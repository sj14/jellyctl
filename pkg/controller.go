package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	jellyapi "github.com/sj14/jellyfin-go/api"
)

type Controller struct {
	ctx    context.Context
	client *jellyapi.APIClient
}

func NewController(ctx context.Context, client *jellyapi.APIClient) *Controller {
	return &Controller{
		ctx:    ctx,
		client: client,
	}
}

func printStruct(v any) {
	j, err := json.MarshalIndent(v, "", "")
	if err != nil {
		log.Fatalln(err)
	}

	output := strings.ReplaceAll(string(j), "\"", "")
	output = strings.Trim(output, "{}")
	output = strings.TrimSpace(output)
	fmt.Println(output)
}
