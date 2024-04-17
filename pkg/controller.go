package pkg

import (
	"context"

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
