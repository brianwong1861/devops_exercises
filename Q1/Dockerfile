FROM python:3.7-alpine3.7
MAINTAINER Brian Wong
ENV APPDIR=/home/app
WORKDIR ${APPDIR}
COPY . .
RUN apk add py3-pip && pip install geoip2


