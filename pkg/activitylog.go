package pkg

import (
	"fmt"
	"net/http"
	"time"
)

func (c *Controller) GetLogEntries(startIdx, limit int32, minDate time.Time) (*http.Response, error) {
	result, resp, err := c.client.ActivityLogAPI.GetLogEntries(c.ctx).StartIndex(startIdx).MinDate(minDate).Limit(limit).Execute()
	if err != nil {
		return resp, err
	}

	for _, i := range result.Items {
		fmt.Printf("%v [%s] %s\n", i.GetDate().Local().Format(time.DateTime), i.GetSeverity(), i.GetName())
	}
	return resp, err
}
