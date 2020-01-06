package test_persistence

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/v3/config"
)

func TestDummyMemoryPersistence(t *testing.T) {
	persister := NewDummyMemoryPersistence()
	persister.Configure(*cconf.NewEmptyConfigParams())

	fixture := NewDummyPersistenceFixture(persister)

	t.Run("DummyMemoryPersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyMemoryPersistence:Batch", fixture.TestBatchOperations)

}
