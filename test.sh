#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o grpcdump cmd/grpcdump/main.go
docker cp grpcdump jzl-grpcdump-dev:/app/
docker exec -ti jzl-grpcdump-dev /app/grpcdump -i eth0 -p 3333 -proto-set /app/leyan-proto.pb -log-level debug
