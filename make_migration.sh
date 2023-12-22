#!/bin/bash
if [ "$1" = "make" ]; then
    migrate create -ext sql -dir db/migrations "$2"
fi

if [ "$1" = "migrate" ]; then
    migrate "$2"
fi