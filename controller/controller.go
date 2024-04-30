package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

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

// Set all fields which are nil to the zero value.
// https://stackoverflow.com/a/61743435/7125878
// TODO: test it, check if it needs to be called recursively on structs
func setZeroValue(v any) {
	rv := reflect.ValueOf(v).Elem()
	for i := 0; i < rv.NumField(); i++ {
		if f := rv.Field(i); !f.IsValid() && f.IsNil() && f.CanSet() {
			f.SetZero()
		}
	}
}
