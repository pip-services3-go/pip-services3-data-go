package persistence

/*
  Interface for data processing components that can create, update and delete data items.
*/
type IWriter interface {

	//  Creates a data item.
	// Parameters:
	//		- correlation_id string
	//    	transaction id to trace execution through call chain.
	// 		- item interface{}
	//      an item to be created.
	// Returns  interface{}, error
	// created item or error.

	Create(correlation_id string, item interface{}) (value *interface{}, err error)

	// Updates a data item.
	// Parameters:
	// 		- correlation_id  string
	//	    transaction id to trace execution through call chain.
	// 		- item interface{}
	//      an item to be updated.
	// Returns: interface{}, error
	// updated item or error.

	Update(correlation_id string, item interface{}) (value *interface{}, err error)

	//  Deleted a data item by it's unique id.
	//	Parameters:
	//  		- correlation_id string
	//		    transaction id to trace execution through call chain.
	//  		- id interface{}
	//			 an id of the item to be deleted
	//  Returns: interface{}, error
	//  deleted item or error.

	DeleteById(correlation_id string, id interface{}) (value *interface{}, err error)
}
