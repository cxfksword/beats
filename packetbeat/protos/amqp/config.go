package amqp

import (
	"github.com/cxfksword/beats/packetbeat/config"
	"github.com/cxfksword/beats/packetbeat/protos"
)

type amqpConfig struct {
	config.ProtocolCommon     `config:",inline"`
	ParseHeaders              bool `config:"parse_headers"`
	ParseArguments            bool `config:"parse_arguments"`
	MaxBodyLength             int  `config:"max_body_length"`
	HideConnectionInformation bool `config:"hide_connection_information"`
}

var (
	defaultConfig = amqpConfig{
		ProtocolCommon: config.ProtocolCommon{
			TransactionTimeout: protos.DefaultTransactionExpiration,
			Ports:              []int{5672},
			SendRequest:        true,
			SendResponse:       true,
		},
		ParseHeaders:              true,
		ParseArguments:            true,
		MaxBodyLength:             1000,
		HideConnectionInformation: true,
	}
)
