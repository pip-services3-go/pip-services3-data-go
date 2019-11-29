package persistence

/*
  Interface for data processing components that can set (create or update) data items.
*/
type ISetter interface {

	// Sets a data item. If the data item exists it updates it,
	// otherwise it create a new data item.
	// Parameters:
	// 		- correlation_id string
	//	    transaction id to trace execution through call chain.
	// 		- item  interface{}
	//      a item to be set.
	// Retruns updated item or error.
	Set(correlation_id string, item interface{}) (value interface{}, err error)
}
