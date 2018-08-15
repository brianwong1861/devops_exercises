#!/bin/sh
#
docker run -d \
        --rm \
        --name webhook \
        --env-file=./env.list \
        -p 3000:3000 \
        brianwong1861/webhook:v1

