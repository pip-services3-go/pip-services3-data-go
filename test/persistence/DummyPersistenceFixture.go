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
	dpf._dummy1 = Dummy{id: "", key: "Key 1", content: "Content 1"}
	dpf._dummy2 = Dummy{id: "", key: "Key 2", content: "Content 2"}
	dpf._persistence = persistence
	return &dpf
}

func (dpf *DummyPersistenceFixture) testCrudOperations(t *testing.T) {
	var dummy1 Dummy
	var dummy2 Dummy

	result, err := dpf._persistence.Create("", dpf._dummy1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy1 = result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1.id)
	assert.Equal(t, dpf._dummy1.key, dummy1.key)
	assert.Equal(t, dpf._dummy1.content, dummy1.content)

	// Create another dummy
	result, err = dpf._persistence.Create("", dpf._dummy2)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy2 = result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2.id)
	assert.Equal(t, dpf._dummy2.key, dummy2.key)
	assert.Equal(t, dpf._dummy2.content, dummy2.content)

	page, errp := dpf._persistence.GetPageByFilter("", *cdata.NewEmptyFilterParams(), *cdata.NewEmptyPagingParams())
	if errp != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 2)

	// Update the dummy
	dummy1.content = "Updated Content 1"
	result, err = dpf._persistence.Update("", dummy1)
	if err != nil {
		t.Errorf("GetPageByFilter method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.id, result.id)
	assert.Equal(t, dummy1.key, result.key)
	assert.Equal(t, dummy1.content, result.content)

	// Partially update the dummy
	result, err = dpf._persistence.UpdatePartially("", dummy1.id, *cdata.NewAnyValueMapFromTuples("content", "Partially Updated Content 1"))
	if err != nil {
		t.Errorf("UpdatePartially method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.id, result.id)
	assert.Equal(t, dummy1.key, result.key)
	assert.Equal(t, "Partially Updated Content 1", result.content)

	// Get the dummy by Id
	result, err = dpf._persistence.GetOneById("", dummy1.id)
	if err != nil {
		t.Errorf("GetOneById method error %v", err)
	}
	// Try to get item
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.id, result.id)
	assert.Equal(t, dummy1.key, result.key)
	assert.Equal(t, "Partially Updated Content 1", result.content)

	// Delete the dummy
	result, err = dpf._persistence.DeleteById("", dummy1.id)
	if err != nil {
		t.Errorf("DeleteById method error %v", err)
	}
	assert.NotNil(t, result)
	assert.Equal(t, dummy1.id, result.id)
	assert.Equal(t, dummy1.key, result.key)
	assert.Equal(t, "Partially Updated Content 1", result.content)

	// Get the deleted dummy
	result, err = dpf._persistence.GetOneById("", dummy1.id)
	if err != nil {
		t.Errorf("GetOneById method error %v", err)
	}
	// Try to get item
	assert.Nil(t, result)
}

func (dpf *DummyPersistenceFixture) testBatchOperations(t *testing.T) {
	var dummy1 Dummy
	var dummy2 Dummy

	// Create one dummy
	result, err := dpf._persistence.Create("", dpf._dummy1)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy1 = result
	assert.NotNil(t, dummy1)
	assert.NotNil(t, dummy1.id)
	assert.Equal(t, dpf._dummy1.key, dummy1.key)
	assert.Equal(t, dpf._dummy1.content, dummy1.content)

	// Create another dummy
	result, err = dpf._persistence.Create("", dpf._dummy2)
	if err != nil {
		t.Errorf("Create method error %v", err)
	}
	dummy2 = result
	assert.NotNil(t, dummy2)
	assert.NotNil(t, dummy2.id)
	assert.Equal(t, dpf._dummy2.key, dummy2.key)
	assert.Equal(t, dpf._dummy2.content, dummy2.content)

	// Read batch
	items, err := dpf._persistence.GetListByIds("", []string{dummy1.id, dummy2.id})
	if err != nil {
		t.Errorf("GetListByIds method error %v", err)
	}
	//assert.isArray(t,items)
	assert.NotNil(t, items)
	assert.Len(t, items, 2)

	// Delete batch
	err = dpf._persistence.DeleteByIds("", []string{dummy1.id, dummy2.id})
	if err != nil {
		t.Errorf("DeleteByIds method error %v", err)
	}
	assert.Nil(t, err)

	// Read empty batch
	items, err = dpf._persistence.GetListByIds("", []string{dummy1.id, dummy2.id})
	if err != nil {
		t.Errorf("GetListByIds method error %v", err)
	}
	//assert.isArray(items)
	assert.NotNil(t, items)
	assert.Len(t, items, 0)

}
