#!/bin/bash
GO_FILE_PATH=/mnt/app
VOLUME_FILE=${GO_FILE_PATH}/goWork/src/github.com/ShowPay/script-decoder
docker stop mvc_script_decoder
docker rm mvc_script_decoder
docker rmi mvc_script_decoder:v1.0
docker build -t mvc_script_decoder:v1.0 -f deploy/Dockerfile .
docker run  --name mvc_script_decoder -p 9030:9030 -d -v ${VOLUME_FILE}:${VOLUME_FILE} --restart=on-failure:3 mvc_script_decoder:v1.0