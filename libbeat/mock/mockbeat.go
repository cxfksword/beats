package mock

import (
	"time"

	"github.com/cxfksword/beats/libbeat/beat"
	"github.com/cxfksword/beats/libbeat/common"
)

///*** Mock Beat Setup ***///

var Version = "9.9.9"
var Name = "mockbeat"

type Mockbeat struct {
	done chan struct{}
}

// Creates beater
func New() *Mockbeat {
	return &Mockbeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (mb *Mockbeat) Config(b *beat.Beat) error {
	return nil
}

func (mb *Mockbeat) Setup(b *beat.Beat) error {
	return nil
}

func (mb *Mockbeat) Run(b *beat.Beat) error {
	// Wait until mockbeat is done
	b.Events.PublishEvent(common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "mock",
		"message":    "Mockbeat is alive!",
	})
	<-mb.done
	return nil
}

func (mb *Mockbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (mb *Mockbeat) Stop() {
	close(mb.done)
}
