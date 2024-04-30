package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

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
