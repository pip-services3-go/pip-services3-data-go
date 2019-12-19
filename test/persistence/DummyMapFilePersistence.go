package test_persistence

import (
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cpersist "github.com/pip-services3-go/pip-services3-data-go/persistence"
)

//  extends DummyMapMemoryPersistence
type DummyMapFilePersistence struct {
	DummyMapMemoryPersistence
	persister *cpersist.JsonFilePersister
}

func NewDummyMapFilePersistence(path string) *DummyMapFilePersistence {
	c := &DummyMapFilePersistence{
		DummyMapMemoryPersistence: *NewDummyMapMemoryPersistence(),
	}

	persister := cpersist.NewJsonFilePersister(path)
	c.persister = persister
	c.Loader = persister
	c.Saver = persister

	return c
}

func (c *DummyMapFilePersistence) Configure(config cconf.ConfigParams) {
	c.DummyMapMemoryPersistence.Configure(config)
	c.persister.Configure(config)
}
