#!/bin/bash

cd cmd/csync

# Clean up any leftover test artifacts
rm -rf .csync __test__ *.txt subdir 2>/dev/null

CSYNC_ENV=test go test -cover -v ./...

# Clean up after tests
rm -rf .csync __test__ *.txt subdir 2>/dev/null
