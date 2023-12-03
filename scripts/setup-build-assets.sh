#!/bin/bash

# Setup GhostScript binary
rm -rf ./assets/bin/**/*
GS_BINARY_URL=https://homepilot-data-public.s3.eu-west-3.amazonaws.com/ghostscript/gsbinary
curl --output ./assets/bin/gs ${GS_BINARY_URL}

# Create config file
default_config=`cat ./assets/config/config.build.json`  
token_value=${LOGTAIL_TOKEN}
rm ./assets/config/config.json
echo "${default_config/LOGTAIL_TOKEN/$token_value}" > ./assets/config/config.json
new_value=`cat ./assets/config/config.json`
echo ${new_value}