package persistence

import (
	"github.com/pip-services3-go/pip-services3-commons-go/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	refl "github.com/pip-services3-go/pip-services3-commons-go/reflect"
	"github.com/pip-services3-go/pip-services3-components-go/log"

	"github.com/jinzhu/copier"

	"math/rand"
	"reflect"
	"time"
)

/*
Abstract persistence component that stores data in memory
and implements a number of CRUD operations over data items with unique ids.
The data items must implement [[https://rawgit.com/pip-services-node/pip-services3-commons-node/master/doc/api/interfaces/data.iidentifiable.html IIdentifiable interface]].

In basic scenarios child classes shall only override [[getPageByFilter]],
[[getListByFilter]] or [[deleteByFilter]] operations with specific filter function.
All other operations can be used out of the box.

In complex scenarios child classes can implement additional operations by
accessing cached items via imp._items property and calling [[save]] method
on updates.

See [[MemoryPersistence]]

 Configuration parameters

- options:
    - max_page_size:       Maximum number of items returned in a single page (default: 100)

 References

- *:logger:*:*:1.0     (optional) ILogger components to pass log messages

 Examples

    class MyMemoryPersistence extends IdentifiableMemoryPersistence<MyData, string> {

        private composeFilter(filter: FilterParams): any {
            filter = filter || new FilterParams();
            let name = filter.getAsNullableString("name");
            return (item) => {
                if (name != null && item.name != name)
                    return false;
                return true;
            };
        }

        func (imp* IdentifiableMemoryPersistence) getPageByFilter(correlationId: string, filter: FilterParams, paging: PagingParams,
                callback: (err: any, page: DataPage<MyData>) => void): void {
            super.getPageByFilter(correlationId, imp.composeFilter(filter), paging, null, null, callback);
        }

    }

    let persistence = new MyMemoryPersistence();

    persistence.create("123", { id: "1", name: "ABC" }, (err, item) => {
        persistence.getPageByFilter(
            "123",
            FilterParams.fromTuples("name", "ABC"),
            null,
            (err, page) => {
                console.log(page.data);          // Result: { id: "1", name: "ABC" }

                persistence.deleteById("123", "1", (err, item) => {
                    ...
                });
            }
        )
    });
*/
//<T extends IIdentifiable<K>, K> extends MemoryPersistence  implements IConfigurable, IWriter<T, K>, IGetter<T, K>, ISetter {

type IdentifiableMemoryPersistence struct {
	MemoryPersistence
	_maxPageSize int
}

// Creates a new instance of the persistence.
// Parameters:
// - loader    (optional) a loader to load items from external datasource.
// - saver     (optional) a saver to save items to external datasource.

func NewEmptyIdentifiableMemoryPersistence() (imp *IdentifiableMemoryPersistence) {
	imp = &IdentifiableMemoryPersistence{}
	imp.MemoryPersistence = *NewEmptyMemoryPersistence()
	imp._logger = *log.NewCompositeLogger()
	imp._maxPageSize = 100
	return imp
}

// Creates a new instance of the persistence.
// Parameters:
// - loader    (optional) a loader to load items from external datasource.
// - saver     (optional) a saver to save items to external datasource.

func NewIdentifiableMemoryPersistence(loader ILoader, saver ISaver) (imp *IdentifiableMemoryPersistence) {
	imp = &IdentifiableMemoryPersistence{}
	imp.MemoryPersistence = *NewEmptyMemoryPersistence()
	imp._loader = loader
	imp._saver = saver
	imp._logger = *log.NewCompositeLogger()
	imp._maxPageSize = 100
	return imp
}

// Configures component by passing configuration parameters.
// Parameters:
// - config    configuration parameters to be set.

func (imp *IdentifiableMemoryPersistence) Configure(config config.ConfigParams) {
	imp._maxPageSize = config.GetAsIntegerWithDefault("options.max_page_size", imp._maxPageSize)
}

// Gets a page of data items retrieved by a given filter and sorted according to sort parameters.
// imp method shall be called by a func (imp* IdentifiableMemoryPersistence) getPageByFilter method from child class that
// receives FilterParams and converts them into a filter function.
// Parameters:
// - correlationId     (optional) transaction id to trace execution through call chain.
// - filter            (optional) a filter function to filter items
// - paging            (optional) paging parameters
// - sort              (optional) sorting parameters implements sort.Interface by providing Less and using the Len and
//  Swap methods
// - select            (optional) projection parameters (not used yet)
// Return cdata.DataPage, error
// data page or error.

func (imp *IdentifiableMemoryPersistence) GetPageByFilter(correlationId string, filter func(interface{}) bool,
	paging cdata.PagingParams, sortWrapper interface{}, sel interface{}) (page cdata.DataPage, err error) {
	var items []interface{}
	// Filter and sort
	if filter != nil {
		for _, v := range imp._items {
			if filter(v) {
				items = append(items, v)
			}
		}
	} else {
		items = make([]interface{}, len(imp._items))
		copier.Copy(items, imp._items)
	}

	if sortWrapper != nil {
		//items = _.sortUniqBy(items, sort)
	}
	// Extract a page
	if &paging == nil {
		paging = *cdata.NewEmptyPagingParams()
	}
	skip := paging.GetSkip(-1)
	take := paging.GetTake((int64)(imp._maxPageSize))

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

	imp._logger.Trace(correlationId, "Retrieved %d items", len(items))

	page = *cdata.NewDataPage(&total, items)
	return page, nil
}

/*
Gets a list of data items retrieved by a given filter and sorted according to sort parameters.

imp method shall be called by a func (imp* IdentifiableMemoryPersistence) getListByFilter method from child class that
receives FilterParams and converts them into a filter function.

- correlationId    (optional) transaction id to trace execution through call chain.
- filter           (optional) a filter function to filter items
- paging           (optional) paging parameters
- sort             (optional) sorting parameters
- select           (optional) projection parameters (not used yet)
- callback         callback function that receives a data list or error.
*/
func (imp *IdentifiableMemoryPersistence) GetListByFilter(correlationId string, filter func(interface{}) bool, sortWrapper interface{}, sel interface{}) (results []interface{}, err error) {

	// Apply filter
	if filter != nil {
		for _, v := range imp._items {
			if filter(v) {
				results = append(results, v)
			}
		}
	} else {
		results = imp._items
	}
	// Apply sorting
	if sortWrapper != nil {
		//items = _.sortUniqBy(items, sort);
	}

	imp._logger.Trace(correlationId, "Retrieved %d items", len(results))

	return results, nil
}

/*
Gets a list of data items retrieved by given unique ids.

- correlationId     (optional) transaction id to trace execution through call chain.
- ids               ids of data items to be retrieved
- callback         callback function that receives a data list or error.
*/
func (imp *IdentifiableMemoryPersistence) GetListByIds(correlationId string, ids []interface{}) (items []interface{}, err error) {
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
	return imp.GetListByFilter(correlationId, filter, nil, nil)
}

/*
Gets a random item from items that match to a given filter.

imp method shall be called by a func (imp* IdentifiableMemoryPersistence) getOneRandom method from child class that
receives FilterParams and converts them into a filter function.

- correlationId     (optional) transaction id to trace execution through call chain.
- filter            (optional) a filter function to filter items.
- callback          callback function that receives a random item or error.
*/
func (imp *IdentifiableMemoryPersistence) GetOneRandom(correlationId string, filter func(interface{}) bool) (item *interface{}, err error) {

	var items []interface{}
	// Apply filter
	if filter != nil {
		for _, v := range imp._items {
			if filter(v) {
				items = append(items, v)
			}
		}
	} else {
		items = imp._items
	}

	rand.Seed(time.Now().UnixNano())

	if len(items) > 0 {
		item = &items[rand.Intn(len(items))]
	}

	if item != nil {
		imp._logger.Trace(correlationId, "Retrieved a random item")
	} else {
		imp._logger.Trace(correlationId, "Nothing to return as random item")
	}
	return item, nil
}

/*
Gets a data item by its unique id.

- correlationId     (optional) transaction id to trace execution through call chain.
- id                an id of data item to be retrieved.
- callback          callback function that receives data item or error.
*/
func (imp *IdentifiableMemoryPersistence) GetOneById(correlationId string, id interface{}) (item *interface{}, err error) {

	var items []interface{}
	for _, v := range imp._items {
		vId := refl.ObjectReader.GetProperty(v, "Id")
		if vId == id {
			items = append(items, v)
		}
	}

	if len(items) > 0 {
		item = &items[0]
	}

	if item != nil {
		imp._logger.Trace(correlationId, "Retrieved item %s", id)
	} else {
		imp._logger.Trace(correlationId, "Cannot find item by %s", id)
	}

	return item, err
}

/*
Creates a data item.

- correlation_id    (optional) transaction id to trace execution through call chain.
- item              an item to be created.
- callback          (optional) callback function that receives created item or error.
*/
func (imp *IdentifiableMemoryPersistence) Create(correlationId string, item interface{}) (result *interface{}, err error) {
	tmp := item
	copier.Copy(tmp, item)
	// Todo: need correct for add uniq id if not exist
	if reflect.ValueOf(tmp).FieldByName("Id").IsValid() {
		refl.ObjectWriter.SetProperty(tmp, "Id", cdata.IdGenerator.NextLong())
	}
	result = &tmp
	imp._items = append(imp._items, tmp)
	imp._logger.Trace(correlationId, "Created item %s", refl.ObjectReader.GetProperty(tmp, "Id"))
	errsave := imp.Save(correlationId)
	return result, errsave
}

/*
Sets a data item. If the data item exists it updates it,
otherwise it create a new data item.

- correlation_id    (optional) transaction id to trace execution through call chain.
- item              a item to be set.
- callback          (optional) callback function that receives updated item or error.
*/
func (imp *IdentifiableMemoryPersistence) Set(correlationId string, item interface{}) (result *interface{}, err error) {
	tmp := item
	copier.Copy(tmp, item)
	// ToDo: Need check
	if reflect.ValueOf(tmp).Elem().FieldByName("Id").IsValid() {
		refl.ObjectWriter.SetProperty(tmp, "Id", cdata.IdGenerator.NextLong())
	}

	var index int = -1
	itemId := refl.ObjectReader.GetProperty(item, "Id")
	for i, v := range imp._items {
		vId := refl.ObjectReader.GetProperty(v, "Id")
		if itemId == vId {
			index = i
			break
		}
	}

	if index < 0 {
		imp._items = append(imp._items, tmp)
	} else {
		imp._items[index] = tmp
	}

	imp._logger.Trace(correlationId, "Set item %s", refl.ObjectReader.GetProperty(tmp, "Id"))

	errsav := imp.Save(correlationId)
	return &tmp, errsav
}

/*
Updates a data item.

- correlation_id    (optional) transaction id to trace execution through call chain.
- item              an item to be updated.
- callback          (optional) callback function that receives updated item or error.
*/
func (imp *IdentifiableMemoryPersistence) Update(correlationId string, item interface{}) (result *interface{}, err error) {
	var index int = -1
	itemId := refl.ObjectReader.GetProperty(item, "Id")
	for i, v := range imp._items {
		vId := refl.ObjectReader.GetProperty(v, "Id")
		if itemId == vId {
			index = i
			break
		}
	}

	if index < 0 {
		imp._logger.Trace(correlationId, "Item %s was not found", refl.ObjectReader.GetProperty(item, "Id"))
		return nil, nil
	}
	tmp := item
	copier.Copy(tmp, item)
	imp._items[index] = item
	imp._logger.Trace(correlationId, "Updated item %s", refl.ObjectReader.GetProperty(item, "Id"))

	errsave := imp.Save(correlationId)
	return &tmp, errsave
}

/*
Updates only few selected fields in a data item.

- correlation_id    (optional) transaction id to trace execution through call chain.
- id                an id of data item to be updated.
- data              a map with fields to be updated.
- callback          callback function that receives updated item or error.
*/
func (imp *IdentifiableMemoryPersistence) UpdatePartially(correlationId string, id interface{}, data cdata.AnyValueMap) (item *interface{}, err error) {
	var index int = -1
	for i, v := range imp._items {
		vId := refl.ObjectReader.GetProperty(v, "Id")
		if vId == id {
			index = i
			break
		}
	}

	if index < 0 {
		imp._logger.Trace(correlationId, "Item %s was not found", id)
		return nil, nil
	}
	tmp := imp._items[index]
	copier.Copy(tmp, imp._items[index])
	// need test!!!
	refl.ObjectWriter.SetProperties(tmp, data.Value())

	copier.Copy(imp._items[index], tmp)
	imp._logger.Trace(correlationId, "Partially updated item %s", id)

	errsave := imp.Save(correlationId)
	return &tmp, errsave
}

/*
  Deleted a data item by it's unique id.

  - correlation_id    (optional) transaction id to trace execution through call chain.
  - id                an id of the item to be deleted
  - callback          (optional) callback function that receives deleted item or error.
*/
func (imp *IdentifiableMemoryPersistence) DeleteById(correlationId string, id interface{}) (item *interface{}, err error) {

	var index int = -1

	for i, v := range imp._items {
		vId := refl.ObjectReader.GetProperty(v, "Id")
		if vId == id {
			index = i
			break
		}
	}

	if index < 0 {
		imp._logger.Trace(correlationId, "Item %s was not found", id)
		return nil, nil
	} else {
		tmp := imp._items[index]
		copier.Copy(tmp, imp._items[index])
		item = &tmp
	}

	imp._items = append(imp._items[:index], imp._items[index+1:])
	imp._logger.Trace(correlationId, "Deleted item by %s", id)

	errsave := imp.Save(correlationId)
	return item, errsave
}

/**
Deletes data items that match to a given filter.
 *
imp method shall be called by a func (imp* IdentifiableMemoryPersistence) deleteByFilter method from child class that
receives FilterParams and converts them into a filter function.
 *
- correlationId     (optional) transaction id to trace execution through call chain.
- filter            (optional) a filter function to filter items.
- callback          (optional) callback function that receives error or null for success.
*/
func (imp *IdentifiableMemoryPersistence) DeleteByFilter(correlationId string, filter func(interface{}) bool) (err error) {
	deleted := 0
	for i, v := range imp._items {
		if filter(v) {
			imp._items = append(imp._items[:i], imp._items[i+1:])
			deleted++
		}
	}

	if deleted == 0 {
		return nil
	}

	imp._logger.Trace(correlationId, "Deleted %s items", deleted)

	errsave := imp.Save(correlationId)
	return errsave
}

/**
Deletes multiple data items by their unique ids.
 *
- correlationId     (optional) transaction id to trace execution through call chain.
- ids               ids of data items to be deleted.
- callback          (optional) callback function that receives error or null for success.
*/
func (imp *IdentifiableMemoryPersistence) DeleteByIds(correlationId string, ids []interface{}) (err error) {
	filter := func(item interface{}) bool {
		var exist bool = false

		itemId := refl.ObjectReader.GetProperty(item, "Id")
		for _, v := range ids {
			vId := refl.ObjectReader.GetProperty(v, "Id")
			if vId == itemId {
				exist = true
				break
			}
		}
		return exist
	}
	return imp.DeleteByFilter(correlationId, filter)
}
