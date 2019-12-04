package test

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	//"github.com/pip-services3-go/pip-services3-data-go/persistence"
)

// extends IdentifiableMemoryPersistence<Dummy, string>
// implements IDummyPersistence {
type DummyMemoryPersistence struct {
	IdentifiableMemoryPersistence
}

func (dmp *DummyMemoryPersistence) NewDummyMemoryPersistence() *DummyMemoryPersistence {
	return DummyMemoryPersistence{}
}

func (dmp *DummyMemoryPersistence) GetPageByFilter(correlationId string, filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage, err error) {

	if &filter == nil {
		filter = *cdata.NewEmptyFilterParams()
	}

	key := filter.GetAsNullableString("key")

	return dmp.IdentifiableMemoryPersistence.getPageByFilter(correlationId, func(item Dummy) bool {
		if key != nil && item.key != *key {
			return false
		}
		return true
	}, paging, nil, nil)
}
