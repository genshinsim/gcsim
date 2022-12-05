#!/bin/bash

protoc --go-grpc_out=. protos/backend/result.proto
protoc --go_out=. protos/backend/result.proto

protoc --go-grpc_out=. protos/backend/db.proto

protoc --go_out=. --proto_path=. \
--go_opt="Mprotos/model/character.proto=github.com/genshinsim/gcsim/pkg/model" \
--go_opt="Mprotos/model/stats.proto=github.com/genshinsim/gcsim/pkg/model" \
protos/backend/db.proto
/


protoc --go_out=. protos/model/character.proto
protoc --go_out=. protos/model/enemy.proto
protoc --go_out=. protos/model/result.proto
protoc --go_out=. protos/model/stats.proto
protoc --go_out=. protos/model/weapon.proto
