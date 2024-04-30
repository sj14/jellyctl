package controller

import (
	"fmt"
	"time"
)

func (c *Controller) KeyCreate(app string) error {
	_, err := c.client.ApiKeyAPI.CreateKey(c.ctx).App(app).Execute()
	return err
}

func (c *Controller) KeyDelete(key string) error {
	_, err := c.client.ApiKeyAPI.RevokeKey(c.ctx, key).Execute()
	return err
}

func (c *Controller) KeyList() error {
	result, _, err := c.client.ApiKeyAPI.GetKeys(c.ctx).Execute()
	if err != nil {
		return err
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
