#!/bin/bash

protoc --go_out=./chat/ --go_opt=paths=source_relative --go-grpc_out=./chat/ --go-grpc_opt=paths=source_relative chat.proto
