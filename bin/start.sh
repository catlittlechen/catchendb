#!/bin/bash

killall ccdb
#ulimit -n 100000
./ccdb-server --config=../etc/config.xml &
