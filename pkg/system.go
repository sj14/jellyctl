package pkg

import (
	"net/http"
)

func (c *Controller) GetSystemInfo() (*http.Response, error) {
	result, resp, err := c.client.SystemAPI.GetSystemInfo(c.ctx).Execute()
	if err != nil {
		return resp, err
	}

	printStruct(result)
	return resp, err
}

func (c *Controller) GetPublicSystemInfo() (*http.Response, error) {
	result, resp, err := c.client.SystemAPI.GetPublicSystemInfo(c.ctx).Execute()
	if err != nil {
		return resp, err
	}

	printStruct(result)
	return resp, err
}

func (c *Controller) SystemShutdown() (*http.Response, error) {
	resp, err := c.client.SystemAPI.ShutdownApplication(c.ctx).Execute()
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c *Controller) SystemRestart() (*http.Response, error) {
	resp, err := c.client.SystemAPI.RestartApplication(c.ctx).Execute()
	if err != nil {
		return resp, err
	}
	return resp, err
}
