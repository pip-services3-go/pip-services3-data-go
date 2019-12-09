package test_persistence

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

type DummyPersistenceFixture struct {
	_dummy1      Dummy
	_dummy2      Dummy
	_persistence IDummyPersistence
}

func NewDummyPersistenceFixture(persistence IDummyPersistence) *DummyPersistenceFixture {
	dpf := DummyPersistenceFixture{}
	dpf._dummy1 = Dummy{Id: "id1", Key: "Key 1", Content: "Content 1"}
	dpf._dummy2 = Dummy{Id: "id2", Key: "Key 2", Content: "Content 2"}
	dpf._persistence = persistence
	return &dpf
}

func (dpf *DummyPersistenceFixture) TestCrudOperations(t *testing.T) {
	var dummy1 Dummy
	var dummy2 Dummy

	result, err := dpf._persistence.Create("", dpf._dummy1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy1 = *result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1.Id)
	assert.Equal(t, dpf._dummy1.Key, dummy1.Key)
	assert.Equal(t, dpf._dummy1.Content, dummy1.Content)

	// Create another dummy
	result, err = dpf._persistence.Create("", dpf._dummy2)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy2 = *result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2.Id)
	assert.Equal(t, dpf._dummy2.Key, dummy2.Key)
	assert.Equal(t, dpf._dummy2.Content, dummy2.Content)

	page, errp := dpf._persistence.GetPageByFilter("", *cdata.NewEmptyFilterParams(), *cdata.NewEmptyPagingParams())
	if errp != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 2)

	// Update the dummy
	dummy1.Content = "Updated Content 1"
	result, err = dpf._persistence.Update("", dummy1)
	if err != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.Id, result.Id)
	assert.Equal(t, dummy1.Key, result.Key)
	assert.Equal(t, dummy1.Content, result.Content)

	// Partially update the dummy
	updateMap := *cdata.NewAnyValueMapFromTuples("Content", "Partially Updated Content 1")
	result, err = dpf._persistence.UpdatePartially("", dummy1.Id, updateMap)
	if err != nil {
		t.Errorf("UpdatePartially method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.Id, result.Id)
	assert.Equal(t, dummy1.Key, result.Key)
	assert.Equal(t, "Partially Updated Content 1", result.Content)

	// Get the dummy by Id
	result, err = dpf._persistence.GetOneById("", dummy1.Id)
	if err != nil {
		t.Errorf("GetOneById method error %v", err)
	}
	// Try to get item
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.Id, result.Id)
	assert.Equal(t, dummy1.Key, result.Key)
	assert.Equal(t, "Partially Updated Content 1", result.Content)

	// Delete the dummy
	result, err = dpf._persistence.DeleteById("", dummy1.Id)
	if err != nil {
		t.Errorf("DeleteById method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.Id, result.Id)
	assert.Equal(t, dummy1.Key, result.Key)
	assert.Equal(t, "Partially Updated Content 1", result.Content)

	// Get the deleted dummy
	result, err = dpf._persistence.GetOneById("", dummy1.Id)
	if err != nil {
		t.Errorf("GetOneById method error %v", err)
	}
	// Try to get item
	assert.Nil(t, result)
}

func (dpf *DummyPersistenceFixture) TestBatchOperations(t *testing.T) {
	var dummy1 Dummy
	var dummy2 Dummy

	// Create one dummy
	result, err := dpf._persistence.Create("", dpf._dummy1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy1 = *result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1.Id)
	assert.Equal(t, dpf._dummy1.Key, dummy1.Key)
	assert.Equal(t, dpf._dummy1.Content, dummy1.Content)

	// Create another dummy
	result, err = dpf._persistence.Create("", dpf._dummy2)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy2 = *result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2.Id)
	assert.Equal(t, dpf._dummy2.Key, dummy2.Key)
	assert.Equal(t, dpf._dummy2.Content, dummy2.Content)

	// Read batch
	items, err := dpf._persistence.GetListByIds("", []string{dummy1.Id, dummy2.Id})
	if err != nil {
		t.Errorf("GetListByIds method error %v", err)
	}
	//assert.isArray(t,items)
	assert.NotNil(t, items)
	assert.Len(t, items, 2)

	// Delete batch
	err = dpf._persistence.DeleteByIds("", []string{dummy1.Id, dummy2.Id})
	if err != nil {
		t.Errorf("DeleteByIds method error %v", err)
	}
	assert.Nil(t, err)

	// Read empty batch
	items, err = dpf._persistence.GetListByIds("", []string{dummy1.Id, dummy2.Id})
	if err != nil {
		t.Errorf("GetListByIds method error %v", err)
	}
	//assert.isArray(items)
	assert.NotNil(t, items)
	assert.Len(t, items, 0)

}
