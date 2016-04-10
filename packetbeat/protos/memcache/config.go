package memcache

import (
	"time"

	"github.com/cxfksword/beats/packetbeat/config"
	"github.com/cxfksword/beats/packetbeat/protos"
)

type memcacheConfig struct {
	config.ProtocolCommon `config:",inline"`
	MaxValues             int
	MaxBytesPerValue      int
	UdpTransactionTimeout time.Duration
	ParseUnknown          bool
}

var (
	defaultConfig = memcacheConfig{
		ProtocolCommon: config.ProtocolCommon{
			Ports:              []int{11211},
			TransactionTimeout: protos.DefaultTransactionExpiration,
		},
		UdpTransactionTimeout: protos.DefaultTransactionExpiration,
	}
)
