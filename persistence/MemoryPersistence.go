package persistence

import (
	"github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-components-go/log"
)

/*
Abstract persistence component that stores data in memory.

This is the most basic persistence component that is only
able to store data items of any type. Specific CRUD operations
over the data items must be implemented in child struct by
accessing _items property and calling Save method.

The component supports loading and saving items from another data source.
That allows to use it as a base struct for file and other types
of persistence components that cache all data in memory.

References

- *:logger:*:*:1.0       (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-go/master/doc/api/interfaces/log.ilogger.html ILogger]] components to pass log messages

Example

    type MyMemoryPersistence struct {
        MemoryPersistence

    }
     func (mmp * MyMemoryPersistence) GetByName(correlationId string, name string)(item interface{}, err error) {
        for _, v := range mmp._items {
            if v.name == name {
                item = v
                break
            }
        }
        return item, nil
    });

    func (mmp * MyMemoryPersistence) Set(correlatonId: string, item: MyData, callback: (err) => void): void {

        mmp._items = append(mmp._items, item);
        mmp.Save(correlationId);
    }

    persistence := NewMyMemoryPersistence();
    err := persistence.Set("123", interface{}({ name: "ABC" }))
    item, err := persistence.GetByName("123", "ABC")
    fmt.Println(item)   // Result: { name: "ABC" }
*/
// implements IReferenceable, IOpenable, ICleanable
type MemoryPersistence struct {
	_logger log.CompositeLogger
	_items  []interface{}
	_loader ILoader
	_saver  ISaver
	_opened bool
}

// Creates a new empty instance of the persistence.
func NewEmptyMemoryPersistence() (mp *MemoryPersistence) {
	mp = &MemoryPersistence{}
	mp._logger = *log.NewCompositeLogger()
	mp._items = make([]interface{}, 0, 10)
	return mp
}

/*
   Creates a new instance of the persistence.

   - loader    (optional) a loader to load items from external datasource.
   - saver     (optional) a saver to save items to external datasource.
*/
func NewMemoryPersistence(loader ILoader, saver ISaver) (mp *MemoryPersistence) {
	mp = &MemoryPersistence{}
	mp._items = make([]interface{}, 0, 5)
	mp._loader = loader
	mp._saver = saver
	mp._logger = *log.NewCompositeLogger()
	return mp
}

/*
	Sets references to dependent components.

	- references 	references to locate the component dependencies.
*/
func (mp *MemoryPersistence) SetReferences(references refer.IReferences) {
	mp._logger.SetReferences(references)
}

/*
	Checks if the component is opened.

	Returns true if the component has been opened and false otherwise.
*/
func (mp *MemoryPersistence) IsOpen() bool {
	return mp._opened
}

/*
	Opens the component.

	- correlationId 	(optional) transaction id to trace execution through call chain.
    - callback 			callback function that receives error or null no errors occured.
*/
func (mp *MemoryPersistence) Open(correlationId string) error {
	err := mp.load(correlationId)
	if err == nil {
		mp._opened = true
	}
	return err
}

func (mp *MemoryPersistence) load(correlationId string) error {
	if mp._loader == nil {
		return nil
	}

	items, err := mp._loader.Load(correlationId)
	if err == nil {
		mp._items = items
		mp._logger.Trace(correlationId, "Loaded %d items", len(mp._items))
	}
	return err
}

/*
	Closes component and frees used resources.

	- correlationId 	(optional) transaction id to trace execution through call chain.
    - callback 			callback function that receives error or null no errors occured.
*/
func (mp *MemoryPersistence) Close(correlationId string) error {
	err := mp.Save(correlationId)
	mp._opened = false
	return err
}

/*
   Saves items to external data source using configured saver component.

   - correlationId     (optional) transaction id to trace execution through call chain.
   - callback          (optional) callback function that receives error or null for success.
*/
func (mp *MemoryPersistence) Save(correlationId string) error {
	if mp._saver == nil {
		return nil
	}

	err := mp._saver.Save(correlationId, mp._items)
	if err == nil {
		mp._logger.Trace(correlationId, "Saved %d items", len(mp._items))
	}
	return err
}

/*
	Clears component state.

	- correlationId 	(optional) transaction id to trace execution through call chain.
    - callback 			callback function that receives error or null no errors occured.
*/
func (mp *MemoryPersistence) Clear(correlationId string) error {
	mp._items = make([]interface{}, 0)
	mp._logger.Trace(correlationId, "Cleared items")
	return mp.Save(correlationId)
}
