package test_persistence

import (
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cpersist "github.com/pip-services3-go/pip-services3-data-go/persistence"
)

//  extends DummyMapMemoryPersistence
type DummyMapFilePersistence struct {
	DummyMapMemoryPersistence
	_persister cpersist.JsonFilePersister
}

func NewDummyMapFilePersistence(path string) *DummyMapFilePersistence {
	c := DummyMapFilePersistence{}
	persister := cpersist.NewJsonFilePersister(path)
	c._persister = *persister
	c.DummyMapMemoryPersistence = *NewDummyMapMemoryPersistence(persister, persister)
	return &c
}

func (c *DummyMapFilePersistence) Configure(config cconf.ConfigParams) {
	c.DummyMapMemoryPersistence.Configure(config)
	c._persister.Configure(config)
}
