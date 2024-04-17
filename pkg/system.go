package pkg

import (
	"fmt"
	"net/http"
)

func (c *Controller) GetPublicSystemInfo() (*http.Response, error) {
	result, resp, err := c.client.SystemAPI.GetPublicSystemInfo(c.ctx).Execute()
	if err != nil {
		return resp, err
	}

	fmt.Printf("%s %s\n", result.GetProductName(), result.GetVersion())
	fmt.Printf("Server: %s (%s)\n", result.GetServerName(), result.GetOperatingSystem())
	fmt.Printf("Setup wizard completed: %v\n", result.GetStartupWizardCompleted())
	return resp, err
}

func (c *Controller) Restart() (*http.Response, error) {
	resp, err := c.client.SystemAPI.RestartApplication(c.ctx).Execute()
	if err != nil {
		return resp, err
	}
	fmt.Println("Jellyfin restart executed")
	return resp, err
}
