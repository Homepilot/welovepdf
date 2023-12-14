#!/bin/bash

# Create config file
default_config=`cat ./assets/config/config.build.json`  
token_value=${LOGTAIL_TOKEN}
rm ./assets/config/config.json
echo "${default_config/LOGTAIL_TOKEN/$token_value}" > ./assets/config/config.json
new_value=`cat ./assets/config/config.json`
echo ${new_value}