package pkg

import (
	"fmt"
	"net/http"

	"github.com/sj14/jellyfin-go/api"
)

func (c *Controller) LibraryScan() (*http.Response, error) {
	return c.client.LibraryAPI.RefreshLibrary(c.ctx).Execute()
}

func (c *Controller) LibraryUnscraped(movies, series, season, episode bool) (*http.Response, error) {
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

	result, resp, err := c.client.ItemsAPI.GetItems(c.ctx).
		Recursive(true).
		IncludeItemTypes(types).
		Filters([]api.ItemFilter{api.ITEMFILTER_IS_NOT_FOLDER}).
		Execute()
	if err != nil {
		return resp, err
	}

	for _, item := range result.Items {
		if !item.ProductionYear.IsSet() {
			fmt.Printf("(%s) %s\n", item.GetType(), item.GetName())
		}
	}
	return resp, err
}
