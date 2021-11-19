package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/latonaio/aion-core/proto/kanbanpb"
	"github.com/latonaio/aion-core/pkg/go-client/msclient"
	"github.com/latonaio/aion-core/pkg/log"
	"github.com/latonaio/salesforce-api-kube/internal/salesforce"
)

func main() {
	// Create Kanban client
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kanbanClient, err := msclient.NewKanbanClient(ctx, msName, kanbanpb.InitializeType_START_SERVICE)
	if err != nil {
		log.Fatalf("failed to get kanban client: %v", err)
	}
	log.Printf("successful get kanban client")
	defer kanbanClient.Close()

	// Create Salesforce client
	oauthClient, err := salesforce.NewOAuthClient()
	if err != nil {
		log.Fatalf("failed to construct OAuthClient: %v", err)
	}
	log.Printf("successful construct OAuth client")

	// Get Kanban channel by Kanban client
	kanbanCh := kanbanClient.GetKanbanCh()
	log.Printf("successful get kanban channel\n")

	signalCh := make(chan os.Signal, 1)
	limit := make(chan struct{},5)
	wg := new(sync.WaitGroup)
	signal.Notify(signalCh, syscall.SIGTERM)
	for {
		select {
		case s := <-signalCh:
			fmt.Printf("received signal: %s", s.String())
			goto END
		case k := <-kanbanCh:
			if k == nil {
				continue
			}
			limit <- struct{}{}
			wg.Add(1)
			go func(k *kanbanpb.StatusKanban) {
				defer func() {
					<-limit
					wg.Done()
				}()
				// Get metadata from Kanban
				fromMetadata, err := msclient.GetMetadataByMap(k)
				if err != nil {
					log.Errorf("failed to get metadata: %v", err)
					return
				}
				log.Printf("got metadata from kanban")
				log.Debugf("metadata: %v\n", fromMetadata)

				ck, ok := fromMetadata["connection_key"].(string)
				if !ok {
					log.Errorf("invalid connection key")
					return
				}
				// Build http request to salesforce
				req, err := salesforce.BuildRequest(fromMetadata, oauthClient)
				if err != nil {
					log.Errorf("failed to build request that send to salesforce api: %v\n", err)
					return
				}

				// Do http request to salesforce
				body, err := salesforce.DoRequest(req)
				if err != nil {
					log.Errorf("failed to do request to salesforce api: %v\n", err)
					return
				}
				log.Printf("successfully do http request to salesforce")

				// Build metadata for Kanban
				toMetadata, err := buildMetadata(fromMetadata, body)
				if err != nil {
					log.Errorf("failed to build metadata to send: %v", err)
					return
				}

				// Write metadata to Kanban
				if err := writeKanban(kanbanClient, toMetadata, ck); err != nil {
					log.Errorf("failed to write kanban: %v", err)
					return
				}
				log.Printf("write metadata to kanban")
				log.Printf("write metadata to kanban: connection_key: %s", ck)
				if _, ok := toMetadata["content"]; ok {
					logMetadata := toMetadata
					logMetadata["content"] = ""
					log.Debugf("metadata: %v\n", logMetadata)
				} else {
					log.Debugf("metadata: %v\n", toMetadata)
				}
			}(k)
		}
	}
END:
	wg.Wait()
}
