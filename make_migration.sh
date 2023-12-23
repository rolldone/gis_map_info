#!/bin/bash
if [ "$1" = "make" ]; then
    migrate create -ext sql -dir db/migrations "$2"
fi

if [ "$1" = "migrate" ]; then
    migrate -database "postgres://root:43lw9rj2@postgres:5432/gis_map_info_db?query&sslmode=disable" -path db/migrations "$2"
fi