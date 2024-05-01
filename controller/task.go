package controller

import (
	"fmt"
	"sort"
	"time"
)

func (c *Controller) TaskList(json bool) error {
	result, _, err := c.client.ScheduledTasksAPI.GetTasks(c.ctx).Execute()
	if err != nil {
		return err
	}

	if json {
		printAsJSON(result)
		return nil
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].GetCategory() < result[j].GetCategory()
	})
	fmt.Printf("ID                               Start Time           State    Category         Name\n")
	fmt.Printf("--------------------------------|--------------------|-------|--------------|---------------------------\n")
	for _, task := range result {
		startTime := task.GetLastExecutionResult().StartTimeUtc
		if startTime == nil {
			startTime = &time.Time{}
		}
		fmt.Printf("%s %s (%s %.f%%) [%s] %s\n",
			task.GetId(),
			startTime.Local().Format(time.DateTime),
			task.GetState(),
			task.GetCurrentProgressPercentage(),
			task.GetCategory(),
			task.GetName(),
		)
	}
	return err
}

func (c *Controller) TaskStart(id string) error {
	_, err := c.client.ScheduledTasksAPI.StartTask(c.ctx, id).Execute()
	return err
}

func (c *Controller) TaskStop(id string) error {
	_, err := c.client.ScheduledTasksAPI.StopTask(c.ctx, id).Execute()
	return err
}
