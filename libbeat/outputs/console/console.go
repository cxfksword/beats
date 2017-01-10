package console

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/cxfksword/beats/libbeat/common"
	"github.com/cxfksword/beats/libbeat/logp"
	"github.com/cxfksword/beats/libbeat/outputs"
    "github.com/shiena/ansicolor"
)

var w = ansicolor.NewAnsiColorWriter(os.Stdout)
func init() {
	outputs.RegisterOutputPlugin("console", New)
}

type console struct {
	config config
}

func New(config *common.Config, _ int) (outputs.Outputer, error) {
	c := &console{config: defaultConfig}
	err := config.Unpack(&c.config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newConsole(pretty bool) *console {
	return &console{config{pretty, ""}}
}

func writeBuffer(buf []byte) error {
	written := 0
	for written < len(buf) {
		n, err := w.Write(buf[written:])
		if err != nil {
			return err
		}

		written += n
	}
	return nil
}

// Implement Outputer
func (c *console) Close() error {
	return nil
}

func (c *console) PublishEvent(
	s outputs.Signaler,
	opts outputs.Options,
	event common.MapStr,
) error {
	var msg []byte
	var err error

	if consoleMsg, exist := event["console"]; exist {
		msg = []byte(consoleMsg.(string))
	} else {
		if c.config.Pretty {
			msg, err = json.MarshalIndent(event, "", "  ")
		} else {
			msg, err = json.Marshal(event)
		}
	}
	if err != nil {
		logp.Err("Fail to convert the event to JSON: %s", err)
		outputs.SignalCompleted(s)
		return nil
	}

	if c.config.Query != "" && !strings.Contains(string(msg), c.config.Query) {
		logp.Debug("console", "[%s] not match query output: %s, will ignore.", string(msg), c.config.Query)
		outputs.SignalCompleted(s)
		return err
	}

	if err = writeBuffer(msg); err != nil {
		goto fail
	}
	if err = writeBuffer([]byte{'\n'}); err != nil {
		goto fail
	}

	outputs.SignalCompleted(s)
	return nil
fail:
	if opts.Guaranteed {
		logp.Critical("Unable to publish events to console: %v", err)
	}
	outputs.SignalFailed(s, err)
	return err
}
