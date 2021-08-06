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

type Owner struct {
	ID    string `json:"id"`
	Asset uint64 `json:"asset,string"`
	Job   uint64 `json:"job,string"`
	Site  uint64 `json:"site,string"`
}

type OwnerGrouping struct {
	Owner
	Version  int32 `json:"version"`
	Modified int64 `json:"Modified,string"`
	Deleted  bool  `json:"Deleted"`
}

type NestedOwnerGroup struct {
	OwnerGrouping
	NestedField string `json:"nested_field,string"`
	Testing     int    `json:"testing,string"`
}

func TestCloneObjectUtils(t *testing.T) {

	now := time.Now().UTC()

	tags := make(map[uint64]*TagV1)

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

	properties := make(map[string]interface{})

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

func TestGenerateObjectIdUtils(t *testing.T) {

	var test interface{} = Owner{
		Asset: 123,
		Job:   456,
		Site:  987,
	}

	persist.GenerateObjectId(&test)
	owner, _ := test.(Owner)
	assert.NotEmpty(t, owner.ID)
	assert.Equal(t, owner.Asset, (uint64)(123))
	assert.Equal(t, owner.Job, (uint64)(456))
	assert.Equal(t, owner.Site, (uint64)(987))

	var test2 interface{} = OwnerGrouping{
		Owner: Owner{
			Asset: 123,
			Job:   456,
			Site:  987,
		},
		Version:  123,
		Modified: 67890,
		Deleted:  false,
	}

	persist.GenerateObjectId(&test2)
	ownerGrouping, _ := test2.(OwnerGrouping)
	assert.NotEmpty(t, ownerGrouping.ID)
	assert.Equal(t, ownerGrouping.Version, (int32)(123))
	assert.Equal(t, ownerGrouping.Modified, (int64)(67890))
	assert.Equal(t, ownerGrouping.Deleted, false)
	assert.Equal(t, ownerGrouping.Asset, (uint64)(123))
	assert.Equal(t, ownerGrouping.Job, (uint64)(456))
	assert.Equal(t, ownerGrouping.Site, (uint64)(987))

	var test3 interface{} = NestedOwnerGroup{
		OwnerGrouping: OwnerGrouping{
			Owner: Owner{
				Asset: 123,
				Job:   456,
				Site:  987,
			},
			Version:  123,
			Modified: 67890,
			Deleted:  false,
		},
		NestedField: "nested 3",
		Testing:     9876,
	}

	persist.GenerateObjectId(&test3)
	nestedOwnerGroup, _ := test3.(NestedOwnerGroup)
	assert.NotEmpty(t, ownerGrouping.ID)
	assert.Equal(t, nestedOwnerGroup.Version, (int32)(123))
	assert.Equal(t, nestedOwnerGroup.Modified, (int64)(67890))
	assert.Equal(t, nestedOwnerGroup.Deleted, false)
	assert.Equal(t, nestedOwnerGroup.Asset, (uint64)(123))
	assert.Equal(t, nestedOwnerGroup.Job, (uint64)(456))
	assert.Equal(t, nestedOwnerGroup.Site, (uint64)(987))
	assert.Equal(t, nestedOwnerGroup.NestedField, "nested 3")
	assert.Equal(t, nestedOwnerGroup.Testing, 9876)
}
