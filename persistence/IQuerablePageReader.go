package persistence

import "github.com/pip-services3-go/pip-services3-commons-go/data"

/*
  Interface for data processing components that can query a page of data items.
*/
type IQuerablePageReader interface {

	//  Gets a page of data items using a query string.
	//  Parameters:
	//  		- correlation_id string
	//			transaction id to trace execution through call chain.
	//  		- query string
	//			 a query string
	//  		- paging data.PagingParams
	//			 paging parameters
	//   		- sort  data.SortParams
	//			 sort parameters
	// Returns interface{}, error
	// receives list of items or error.

	GetPageByQuery(correlation_id string, query string, paging data.PagingParams, sort data.SortParams) (page interface{}, err error)
}
