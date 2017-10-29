#!/usr/bin/env bash
set -e
rm -f BENCHMARK.md
touch BENCHMARK.md

server_types=(simple-server)

for i in ${server_types[@]}; do
  go run "$i/server.go" &
  server_pid=$!
  sleep 3
  echo "# ${i}" >> BENCHMARK.md
  echo "\`\`\`" >> BENCHMARK.md
  wrk -t12 -c400 -d30s http://localhost:8081/index.html >> BENCHMARK.md
  echo "\`\`\`" >> BENCHMARK.md
  pkill server
done
