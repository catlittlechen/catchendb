#!/bin/bash

ulimit -n 100000
./main --config=../etc/config.xml &
