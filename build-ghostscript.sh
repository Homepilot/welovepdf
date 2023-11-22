# Fetch GhostScript source code
curl https://github.com/ArtifexSoftware/ghostpdl-downloads/releases/download/gs10021/ghostpdl-10.02.1.tar.gz --output ghostscript-source.tar.gz
gzip -d -c ghostscript-source.tar.gz | tar -xf -
cd ghostscript-source
# Build binary
./configure
make O=./binary/temp
mv ./binary/temp/gs ./binary
# Remove temp files
rm -rf ./binary/temp
rm -rf ghostscript-source
rm ghostscript-source.tar.gz