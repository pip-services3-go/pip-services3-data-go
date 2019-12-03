package persistence

import (
	"io/ioutil"
	"os"

	"github.com/pip-services3-go/pip-services3-commons-go/config"
	"github.com/pip-services3-go/pip-services3-commons-go/convert"
	"github.com/pip-services3-go/pip-services3-commons-go/errors"
)

/*
Persistence component that loads and saves data from/to flat file.

It is used by FilePersistence, but can be useful on its own.

 Configuration parameters

- path:          path to the file where data is stored

 Example

    persister := NewJsonFilePersister("./data/data.json");

	err_sav := persister.Save("123", ["A", "B", "C"])
	if err_sav == nil {
		items, err_lod := persister.Load("123")
		if err_lod == nil {
			fmt.Println(items);// Result: ["A", "B", "C"]
		}
*/
// implements ILoader, ISaver, IConfigurable
type JsonFilePersister struct {
	_path string
}

// Creates a new instance of the persistence.
// Parameters:
//		- path  string
//		(optional) a path to the file where data is stored.
func NewJsonFilePersister(path string) (jfp *JsonFilePersister) {
	jfp = &JsonFilePersister{_path: path}
	return jfp
}

// Gets the file path where data is stored.
// Returns the file path where data is stored.
func (jfp *JsonFilePersister) Path() string {
	return jfp._path
}

// Sets the file path where data is stored.
// Parameters:
//		- value  string
//	    the file path where data is stored.
func (jfp *JsonFilePersister) SetPath(value string) {
	jfp._path = value
}

// Configures component by passing configuration parameters.
// Parameters:
//		- config  config.ConfigParams
//		parameters to be set.
func (jfp *JsonFilePersister) Configure(config config.ConfigParams) {
	jfp._path = config.GetAsStringWithDefault("path", jfp._path)
}

// Loads data items from external JSON file.
//		- correlation_id  string
//		transaction id to trace execution through call chain.
// Returns []interface{}, error
// loaded items or error.
func (jfp *JsonFilePersister) Load(correlation_id string) (data []interface{}, err error) {
	if jfp._path == "" {
		data = nil
		err = errors.NewConfigError("", "NO_PATH", "Data file path is not set")
		return
	}

	_, fserr := os.Stat(jfp._path)
	if os.IsNotExist(fserr) {
		data = nil
		err = nil
		return
	}

	json, jsonerr := ioutil.ReadFile(jfp._path)
	if jsonerr != nil {
		err = errors.NewFileError(correlation_id, "READ_FAILED", "Failed to read data file: "+jfp._path).WithCause(jsonerr)
		data = nil
		return
	}
	list := convert.JsonConverter.ToNullableMap((string)(json))
	data = convert.ArrayConverter.ListToArray(list)
	err = nil
	return
}

//Saves given data items to external JSON file.
// Parameters:
//		- correlation_id string
//	    transaction id to trace execution through call chain.
// 		- items []interface[]
//      list if data items to save
//  Retruns error
//  error or nil for success.
func (jfp *JsonFilePersister) Save(correlation_id string, items []interface{}) error {
	json, jsonerr := convert.ToJson(items)
	if jsonerr != nil {
		err := errors.NewInternalError(correlation_id, "CANT CONVERT", "Failed convert to JSON")
		return err
	}
	werr := ioutil.WriteFile(jfp._path, ([]byte)(json), 0777)
	if werr != nil {
		err := errors.NewFileError(correlation_id, "WRITE_FAILED", "Failed to write data file: "+jfp._path).WithCause(werr)
		return err
	}
	return nil
}
