package persistence

/*
  Interface for data processing components that can get data items.
*/
// <T extends IIdentifiable<K>, K>
type IGetter interface {

	//  Gets a data items by its unique id.
	//  Parameters:
	//  - correlation_id    (optional) transaction id to trace execution through call chain.
	//  - id                an id of item to be retrieved.
	//  Return interface{}, error
	// item or error
	GetOneById(correlation_id string, id interface{}) (item interface{}, err error)
}
