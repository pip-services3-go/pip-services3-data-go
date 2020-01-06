package test_persistence

import (
	"reflect"
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/v3/config"
	cpersist "github.com/pip-services3-go/pip-services3-data-go/v3/persistence"
	"github.com/stretchr/testify/assert"
)

func TestJsonFilePersister(t *testing.T) {
	var p interface{}
	persistence := cpersist.NewJsonFilePersister(reflect.TypeOf(p), "")
	persistence.Configure(*cconf.NewEmptyConfigParams())
	fileName := "../JsonFilePersisterTest"
	persistence.Configure(*cconf.NewConfigParamsFromTuples("path", fileName))
	assert.Equal(t, fileName, persistence.Path())
}
