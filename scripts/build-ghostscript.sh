## Run from parent folder (welovepdf) as such : ./scripts/build-ghostcript.sh

# Set target version
GHOSTPDL_VERSION="10.02.1"
TEMP_DIR="$PWD/temp"
GS_SRC_DIR="${TEMP_DIR}/ghostpdl-${GHOSTPDL_VERSION}"
TEMP_BIN_DIR="${TEMP_DIR}/bin"
OUTPUT_DIR="$PWD/"
TEMP_ARCHIVE_FILE=""


echo "Start GhostScript setup v${GHOSTPDL_VERSION} using temp folder : ${GS_SRC_DIR}"
# mkdir ${TEMP_DIR}
# rm -rf ${GS_SRC_DIR}
rm -rf ${TEMP_BIN_DIR}
mkdir ${TEMP_BIN_DIR}

# # Fetch & extract GhostScript source code
# # curl --output ghostscript-source.tar.gz "https://homepilot-data-public.s3.eu-west-3.amazonaws.com/ghostpdl-10.02.1.tar.gz" # --output ghostscript-source.tar.gz
# curl --output ${GS_SRC_DIR}.tar.gz "https://homepilot-data-public.s3.eu-west-3.amazonaws.com/ghostpdl-${GHOSTPDL_VERSION}.tar.gz" # --output ghostscript-source.tar.gz
# cd "${TEMP_DIR}"
# gzip -d -c "${GS_SRC_DIR}.tar.gz" | tar -xf -

# Build binary
cd ${GS_SRC_DIR}
echo "\n\n\nBuilding binary in dir : $PWD\n\n"
# ./configure
make O=../bin

# cd ${TEMP_DIR}/..
# mv ./binary/temp/gs ./binary

# Remove temp & unnecessary files
# echo "\n\n\nRemoving temp files\n\n"
# rm -rf ./binary/temp
# rm -rf ${GS_SRC_DIR}
# rm "${GS_SRC_DIR}.tar.gz"