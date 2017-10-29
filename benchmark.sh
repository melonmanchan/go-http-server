#!/usr/bin/env bash
set -e
rm -f BENCHMARK.md
touch BENCHMARK.md

server_types=(simple-server)

for i in ${server_types[@]}; do
  go run "$i/$i.go" &
  server_pid=$!
  sleep 5
  echo "# ${i}" >> BENCHMARK.md
  wrk -t12 -c400 -d30s http://localhost:8081/index.html >> BENCHMARK.md
  kill $server_pid
done
