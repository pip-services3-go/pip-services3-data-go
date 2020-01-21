package persistence

import (
	"io/ioutil"
	"os"
	"reflect"

	"github.com/pip-services3-go/pip-services3-commons-go/v3/config"
	"github.com/pip-services3-go/pip-services3-commons-go/v3/convert"
	"github.com/pip-services3-go/pip-services3-commons-go/v3/errors"
)

/*
Persistence component that loads and saves data from/to flat file.

It is used by FilePersistence, but can be useful on its own.

 Configuration parameters

- path:          path to the file where data is stored

 Example

    persister := NewJsonFilePersister(reflect.TypeOf(MyData{}), "./data/data.json");

	err_sav := persister.Save("123", ["A", "B", "C"])
	if err_sav == nil {
		items, err_lod := persister.Load("123")
		if err_lod == nil {
			fmt.Println(items);// Result: ["A", "B", "C"]
		}
*/
// implements ILoader, ISaver, IConfigurable
type JsonFilePersister struct {
	path      string
	Prototype reflect.Type
}

// Creates a new instance of the persistence.
// Parameters:
//		- path  string
//		(optional) a path to the file where data is stored.
func NewJsonFilePersister(prototype reflect.Type, path string) *JsonFilePersister {
	var c = &JsonFilePersister{path: path, Prototype: prototype}
	return c
}

// Gets the file path where data is stored.
// Returns the file path where data is stored.
func (c *JsonFilePersister) Path() string {
	return c.path
}

// Sets the file path where data is stored.
// Parameters:
//		- value  string
//	    the file path where data is stored.
func (c *JsonFilePersister) SetPath(value string) {
	c.path = value
}

// Configures component by passing configuration parameters.
// Parameters:
//		- config  config.ConfigParams
//		parameters to be set.
func (c *JsonFilePersister) Configure(config *config.ConfigParams) {
	c.path = config.GetAsStringWithDefault("path", c.path)
}

// Loads data items from external JSON file.
// Parameters:
//		- correlation_id  string
//		transaction id to trace execution through call chain.
// Returns []interface{}, error
// loaded items or error.
func (c *JsonFilePersister) Load(correlation_id string) (data []interface{}, err error) {
	if c.path == "" {
		data = nil
		err = errors.NewConfigError("", "NO_PATH", "Data file path is not set")
		return data, err
	}

	_, fserr := os.Stat(c.path)
	if os.IsNotExist(fserr) {
		data = nil
		err = nil
		return data, err
	}

	jsonStr, jsonerr := ioutil.ReadFile(c.path)
	if jsonerr != nil {
		err = errors.NewFileError(correlation_id, "READ_FAILED", "Failed to read data file: "+c.path).WithCause(jsonerr)
		data = nil
		return data, err
	}

	if len(jsonStr) == 0 {
		return nil, nil
	}

	list, err := convert.FromJson((string)(jsonStr))
	if list == nil {
		data = nil
		return data, err
	}
	data = convert.ArrayConverter.ListToArray(list)
	err = nil
	return data, err
}

// Saves given data items to external JSON file.
// Parameters:
//		- correlation_id string
//	    transaction id to trace execution through call chain.
// 		- items []interface[]
//      list of data items to save
//  Retruns error
//  error or nil for success.
func (c *JsonFilePersister) Save(correlationId string, items []interface{}) error {
	json, jsonerr := convert.ToJson(items)
	if jsonerr != nil {
		err := errors.NewInternalError(correlationId, "CAN'T_CONVERT", "Failed convert to JSON")
		return err
	}
	werr := ioutil.WriteFile(c.path, ([]byte)(json), 0777)
	if werr != nil {
		err := errors.NewFileError(correlationId, "WRITE_FAILED", "Failed to write data file: "+c.path).WithCause(werr)
		return err
	}
	return nil
}
