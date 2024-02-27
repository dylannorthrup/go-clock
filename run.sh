#!/bin/bash

go mod tidy && \
  go build -o go_clock && \
  ./go_clock
