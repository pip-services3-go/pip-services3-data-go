package persistence

import "github.com/pip-services3-go/pip-services3-commons-go/data"

/*
  Interface for data processing components that can retrieve a list of data items by filter.
*/
type IFilteredReader interface {

	// Gets a list of data items using filter
	// Parameters:
	// 			- correlation_id string
	//		    transaction id to trace execution through call chain.
	// 			- filter data.FilterParams
	//          filter parameters
	// 			- sort  data.SortParams
	//			sort parameters
	// Returns []interfcace{}, error
	// receives list of items or error.

	getListByFilter(correlation_id string, filter data.FilterParams, sort data.SortParams) (items []interface{}, err error)
}
