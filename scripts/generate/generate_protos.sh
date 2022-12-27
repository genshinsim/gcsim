#!/bin/bash

protoc --go-grpc_out=. protos/backend/result.proto
protoc --go_out=. protos/backend/result.proto

protoc --go-grpc_out=. protos/backend/compute.proto
protoc --go_out=. protos/backend/compute.proto

protoc --go-grpc_out=. protos/backend/submission.proto
protoc --go_out=. protos/backend/submission.proto

protoc --go-grpc_out=. protos/backend/db.proto
protoc --go_out=. --proto_path=. \
--go_opt="Mprotos/model/result.proto=github.com/genshinsim/gcsim/pkg/model" \
--go_opt="Mprotos/model/sim.proto=github.com/genshinsim/gcsim/pkg/model" \
--go_opt="Mprotos/model/db.proto=github.com/genshinsim/gcsim/pkg/model" \
protos/backend/db.proto


protoc --go_out=. protos/model/result.proto
protoc --go_out=. protos/model/sim.proto
protoc --go_out=. protos/model/db.proto