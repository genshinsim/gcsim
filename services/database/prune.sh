#!/bin/bash

docker container kill docker-postgrest-1
docker container kill docker-postgres-1
docker container rm docker-postgrest-1
docker container rm docker-postgres-1
docker volume rm docker_postgres-data
docker-compose up -d