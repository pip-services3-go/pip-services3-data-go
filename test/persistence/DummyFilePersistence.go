package test_persistence

import (
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cpersist "github.com/pip-services3-go/pip-services3-data-go/persistence"
)

//  extends DummyMemoryPersistence
type DummyFilePersistence struct {
	DummyMemoryPersistence
	persister cpersist.JsonFilePersister
}

func NewDummyFilePersistence(path string) *DummyFilePersistence {
	c := DummyFilePersistence{}
	persister := cpersist.NewJsonFilePersister(path)
	c.persister = *persister
	c.DummyMemoryPersistence = *NewDummyMemoryPersistence(persister, persister)
	return &c
}

func (c *DummyFilePersistence) Configure(config cconf.ConfigParams) {
	c.DummyMemoryPersistence.Configure(config)
	c.persister.Configure(config)
}
