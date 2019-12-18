package persistence

import (
	"github.com/jinzhu/copier"
	"github.com/pip-services3-go/pip-services3-commons-go/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	"github.com/pip-services3-go/pip-services3-commons-go/errors"
	refl "github.com/pip-services3-go/pip-services3-commons-go/reflect"
	"github.com/pip-services3-go/pip-services3-components-go/log"
	"math/rand"
	"reflect"
	"sort"
	"time"
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
        if &filter == nil {
			filter = NewFilterParams()
		}
        name := filter.getAsNullableString("Name");
        return func(item interface{}) bool {
			dummy, ok := item.(Dummy)
            if (*name != "" && ok && item.Name != *name)
                return false;
            return true;
        };
    }

    func (mmp * MyMemoryPersistence) GetPageByFilter(correlationId string, filter FilterParams, paging PagingParams) (page DataPage, err error) {
        return mmp.IdentifiableMemoryPersistence.GetPageByFilter(correlationId, c.composeFilter(filter), paging, nil, nil)
    }

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
//  extends MemoryPersistence  implements IConfigurable, IWriter, IGetter, ISetter
type IdentifiableMemoryPersistence struct {
	MemoryPersistence
	_maxPageSize int
}

// ================================= Service methods ====================================
// Generate new id if "Id" field in item is zero or empty
// result saves in item
func setIdIfEmpty(correlationId string, item *interface{}) error {
	value := *item
	idField := refl.ObjectReader.GetProperty(value, "Id")
	if idField != nil {
		if reflect.ValueOf(idField).IsZero() {
			if reflect.ValueOf(value).Kind() == reflect.Map {
				refl.ObjectWriter.SetProperty(value, "Id", cdata.IdGenerator.NextLong())
			} else {
				typePointer := reflect.New(reflect.TypeOf(value))
				typePointer.Elem().Set(reflect.ValueOf(value))
				interfacePointer := typePointer.Interface()
				refl.ObjectWriter.SetProperty(interfacePointer, "Id", cdata.IdGenerator.NextLong())
				*item = reflect.ValueOf(interfacePointer).Elem().Interface()
			}
		}
	} else {
		return errors.NewInternalError(correlationId, "ID_FIELD_NOT_EXIST", "Id field doesn't exist")
	}
	return nil
}

// Get index by "Id" field
// return index number
func (c *IdentifiableMemoryPersistence) getIndexById(id interface{}) int {
	var index int = -1
	for i, v := range c.Items {
		vId := refl.ObjectReader.GetProperty(v, "Id")
		if vId == id {
			index = i
			break
		}
	}
	return index
}

// Clone object
// Parameters:
//		- dst *interface{}
// 		- src *interface{}
func cloneObject(dst *interface{}, src interface{}) {
	srcValue := src
	if reflect.TypeOf(srcValue).Kind() == reflect.Ptr {
		srcValue = reflect.ValueOf(srcValue).Elem().Interface()
	}
	if reflect.ValueOf(srcValue).Kind() == reflect.Map {
		itemType := reflect.TypeOf(srcValue)
		itemValue := reflect.ValueOf(srcValue)
		mapType := reflect.MapOf(itemType.Key(), itemType.Elem())
		newMap := reflect.MakeMap(mapType)
		for _, k := range itemValue.MapKeys() {
			v := itemValue.MapIndex(k)
			newMap.SetMapIndex(k, v)
		}
		// inflate map
		typePointer := reflect.New(itemType)
		typePointer.Elem().Set(newMap)
		interfacePointer := typePointer.Interface()
		*dst = reflect.ValueOf(interfacePointer).Elem().Interface()

	} else {
		copier.Copy(dst, &srcValue)
	}
}

// Creates a new empty instance of the persistence.
// Parameters:
// 		- prototype reflect.Type
//		data type of contains items
// Return * IdentifiableMemoryPersistence
// created empty IdentifiableMemoryPersistence
func NewEmptyIdentifiableMemoryPersistence(prototype reflect.Type) (c *IdentifiableMemoryPersistence) {
	c = &IdentifiableMemoryPersistence{}
	c.MemoryPersistence = *NewEmptyMemoryPersistence(prototype)
	c.Logger = *log.NewCompositeLogger()
	c._maxPageSize = 100
	return c
}

// Creates a new instance of the persistence.
// Parameters:
// 		- prototype reflect.Type
//		data type of contains items
// 		- loader ILoader
//	    a loader to load items from external datasource.
// 		- saver  ISaver
//		a saver to save items to external datasource.
// Return * IdentifiableMemoryPersistence
// created empty IdentifiableMemoryPersistence
func NewIdentifiableMemoryPersistence(prototype reflect.Type, loader ILoader, saver ISaver) (c *IdentifiableMemoryPersistence) {
	c = &IdentifiableMemoryPersistence{}
	c.MemoryPersistence = *NewEmptyMemoryPersistence(prototype)
	c.Loader = loader
	c.Saver = saver
	c.Logger = *log.NewCompositeLogger()
	c._maxPageSize = 100
	return c
}

// Configures component by passing configuration parameters.
// Parameters:
// 		- config  config.ConfigParams
//		 configuration parameters to be set.
func (c *IdentifiableMemoryPersistence) Configure(config config.ConfigParams) {
	c._maxPageSize = config.GetAsIntegerWithDefault("options.max_page_size", c._maxPageSize)
}

// Gets a page of data items retrieved by a given filter and sorted according to sort parameters.
// cmethod shall be called by a func (imp* IdentifiableMemoryPersistence) getPageByFilter method from child struct that
// receives FilterParams and converts them into a filter function.
// Parameters:
// 		- correlationId string
//	     transaction id to trace execution through call chain.
// 		- filter func(interface{}) bool
//      (optional) a filter function to filter items
// 		- paging cdata.PagingParams
//      (optional) paging parameters
// 		- sortFunc func(a, b interface{}) bool
//      (optional) sorting compare function func Less (a, b interface{}) bool  see sort.Interface Less function
// 		- selectFunc func(in interface{}) (out interface{})
//      (optional) projection parameters
// Return cdata.DataPage, error
// data page or error.
func (c *IdentifiableMemoryPersistence) GetPageByFilter(correlationId string, filter func(interface{}) bool,
	paging cdata.PagingParams, sortFunc func(a, b interface{}) bool, selectFunc func(in interface{}) (out interface{})) (page cdata.DataPage, err error) {
	var items []interface{}
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	// Filter and sort
	if filter != nil {
		for _, v := range c.Items {
			if filter(v) {
				items = append(items, v)
			}
		}
	} else {
		items = make([]interface{}, len(c.Items))
		copy(items, c.Items)
	}
	if sortFunc != nil {
		if sortFunc != nil {
			localSort := sorter{items: items, compFunc: sortFunc}
			sort.Sort(localSort)
		}
	}
	// Extract a page
	if &paging == nil {
		paging = *cdata.NewEmptyPagingParams()
	}
	skip := paging.GetSkip(-1)
	take := paging.GetTake((int64)(c._maxPageSize))
	var total int64
	if paging.Total {
		total = (int64)(len(items))
	}
	if skip > 0 {
		items = items[skip:]
	}
	if (int64)(len(items)) >= take {
		items = items[:take]
	}
	// Get projection
	if selectFunc != nil {
		for i, v := range items {
			items[i] = selectFunc(v)
		}
	}
	c.Logger.Trace(correlationId, "Retrieved %d items", len(items))
	page = *cdata.NewDataPage(&total, items)
	return page, nil
}

// Gets a list of data items retrieved by a given filter and sorted according to sort parameters.
// This method shall be called by a func (c * IdentifiableMemoryPersistence) GetListByFilter method from child struct that
// receives FilterParams and converts them into a filter function.
// Parameters:
// 		- correlationId string
//      (optional) transaction id to trace execution through call chain.
// 		- filter func(interface{}) bool
//      (optional) a filter function to filter items
// 		- sortFunc func(a, b interface{}) bool
//      (optional) sorting compare function func Less (a, b interface{}) bool  see sort.Interface Less function
// 		- selectFunc func(in interface{}) (out interface{})
//      (optional) projection parameters
// Returns  []interface{},  error
// array of items and error
func (c *IdentifiableMemoryPersistence) GetListByFilter(correlationId string, filter func(interface{}) bool,
	sortFunc func(a, b interface{}) bool, selectFunc func(in interface{}) (out interface{})) (results []interface{}, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	// Apply filter
	if filter != nil {
		for _, v := range c.Items {
			if filter(v) {
				results = append(results, v)
			}
		}
	} else {
		copy(results, c.Items)
	}
	// Apply sorting
	if sortFunc != nil {
		localSort := sorter{items: results, compFunc: sortFunc}
		sort.Sort(localSort)
	}
	// Get projection
	if selectFunc != nil {
		for i, v := range results {
			results[i] = selectFunc(v)
		}
	}
	c.Logger.Trace(correlationId, "Retrieved %d items", len(results))
	return results, nil
}

// Gets a list of data items retrieved by given unique ids.
// Parameters:
// 		- correlationId string
//   	(optional) transaction id to trace execution through call chain.
// 		- ids  []interface{}
//      ids of data items to be retrieved
// Returns  []interface{}, error
// data list or error.
func (c *IdentifiableMemoryPersistence) GetListByIds(correlationId string, ids []interface{}) (result []interface{}, err error) {
	filter := func(item interface{}) bool {
		var exist bool = false
		itemId := refl.ObjectReader.GetProperty(item, "Id")
		for _, v := range ids {
			vId := refl.ObjectReader.GetValue(v)
			if itemId == vId {
				exist = true
				break
			}
		}
		return exist
	}
	return c.GetListByFilter(correlationId, filter, nil, nil)
}

// Gets a random item from items that match to a given filter.
// This method shall be called by a func (c* IdentifiableMemoryPersistence) GetOneRandom method from child type that
// receives FilterParams and converts them into a filter function.
// Parameters:
// 		- correlationId string
//     (optional) transaction id to trace execution through call chain.
// 		- filter   func(interface{}) bool
//     (optional) a filter function to filter items.
// Returns: *interface{}, error
// random item or error.
func (c *IdentifiableMemoryPersistence) GetOneRandom(correlationId string, filter func(interface{}) bool) (item interface{}, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	var items []interface{}
	// Apply filter
	if filter != nil {
		for _, v := range c.Items {
			if filter(v) {
				items = append(items, v)
			}
		}
	} else {
		copy(items, c.Items)
	}
	rand.Seed(time.Now().UnixNano())
	if len(items) > 0 {
		item = items[rand.Intn(len(items))]
	}
	if item != nil {
		c.Logger.Trace(correlationId, "Retrieved a random item")
	} else {
		c.Logger.Trace(correlationId, "Nothing to return as random item")
	}
	return item, nil
}

// Gets a data item by its unique id.
// Parameters:
// 		- correlationId  string
//   	(optional) transaction id to trace execution through call chain.
// 		- id interface{}
//      an id of data item to be retrieved.
// Returns:  *interface{}, error
// data item or error.
func (c *IdentifiableMemoryPersistence) GetOneById(correlationId string, id interface{}) (item interface{}, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	var items []interface{}
	for _, v := range c.Items {
		vId := refl.ObjectReader.GetProperty(v, "Id")
		if vId == id {
			items = append(items, v)
		}
	}
	if len(items) > 0 {
		item = items[0]
	}
	if item != nil {
		c.Logger.Trace(correlationId, "Retrieved item %s", id)
	} else {
		c.Logger.Trace(correlationId, "Cannot find item by %s", id)
	}
	return item, err
}

// Creates a data item.
// Returns:
// 	 - correlation_id string
//   (optional) transaction id to trace execution through call chain.
// 	 - item  string
//   an item to be created.
// Returns:  *interface{}, error
// created item or error.
func (c *IdentifiableMemoryPersistence) Create(correlationId string, item interface{}) (result interface{}, err error) {
	c.Lock.Lock()
	var tmp interface{}
	cloneObject(&tmp, item)
	err = setIdIfEmpty(correlationId, &tmp)
	if err != nil {
		return nil, err
	}
	//result = tmp
	cloneObject(&result, tmp)
	c.Items = append(c.Items, tmp)
	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Created item %s", refl.ObjectReader.GetProperty(tmp, "Id"))
	errsave := c.Save(correlationId)
	return result, errsave
}

// Sets a data item. If the data item exists it updates it,
// otherwise it create a new data item.
// Parameters:
// 		- correlation_id string
//	    (optional) transaction id to trace execution through call chain.
// 		- item  interface{}
//      a item to be set.
// Returns:  *interface{}, error
// updated item or error.
func (c *IdentifiableMemoryPersistence) Set(correlationId string, item interface{}) (result interface{}, err error) {
	c.Lock.Lock()
	var tmp interface{}
	cloneObject(&tmp, item)
	err = setIdIfEmpty(correlationId, &tmp)
	if err != nil {
		return nil, err
	}
	id := refl.ObjectReader.GetProperty(item, "Id")
	index := c.getIndexById(id)
	if index < 0 {
		c.Items = append(c.Items, tmp)
	} else {
		//c.Items[index] = tmp
		cloneObject(&c.Items[index], tmp)
	}
	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Set item %s", refl.ObjectReader.GetProperty(tmp, "Id"))
	errsav := c.Save(correlationId)
	//result = new(interface{})
	cloneObject(&result, tmp)
	return result, errsav
}

// Updates a data item.
// Parameters:
// 		- correlation_id string
//  	(optional) transaction id to trace execution through call chain.
// 		- item  interface{}
//      an item to be updated.
// Returns:   *interface{}, error
// updated item or error.
func (c *IdentifiableMemoryPersistence) Update(correlationId string, item interface{}) (result interface{}, err error) {
	c.Lock.Lock()
	id := refl.ObjectReader.GetProperty(item, "Id")
	index := c.getIndexById(id)
	if index < 0 {
		c.Logger.Trace(correlationId, "Item %s was not found", refl.ObjectReader.GetProperty(item, "Id"))
		return nil, nil
	}
	var tmp interface{}
	cloneObject(&tmp, item)
	cloneObject(&c.Items[index], item)
	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Updated item %s", refl.ObjectReader.GetProperty(item, "Id"))
	errsave := c.Save(correlationId)
	return tmp, errsave
}

// Updates only few selectFuncected fields in a data item.
// Parameters:
// 		- correlation_id string
//    	(optional) transaction id to trace execution through call chain.
// 		- id interface{}
//      an id of data item to be updated.
// 		- data  cdata.AnyValueMap
//      a map with fields to be updated.
// Returns: *interface{}, error
// updated item or error.
func (c *IdentifiableMemoryPersistence) UpdatePartially(correlationId string, id interface{}, data cdata.AnyValueMap) (item interface{}, err error) {
	c.Lock.Lock()
	index := c.getIndexById(id)
	if index < 0 {
		c.Logger.Trace(correlationId, "Item %s was not found", id)
		return nil, nil
	}
	var tmp interface{}
	cloneObject(&tmp, c.Items[index])
	if reflect.ValueOf(tmp).Kind() == reflect.Map {
		refl.ObjectWriter.SetProperties(tmp, data.Value())
	} else {
		objPointer := reflect.New(reflect.TypeOf(tmp))
		objPointer.Elem().Set(reflect.ValueOf(tmp))
		intPointer := objPointer.Interface()
		refl.ObjectWriter.SetProperties(intPointer, data.Value())
		tmp = reflect.ValueOf(intPointer).Elem().Interface()
	}
	cloneObject(&c.Items[index], tmp)
	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Partially updated item %s", id)
	errsave := c.Save(correlationId)
	return tmp, errsave
}

// Deleted a data item by it's unique id.
// Parameters:
// 		- correlation_id string
//	    (optional) transaction id to trace execution through call chain.
//  	- id interface{}
//      an id of the item to be deleted
// Retruns:
// deleted item or error.
func (c *IdentifiableMemoryPersistence) DeleteById(correlationId string, id interface{}) (item interface{}, err error) {
	c.Lock.Lock()
	index := c.getIndexById(id)
	if index < 0 {
		c.Logger.Trace(correlationId, "Item %s was not found", id)
		return nil, nil
	} else {
		//var tmp interface{}
		cloneObject(&item, c.Items[index])
		//item = &tmp
	}
	if index == len(c.Items) {
		c.Items = c.Items[:index-1]
	} else {
		c.Items = append(c.Items[:index], c.Items[index+1:]...)
	}
	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Deleted item by %s", id)
	errsave := c.Save(correlationId)
	return item, errsave
}

// Deletes data items that match to a given filter.
// this method shall be called by a func (c* IdentifiableMemoryPersistence) DeleteByFilter method from child struct that
// receives FilterParams and converts them into a filter function.
// Parameters:
// 		- correlationId  string
//		(optional) transaction id to trace execution through call chain.
// 		- filter  filter func(interface{}) bool
//      (optional) a filter function to filter items.
// Retruns: error
// error or nil for success.
func (c *IdentifiableMemoryPersistence) DeleteByFilter(correlationId string, filter func(interface{}) bool) (err error) {
	c.Lock.Lock()
	deleted := 0

	for i := 0; i < len(c.Items); {
		if filter(c.Items[i]) {
			if i == len(c.Items)-1 {
				c.Items = c.Items[:i]
			} else {
				c.Items = append(c.Items[:i], c.Items[i+1:]...)
			}
			deleted++
		} else {
			i++
		}
	}
	if deleted == 0 {
		return nil
	}
	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Deleted %s items", deleted)
	errsave := c.Save(correlationId)
	return errsave
}

// Deletes multiple data items by their unique ids.
// Parameters:
// 		- correlationId  string
//     	(optional) transaction id to trace execution through call chain.
// 		- ids []interface{}
//     	ids of data items to be deleted.
// Returns: error
// error or null for success.
func (c *IdentifiableMemoryPersistence) DeleteByIds(correlationId string, ids []interface{}) (err error) {
	filter := func(item interface{}) bool {
		var exist bool = false
		itemId := refl.ObjectReader.GetProperty(item, "Id")
		for _, v := range ids {
			if v == itemId {
				exist = true
				break
			}
		}
		return exist
	}
	return c.DeleteByFilter(correlationId, filter)
}
