package persistence

import (
	"encoding/json"
	"github.com/pip-services3-go/pip-services3-commons-go/convert"
	"github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-components-go/log"
	"reflect"
	"sync"
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

- *:logger:*:*:1.0       (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-go/master/doc/api/interfaces/log.ilogger.html ILogger]] components to pass log messages

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
    err := persistence.Set("123", interface{}({ name: "ABC" }))
    item, err := persistence.GetByName("123", "ABC")
    fmt.Println(item)   // Result: { name: "ABC" }
*/
// implements IReferenceable, IOpenable, ICleanable
type MemoryPersistence struct {
	Logger    *log.CompositeLogger
	Items     []interface{}
	Loader    ILoader
	Saver     ISaver
	opened    bool
	Prototype reflect.Type
	Lock      sync.RWMutex
}

// Creates a new instance of the MemoryPersistence
// Parameters:
// 	- prototype reflect.Type
//    type of contained data
// Return *MemoryPersistence
// a MemoryPersistence
func NewMemoryPersistence(prototype reflect.Type) *MemoryPersistence {
	if prototype == nil {
		return nil
	}
	c := &MemoryPersistence{}
	c.Prototype = prototype
	c.Logger = log.NewCompositeLogger()
	c.Items = make([]interface{}, 0, 10)
	return c
}

//  Sets references to dependent components.
//  Parameters:
// 	- references refer.IReferences
//	references to locate the component dependencies.
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
// 		- correlationId  string
// 		(optional) transaction id to trace execution through call chain.
// Returns  error or null no errors occured.
func (c *MemoryPersistence) Open(correlationId string) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	err := c.Load(correlationId)
	if err == nil {
		c.opened = true
	}
	return err
}

func (c *MemoryPersistence) Load(correlationId string) error {
	if c.Loader == nil {
		return nil
	}

	items, err := c.Loader.Load(correlationId)
	if err == nil && items != nil {
		c.Items = make([]interface{}, len(items))
		for i, v := range items {
			item := convert.MapConverter.ToNullableMap(v)
			jsonMarshalStr, errJson := json.Marshal(item)
			if errJson != nil {
				panic("MemoryPersistence.Load Error can't convert from Json to any type")
			}
			value := reflect.New(c.Prototype).Interface()
			json.Unmarshal(jsonMarshalStr, value)
			c.Items[i] = reflect.ValueOf(value).Elem().Interface()
		}
		length := len(c.Items)
		c.Logger.Trace(correlationId, "Loaded %d items", length)
	}
	return err
}

// Closes component and frees used resources.
// Parameters:
// 	- correlationId 	(optional) transaction id to trace execution through call chain.
// Retruns: error or null no errors occured.
func (c *MemoryPersistence) Close(correlationId string) error {
	err := c.Save(correlationId)
	c.opened = false
	return err
}

// Saves items to external data source using configured saver component.
//    - correlationId string
//     (optional) transaction id to trace execution through call chain.
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
// 	- correlationId 	(optional) transaction id to trace execution through call chain.
//  Returns error or null no errors occured.
func (c *MemoryPersistence) Clear(correlationId string) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Items = make([]interface{}, 0, 5)
	c.Logger.Trace(correlationId, "Cleared items")
	return c.Save(correlationId)
}
