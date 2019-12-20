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

func NewDummyMapMemoryPersistence() *DummyMapMemoryPersistence {
	var t map[string]interface{}
	proto := reflect.TypeOf(t)
	return &DummyMapMemoryPersistence{*cpersist.NewIdentifiableMemoryPersistence(proto)}
}

func (c *DummyMapMemoryPersistence) Create(correlationId string, item map[string]interface{}) (result map[string]interface{}, err error) {
	value, err := c.IdentifiableMemoryPersistence.Create(correlationId, item)
	if value != nil {
		val, _ := value.(map[string]interface{})
		result = val
	}
	return result, err
}

func (c *DummyMapMemoryPersistence) GetListByIds(correlationId string, ids []string) (items []map[string]interface{}, err error) {
	convIds := make([]interface{}, len(ids))
	for i, v := range ids {
		convIds[i] = v
	}
	result, err := c.IdentifiableMemoryPersistence.GetListByIds(correlationId, convIds)
	items = make([]map[string]interface{}, len(result))
	for i, v := range result {
		val, _ := v.(map[string]interface{})
		items[i] = val
	}
	return items, err
}

func (c *DummyMapMemoryPersistence) GetOneById(correlationId string, id string) (item map[string]interface{}, err error) {
	result, err := c.IdentifiableMemoryPersistence.GetOneById(correlationId, id)

	if result != nil {
		val, _ := result.(map[string]interface{})
		item = val
	}
	return item, err
}

func (c *DummyMapMemoryPersistence) Update(correlationId string, item map[string]interface{}) (result map[string]interface{}, err error) {
	value, err := c.IdentifiableMemoryPersistence.Update(correlationId, item)

	if value != nil {
		val, _ := value.(map[string]interface{})
		result = val
	}
	return result, err
}

func (c *DummyMapMemoryPersistence) UpdatePartially(correlationId string, id string, data cdata.AnyValueMap) (item map[string]interface{}, err error) {
	result, err := c.IdentifiableMemoryPersistence.UpdatePartially(correlationId, id, data)

	if result != nil {
		val, _ := result.(map[string]interface{})
		item = val
	}
	return item, err
}

func (c *DummyMapMemoryPersistence) DeleteById(correlationId string, id string) (item map[string]interface{}, err error) {
	result, err := c.IdentifiableMemoryPersistence.DeleteById(correlationId, id)

	if result != nil {
		val, _ := result.(map[string]interface{})
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

func (c *DummyMapMemoryPersistence) GetPageByFilter(correlationId string, filter cdata.FilterParams, paging cdata.PagingParams) (page MapPage, err error) {

	if &filter == nil {
		filter = *cdata.NewEmptyFilterParams()
	}

	key := filter.GetAsNullableString("Key")

	tempPage, err := c.IdentifiableMemoryPersistence.GetPageByFilter(correlationId, func(item interface{}) bool {
		dummy, ok := item.(map[string]interface{})
		if *key != "" && ok && dummy["Key"] != *key {
			return false
		}
		return true
	}, paging,
		func(a, b interface{}) bool {
			_a, _ := a.(map[string]interface{})
			_b, _ := b.(map[string]interface{})
			return len(_a["Key"].(string)) < len(_b["Key"].(string))
		}, nil)
	dataLen := int64(len(tempPage.Data))
	data := make([]map[string]interface{}, dataLen)
	for i, v := range tempPage.Data {
		data[i] = v.(map[string]interface{})
	}
	dataPage := *NewMapPage(&dataLen, data)
	return dataPage, err
}
