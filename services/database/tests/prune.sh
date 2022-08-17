#!/bin/bash

docker container kill postgrest 
docker container kill postgres
docker container rm postgrest
docker container rm postgres
docker volume rm tests_postgres-data 
docker-compose up -d
