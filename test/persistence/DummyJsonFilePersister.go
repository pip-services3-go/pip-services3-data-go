package test_persistence

import (
	"encoding/json"
	cpersist "github.com/pip-services3-go/pip-services3-data-go/persistence"

	"github.com/pip-services3-go/pip-services3-commons-go/convert"
)

type DummyJsonFilePersister struct {
	cpersist.JsonFilePersister
}

func NewDummyJsonFilePersister(path string) *DummyJsonFilePersister {
	res := DummyJsonFilePersister{*cpersist.NewJsonFilePersister(path)}
	return &res
}

func (c *DummyJsonFilePersister) Load(correlation_id string) (data []interface{}, err error) {
	result, err := c.JsonFilePersister.Load(correlation_id)
	if err != nil {
		return result, err
	}
	data = make([]interface{}, len(result))
	for i, v := range result {
		item := convert.MapConverter.ToNullableMap(v)
		jsonMarshalStr, errJson := json.Marshal(item)
		if errJson != nil {
			panic("NewDummyJsonFilePersister.Load Error can't convert from Json to Dummy type")
		}
		value := Dummy{}
		json.Unmarshal(jsonMarshalStr, &value)
		data[i] = value
	}
	return
}
