#!/bin/bash
if [ "$1" = "make" ]; then
    migrate create -ext sql -dir db/migrations "$2"
fi

if [ "$1" = "migrate" ]; then
    # Load variables from .env file
    source .env
    # $2 is up or down
    if [ "$2" = "down" ]; then
        if [ "$3" = "all" ]; then
        migrate -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?query&sslmode=$DB_SSLMODE" -path db/migrations down
        else
        migrate -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?query&sslmode=$DB_SSLMODE" -path db/migrations down 1
        fi
    elif [ "$2" = "up" ]; then
        migrate -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?query&sslmode=$DB_SSLMODE" -path db/migrations up
    fi
fi