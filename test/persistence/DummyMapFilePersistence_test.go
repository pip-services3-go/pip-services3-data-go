package test_persistence

import (
	"os"
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

func TestDummyMapFilePersistence(t *testing.T) {
	filename := "../../data/dummies_map.json"

	//cleaning file before testing
	f, err := os.Create(filename)
	if err != nil {
		t.Error("Can't clean file: ", filename)
	}
	f.Close()

	persistence := NewDummyMapFilePersistence(filename)
	persistence.Configure(*cconf.NewEmptyConfigParams())

	defer persistence.Close("")

	fixture := NewDummyMapPersistenceFixture(persistence)
	persistence.Open("")

	t.Run("DummyMapFilePersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyMapFilePersistence:Batch", fixture.TestBatchOperations)

}
