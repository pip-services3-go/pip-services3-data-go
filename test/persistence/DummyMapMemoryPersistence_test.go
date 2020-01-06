package test_persistence

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/v3/config"
)

func TestDummyMapMemoryPersistence(t *testing.T) {
	persister := NewDummyMapMemoryPersistence()
	persister.Configure(*cconf.NewEmptyConfigParams())

	fixture := NewDummyMapPersistenceFixture(persister)

	t.Run("DummyMapMemoryPersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyMapMemoryPersistence:Batch", fixture.TestBatchOperations)

}
