package main

import (
	"os"

	"github.com/cxfksword/beats/metricbeat/beater"
	_ "github.com/cxfksword/beats/metricbeat/include"

	"github.com/cxfksword/beats/libbeat/beat"
)

var Name = "metricbeat"

func main() {
	if err := beat.Run(Name, "", beater.New()); err != nil {
		os.Exit(1)
	}
}
