package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	jellyapi "github.com/sj14/jellyfin-go/api"
)

type Controller struct {
	ctx    context.Context
	client *jellyapi.APIClient
}

func New(ctx context.Context, client *jellyapi.APIClient) *Controller {
	return &Controller{
		ctx:    ctx,
		client: client,
	}
}

func pointer[T any](v T) *T {
	return &v
}

func printAsJSON(v any) {
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(j))
}
