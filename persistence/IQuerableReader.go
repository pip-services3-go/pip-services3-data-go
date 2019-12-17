package persistence

import "github.com/pip-services3-go/pip-services3-commons-go/data"

/*
  Interface for data processing components that can query a list of data items.
*/
type IQuerableReader interface {

	//  Gets a list of data items using a query string.
	//  Prameters:
	//  		- correlation_id  string
	//			transaction id to trace execution through call chain.
	//  		- query string
	//			a query string
	//  		- sort data.SortParams
	//          sort parameters
	// Returns []interface{}, error
	// list of items or error.

	GetListByQuery(correlation_id string, query string, sort data.SortParams) (items []interface{}, err error)
}
