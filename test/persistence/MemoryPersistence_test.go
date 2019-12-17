package test_persistence

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

func TestMemoryPersistence(t *testing.T) {
	persister := NewEmptyDummyMemoryPersistence()
	persister.Configure(*cconf.NewEmptyConfigParams())

	fixture := NewDummyPersistenceFixture(persister)

	t.Run("MemoryPersistence:CRUD", fixture.TestCrudOperations)
	t.Run("MemoryPersistence:Batch", fixture.TestBatchOperations)

}
