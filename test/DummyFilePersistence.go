package test_persistence

import (
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
)

//  extends DummyMemoryPersistence
type DummyFilePersistence struct {
	DummyMemoryPersistence
	_persister JsonFilePersister
}

func NewDummyFilePersistence(path string) *DummyFilePersistence {
	dfp := DummyFilePersistence{}
	dfp._persister = NewJsonFilePersister(path)
	dfp._loader = dfp._persister
	dfp._saver = dfp._persister
	return &dfp
}

func (dfp *DummyFilePersistence) Configure(config cconf.ConfigParams) {
	dfp.DummyMemoryPersistence.Configure(config)
	dfp._persister.Configure(config)
}
