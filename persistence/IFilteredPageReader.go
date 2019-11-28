package persistence

import "github.com/pip-services3-go/pip-services3-commons-go/data"

/*
  Interface for data processing components that can retrieve a page of data items by a filter.
*/
type IFilteredPageReader interface {

	// Gets a page of data items using filter
	// Parameters
	//   	- correlation_id
	//		transaction id to trace execution through call chain.
	//   	- filter  data.FilterParams
	//		filter parameters
	//   	- paging data.PagingParams
	//      paging parameters
	//   	- sort data.SortParams
	//      sort parameters
	// Retrun  interface{}, error
	// list of items or error.

	getPageByFilter(correlation_id string, filter data.FilterParams, paging data.PagingParams, sort data.SortParams) (page interface{}, err error)
}
