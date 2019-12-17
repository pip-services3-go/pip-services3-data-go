package test_persistence

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

func TestDummyMapMemoryPersistence(t *testing.T) {
	persister := NewEmptyDummyMapMemoryPersistence()
	persister.Configure(*cconf.NewEmptyConfigParams())

	fixture := NewDummyMapPersistenceFixture(persister)

	t.Run("MemoryPersistence:CRUD", fixture.TestCrudOperations)
	t.Run("MemoryPersistence:Batch", fixture.TestBatchOperations)

}
