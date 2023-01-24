package persistence

import (
	"reflect"

	"github.com/pip-services3-go/pip-services3-commons-go/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	refl "github.com/pip-services3-go/pip-services3-commons-go/reflect"
	"github.com/pip-services3-go/pip-services3-components-go/log"
)

/*
Abstract persistence component that stores data in memory
and implements a number of CRUD operations over data items with unique ids.
The data items must have Id field.

In basic scenarios child structs shall only override GetPageByFilter,
GetListByFilter or DeleteByFilter operations with specific filter function.
All other operations can be used out of the box.

In complex scenarios child structes can implement additional operations by
accessing cached items via c.Items property and calling Save method
on updates.

See MemoryPersistence

Configuration parameters

- options:
    - max_page_size:       Maximum number of items returned in a single page (default: 100)

 References

- *:logger:*:*:1.0     (optional) ILogger components to pass log messages

 Examples

  type MyMemoryPersistence struct{
  	IdentifiableMemoryPersistence
  }
      func composeFilter(filter: FilterParams) (func (item interface{}) bool ) {
          if filter == nil {
  			filter = NewFilterParams()
  		}
          name := filter.getAsNullableString("Name");
          return func(item interface{}) bool {
  			dummy, ok := item.(MyData)
              if (*name != "" && ok && item.Name != *name)
                  return false;
              return true;
          };
      }

      func (mmp * MyMemoryPersistence) GetPageByFilter(correlationId string, filter FilterParams, paging PagingParams) (page DataPage, err error) {
          tempPage, err := c.GetPageByFilter(correlationId, composeFilter(filter), paging, nil, nil)
  		dataLen := int64(len(tempPage.Data))
  		data := make([]MyData, dataLen)
  		for i, v := range tempPage.Data {
  			data[i] = v.(MyData)
  		}
  		page = *NewMyDataPage(&dataLen, data)
  		return page, err}

      persistence := NewMyMemoryPersistence();

  	item, err := persistence.Create("123", { Id: "1", Name: "ABC" })
  	...
  	page, err := persistence.GetPageByFilter("123", NewFilterParamsFromTuples("Name", "ABC"), nil)
  	if err != nil {
  		panic("Error can't get data")
  	}
      fmt.Prnitln(page.data)         // Result: { Id: "1", Name: "ABC" }
  	item, err := persistence.DeleteById("123", "1")
  	...

*/
// extends MemoryPersistence  implements IConfigurable, IWriter, IGetter, ISetter
type IdentifiableMemoryPersistence struct {
	MemoryPersistence
}

// Creates a new empty instance of the persistence.
// Parameters:
//  - prototype reflect.Type
//  data type of contains items
// Return * IdentifiableMemoryPersistence
// created empty IdentifiableMemoryPersistence
func NewIdentifiableMemoryPersistence(prototype reflect.Type) (c *IdentifiableMemoryPersistence) {
	c = &IdentifiableMemoryPersistence{}
	c.MemoryPersistence = *NewMemoryPersistence(prototype)
	c.Logger = log.NewCompositeLogger()
	c.MaxPageSize = 100
	return c
}

// Configures component by passing configuration parameters.
// Parameters:
//  - config  *config.ConfigParams
//  configuration parameters to be set.
func (c *IdentifiableMemoryPersistence) Configure(config *config.ConfigParams) {
	c.MaxPageSize = config.GetAsIntegerWithDefault("options.max_page_size", c.MaxPageSize)
}

// Gets a list of data items retrieved by given unique ids.
// Parameters:
//   - correlationId string
//   (optional) transaction id to trace execution through call chain.
//   - ids  []interface{}
//   ids of data items to be retrieved
// Returns  []interface{}, error
// data list or error.
func (c *IdentifiableMemoryPersistence) GetListByIds(correlationId string, ids []interface{}) (result []interface{}, err error) {
	filter := func(item interface{}) bool {
		exist := false
		id := GetObjectId(item)
		for _, v := range ids {
			vId := refl.ObjectReader.GetValue(v)
			if CompareValues(id, vId) {
				exist = true
				break
			}
		}
		return exist
	}
	return c.GetListByFilter(correlationId, filter, nil, nil)
}

// Gets a data item by its unique id.
// Parameters:
//   - correlationId  string
//   (optional) transaction id to trace execution through call chain.
//   - id interface{}
//   an id of data item to be retrieved.
// Returns:  interface{}, error
// data item or error.
func (c *IdentifiableMemoryPersistence) GetOneById(correlationId string, id interface{}) (result interface{}, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	var items []interface{}
	for _, v := range c.Items {
		vId := GetObjectId(v)
		if CompareValues(vId, id) {
			items = append(items, v)
		}
	}

	var item interface{} = nil
	if len(items) > 0 {
		//item = CloneObject(items[0])
		item = CloneObjectForResult(items[0], c.Prototype)
	}
	if item != nil {
		c.Logger.Trace(correlationId, "Retrieved item %s", id)
	} else {
		c.Logger.Trace(correlationId, "Cannot find item by %s", id)
	}
	return item, err
}

// Get index by "Id" field
// return index number
func (c *IdentifiableMemoryPersistence) GetIndexById(id interface{}) int {
	var index int = -1
	for i, v := range c.Items {
		vId := GetObjectId(v)
		if CompareValues(vId, id) {
			index = i
			break
		}
	}
	return index
}

// Creates a data item.
// Returns:
//   - correlation_id string
//   (optional) transaction id to trace execution through call chain.
//   - item  string
//   an item to be created.
// Returns:  interface{}, error
// created item or error.
func (c *IdentifiableMemoryPersistence) Create(correlationId string, item interface{}) (result interface{}, err error) {
	c.Lock.Lock()

	newItem := CloneObject(item, c.Prototype)
	GenerateObjectId(&newItem)
	id := GetObjectId(newItem)
	c.Items = append(c.Items, newItem)

	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Created item %s", id)

	errsave := c.Save(correlationId)
	result = CloneObjectForResult(newItem, c.Prototype)

	return result, errsave
}

// Sets a data item. If the data item exists it updates it,
// otherwise it create a new data item.
// Parameters:
//   - correlation_id string
//   (optional) transaction id to trace execution through call chain.
//   - item  interface{}
//   a item to be set.
// Returns:  interface{}, error
// updated item or error.
func (c *IdentifiableMemoryPersistence) Set(correlationId string, item interface{}) (result interface{}, err error) {
	c.Lock.Lock()

	newItem := CloneObject(item, c.Prototype)
	GenerateObjectId(&newItem)

	id := GetObjectId(item)
	index := c.GetIndexById(id)
	if index < 0 {
		c.Items = append(c.Items, newItem)
	} else {
		c.Items[index] = newItem
	}

	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Set item %s", id)

	errsav := c.Save(correlationId)

	result = CloneObjectForResult(newItem, c.Prototype)
	return result, errsav
}

// Updates a data item.
// Parameters:
//   - correlation_id string
//   (optional) transaction id to trace execution through call chain.
//   - item  interface{}
//   an item to be updated.
// Returns:   interface{}, error
// updated item or error.
func (c *IdentifiableMemoryPersistence) Update(correlationId string, item interface{}) (result interface{}, err error) {
	c.Lock.Lock()

	id := GetObjectId(item)
	index := c.GetIndexById(id)
	if index < 0 {
		c.Logger.Trace(correlationId, "Item %s was not found", id)
		c.Lock.Unlock()
		return nil, nil
	}
	newItem := CloneObject(item, c.Prototype)
	c.Items[index] = newItem

	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Updated item %s", id)

	errsave := c.Save(correlationId)

	result = CloneObjectForResult(newItem, c.Prototype)
	return result, errsave
}

// Updates only few selectFuncected fields in a data item.
// Parameters:
//   - correlation_id string
//   (optional) transaction id to trace execution through call chain.
//   - id interface{}
//   an id of data item to be updated.
//   - data  cdata.AnyValueMap
//   a map with fields to be updated.
// Returns: interface{}, error
// updated item or error.
func (c *IdentifiableMemoryPersistence) UpdatePartially(correlationId string, id interface{}, data *cdata.AnyValueMap) (result interface{}, err error) {
	c.Lock.Lock()

	index := c.GetIndexById(id)
	if index < 0 {
		c.Logger.Trace(correlationId, "Item %s was not found", id)
		c.Lock.Unlock()
		return nil, nil
	}

	newItem := CloneObject(c.Items[index], c.Prototype)

	if reflect.ValueOf(newItem).Kind() == reflect.Map {
		refl.ObjectWriter.SetProperties(newItem, data.Value())
	} else {
		objPointer := reflect.New(reflect.TypeOf(newItem))
		objPointer.Elem().Set(reflect.ValueOf(newItem))
		intPointer := objPointer.Interface()
		refl.ObjectWriter.SetProperties(intPointer, data.Value())
		newItem = reflect.ValueOf(intPointer).Elem().Interface()
	}

	c.Items[index] = newItem

	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Partially updated item %s", id)

	errsave := c.Save(correlationId)

	result = CloneObjectForResult(newItem, c.Prototype)
	return result, errsave
}

// Deleted a data item by it's unique id.
// Parameters:
//   - correlation_id string
//   (optional) transaction id to trace execution through call chain.
//   - id interface{}
//   an id of the item to be deleted
// Retruns:  interface{}, error
// deleted item or error.
func (c *IdentifiableMemoryPersistence) DeleteById(correlationId string, id interface{}) (result interface{}, err error) {
	c.Lock.Lock()

	index := c.GetIndexById(id)
	if index < 0 {
		c.Logger.Trace(correlationId, "Item %s was not found", id)
		c.Lock.Unlock()
		return nil, nil
	}

	oldItem := c.Items[index]

	if index == len(c.Items) {
		c.Items = c.Items[:index-1]
	} else {
		c.Items = append(c.Items[:index], c.Items[index+1:]...)
	}

	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Deleted item by %s", id)

	errsave := c.Save(correlationId)
	//result = CloneObject(oldItem)
	result = CloneObjectForResult(oldItem, c.Prototype)
	return result, errsave
}

// Deletes multiple data items by their unique ids.
// Parameters:
//   - correlationId  string
//   (optional) transaction id to trace execution through call chain.
//   - ids []interface{}
//   ids of data items to be deleted.
// Returns: error
// error or null for success.
func (c *IdentifiableMemoryPersistence) DeleteByIds(correlationId string, ids []interface{}) (err error) {
	filterFunc := func(item interface{}) bool {
		exist := false
		itemId := GetObjectId(item)
		for _, v := range ids {
			if CompareValues(v, itemId) {
				exist = true
				break
			}
		}
		return exist
	}

	return c.DeleteByFilter(correlationId, filterFunc)
}
