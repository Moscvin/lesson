#!/bin/bash
cat $1 | redis-cli -h $HOST -p $PORT --pipe