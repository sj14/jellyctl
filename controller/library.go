package controller

import (
	"fmt"

	"github.com/sj14/jellyfin-go/api"
)

func (c *Controller) LibraryScan() error {
	_, err := c.client.LibraryAPI.RefreshLibrary(c.ctx).Execute()
	return err
}

func (c *Controller) LibraryUnscraped(movies, series, seasons, episodes bool) error {
	// Determine based on missing production date
	// TODO: look for a better endpoints/approach.

	var types []api.BaseItemKind
	if movies {
		types = append(types, api.BASEITEMKIND_MOVIE)
	}
	if series {
		types = append(types, api.BASEITEMKIND_SERIES)
	}
	if seasons {
		types = append(types, api.BASEITEMKIND_SEASON)
	}
	if episodes {
		types = append(types, api.BASEITEMKIND_EPISODE)
	}

	result, _, err := c.client.ItemsAPI.GetItems(c.ctx).
		Recursive(true).
		IncludeItemTypes(types).
		Filters([]api.ItemFilter{api.ITEMFILTER_IS_NOT_FOLDER}).
		Execute()
	if err != nil {
		return err
	}

	for _, item := range result.Items {
		if !item.ProductionYear.IsSet() {
			fmt.Printf("(%s) %s\n", item.GetType(), item.GetName())
		}
	}
	return err
}

func (c *Controller) LibrarySearch(term string, movies, series, seasons, episodes bool) error {
	var types []api.BaseItemKind
	if movies {
		types = append(types, api.BASEITEMKIND_MOVIE)
	}
	if series {
		types = append(types, api.BASEITEMKIND_SERIES)
	}
	if seasons {
		types = append(types, api.BASEITEMKIND_SEASON)
	}
	if episodes {
		types = append(types, api.BASEITEMKIND_EPISODE)
	}

	results, _, err := c.client.ItemsAPI.GetItems(c.ctx).
		SearchTerm(term).
		IncludeItemTypes(types).
		Recursive(true).
		Execute()
	if err != nil {
		return err
	}

	for _, result := range results.Items {
		fmt.Printf("(%s) %s (%d)\n", result.GetType(), result.GetName(), result.GetProductionYear())
	}
	return nil
}
