package persistence

import (
	"github.com/pip-services3-go/pip-services3-commons-go/v3/config"
	"reflect"
)

/*
Abstract persistence component that stores data in flat files
and implements a number of CRUD operations over data items with unique ids.
The data items must implement
 IIdentifiable interface

In basic scenarios child classes shall only override GetPageByFilter,
GetListByFilter or DeleteByFilter operations with specific filter function.
All other operations can be used out of the box.

In complex scenarios child classes can implement additional operations by
accessing cached items via IdentifiableFilePersistence._items property and calling Save method
on updates.

See JsonFilePersister
See MemoryPersistence

Configuration parameters

- path:                    path to the file where data is stored
- options:
    - max_page_size:       Maximum number of items returned in a single page (default: 100)

 References

- *:logger:*:*:1.0      (optional)  ILogger components to pass log messages

 Examples
type MyFilePersistence  struct {
	IdentifiableFilePersistence
}
    func NewMyFilePersistence(path string)(mfp *MyFilePersistence) {
		mfp = MyFilePersistence{}
		prototype := reflect.TypeOf(MyData{})
		mfp.IdentifiableFilePersistence = *NewJsonPersister(prototype,path)
		return mfp
    }

    func composeFilter(filter cdata.FilterParams)(func (item interface{})bool) {
		if &filter == nil {
			filter = NewFilterParams()
		}
        name := filter.GetAsNullableString("name");
        return func (item interface) bool {
            dummy, ok := item.(MyData)
			if *name != "" && ok && dummy.Name != *name {
				return false
			}
            return true
        }
    }

    func (c *MyFilePersistence ) GetPageByFilter(correlationId string, filter FilterParams, paging PagingParams)(page cdata.MyDataPage, err error){
		tempPage, err := c.GetPageByFilter(correlationId, composeFilter(filter), paging, nil, nil)
		dataLen := int64(len(tempPage.Data))
		data := make([]MyData, dataLen)
		for i, v := range tempPage.Data {
			data[i] = v.(MyData)
		}
		page = *NewMyDataPage(&dataLen, data)
		return page, err
    }

    persistence := NewMyFilePersistence("./data/data.json")

	_, errc := persistence.Create("123", { Id: "1", Name: "ABC" })
	if (errc != nil) {
		panic()
	}
    page, errg := persistence.GetPageByFilter("123", NewFilterParamsFromTuples("Name", "ABC"), nil)
    if errg != nil {
		panic("Error")
	}
    fmt.Println(page.Data)         // Result: { Id: "1", Name: "ABC" )
    persistence.DeleteById("123", "1")
*/
type IdentifiableFilePersistence struct {
	IdentifiableMemoryPersistence
	Persister *JsonFilePersister
}

// Creates a new instance of the persistence.
// Parameters:
// 		- prototype reflect.Type
//    	type of contained data
//  	- persister    (optional) a persister component that loads and saves data from/to flat file.
// Return *IdentifiableFilePersistence
// pointer on new IdentifiableFilePersistence
func NewIdentifiableFilePersistence(prototype reflect.Type, persister *JsonFilePersister) *IdentifiableFilePersistence {
	c := &IdentifiableFilePersistence{}
	if persister == nil {
		persister = NewJsonFilePersister(prototype, "")
	}
	c.IdentifiableMemoryPersistence = *NewIdentifiableMemoryPersistence(prototype)
	c.Loader = persister
	c.Saver = persister
	c.Persister = persister
	return c
}

// Configures component by passing configuration parameters.
// Parameters:
// 		- config    configuration parameters to be set.
func (c *IdentifiableFilePersistence) Configure(config config.ConfigParams) {
	c.Configure(config)
	c.Persister.Configure(config)
}
