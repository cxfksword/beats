package web

import (
	"fmt"
	"os"

	"github.com/cxfksword/beats/libbeat/common"
	_ "github.com/cxfksword/beats/libbeat/logp"
	"github.com/cxfksword/beats/libbeat/outputs"
)

func init() {
	outputs.RegisterOutputPlugin("web", New)
}

type web struct {
	config config
	server *WebServer
}

func New(config *common.Config, _ int) (outputs.Outputer, error) {
	c := &web{config: defaultConfig}
	err := config.Unpack(&c.config)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf(":%d", c.config.Port)
	c.server = NewWebServer(addr)
	go c.server.Start()
	return c, nil
}

func newWeb(port int) *web {
	return &web{config: config{port}}
}

func writeBuffer(buf []byte) error {
	written := 0
	for written < len(buf) {
		n, err := os.Stdout.Write(buf[written:])
		if err != nil {
			return err
		}

		written += n
	}
	return nil
}

// Implement Outputer
func (c *web) Close() error {
	return nil
}

func (c *web) PublishEvent(
	s outputs.Signaler,
	opts outputs.Options,
	event common.MapStr,
) error {

	c.server.SendEvent(event)
	outputs.SignalCompleted(s)
	return nil
}
