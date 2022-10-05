#!/bin/bash

export LISTEN_ADDR=127.0.0.1:4400
export DEBUG=true 
export FILE_PATH=bot.txt
export FILE_SCAN_TIMEOUT=10


go build -o ./debug/ip-validator ./app
./debug/ip-validator