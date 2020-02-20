#!/bin/bash

go run . $@ > $1.OUT && echo "DONE WITH $@" &
