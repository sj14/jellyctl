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
	// Determine based on missing production date or missing community rating
	// TODO: look for a better endpoints/approach.

	var t []api.BaseItemKind
	for _, ty := range types {
		t = append(t, api.BaseItemKind(ty))
	}

	allItems, _, err := c.client.ItemsAPI.GetItems(c.ctx).
		Recursive(true).
		IncludeItemTypes(t).
		Filters([]api.ItemFilter{api.ITEMFILTER_IS_NOT_FOLDER}).
		Execute()
	if err != nil {
		return err
	}

	var jsonResult []api.BaseItemDto
	for _, item := range allItems.Items {
		if !item.ProductionYear.IsSet() || !item.CommunityRating.IsSet() {
			if json {
				jsonResult = append(jsonResult, item)
				continue
			}
			fmt.Printf("(%s) %s\n", item.GetType(), item.GetName())
		}
	}

	if json {
		printAsJSON(jsonResult)
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

func (c *Controller) LibraryDuplicates(term string, types []string, json bool) error {
	var t []api.BaseItemKind
	for _, ty := range types {
		t = append(t, api.BaseItemKind(ty))
	}

	response, _, err := c.client.ItemsAPI.GetItems(c.ctx).
		SearchTerm(term).
		IncludeItemTypes(t).
		Recursive(true).
		Execute()
	if err != nil {
		return err
	}

	var result []api.BaseItemDto

	for idx1, item1 := range response.Items {
		for idx2, item2 := range response.Items {
			if idx1 == idx2 {
				// Do not compare exactly the same entry.
				continue
			}

			// Check if we have a duplicate.
			if item1.GetName() == item2.GetName() && item1.GetProductionYear() == item2.GetProductionYear() {
				found := false
				for _, r := range result {
					// Only add the duplicate to the result list when the same entry is not already listed.
					if r.GetName() == item1.GetName() && r.GetProductionYear() == item1.GetProductionYear() {
						found = true
						break
					}
				}

				if !found {
					result = append(result, item1)
				}
			}
		}
	}

	if json {
		printAsJSON(result)
		return nil
	}

	for _, item := range result {
		fmt.Printf("(%s) %s (%d)\n", item.GetType(), item.GetName(), item.GetProductionYear())
	}

	return nil
}
