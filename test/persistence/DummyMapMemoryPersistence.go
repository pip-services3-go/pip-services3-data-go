package test_persistence

import (
	"reflect"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cpersist "github.com/pip-services3-go/pip-services3-data-go/persistence"
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

func (c *DummyMapMemoryPersistence) toPublicPage(page *cdata.DataPage) *MapPage {
	if page == nil {
		return nil
	}

	dataLen := int64(len(page.Data))
	data := make([]map[string]interface{}, dataLen)
	for i, v := range page.Data {
		data[i] = cpersist.ToPublicMap(v)
	}
	dataPage := NewMapPage(&dataLen, data)

	return dataPage
}

func (c *DummyMapMemoryPersistence) Create(correlationId string, item map[string]interface{}) (result map[string]interface{}, err error) {
	value, err := c.IdentifiableMemoryPersistence.Create(correlationId, item)
	result = cpersist.ToPublicMap(value)
	return result, err
}

func (c *DummyMapMemoryPersistence) GetListByIds(correlationId string, ids []string) (result []map[string]interface{}, err error) {
	convIds := cpersist.FromIds(ids)
	values, err := c.IdentifiableMemoryPersistence.GetListByIds(correlationId, convIds)
	result = cpersist.ToPublicArray(values)
	return result, err
}

func (c *DummyMapMemoryPersistence) GetOneById(correlationId string, id string) (result map[string]interface{}, err error) {
	value, err := c.IdentifiableMemoryPersistence.GetOneById(correlationId, id)
	result = cpersist.ToPublicMap(value)
	return result, err
}

func (c *DummyMapMemoryPersistence) Update(correlationId string, item map[string]interface{}) (result map[string]interface{}, err error) {
	value, err := c.IdentifiableMemoryPersistence.Update(correlationId, item)
	result = cpersist.ToPublicMap(value)
	return result, err
}

func (c *DummyMapMemoryPersistence) UpdatePartially(correlationId string, id string, data *cdata.AnyValueMap) (result map[string]interface{}, err error) {
	value, err := c.IdentifiableMemoryPersistence.UpdatePartially(correlationId, id, data)
	result = cpersist.ToPublicMap(value)
	return result, err
}

func (c *DummyMapMemoryPersistence) DeleteById(correlationId string, id string) (result map[string]interface{}, err error) {
	value, err := c.IdentifiableMemoryPersistence.DeleteById(correlationId, id)
	result = cpersist.ToPublicMap(value)
	return result, err
}

func (c *DummyMapMemoryPersistence) DeleteByIds(correlationId string, ids []string) (err error) {
	convIds := cpersist.FromIds(ids)
	return c.IdentifiableMemoryPersistence.DeleteByIds(correlationId, convIds)
}

func filterFunc(filter *cdata.FilterParams) func(interface{}) bool {
	if &filter == nil {
		filter = cdata.NewEmptyFilterParams()
	}

	key := filter.GetAsNullableString("Key")

	return func(value interface{}) bool {
		dummy, ok := value.(map[string]interface{})
		if *key != "" && ok && dummy["Key"] != *key {
			return false
		}
		return true
	}
}

func sortFunc(value1, value2 interface{}) bool {
	v1, _ := value1.(map[string]interface{})
	v2, _ := value2.(map[string]interface{})
	// Todo: Why len instead of string comparison?
	return len(v1["Key"].(string)) < len(v2["Key"].(string))
}

func (c *DummyMapMemoryPersistence) GetPageByFilter(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *MapPage, err error) {
	page, err := c.IdentifiableMemoryPersistence.GetPageByFilter(correlationId, filterFunc(filter), paging, sortFunc, nil)
	result = c.toPublicPage(page)
	return result, err
}

func (c *DummyMapMemoryPersistence) GetCountByFilter(correlationId string, filter *cdata.FilterParams) (count int64, err error) {
	if &filter == nil {
		filter = cdata.NewEmptyFilterParams()
	}
	key := filter.GetAsNullableString("Key")
	count, err = c.IdentifiableMemoryPersistence.GetCountByFilter(correlationId, func(item interface{}) bool {
		dummy, ok := item.(Dummy)
		if *key != "" && ok && dummy.Key != *key {
			return false
		}
		return true
	})
	return count, err
}
