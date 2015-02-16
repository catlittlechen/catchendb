#!/bin/bash

killall ccdb
ulimit -n 100000
./ccdb --config=../etc/config.xml &
