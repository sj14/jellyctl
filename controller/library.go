package controller

import (
	"fmt"

	"github.com/sj14/jellyfin-go/api"
)

func (c *Controller) LibraryScan() error {
	_, err := c.client.LibraryAPI.RefreshLibrary(c.ctx).Execute()
	return err
}

func (c *Controller) LibraryUnscraped(movies, series, season, episode bool) error {
	// Determine based on missing production date
	// TODO: look for a better endpoints/approach.

	var types []api.BaseItemKind
	if movies {
		types = append(types, api.BASEITEMKIND_MOVIE)
	}
	if series {
		types = append(types, api.BASEITEMKIND_SERIES)
	}
	if season {
		types = append(types, api.BASEITEMKIND_SEASON)
	}
	if episode {
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
