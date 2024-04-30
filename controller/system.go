package controller

func (c *Controller) GetSystemInfo() error {
	result, _, err := c.client.SystemAPI.GetSystemInfo(c.ctx).Execute()
	if err != nil {
		return err
	}

	printAsJSON(result)
	return err
}

func (c *Controller) GetPublicSystemInfo() error {
	result, _, err := c.client.SystemAPI.GetPublicSystemInfo(c.ctx).Execute()
	if err != nil {
		return err
	}

	printAsJSON(result)
	return err
}

func (c *Controller) SystemShutdown() error {
	_, err := c.client.SystemAPI.ShutdownApplication(c.ctx).Execute()
	return err
}

func (c *Controller) SystemRestart() error {
	_, err := c.client.SystemAPI.RestartApplication(c.ctx).Execute()
	return err
}
