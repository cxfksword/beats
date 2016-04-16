package mongodb

import (
	"github.com/cxfksword/beats/packetbeat/config"
	"github.com/cxfksword/beats/packetbeat/protos"
)

type mongodbConfig struct {
	config.ProtocolCommon `config:",inline"`
	MaxDocLength          int `config:"max_doc_length"`
	MaxDocs               int `config:"max_docs"`
}

var (
	defaultConfig = mongodbConfig{
		ProtocolCommon: config.ProtocolCommon{
			TransactionTimeout: protos.DefaultTransactionExpiration,
			Ports:              []int{27017},
			SendRequest:        true,
			SendResponse:       true,
		},
		MaxDocLength: 5000,
		MaxDocs:      10,
	}
)
