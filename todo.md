
## Must have
- prompt 3 compression modes
- style 2 (Abde ?) (spinner, backdrop, etc)
- add logger
- keep mem usage in check
- convert images to A4 format...
- build pipeline

## Nice to have
- open dir after actions
- optmize after compression ?
- setup husky
- remove unused dependencies

# Compressing
## 3 compression modes
### Optimization
Simply optimize the PDF
### Compression
Divide by 4 the image quality 
### Extreme Compression
Divide by 8 the image quality
#### Compress to 19MB
Automatically compress the given file to 19MB


## Steps
1. Convert PDF to JPEG files (depends on ghostscript)
2. Lower image quality with given ratio for each page
3. Convert all pages to pdf
4. Merge all pages in the right order