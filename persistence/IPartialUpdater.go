package persistence

import "github.com/pip-services3-go/pip-services3-commons-go/data"

/*
  Interface for data processing components to update data items partially.
*/
//<T, K>
type IPartialUpdater interface {

	// Updates only few selected fields in a data item.
	// Parameters:
	//   - correlation_id string
	//   transaction id to trace execution through call chain.
	//   - id interface{}
	//   an id of data item to be updated.
	//   - data data.AnyValueMap
	//   a map with fields to be updated.
	// Returns interface{}, error
	// updated item or error.
	UpdatePartially(correlation_id string, id interface{}, data *data.AnyValueMap) (item interface{}, err error)
}
