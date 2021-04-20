package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
	"bitbucket.org/latonaio/aion-core/pkg/log"
	"bitbucket.org/latonaio/salesforce-api-kube/internal/salesforce"
)

func main() {
	// Create Kanban client
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kanbanClient, err := msclient.NewKanbanClient(ctx, msName)
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
	kanbanCh, err := kanbanClient.GetKanbanCh()
	if err != nil {
		log.Fatalf("failed to get kanban channel: %v", err)
	}
	log.Printf("successful get kanban channel\n")

	signalCh := make(chan os.Signal, 1)
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
			limit := make(chan struct{},5)
			go func(k *msclient.WrapKanban) {
				limit <- struct{}{}
				// Get metadata from Kanban
				fromMetadata, err := k.GetMetadataByMap()
				if err != nil {
					log.Printf("failed to get metadata: %v", err)
					return
				}
				log.Printf("got metadata from kanban")
				log.Printf("metadata: %v\n", fromMetadata)

				ck, ok := fromMetadata["connection_key"].(string)
				if !ok {
					log.Printf("invalid connection key")
					return
				}
				// Build http request to salesforce
				req, err := salesforce.BuildRequest(fromMetadata, oauthClient)
				if err != nil {
					log.Printf("failed to build request that send to salesforce api: %v\n", err)
					return
				}

				// Do http request to salesforce
				body, err := salesforce.DoRequest(req)
				if err != nil {
					log.Printf("failed to do request to salesforce api: %v\n", err)
					return
				}
				log.Printf("successfully do http request to salesforce")

				// Build metadata for Kanban
				toMetadata, err := buildMetadata(fromMetadata, body)
				if err != nil {
					log.Printf("failed to build metadata to send: %v", err)
					return
				}

				// Write metadata to Kanban
				if err := writeKanban(kanbanClient, toMetadata, ck); err != nil {
					log.Printf("failed to write kanban: %v", err)
					return
				}
				log.Printf("write metadata to kanban")
				log.Printf("write metadata to kanban: connection_key: %s", ck)
				if _, ok := toMetadata["content"]; ok {
					logMetadata := toMetadata
					logMetadata["content"] = ""
					log.Printf("metadata: %v\n", logMetadata)
				} else {
					log.Printf("metadata: %v\n", toMetadata)
				}
				<-limit
				return
			}(k)

		}
	}
END:
}
