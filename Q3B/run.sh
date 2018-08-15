#!/bin/sh
#

docker run -d \
	--rm \
	--name urlshortener \
	-p 8000:8000 \
	brianwong1861/urlshorten:v1
