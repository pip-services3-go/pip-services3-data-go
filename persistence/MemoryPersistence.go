package persistence

import (
	"github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-components-go/log"
)

/*
Abstract persistence component that stores data in memory.

This is the most basic persistence component that is only
able to store data items of any type. Specific CRUD operations
over the data items must be implemented in child classes by
accessing <code>this._items</code> property and calling [[save]] method.

The component supports loading and saving items from another data source.
That allows to use it as a base class for file and other types
of persistence components that cache all data in memory.

### References ###

- <code>\*:logger:\*:\*:1.0</code>       (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/log.ilogger.html ILogger]] components to pass log messages

### Example ###

    class MyMemoryPersistence extends MemoryPersistence<MyData> {

         getByName(correlationId: string, name: string, callback: (err, item) => void): void {
            let item = _.find(this._items, (d) => d.name == name);
            callback(null, item);
        });

         set(correlatonId: string, item: MyData, callback: (err) => void): void {
            this._items = _.filter(this._items, (d) => d.name != name);
            this._items.push(item);
            this.save(correlationId, callback);
        }

    }

    let persistence = new MyMemoryPersistence();

    persistence.set("123", { name: "ABC" }, (err) => {
        persistence.getByName("123", "ABC", (err, item) => {
            console.log(item);                   // Result: { name: "ABC" }
        });
    });
*/
//<T> implements IReferenceable, IOpenable, ICleanable
type MemoryPersistence struct {
	_logger log.CompositeLogger
	_items  []interface{}
	_loader ILoader
	_saver  ISaver
	_opened bool
}

func (mp *MemoryPersistence) NewEmptyMemoryPersistence() {
	mp._logger = *log.NewCompositeLogger()
}

/*
   Creates a new instance of the persistence.

   @param loader    (optional) a loader to load items from external datasource.
   @param saver     (optional) a saver to save items to external datasource.
*/
func (mp *MemoryPersistence) NewMemoryPersistence(loader ILoader, saver ISaver) {
	mp._loader = loader
	mp._saver = saver
	mp._logger = *log.NewCompositeLogger()
}

/*
	Sets references to dependent components.

	@param references 	references to locate the component dependencies.
*/
func (mp *MemoryPersistence) SetReferences(references refer.IReferences) {
	mp._logger.SetReferences(references)
}

/*
	Checks if the component is opened.

	@returns true if the component has been opened and false otherwise.
*/
func (mp *MemoryPersistence) IsOpen() bool {
	return mp._opened
}

/*
	Opens the component.

	@param correlationId 	(optional) transaction id to trace execution through call chain.
    @param callback 			callback function that receives error or null no errors occured.
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

	@param correlationId 	(optional) transaction id to trace execution through call chain.
    @param callback 			callback function that receives error or null no errors occured.
*/
func (mp *MemoryPersistence) Close(correlationId string) error {
	err := mp.Save(correlationId)
	mp._opened = false
	return err
}

/*
   Saves items to external data source using configured saver component.

   @param correlationId     (optional) transaction id to trace execution through call chain.
   @param callback          (optional) callback function that receives error or null for success.
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

	@param correlationId 	(optional) transaction id to trace execution through call chain.
    @param callback 			callback function that receives error or null no errors occured.
*/
func (mp *MemoryPersistence) Clear(correlationId string) error {
	mp._items = make([]interface{}, 0)
	mp._logger.Trace(correlationId, "Cleared items")
	return mp.Save(correlationId)
}
