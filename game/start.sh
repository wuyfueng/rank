#!/usr/bin/env bash

set -e

go run main.go \
-redis.host=127.0.0.1 \
-redis.port=6379 \
-redis.password=
