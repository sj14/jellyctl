package controller

import (
	"fmt"

	"github.com/sj14/jellyfin-go/api"
)

func (c *Controller) LibraryScan() error {
	_, err := c.client.LibraryAPI.RefreshLibrary(c.ctx).Execute()
	return err
}

func (c *Controller) LibraryUnscraped(types []string, json bool) error {
	// Determine based on missing production date
	// TODO: look for a better endpoints/approach.

	var t []api.BaseItemKind
	for _, ty := range types {
		t = append(t, api.BaseItemKind(ty))
	}

	result, _, err := c.client.ItemsAPI.GetItems(c.ctx).
		Recursive(true).
		IncludeItemTypes(t).
		Filters([]api.ItemFilter{api.ITEMFILTER_IS_NOT_FOLDER}).
		Execute()
	if err != nil {
		return err
	}

	if json {
		printAsJSON(result)
		return nil
	}

	for _, item := range result.Items {
		if !item.ProductionYear.IsSet() {
			fmt.Printf("(%s) %s\n", item.GetType(), item.GetName())
		}
	}
	return err
}

func (c *Controller) LibrarySearch(term string, types []string, json bool) error {
	var t []api.BaseItemKind
	for _, ty := range types {
		t = append(t, api.BaseItemKind(ty))
	}

	results, _, err := c.client.ItemsAPI.GetItems(c.ctx).
		SearchTerm(term).
		IncludeItemTypes(t).
		Recursive(true).
		Execute()
	if err != nil {
		return err
	}

	if json {
		printAsJSON(results)
		return nil
	}

	for _, result := range results.Items {
		fmt.Printf("(%s) %s (%d)\n", result.GetType(), result.GetName(), result.GetProductionYear())
	}
	return nil
}
