#!/usr/bin/env bash

set -e

# 生成pb
protoc ./game.proto --go_out=.
