#!/bin/bash

cd cmd/nexio

# Clean up any leftover test artifacts
rm -rf .nexio __test__ *.txt subdir 2>/dev/null

NEXIO_ENV=test go test -cover -v ./...

# Clean up after tests
rm -rf .nexio __test__ *.txt subdir 2>/dev/null
