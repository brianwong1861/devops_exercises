#!/bin/sh
#
#
#if [ -z "$STARTAT" ] && [ -z "$ENDAT" ]
#  then
#    export STARTAT="2017-06-10 00:00:00"
#    export ENDAT="2017-05-19 23:59:59" 
#fi
docker run --rm -v `pwd`:/home/app --env-file ./env.list brianwong1861/loganalyzer:v1 python log_analyzer.py $@
