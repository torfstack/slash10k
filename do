#! /bin/bash

cd backend || return
CGO_ENABLED=0 GOOS=linux go build -o bin/scurvy10k-backend src/main.go
docker buildx build . -t scurvy10k-backend
