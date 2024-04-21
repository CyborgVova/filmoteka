#!/bin/sh

go mod tidy
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir db/migrations up
cd cmd/app/
go build app.go && ./app