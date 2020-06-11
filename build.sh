#!/bin/bash

mkdir -p output/bin
go build -a -o output/bin/lambow

chmod +x output/bin/lambow
