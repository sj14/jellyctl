package controller

import (
	"fmt"
	"time"
)

func (c *Controller) GetLogEntries(startIdx, limit int32, minDate time.Time, json bool) error {
	result, _, err := c.client.ActivityLogAPI.GetLogEntries(c.ctx).
		StartIndex(startIdx).
		MinDate(minDate).
		Limit(limit).
		Execute()
	if err != nil {
		return err
	}

	if json {
		printAsJSON(result)
		return nil
	}

	for _, i := range result.Items {
		fmt.Printf("%v [%s] %s\n",
			i.GetDate().Local().Format(time.DateTime),
			i.GetSeverity(),
			i.GetName())
	}
	return err
}
