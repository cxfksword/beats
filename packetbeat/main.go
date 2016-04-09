package main

import (
	"os"

	"github.com/cxfksword/beats/libbeat/beat"
	"github.com/cxfksword/beats/packetbeat/beater"

	// import support protocol modules
	_ "github.com/cxfksword/beats/packetbeat/protos/amqp"
	_ "github.com/cxfksword/beats/packetbeat/protos/dns"
	_ "github.com/cxfksword/beats/packetbeat/protos/http"
	_ "github.com/cxfksword/beats/packetbeat/protos/memcache"
	_ "github.com/cxfksword/beats/packetbeat/protos/mongodb"
	_ "github.com/cxfksword/beats/packetbeat/protos/mysql"
	_ "github.com/cxfksword/beats/packetbeat/protos/nfs"
	_ "github.com/cxfksword/beats/packetbeat/protos/pgsql"
	_ "github.com/cxfksword/beats/packetbeat/protos/redis"
	_ "github.com/cxfksword/beats/packetbeat/protos/thrift"
)

var Name = "packetbeat"

// Setups and Runs Packetbeat
func main() {
	if err := beat.Run(Name, "", beater.New()); err != nil {
		os.Exit(1)
	}
}
