#!/bin/bash

NAME="prod"
LOOP=true
# trap ctrl-c and call ctrl_c()
trap ctrl_c INT

function clear_logs(){
    sleep 2
    pm2 flush all
}

function ctrl_c() {
    LOOP=false
    pm2 stop $NAME && \
    clear_logs
}

pm2 stop all
pm2 delete all
go build main.go
pm2 start ./main --name $NAME

while $LOOP == true;
do
    pm2 logs
done