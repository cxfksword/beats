#!/bin/bash

BEATNAME=metricbeat
SYSTEM_TESTS=true
TEST_ENVIRONMENT?=true
GOPACKAGES=$(shell go list ${BEAT_DIR}/${BEATNAME}/... | grep -v /vendor/)

include ../libbeat/scripts/Makefile

# Collects all module dashboards
.PHONY: kibana
kibana:
	cp -r module/*/.kibana/ etc/kibana


# Collects all module and metricset fields
.PHONY: fields
fields:
	cat etc/fields_base.yml > etc/fields.yml
	cat module/*/*/fields.yml >> etc/fields.yml

# Collects all module docs
.PHONY: docs
docs:
	cp -r module/*/*.asciidoc docs/modules


# This is called by the beats packer before building starts
.PHONY: before-build
before-build:

# Runs all collection steps and updates afterwards
.PHONY: collect
collect: kibana fields docs update
