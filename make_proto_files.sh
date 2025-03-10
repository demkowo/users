#!/bin/bash

protoc \
  -I=proto \
  --go_out=internal/generated \
  --go_opt=paths=source_relative \
  --go-grpc_out=internal/generated \
  --go-grpc_opt=paths=source_relative \
  proto/users.proto
