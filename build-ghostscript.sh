# Set target version
GHOSTPDL_VERSION="10.02.1"

echo "Start GhostScript setup v${GHOSTPDL_VERSION}"
# Fetch & extract GhostScript source code
# curl "https://github.com/ArtifexSoftware/ghostpdl-downloads/releases/download/gs10021/ghostpdl--${GHOSTPDL_VERSION}.tar.gz" # --output ghostscript-source.tar.gz
gzip -d -c "ghostpdl-${GHOSTPDL_VERSION}.tar.gz | tar -xf -"
cd "./ghostpdl--${GHOSTPDL_VERSION}"

# Build binary
echo "\n\n\nBuilding binary\n\n"
mkdir ./binary/temp
./configure
make O=../binary/temp

cd ..
mv ./binary/temp/gs ./binary

# Remove temp & unnecessary files
echo "\n\n\nRemoving temp files\n\n"
rm -rf ./binary/temp
rm -rf "ghostpdl--${GHOSTPDL_VERSION}"
# rm "ghostpdl--${GHOSTPDL_VERSION}.tar.gz"