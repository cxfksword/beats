package nfs

import (
	"github.com/cxfksword/beats/packetbeat/config"
	"github.com/cxfksword/beats/packetbeat/protos"
)

type rpcConfig struct {
	config.ProtocolCommon `config:",inline"`
}

var (
	defaultConfig = rpcConfig{
		ProtocolCommon: config.ProtocolCommon{
			TransactionTimeout: protos.DefaultTransactionExpiration,
			Ports:              []int{2049},
		},
	}
)
