package test_persistence

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

func TestItemMemoryPersistence(t *testing.T) {
	persister := NewItemMemoryPersistence()
	persister.Configure(cconf.NewEmptyConfigParams())

	fixture := NewItemPersistenceFixture(persister)

	t.Run("ItemMemoryPersistence:CRUD", fixture.TestCrudOperations)
}
