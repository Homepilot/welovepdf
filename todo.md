## Must have                          
| Must have 1 - 4h | 
| ---------------------| 
- [ ] SIGN BUILDS                
- [ ] resize images before converting to PDF 
- [ ] 30m - per package constants files
- [ ] 30m - FileCard components,  alignment & overflow
- [ ] 30m -  config file in build pipeline
- [ ] 10m - style toasts
- [ ] 10m - Rotate Logtail token (leaked in source code)
- [ ] ?? - search through whole user home dir
-----
- [ ] check new mem usage       
<br/>

| Known bugs |                                 
| ---------------------|
| Run GS commands |
<br/>

| V2 |                                3h30 + ?? |
| ---------------------| ----------|
- [ ] 30m -  Logger - send logs by batches ?           
| improved logger ||
| hide file names in logs (personal info) replace only letters for easier debugging ||
| Backend logs - 🏆 unexpected errors 🏆
| add tests          |                         2h30|
| add tests backend   |    2h|
| add tests frontend   |   2h30|
| test pipeline         |                      30min|
<br/>

| Nice to have                 | 2h30 + ??? |
| --------------------- | ----------|
| different loader messages     | 30min |
| setup husky                   | 30min |
| open dir after actions        | ??? |
| paralellize when several files... | |
| Split PDF ?                   | 30min |

<br/>
<br/>

## Tests
### Merge
- [x] Merge w/o resize
- [x] Merge w/ resize
### Compression
- [x] Optimisaton
- [x] Compression
- [x] Compression Extreme
### Convert image
- [x] Convert image w/o Resize
- [x] Convert image w Resize
### Format
- [x] Format to A4
### Drag n Drop
- [x] Drag n Drop files consecutively several times

## Memory Tests
| Operation | Mem before | Mem after | file size | Leaked Mem | Solution |
| ------------ | ----- | ----- | ----- | ----- | ----- |
| Start (Idle) | N/A | 28 MB | N/A | N/A | N/A |
| Format to A4 | 38 MB | >330MB | > 300 MB | 0 | N/A |
| Format to A4 | 33 MB | 33MB | 6.8 MB | 0 | N/A |
| Merge w resize | 34 | 80 | 26.4 MB | 0 | N/A |
| Merge w/o resize | 34 | 84 MB| 26.4 MB | 0 | N/A |
| Compression | 34 MB| 57 MB | 25.8MB | 23 MB | pdfcpu out in convert |
| Compression Extreme | 57 MB | 71 MB | 12.7 MB | 13 MB | pdfcpu out in convert |
| Convert image w/o Resize | 34 MB | 114 MB | 12.7 MB| 80 MB | pdfcpu out in convert |
| Convert image w Resize | 34 | 101 | 12.7MB | 67 MB | pdfcpu out in convert |
