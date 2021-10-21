package test_persistence

import cdata "github.com/pip-services3-go/pip-services3-commons-go/data"

// extends IGetter<Item, String>, IWriter<Item, String>, IPartialUpdater<Item, String> {
type IItemPersistence interface {
	GetOneById(correlationId string, id string) (item Item, err error)
	Create(correlationId string, item Item) (result Item, err error)
	Update(correlationId string, item Item) (result Item, err error)
	UpdatePartially(correlationId string, id string, data *cdata.AnyValueMap) (item Item, err error)
	DeleteById(correlationId string, id string) (item Item, err error)
}
