package test_persistence

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cpersist "github.com/pip-services3-go/pip-services3-data-go/persistence"
	"reflect"
)

// extends IdentifiableMemoryPersistence<DummyMap, string>
// implements IDummyMapPersistence {
type DummyMapMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence
}

func NewEmptyDummyMapMemoryPersistence() *DummyMapMemoryPersistence {
	proto := reflect.TypeOf(DummyMap{})
	return &DummyMapMemoryPersistence{*cpersist.NewEmptyIdentifiableMemoryPersistence(proto)}
}

func NewDummyMapMemoryPersistence(loader cpersist.ILoader, saver cpersist.ISaver) *DummyMapMemoryPersistence {
	proto := reflect.TypeOf(DummyMap{})
	return &DummyMapMemoryPersistence{*cpersist.NewIdentifiableMemoryPersistence(proto, loader, saver)}
}

func (c *DummyMapMemoryPersistence) Create(correlationId string, item interface{}) (result DummyMap, err error) {
	value, err := c.IdentifiableMemoryPersistence.Create(correlationId, item)
	//result = nil
	if value != nil {
		val, _ := value.(DummyMap)
		result = val
	}
	return result, err
}

func (c *DummyMapMemoryPersistence) GetListByIds(correlationId string, ids []string) (items []DummyMap, err error) {
	convIds := make([]interface{}, len(ids))
	for i, v := range ids {
		convIds[i] = v
	}
	result, err := c.IdentifiableMemoryPersistence.GetListByIds(correlationId, convIds)
	items = make([]DummyMap, len(result))
	for i, v := range result {
		val, _ := v.(DummyMap)
		items[i] = val
	}
	return items, err
}

func (c *DummyMapMemoryPersistence) GetOneById(correlationId string, id string) (item DummyMap, err error) {
	result, err := c.IdentifiableMemoryPersistence.GetOneById(correlationId, id)
	//item = nil
	if result != nil {
		val, _ := result.(DummyMap)
		item = val
	}
	return item, err
}

func (c *DummyMapMemoryPersistence) Update(correlationId string, item interface{}) (result DummyMap, err error) {
	value, err := c.IdentifiableMemoryPersistence.Update(correlationId, item)
	//result = nil
	if value != nil {
		val, _ := value.(DummyMap)
		result = val
	}
	return result, err
}

func (c *DummyMapMemoryPersistence) UpdatePartially(correlationId string, id string, data cdata.AnyValueMap) (item DummyMap, err error) {
	result, err := c.IdentifiableMemoryPersistence.UpdatePartially(correlationId, id, data)
	//item = nil
	if result != nil {
		val, _ := result.(DummyMap)
		item = val
	}
	return item, err
}

func (c *DummyMapMemoryPersistence) DeleteById(correlationId string, id string) (item DummyMap, err error) {
	result, err := c.IdentifiableMemoryPersistence.DeleteById(correlationId, id)
	//item = nil
	if result != nil {
		val, _ := result.(DummyMap)
		item = val
	}
	return item, err
}

func (c *DummyMapMemoryPersistence) DeleteByIds(correlationId string, ids []string) (err error) {
	convIds := make([]interface{}, len(ids))
	for i, v := range ids {
		convIds[i] = v
	}
	return c.IdentifiableMemoryPersistence.DeleteByIds(correlationId, convIds)
}

func (c *DummyMapMemoryPersistence) GetPageByFilter(correlationId string, filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage, err error) {

	if &filter == nil {
		filter = *cdata.NewEmptyFilterParams()
	}

	key := filter.GetAsNullableString("Key")

	return c.IdentifiableMemoryPersistence.GetPageByFilter(correlationId, func(item interface{}) bool {
		dummy, ok := item.(DummyMap)
		if *key != "" && ok && dummy["Key"] != *key {
			return false
		}
		return true
	}, paging,
		func(a, b interface{}) bool {
			_a, _ := a.(DummyMap)
			_b, _ := b.(DummyMap)
			return len(_a["Key"]) < len(_b["Key"])
		}, nil)
}
