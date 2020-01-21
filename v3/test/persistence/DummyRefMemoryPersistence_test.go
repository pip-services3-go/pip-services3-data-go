package test_persistence

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/v3/config"
)

func TestDummyRefMemoryPersistence(t *testing.T) {
	persister := NewDummyRefMemoryPersistence()
	persister.Configure(cconf.NewEmptyConfigParams())

	fixture := NewDummyRefPersistenceFixture(persister)

	t.Run("DummyRefMemoryPersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyRefMemoryPersistence:Batch", fixture.TestBatchOperations)

}
