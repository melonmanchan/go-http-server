#!/usr/bin/env bash
set -e

go run "./simple-server/simple-server.go" &
sleep 3
server_pid=$!
wrk -t12 -c400 -d3s http://localhost:8081/index.html
pkill simple-server
