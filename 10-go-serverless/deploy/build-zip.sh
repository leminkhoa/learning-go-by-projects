#!/bin/bash
set -e

GOOS=linux GOARCH=arm64 go build -o build/bootstrap -tags lambda.norpc cmd/main.go
zip -jrm build/go-serverless.zip build/bootstrap
