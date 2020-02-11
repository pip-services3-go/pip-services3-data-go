package test_persistence

import (
	"reflect"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cpersist "github.com/pip-services3-go/pip-services3-data-go/persistence"
)

// extends IdentifiableMemoryPersistence<Dummy, string>
// implements IDummyPersistence {
type DummyRefMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence
}

func NewDummyRefMemoryPersistence() *DummyRefMemoryPersistence {
	var t *Dummy
	proto := reflect.TypeOf(t)
	return &DummyRefMemoryPersistence{*cpersist.NewIdentifiableMemoryPersistence(proto)}
}

func (c *DummyRefMemoryPersistence) Create(correlationId string, item *Dummy) (result *Dummy, err error) {
	value, err := c.IdentifiableMemoryPersistence.Create(correlationId, item)

	if value != nil {
		val, _ := value.(*Dummy)
		result = val
	}
	return result, err
}

func (c *DummyRefMemoryPersistence) GetListByIds(correlationId string, ids []string) (items []*Dummy, err error) {
	convIds := make([]interface{}, len(ids))
	for i, v := range ids {
		convIds[i] = v
	}
	result, err := c.IdentifiableMemoryPersistence.GetListByIds(correlationId, convIds)
	items = make([]*Dummy, len(result))
	for i, v := range result {
		val, _ := v.(*Dummy)
		items[i] = val
	}
	return items, err
}

func (c *DummyRefMemoryPersistence) GetOneById(correlationId string, id string) (item *Dummy, err error) {
	result, err := c.IdentifiableMemoryPersistence.GetOneById(correlationId, id)
	if result != nil {
		val, _ := result.(*Dummy)
		item = val
	}
	return item, err
}

func (c *DummyRefMemoryPersistence) Update(correlationId string, item *Dummy) (result *Dummy, err error) {
	value, err := c.IdentifiableMemoryPersistence.Update(correlationId, item)
	if value != nil {
		val, _ := value.(*Dummy)
		result = val
	}
	return result, err
}

func (c *DummyRefMemoryPersistence) UpdatePartially(correlationId string, id string, data *cdata.AnyValueMap) (item *Dummy, err error) {
	result, err := c.IdentifiableMemoryPersistence.UpdatePartially(correlationId, id, data)

	if result != nil {
		val, _ := result.(*Dummy)
		item = val
	}
	return item, err
}

func (c *DummyRefMemoryPersistence) DeleteById(correlationId string, id string) (item *Dummy, err error) {
	result, err := c.IdentifiableMemoryPersistence.DeleteById(correlationId, id)
	if result != nil {
		val, _ := result.(*Dummy)
		item = val
	}
	return item, err
}

func (c *DummyRefMemoryPersistence) DeleteByIds(correlationId string, ids []string) (err error) {
	convIds := make([]interface{}, len(ids))
	for i, v := range ids {
		convIds[i] = v
	}
	return c.IdentifiableMemoryPersistence.DeleteByIds(correlationId, convIds)
}

func (c *DummyRefMemoryPersistence) GetPageByFilter(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (page *DummyRefPage, err error) {

	if &filter == nil {
		filter = cdata.NewEmptyFilterParams()
	}

	key := filter.GetAsNullableString("Key")

	tempPage, err := c.IdentifiableMemoryPersistence.GetPageByFilter(correlationId, func(item interface{}) bool {
		dummy, ok := item.(Dummy)
		if *key != "" && ok && (dummy).Key != *key {
			return false
		}
		return true
	}, paging,
		func(a, b interface{}) bool {
			_a, _ := a.(Dummy)
			_b, _ := b.(Dummy)
			return len((_a).Key) < len((_b).Key)
		}, nil)
	// Convert to DummyRefPage
	dataLen := int64(len(tempPage.Data)) // For full release tempPage and delete this by GC
	data := make([]*Dummy, dataLen)
	for i := range tempPage.Data {
		temp := tempPage.Data[i].(*Dummy)
		data[i] = temp
	}
	page = NewDummyRefPage(&dataLen, data)
	return page, err
}
