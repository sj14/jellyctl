package controller

import (
	"errors"
	"fmt"
	"time"
)

func (c *Controller) KeyCreate(app string) error {
	if app == "" {
		return errors.New("missing app name")
	}
	_, err := c.client.ApiKeyAPI.CreateKey(c.ctx).App(app).Execute()
	return err
}

func (c *Controller) KeyDelete(key string) error {
	_, err := c.client.ApiKeyAPI.RevokeKey(c.ctx, key).Execute()
	return err
}

func (c *Controller) KeyList(json bool) error {
	result, _, err := c.client.ApiKeyAPI.GetKeys(c.ctx).Execute()
	if err != nil {
		return err
	}

	if json {
		printAsJSON(result)
		return nil
	}

	fmt.Printf("Token                             Created               Name\n")
	fmt.Printf("---------------------------------|---------------------|-----\n")
	for _, key := range result.Items {
		fmt.Printf("%s   %s   %s\n",
			key.GetAccessToken(),
			key.GetDateCreated().Local().Format(time.DateTime),
			key.GetAppName(),
		)
	}

	return nil
}
