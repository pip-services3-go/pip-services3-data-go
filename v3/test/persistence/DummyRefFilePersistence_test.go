package test_persistence

import (
	"os"
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/v3/config"
)

func TestDummyRefFilePersistence(t *testing.T) {
	filename := "../../data/dummies_ref.json"

	//cleaning file before testing
	f, err := os.Create(filename)
	if err != nil {
		t.Error("Can't clean file: ", filename)
	}
	f.Close()

	persistence := NewDummyRefFilePersistence(filename)
	persistence.Configure(cconf.NewEmptyConfigParams())

	defer persistence.Close("")

	fixture := NewDummyRefPersistenceFixture(persistence)
	persistence.Open("")

	t.Run("DummyFilePersistence:CRUD", fixture.TestCrudOperations)
	t.Run("DummyFilePersistence:Batch", fixture.TestBatchOperations)

}
