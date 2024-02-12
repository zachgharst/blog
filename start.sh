#!/bin/bash

cd -- "$( dirname -- "${BASH_SOURCE[0]}" )"
go run cmd/server/main.go
