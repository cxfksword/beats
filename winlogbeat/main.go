package main

import (
	"os"

	"github.com/cxfksword/beats/libbeat/beat"
	"github.com/cxfksword/beats/winlogbeat/beater"
)

// Name of this beat.
var Name = "winlogbeat"

func main() {
	if err := beat.Run(Name, "", beater.New()); err != nil {
		os.Exit(1)
	}
}
