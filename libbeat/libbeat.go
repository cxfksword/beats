package main

import (
	"os"

	"github.com/cxfksword/beats/libbeat/beat"
	"github.com/cxfksword/beats/libbeat/mock"
)

func main() {
	if err := beat.Run(mock.Name, mock.Version, mock.New()); err != nil {
		os.Exit(1)
	}
}
