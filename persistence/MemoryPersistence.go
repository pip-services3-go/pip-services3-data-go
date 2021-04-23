package persistence

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/pip-services3-go/pip-services3-commons-go/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	"github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-components-go/log"
)

/*
Abstract persistence component that stores data in memory.

This is the most basic persistence component that is only
able to store data items of any type. Specific CRUD operations
over the data items must be implemented in child struct by
accessing Items property and calling Save method.

The component supports loading and saving items from another data source.
That allows to use it as a base struct for file and other types
of persistence components that cache all data in memory.

References

- *:logger:*:*:1.0    ILogger components to pass log messages

Example

    type MyMemoryPersistence struct {
        MemoryPersistence

    }
     func (c * MyMemoryPersistence) GetByName(correlationId string, name string)(item interface{}, err error) {
        for _, v := range c.Items {
            if v.name == name {
                item = v
                break
            }
        }
        return item, nil
    });

    func (c * MyMemoryPersistence) Set(correlatonId: string, item: MyData, callback: (err) => void): void {
        c.Items = append(c.Items, item);
        c.Save(correlationId);
    }

    persistence := NewMyMemoryPersistence();
    err := persistence.Set("123", MyData{ name: "ABC" })
    item, err := persistence.GetByName("123", "ABC")
    fmt.Println(item)   // Result: { name: "ABC" }
*/
// implements IReferenceable, IOpenable, ICleanable
type MemoryPersistence struct {
	Logger      *log.CompositeLogger
	Items       []interface{}
	Loader      ILoader
	Saver       ISaver
	opened      bool
	Prototype   reflect.Type
	Lock        sync.RWMutex
	MaxPageSize int
}

// Creates a new instance of the MemoryPersistence
// Parameters:
//  - prototype reflect.Type
//   type of contained data
// Return *MemoryPersistence
// a MemoryPersistence
func NewMemoryPersistence(prototype reflect.Type) *MemoryPersistence {
	if prototype == nil {
		panic("Prototype cannot be nil")
	}
	c := &MemoryPersistence{}
	c.Prototype = prototype
	c.Logger = log.NewCompositeLogger()
	c.Items = make([]interface{}, 0, 10)
	return c
}

//  Sets references to dependent components.
//  Parameters:
//   - references refer.IReferences
//   references to locate the component dependencies.
func (c *MemoryPersistence) SetReferences(references refer.IReferences) {
	c.Logger.SetReferences(references)
}

//  Checks if the component is opened.
//  Returns true if the component has been opened and false otherwise.
func (c *MemoryPersistence) IsOpen() bool {
	return c.opened
}

// Opens the component.
// Parameters:
//   - correlationId  string
//   (optional) transaction id to trace execution through call chain.
// Returns  error or null no errors occured.
func (c *MemoryPersistence) Open(correlationId string) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	err := c.load(correlationId)
	if err == nil {
		c.opened = true
	}
	return err
}

func (c *MemoryPersistence) load(correlationId string) error {
	if c.Loader == nil {
		return nil
	}

	items, err := c.Loader.Load(correlationId)
	if err == nil && items != nil {
		c.Items = items
		c.Items = make([]interface{}, len(items))
		for i, v := range items {
			item := convert.MapConverter.ToNullableMap(v)
			jsonMarshalStr, errJson := json.Marshal(item)
			if errJson != nil {
				panic("MemoryPersistence.Load Error can't convert from Json to type")
			}
			value := reflect.New(c.Prototype).Interface()
			json.Unmarshal(jsonMarshalStr, value)
			c.Items[i] = reflect.ValueOf(value).Elem().Interface() // load value
		}
		length := len(c.Items)
		c.Logger.Trace(correlationId, "Loaded %d items", length)
	}
	return err
}

// Closes component and frees used resources.
// Parameters:
//  - correlationId string
//  (optional) transaction id to trace execution through call chain.
// Retruns: error or nil if no errors occured.
func (c *MemoryPersistence) Close(correlationId string) error {
	err := c.Save(correlationId)
	c.opened = false
	return err
}

// Saves items to external data source using configured saver component.
// Parameters:
//   - correlationId string
//   (optional) transaction id to trace execution through call chain.
// Return error or null for success.
func (c *MemoryPersistence) Save(correlationId string) error {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	if c.Saver == nil {
		return nil
	}

	err := c.Saver.Save(correlationId, c.Items)
	if err == nil {
		length := len(c.Items)
		c.Logger.Trace(correlationId, "Saved %d items", length)
	}
	return err
}

// Clears component state.
// Parameters:
//  - correlationId string
//  (optional) transaction id to trace execution through call chain.
//  Returns error or null no errors occured.
func (c *MemoryPersistence) Clear(correlationId string) error {
	c.Lock.Lock()

	c.Items = make([]interface{}, 0, 5)
	c.Logger.Trace(correlationId, "Cleared items")

	c.Lock.Unlock()
	return c.Save(correlationId)
}

// Gets a page of data items retrieved by a given filter and sorted according to sort parameters.
// cmethod shall be called by a func (imp* IdentifiableMemoryPersistence) getPageByFilter method from child struct that
// receives FilterParams and converts them into a filter function.
// Parameters:
//   - correlationId string
//   transaction id to trace execution through call chain.
//   - filter func(interface{}) bool
//   (optional) a filter function to filter items
//   - paging *cdata.PagingParams
//   (optional) paging parameters
//   - sortFunc func(a, b interface{}) bool
//   (optional) sorting compare function func Less (a, b interface{}) bool  see sort.Interface Less function
//   - selectFunc func(in interface{}) (out interface{})
// (optional) projection parameters
// Return cdata.DataPage, error
// data page or error.
func (c *MemoryPersistence) GetPageByFilter(correlationId string, filterFunc func(interface{}) bool,
	paging *cdata.PagingParams, sortFunc func(a, b interface{}) bool, selectFunc func(in interface{}) (out interface{})) (page *cdata.DataPage, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	var items []interface{}

	// Apply filtering
	if filterFunc != nil {
		for _, v := range c.Items {
			if filterFunc(v) {
				items = append(items, v)
			}
		}
	} else {
		items = make([]interface{}, len(c.Items))
		copy(items, c.Items)
	}

	// Apply sorting
	if sortFunc != nil {
		localSort := sorter{items: items, compFunc: sortFunc}
		sort.Sort(localSort)
	}

	// Extract a page
	if paging == nil {
		paging = cdata.NewEmptyPagingParams()
	}
	skip := paging.GetSkip(-1)
	take := paging.GetTake((int64)(c.MaxPageSize))
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
	// W!
	for i := 0; i < len(items); i++ {
		//items[i] = CloneObject(items[i])
		items[i] = CloneObjectForResult(items[i], c.Prototype)
	}

	page = cdata.NewDataPage(&total, items)
	return page, nil
}

// Gets a list of data items retrieved by a given filter and sorted according to sort parameters.
// This method shall be called by a func (c * IdentifiableMemoryPersistence) GetListByFilter method from child struct that
// receives FilterParams and converts them into a filter function.
// Parameters:
//   - correlationId string
//   (optional) transaction id to trace execution through call chain.
//   - filter func(interface{}) bool
//   (optional) a filter function to filter items
//   - sortFunc func(a, b interface{}) bool
//   (optional) sorting compare function func Less (a, b interface{}) bool  see sort.Interface Less function
//   - selectFunc func(in interface{}) (out interface{})
//   (optional) projection parameters
// Returns  []interface{},  error
// array of items and error
func (c *MemoryPersistence) GetListByFilter(correlationId string, filterFunc func(interface{}) bool,
	sortFunc func(a, b interface{}) bool, selectFunc func(in interface{}) (out interface{})) (results []interface{}, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	// Apply filter
	if filterFunc != nil {
		for _, v := range c.Items {
			if filterFunc(v) {
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
	//W!
	for i := 0; i < len(results); i++ {
		//results[i] = CloneObject(results[i])
		results[i] = CloneObjectForResult(results[i], c.Prototype)
	}
	return results, nil
}

// Gets a random item from items that match to a given filter.
// This method shall be called by a func (c* IdentifiableMemoryPersistence) GetOneRandom method from child type that
// receives FilterParams and converts them into a filter function.
// Parameters:
//   - correlationId string
//   (optional) transaction id to trace execution through call chain.
//   - filter   func(interface{}) bool
//   (optional) a filter function to filter items.
// Returns: interface{}, error
// random item or error.
func (c *MemoryPersistence) GetOneRandom(correlationId string, filterFunc func(interface{}) bool) (result interface{}, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	var items []interface{}

	// Apply filter
	if filterFunc != nil {
		for _, v := range c.Items {
			if filterFunc(v) {
				items = append(items, v)
			}
		}
	} else {
		copy(items, c.Items)
	}
	rand.Seed(time.Now().UnixNano())

	var item interface{} = nil
	if len(items) > 0 {
		item = items[rand.Intn(len(items))]
	}

	if item != nil {
		c.Logger.Trace(correlationId, "Retrieved a random item")
	} else {
		c.Logger.Trace(correlationId, "Nothing to return as random item")
	}
	//result = CloneObject(item)
	result = CloneObjectForResult(item, c.Prototype)
	return result, nil
}

// Creates a data item.
// Returns:
//   - correlation_id string
//   (optional) transaction id to trace execution through call chain.
//   - item  string
//   an item to be created.
// Returns:  interface{}, error
// created item or error.
func (c *MemoryPersistence) Create(correlationId string, item interface{}) (result interface{}, err error) {
	c.Lock.Lock()

	newItem := CloneObject(item, c.Prototype)
	c.Items = append(c.Items, newItem)

	c.Lock.Unlock()
	c.Logger.Trace(correlationId, "Created item")

	errsave := c.Save(correlationId)
	result = CloneObjectForResult(newItem, c.Prototype)

	return result, errsave
}

// Deletes data items that match to a given filter.
// this method shall be called by a func (c* IdentifiableMemoryPersistence) DeleteByFilter method from child struct that
// receives FilterParams and converts them into a filter function.
// Parameters:
//   - correlationId  string
//   (optional) transaction id to trace execution through call chain.
//   - filter  filter func(interface{}) bool
//   (optional) a filter function to filter items.
// Retruns: error
// error or nil for success.
func (c *MemoryPersistence) DeleteByFilter(correlationId string, filterFunc func(interface{}) bool) (err error) {
	c.Lock.Lock()

	deleted := 0
	for i := 0; i < len(c.Items); {
		if filterFunc(c.Items[i]) {
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
	c.Lock.Unlock()

	if deleted == 0 {
		return nil
	}

	c.Logger.Trace(correlationId, "Deleted %s items", deleted)

	errsave := c.Save(correlationId)
	return errsave
}

// Gets a count of data items retrieved by a given filter.
// this method shall be called by a func (imp* IdentifiableMemoryPersistence) getCountByFilter method from child struct that
// receives FilterParams and converts them into a filter function.
// Parameters:
//  - correlationId string
//  transaction id to trace execution through call chain.
//  - filter func(interface{}) bool
//  (optional) a filter function to filter items
// Return int, error
// data count or error.
func (c *MemoryPersistence) GetCountByFilter(correlationId string, filterFunc func(interface{}) bool) (count int64, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	// Apply filtering
	if filterFunc != nil {
		for _, v := range c.Items {
			if filterFunc(v) {
				count++
			}
		}
	} else {
		count = 0
	}
	c.Logger.Trace(correlationId, "Find %d items", count)
	return count, nil
}
