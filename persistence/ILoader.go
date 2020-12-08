package persistence

/*
  Interface for data processing components that load data items.
*/
type ILoader interface {

	// Loads data items.
	// Parameters:
	//   - correlation_id string
	//   transaction id to trace execution through call chain.
	// Retruns []interface{}, error
	// a list of data items or error.
	Load(correlation_id string) (items []interface{}, err error)
}
