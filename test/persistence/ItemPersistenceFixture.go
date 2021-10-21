package test_persistence

import (
	"testing"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	"github.com/stretchr/testify/assert"
)

type ItemPersistenceFixture struct {
	item1       Item
	persistence IItemPersistence
}

func NewItemPersistenceFixture(persistence IItemPersistence) *ItemPersistenceFixture {
	c := ItemPersistenceFixture{}
	c.item1 = Item{
		Id:                             "",
		Failing_to_update_this_field_1: 1434824722285792000,
		Failing_to_update_this_field_2: 1234757257822780121,
		Updated_by:                     "y76bh72c1",
	}
	c.persistence = persistence
	return &c
}

func (c *ItemPersistenceFixture) TestCrudOperations(t *testing.T) {
	var item1 Item

	result, err := c.persistence.Create("", c.item1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	item1 = result
	assert.NotNil(t, item1)
	assert.NotNil(t, item1.Id)
	assert.Equal(t, c.item1.Failing_to_update_this_field_1, item1.Failing_to_update_this_field_1)
	assert.Equal(t, c.item1.Failing_to_update_this_field_2, item1.Failing_to_update_this_field_2)
	assert.Equal(t, c.item1.Updated_by, item1.Updated_by)

	// Partially update the item
	updateMap := cdata.NewAnyValueMapFromTuples(
		"failing_to_update_this_field_1", (int64)(1234757257822780121),
		"updated_by", "upd_y76bh72c1",
		"failing_to_update_this_field_2", (int64)(1434824722285792000),
	)
	result, err = c.persistence.UpdatePartially("", item1.Id, updateMap)
	if err != nil {
		t.Errorf("UpdatePartially method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, item1.Id, result.Id)
	assert.Equal(t, (int64)(1234757257822780121), result.Failing_to_update_this_field_1)
	assert.Equal(t, (int64)(1434824722285792000), result.Failing_to_update_this_field_2)
	assert.Equal(t, "upd_y76bh72c1", result.Updated_by)

}
