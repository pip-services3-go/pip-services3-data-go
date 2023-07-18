package test_persistence

import (
	"testing"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	"github.com/stretchr/testify/assert"
)

type ItemPersistenceFixture struct {
	item1       Item
	item2       Item
	item3       Item
	item4       Item
	persistence IItemPersistence
}

func NewItemPersistenceFixture(persistence IItemPersistence) *ItemPersistenceFixture {
	c := ItemPersistenceFixture{}
	c.item1 = Item{
		Id:                        "",
		FailingToUpdateThisField1: 0,
		FailingToUpdateThisField2: 0,
		UpdatedBy:                 "y76bh72c1",
	}

	c.item2 = Item{
		Id:                        "",
		FailingToUpdateThisField1: 0,
		FailingToUpdateThisField2: 0,
		UpdatedBy:                 "y76bh72c12",
	}

	c.item3 = Item{
		Id:                        "",
		FailingToUpdateThisField1: 0,
		FailingToUpdateThisField2: 0,
		UpdatedBy:                 "y76bh72c123",
	}

	c.item4 = Item{
		Id:                        "",
		FailingToUpdateThisField1: 0,
		FailingToUpdateThisField2: 0,
		UpdatedBy:                 "y76bh72c1234",
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
	assert.Equal(t, c.item1.FailingToUpdateThisField1, item1.FailingToUpdateThisField1)
	assert.Equal(t, c.item1.FailingToUpdateThisField2, item1.FailingToUpdateThisField2)
	assert.Equal(t, c.item1.UpdatedBy, item1.UpdatedBy)

	intvar1 := (int64)(1234757257822780121)
	intvar2 := (int64)(1434824722285792000)
	// Partially update the item
	updateMap := cdata.NewAnyValueMapFromTuples(
		"FailingToUpdateThisField1", intvar1,
		"updated_by", "upd_y76bh72c1",
		"failing_to_update_this_field_2", intvar2,
	)
	result, err = c.persistence.UpdatePartially("", item1.Id, updateMap)
	if err != nil {
		t.Errorf("UpdatePartially method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, item1.Id, result.Id)
	assert.Equal(t, (int64)(1234757257822780121), result.FailingToUpdateThisField1)
	assert.Equal(t, (int64)(1434824722285792000), result.FailingToUpdateThisField2)
	assert.Equal(t, "upd_y76bh72c1", result.UpdatedBy)

}

func (c *ItemPersistenceFixture) TestSortOperations(t *testing.T) {

	// first item
	item, err := c.persistence.Create("", c.item1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}

	assert.NotNil(t, item)
	assert.NotNil(t, item.Id)
	assert.Equal(t, c.item1.FailingToUpdateThisField1, item.FailingToUpdateThisField1)
	assert.Equal(t, c.item1.FailingToUpdateThisField2, item.FailingToUpdateThisField2)
	assert.Equal(t, c.item1.UpdatedBy, item.UpdatedBy)

	// second item
	item, err = c.persistence.Create("", c.item2)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}

	assert.NotNil(t, item)
	assert.NotNil(t, item.Id)
	assert.Equal(t, c.item2.FailingToUpdateThisField1, item.FailingToUpdateThisField1)
	assert.Equal(t, c.item2.FailingToUpdateThisField2, item.FailingToUpdateThisField2)
	assert.Equal(t, c.item2.UpdatedBy, item.UpdatedBy)

	// third item
	item, err = c.persistence.Create("", c.item3)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}

	assert.NotNil(t, item)
	assert.NotNil(t, item.Id)
	assert.Equal(t, c.item3.FailingToUpdateThisField1, item.FailingToUpdateThisField1)
	assert.Equal(t, c.item3.FailingToUpdateThisField2, item.FailingToUpdateThisField2)
	assert.Equal(t, c.item3.UpdatedBy, item.UpdatedBy)

	// forth item
	item, err = c.persistence.Create("", c.item4)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}

	assert.NotNil(t, item)
	assert.NotNil(t, item.Id)
	assert.Equal(t, c.item4.FailingToUpdateThisField1, item.FailingToUpdateThisField1)
	assert.Equal(t, c.item4.FailingToUpdateThisField2, item.FailingToUpdateThisField2)
	assert.Equal(t, c.item4.UpdatedBy, item.UpdatedBy)

	// get all without sorting
	page, errp := c.persistence.GetPageByFilter("", cdata.NewEmptyFilterParams(), cdata.NewEmptyPagingParams(), false)
	if errp != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 4)

	assert.Equal(t, c.item1.UpdatedBy, page.Data[0].UpdatedBy)
	assert.Equal(t, c.item2.UpdatedBy, page.Data[1].UpdatedBy)
	assert.Equal(t, c.item3.UpdatedBy, page.Data[2].UpdatedBy)
	assert.Equal(t, c.item4.UpdatedBy, page.Data[3].UpdatedBy)

	// get all with sort
	page, errp = c.persistence.GetPageByFilter("", cdata.NewEmptyFilterParams(), cdata.NewEmptyPagingParams(), true)
	if errp != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 4)

	assert.Equal(t, c.item4.UpdatedBy, page.Data[0].UpdatedBy)
	assert.Equal(t, c.item3.UpdatedBy, page.Data[1].UpdatedBy)
	assert.Equal(t, c.item2.UpdatedBy, page.Data[2].UpdatedBy)
	assert.Equal(t, c.item1.UpdatedBy, page.Data[3].UpdatedBy)

}
