package test_persistence

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	"github.com/pip-services3-go/pip-services3-data-go/persistence"
)

// extends IdentifiableMemoryPersistence<Dummy, string>
// implements IDummyPersistence {
type DummyMemoryPersistence struct {
	persistence.IdentifiableMemoryPersistence
}

func NewDummyMemoryPersistence() *DummyMemoryPersistence {
	return &DummyMemoryPersistence{*persistence.NewEmptyIdentifiableMemoryPersistence()}
}

func (dmp *DummyMemoryPersistence) Create(correlationId string, item Dummy) (result *Dummy, err error) {
	value, err := dmp.IdentifiableMemoryPersistence.Create(correlationId, item)
	result = nil
	if value != nil {
		val, _ := (*value).(Dummy)
		result = &val
	}
	return
}

func (dmp *DummyMemoryPersistence) GetListByIds(correlationId string, ids []string) (items []Dummy, err error) {
	convIds := make([]interface{}, len(ids))
	for i, v := range ids {
		convIds[i] = v
	}
	result, err := dmp.IdentifiableMemoryPersistence.GetListByIds(correlationId, convIds)
	items = make([]Dummy, len(result))
	for i, v := range result {
		val, _ := v.(Dummy)
		items[i] = val
	}
	return
}

func (dmp *DummyMemoryPersistence) GetOneById(correlationId string, id string) (item *Dummy, err error) {
	result, err := dmp.IdentifiableMemoryPersistence.GetOneById(correlationId, id)
	item = nil
	if result != nil {
		val, _ := (*result).(Dummy)
		item = &val
	}
	return
}

func (dmp *DummyMemoryPersistence) Update(correlationId string, item Dummy) (result *Dummy, err error) {
	value, err := dmp.IdentifiableMemoryPersistence.Update(correlationId, item)
	result = nil
	if value != nil {
		val, _ := (*value).(Dummy)
		result = &val
	}
	return
}

func (dmp *DummyMemoryPersistence) UpdatePartially(correlationId string, id string, data cdata.AnyValueMap) (item *Dummy, err error) {
	result, err := dmp.IdentifiableMemoryPersistence.UpdatePartially(correlationId, id, data)
	item = nil
	if result != nil {
		val, _ := (*result).(Dummy)
		item = &val
	}
	return
}

func (dmp *DummyMemoryPersistence) DeleteById(correlationId string, id string) (item *Dummy, err error) {
	result, err := dmp.IdentifiableMemoryPersistence.DeleteById(correlationId, id)
	item = nil
	if result != nil {
		val, _ := (*result).(Dummy)
		item = &val
	}
	return
}

func (dmp *DummyMemoryPersistence) DeleteByIds(correlationId string, ids []string) (err error) {
	convIds := make([]interface{}, len(ids))
	for i, v := range ids {
		convIds[i] = v
	}
	return dmp.IdentifiableMemoryPersistence.DeleteByIds(correlationId, convIds)
}

func (dmp *DummyMemoryPersistence) GetPageByFilter(correlationId string, filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage, err error) {

	if &filter == nil {
		filter = *cdata.NewEmptyFilterParams()
	}

	key := filter.GetAsNullableString("Key")

	return dmp.IdentifiableMemoryPersistence.GetPageByFilter(correlationId, func(item interface{}) bool {
		dummy, ok := item.(Dummy)
		if *key != "" && ok && dummy.Key != *key {
			return false
		}
		return true
	}, paging, nil, nil)
}
