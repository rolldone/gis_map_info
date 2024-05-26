#!/bin/bash

NAME="dev"
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

rm main || true

pm2 stop all
pm2 delete all
pm2 start go --name $NAME -- run main.go

while $LOOP == true; 
do
    pm2 logs
done