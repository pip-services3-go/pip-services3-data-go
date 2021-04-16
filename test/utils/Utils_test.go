package test_utils

import (
	"reflect"
	"testing"
	"time"

	persist "github.com/pip-services3-go/pip-services3-data-go/persistence"
	"github.com/stretchr/testify/assert"
)

type AttributeV1 struct {
	Id          uint64                 `json:"id,string"`
	DisplayName string                 `json:"display_name"`
	Description string                 `json:"description"`
	AssetId     uint64                 `json:"asset_id,string"`
	TagMap      map[uint64]*TagV1      `json:"tag_map"`
	Properties  map[string]interface{} `json:"properties"`
}

type TagV1 struct {
	Id        uint64    `json:"id,string"`
	ValidFrom time.Time `json:"valid_from"`
	UoM       int64     `json:"uom,string"`
}

func TestCloneObjectUtils(t *testing.T) {

	now := time.Now().UTC()

	tags := make(map[uint64]*TagV1, 0)

	tags[123] = &TagV1{
		Id:        456,
		ValidFrom: now,
		UoM:       3456,
	}

	tags[456] = &TagV1{
		Id:        123,
		ValidFrom: now,
		UoM:       987,
	}

	properties := make(map[string]interface{}, 0)

	atribute := AttributeV1{
		Id:          12345,
		DisplayName: "Display Name",
		Description: "Description",
		AssetId:     890,
		TagMap:      tags,
		Properties:  properties,
	}

	copy := persist.CloneObject(atribute, reflect.TypeOf(atribute))

	assert.NotNil(t, copy)
	copyAttribute, _ := copy.(AttributeV1)

	assert.Equal(t, copyAttribute.Id, atribute.Id)
	assert.Equal(t, copyAttribute.DisplayName, atribute.DisplayName)
	assert.Equal(t, copyAttribute.Description, atribute.Description)
	assert.Equal(t, copyAttribute.AssetId, atribute.AssetId)
	assert.NotNil(t, copyAttribute.TagMap)

	tag := copyAttribute.TagMap[123]
	assert.NotNil(t, tag)
	assert.Equal(t, tag.Id, (uint64)(456))
	assert.Equal(t, tag.ValidFrom, now)
	assert.Equal(t, tag.UoM, (int64)(3456))

	tag = copyAttribute.TagMap[456]
	assert.NotNil(t, tag)
	assert.Equal(t, tag.Id, (uint64)(123))
	assert.Equal(t, tag.ValidFrom, now)
	assert.Equal(t, tag.UoM, (int64)(987))

	assert.NotNil(t, copyAttribute.Properties)

	// atribute.TagMap[456].Id += 1
	// assert.Equal(t, atribute.TagMap[456].Id, (uint64)(124))
	// assert.Equal(t, tag.Id, (uint64)(123))

}
