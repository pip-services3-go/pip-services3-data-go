package persistence

/*
  Interface for data processing components that save data items.
*/
type ISaver interface {

	// Saves given data items.
	// Parameters:
	//		- correlation_id string
	//    	transaction id to trace execution through call chain.
	// 		- items []interface{}
	//      a list of items to save.
	// Retuirns error or nil for success.

	Save(correlation_id string, items []interface{}) error
}
