#!/bin/bash

rm log.txt
go build -o ../../terraform-provider-relyt ../../
TF_LOG=TRACE TF_LOG_PATH=./log.txt terraform apply