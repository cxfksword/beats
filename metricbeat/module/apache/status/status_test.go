// +build integration

package status

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cxfksword/beats/libbeat/common"
	"github.com/cxfksword/beats/metricbeat/helper"
	"github.com/cxfksword/beats/metricbeat/module/apache"
)

func TestConnect(t *testing.T) {

	config, _ := getApacheModuleConfig()

	module, mErr := helper.NewModule(config, apache.New)
	assert.NoError(t, mErr)
	ms, msErr := helper.NewMetricSet("status", New, module)
	assert.NoError(t, msErr)

	// Setup metricset and metricseter
	err := ms.Setup()
	assert.NoError(t, err)
	err = ms.MetricSeter.Setup(ms)
	assert.NoError(t, err)

	// Check that host is correctly set
	assert.Equal(t, apache.GetApacheEnvHost(), ms.Config.Hosts[0])

	data, err := ms.MetricSeter.Fetch(ms, ms.Config.Hosts[0])
	assert.NoError(t, err)

	// Check fields
	assert.Equal(t, 13, len(data))
}

type ApacheModuleConfig struct {
	Hosts  []string `config:"hosts"`
	Module string   `config:"module"`
}

func getApacheModuleConfig() (*common.Config, error) {
	return common.NewConfigFrom(ApacheModuleConfig{
		Module: "apache",
		Hosts:  []string{apache.GetApacheEnvHost()},
	})
}
