package test_persistence

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

func TestMemoryPersistence(t *testing.T) {
	persister := NewDummyMemoryPersistence()
	persister.Configure(*cconf.NewEmptyConfigParams())

	fixture := NewDummyPersistenceFixture(persister)

	t.Run("MemoryPersistence", fixture.TestCrudOperations)
	t.Run("MemoryPersistence", fixture.TestBatchOperations)

}
