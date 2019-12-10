package test_persistence

import (
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cpersist "github.com/pip-services3-go/pip-services3-data-go/persistence"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonFilePersister(t *testing.T) {
	persistence := cpersist.NewJsonFilePersister("")
	persistence.Configure(*cconf.NewEmptyConfigParams())
	fileName := "../JsonFilePersisterTest"
	persistence.Configure(*cconf.NewConfigParamsFromTuples("path", fileName))
	assert.Equal(t, fileName, persistence.Path())
}
