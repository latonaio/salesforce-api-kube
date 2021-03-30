package main

import (
	"errors"
	"fmt"

	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
)

const msName = "salesforce-api-kube"

func writeKanban(kanbanClient msclient.MicroserviceClient, data map[string]interface{}) error {
	var options []msclient.Option
	options = append(options, msclient.SetMetadata(data))
	options = append(options, msclient.SetProcessNumber(kanbanClient.GetProcessNumber()))
	req, err := msclient.NewOutputData(options...)
	if err != nil {
		return fmt.Errorf("failed to construct output request: %v", err)
	}
	if err := kanbanClient.OutputKanban(req); err != nil {
		return fmt.Errorf("failed to output to kanban: %v", err)
	}
	return nil
}

func buildMetadata(metadata map[string]interface{}, body string) (map[string]interface{}, error) {
	object, ok := metadata["object"]
	if !ok {
		return nil, errors.New("invalid metadata: object not found")
	}
	objectStr, ok := object.(string)
	if !ok {
		return nil, errors.New("failed to convert interface{} to string")
	}
	return map[string]interface{}{"key": objectStr, "content": body}, nil
}
