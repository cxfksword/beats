package main

import (
	"os"

	"github.com/cxfksword/beats/libbeat/beat"
	"github.com/cxfksword/beats/topbeat/beater"
)

var Name = "topbeat"

func main() {
	if err := beat.Run(Name, "", beater.New()); err != nil {
		os.Exit(1)
	}
}
