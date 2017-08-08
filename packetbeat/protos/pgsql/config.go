package pgsql

import (
	"github.com/cxfksword/beats/packetbeat/config"
	"github.com/cxfksword/beats/packetbeat/protos"
)

type pgsqlConfig struct {
	config.ProtocolCommon `config:",inline"`
	MaxRowLength          int `config:"max_row_length"`
	MaxRows               int `config:"max_rows"`
}

var (
	defaultConfig = pgsqlConfig{
		ProtocolCommon: config.ProtocolCommon{
			TransactionTimeout: protos.DefaultTransactionExpiration,
			Ports:              []int{5432},
			SendRequest:        true,
			SendResponse:       true,
		},
		MaxRowLength: 1024,
		MaxRows:      10,
	}
)