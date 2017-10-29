#!/usr/bin/env bash
fswatch -o ./**/*.go | xargs -n1 ./bin/fast-benchmark.sh
