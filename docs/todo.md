
## Must have
- (spinner, backdrop, etc)
- add logger
- keep mem usage in check
- convert images to A4 format...
- build pipeline
- prerequisites.md
- replace icon in Dialog

## Nice to have
- open dir after actions
- optmize after compression ?
- setup husky
- remove unused dependencies
- drag files into window (cf wails options CSSDragProperty: "--wails-draggable" ,CSSDragValue: "drag",)

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