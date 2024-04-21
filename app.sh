#!/bin/sh

go mod tidy
go install github.com/pressly/goose/v3/cmd/goose@latest
goose up
go build cmd/app/app.go && ./app