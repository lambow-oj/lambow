#!/bin/bash

mkdir -p output/bin
export GO111MODULE=on && go build -a -o output/bin/lambow

chmod +x output/bin/lambow
