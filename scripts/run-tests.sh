#!/bin/bash

cd cmd/csync
CSYNC_ENV=test go test -cover -v ./...
