package test_persistence

import (
	"reflect"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cpersist "github.com/pip-services3-go/pip-services3-data-go/persistence"
)

// extends IdentifiableMemoryPersistence
// implements IItemPersistence
type ItemMemoryPersistence struct {
	cpersist.IdentifiableMemoryPersistence
}

func NewItemMemoryPersistence() *ItemMemoryPersistence {
	proto := reflect.TypeOf(Item{})
	return &ItemMemoryPersistence{*cpersist.NewIdentifiableMemoryPersistence(proto)}
}

func (c *ItemMemoryPersistence) Create(correlationId string, item Item) (result Item, err error) {
	value, err := c.IdentifiableMemoryPersistence.Create(correlationId, item)
	if value != nil {
		val, _ := value.(Item)
		result = val
	}
	return result, err
}

func (c *ItemMemoryPersistence) GetOneById(correlationId string, id string) (item Item, err error) {
	result, err := c.IdentifiableMemoryPersistence.GetOneById(correlationId, id)
	if result != nil {
		val, _ := result.(Item)
		item = val
	}
	return item, err
}

func (c *ItemMemoryPersistence) Update(correlationId string, item Item) (result Item, err error) {
	value, err := c.IdentifiableMemoryPersistence.Update(correlationId, item)
	if value != nil {
		val, _ := value.(Item)
		result = val
	}
	return result, err
}

func (c *ItemMemoryPersistence) UpdatePartially(correlationId string, id string, data *cdata.AnyValueMap) (item Item, err error) {
	result, err := c.IdentifiableMemoryPersistence.UpdatePartially(correlationId, id, data)

	if result != nil {
		val, _ := result.(Item)
		item = val
	}
	return item, err
}

func (c *ItemMemoryPersistence) DeleteById(correlationId string, id string) (item Item, err error) {
	result, err := c.IdentifiableMemoryPersistence.DeleteById(correlationId, id)
	if result != nil {
		val, _ := result.(Item)
		item = val
	}
	return item, err
}

func (c *ItemMemoryPersistence) GetPageByFilter(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams, sort bool) (page *ItemPage, err error) {

	if filter == nil {
		filter = cdata.NewEmptyFilterParams()
	}

	updatedBy := filter.GetAsNullableString("UpdatedBy")

	filterFn := func(item interface{}) bool {
		Item, ok := item.(Item)
		if updatedBy != nil && ok && Item.UpdatedBy != *updatedBy {
			return false
		}
		return true
	}

	sortFn := func(a interface{}, b interface{}) bool {

		if a == nil || b == nil {
			return false
		}

		if len(a.(Item).UpdatedBy) > len(b.(Item).UpdatedBy) {
			return true
		}
		return false
	}

	if !sort {
		sortFn = nil
	}

	tempPage, err := c.IdentifiableMemoryPersistence.GetPageByFilter(correlationId, filterFn, paging, sortFn, nil)
	// Convert to ItemPage
	dataLen := int64(len(tempPage.Data)) // For full release tempPage and delete this by GC
	data := make([]Item, dataLen)
	for i, v := range tempPage.Data {
		data[i] = v.(Item)
	}
	total := (int64)(0)
	if tempPage.Total != nil {
		total = *tempPage.Total
	}
	page = NewItemPage(&total, data)
	return page, err
}
