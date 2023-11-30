GS_BINARY_URL=https://homepilot-data-public.s3.eu-west-3.amazonaws.com/ghostscript/gsbinary
curl --output ./assets/bin/gs ${GS_BINARY_URL}

JSON_CONFIG_STR="{\"debugMode\": false, \"logger\": { \"logtailToken\": \"${LOGTAIL_TOKEN}\" }}"
rm ./assets/config/config.json
echo ${JSON_CONFIG_STR} > ./assets/config/config.json
rm -rf ./assets/bin/**/*