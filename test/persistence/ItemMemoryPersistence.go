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
