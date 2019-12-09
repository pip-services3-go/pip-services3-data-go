package test_persistence

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cpersist "github.com/pip-services3-go/pip-services3-data-go/persistence"
)

// extends IdentifiableMemoryPersistence<Dummy, string>
// implements IDummyPersistence {
type DummyMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence
}

func (c *DummyMemoryPersistence) NewDummyMemoryPersistence() *DummyMemoryPersistence {
	return &DummyMemoryPersistence{}
}

func (c *DummyMemoryPersistence) GetPageByFilter(correlationId string, filter cdata.FilterParams, paging cdata.PagingParams) (page cdata.DataPage, err error) {

	if &filter == nil {
		filter = *cdata.NewEmptyFilterParams()
	}

	key := filter.GetAsNullableString("key")

	return c.IdentifiableMemoryPersistence.getPageByFilter(correlationId, func(item Dummy) bool {
		if key != nil && item.key != *key {
			return false
		}
		return true
	}, paging, nil, nil)
}
