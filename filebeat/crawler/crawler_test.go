// +build !integration

package crawler

import (
	"testing"

	"github.com/cxfksword/beats/filebeat/config"
	"github.com/cxfksword/beats/filebeat/input"
	"github.com/stretchr/testify/assert"
)

func TestCrawlerStartError(t *testing.T) {
	crawler := Crawler{}
	channel := make(chan *input.FileEvent, 1)
	prospectorConfigs := []config.ProspectorConfig{}

	error := crawler.Start(prospectorConfigs, channel)

	assert.Error(t, error)
}
