package test_persistence

import (
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

func TestDummyFilePersistence(t *testing.T) {
	persistence := NewDummyFilePersistence("../../data/dummies.json")
	persistence.Configure(*cconf.NewEmptyConfigParams())

	defer persistence.Close("")

	fixture := NewDummyPersistenceFixture(persistence)
	persistence.Open("")

	t.Run("DummyFilePersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyFilePersistence:Batch", fixture.TestBatchOperations)

}
