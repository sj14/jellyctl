package pkg

import (
	"net/http"
)

func (c *Controller) LibraryScan() (*http.Response, error) {
	return c.client.LibraryAPI.RefreshLibrary(c.ctx).Execute()
}
