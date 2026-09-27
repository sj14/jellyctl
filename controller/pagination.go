package controller

import (
	"fmt"
	"net/http"

	"github.com/sj14/jellyfin-go/api"
)

const itemPageSize int32 = 100

type itemPageRequest[T any] interface {
	StartIndex(int32) T
	Limit(int32) T
	Execute() (*api.BaseItemDtoQueryResult, *http.Response, error)
}

func getAllRequestItems[T itemPageRequest[T]](request T) (*api.BaseItemDtoQueryResult, error) {
	return getAllItems(func(startIndex, limit int32) (*api.BaseItemDtoQueryResult, error) {
		result, _, err := request.StartIndex(startIndex).Limit(limit).Execute()
		return result, err
	})
}

func getAllItems(fetch func(int32, int32) (*api.BaseItemDtoQueryResult, error)) (*api.BaseItemDtoQueryResult, error) {
	all := &api.BaseItemDtoQueryResult{Items: []api.BaseItemDto{}}
	for startIndex := int32(0); ; {
		page, err := fetch(startIndex, itemPageSize)
		if err != nil {
			return nil, err
		}
		if page == nil {
			return nil, fmt.Errorf("items page at index %d has no result", startIndex)
		}
		if startIndex == 0 {
			all.SetStartIndex(0)
		}
		if page.TotalRecordCount != nil {
			all.TotalRecordCount = page.TotalRecordCount
		}
		if len(page.Items) == 0 {
			if total, ok := all.GetTotalRecordCountOk(); ok && startIndex < *total {
				return nil, fmt.Errorf("items page at index %d is empty before total count %d", startIndex, *total)
			}
			break
		}

		all.Items = append(all.Items, page.Items...)
		startIndex += int32(len(page.Items))
		if total, ok := all.GetTotalRecordCountOk(); ok {
			if startIndex >= *total {
				break
			}
		}
	}
	if all.TotalRecordCount == nil {
		all.SetTotalRecordCount(int32(len(all.Items)))
	}
	return all, nil
}
