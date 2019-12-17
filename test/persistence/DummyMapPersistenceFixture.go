package test_persistence

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

type DummyMapPersistenceFixture struct {
	_dummy1      DummyMap
	_dummy2      DummyMap
	_persistence IDummyMapPersistence
}

func NewDummyMapPersistenceFixture(persistence IDummyMapPersistence) *DummyMapPersistenceFixture {
	c := DummyMapPersistenceFixture{}
	c._dummy1 = DummyMap{"Id": "", "Key": "Key 11", "Content": "Content 1"}
	c._dummy2 = DummyMap{"Id": "", "Key": "Key 2", "Content": "Content 2"}
	c._persistence = persistence
	return &c
}

func (c *DummyMapPersistenceFixture) TestCrudOperations(t *testing.T) {
	var dummy1 DummyMap
	var dummy2 DummyMap

	result, err := c._persistence.Create("", c._dummy1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy1 = *result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1["Id"])
	assert.Equal(t, c._dummy1["Key"], dummy1["Key"])
	assert.Equal(t, c._dummy1["Content"], dummy1["Content"])

	// Create another dummy
	result, err = c._persistence.Create("", c._dummy2)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy2 = *result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2["Id"])
	assert.Equal(t, c._dummy2["Key"], dummy2["Key"])
	assert.Equal(t, c._dummy2["Content"], dummy2["Content"])

	page, errp := c._persistence.GetPageByFilter("", *cdata.NewEmptyFilterParams(), *cdata.NewEmptyPagingParams())
	if errp != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 2)
	//Testing default sorting by Key field len

	item1 := page.Data[0].(DummyMap)
	assert.Equal(t, item1["Key"], dummy2["Key"])
	item2 := page.Data[1].(DummyMap)
	assert.Equal(t, item2["Key"], dummy1["Key"])

	// Update the dummy
	dummy1["Content"] = "Updated Content 1"
	result, err = c._persistence.Update("", dummy1)
	if err != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1["Id"], (*result)["Id"])
	assert.Equal(t, dummy1["Key"], (*result)["Key"])
	assert.Equal(t, dummy1["Content"], (*result)["Content"])

	// Partially update the dummy
	updateMap := *cdata.NewAnyValueMapFromTuples("Content", "Partially Updated Content 1")
	result, err = c._persistence.UpdatePartially("", dummy1["Id"], updateMap)
	if err != nil {
		t.Errorf("UpdatePartially method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1["Id"], (*result)["Id"])
	assert.Equal(t, dummy1["Key"], (*result)["Key"])
	assert.Equal(t, "Partially Updated Content 1", (*result)["Content"])

	// Get the dummy by Id
	result, err = c._persistence.GetOneById("", dummy1["Id"])
	if err != nil {
		t.Errorf("GetOneById method error %v", err)
	}
	// Try to get item
	assert.NotNil(t, result)
	assert.Equal(t, dummy1["Id"], (*result)["Id"])
	assert.Equal(t, dummy1["Key"], (*result)["Key"])
	assert.Equal(t, "Partially Updated Content 1", (*result)["Content"])

	// Delete the dummy
	result, err = c._persistence.DeleteById("", dummy1["Id"])
	if err != nil {
		t.Errorf("DeleteById method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1["Id"], (*result)["Id"])
	assert.Equal(t, dummy1["Key"], (*result)["Key"])
	assert.Equal(t, "Partially Updated Content 1", (*result)["Content"])

	// Get the deleted dummy
	result, err = c._persistence.GetOneById("", dummy1["Id"])
	if err != nil {
		t.Errorf("GetOneById method error %v", err)
	}
	// Try to get item
	assert.Nil(t, result)
}

func (c *DummyMapPersistenceFixture) TestBatchOperations(t *testing.T) {
	var dummy1 DummyMap
	var dummy2 DummyMap

	// Create one dummy
	result, err := c._persistence.Create("", c._dummy1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy1 = *result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1["Id"])
	assert.Equal(t, c._dummy1["Key"], dummy1["Key"])
	assert.Equal(t, c._dummy1["Content"], dummy1["Content"])

	// Create another dummy
	result, err = c._persistence.Create("", c._dummy2)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy2 = *result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2["Id"])
	assert.Equal(t, c._dummy2["Key"], dummy2["Key"])
	assert.Equal(t, c._dummy2["Content"], dummy2["Content"])

	// Read batch
	items, err := c._persistence.GetListByIds("", []string{dummy1["Id"], dummy2["Id"]})
	if err != nil {
		t.Errorf("GetListByIds method error %v", err)
	}
	//assert.isArray(t,items)
	assert.NotNil(t, items)
	assert.Len(t, items, 2)

	// Delete batch
	err = c._persistence.DeleteByIds("", []string{dummy1["Id"], dummy2["Id"]})
	if err != nil {
		t.Errorf("DeleteByIds method error %v", err)
	}
	assert.Nil(t, err)

	// Read empty batch
	items, err = c._persistence.GetListByIds("", []string{dummy1["Id"], dummy2["Id"]})
	if err != nil {
		t.Errorf("GetListByIds method error %v", err)
	}
	assert.NotNil(t, items)
	assert.Len(t, items, 0)

}
