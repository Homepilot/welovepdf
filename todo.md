- compress
- style 2 (Abde ?)
- add logger
- open dir after actions
- keep mem usage in check
- Go : make PDF methods independant from app.ctx (used for logger)

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