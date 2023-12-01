## Must have                          
| Must have 1 - 4h | 
| ---------------------| 
- [ ] Test img resize
- [ ] Test img rotate
- [ ] 30m - FileCard components,  alignment & overflow
- [ ] SIGN BUILDS                
- [ ] new logo by iad ?
- [ ] batch Ids in logs
- [ ] page count in compress
- [ ] restore CompresImage in PdfService
- [ ] disclaimer about allowed image file types
- [ ] Toast if wrong file type
-----
- [ ] check new mem usage       
<br/>

| V2 |                                3h30 + ?? |
| ---------------------| ----------|
| batch Ids in logs ||
|  Logger - send logs by batches ?   | 30m |
| improved logger ||
| hide file names in logs (personal info) replace only letters for easier debugging ||
| Backend logs - 🏆 unexpected errors 🏆
| add tests          |                         2h30|
| add tests backend   |    2h |
| add tests frontend   |   2h30 |
| lint + test pipeline         | 30min |
<br/>

| Nice to have                 | 2h30 + ??? |
| --------------------- | ----------| 
| cache gs binary in build pipeline     | 30min |
| cache node modules & go in test/build pipelines for sharing ?   | 30min |
| format frontend code     | 30min |
| different loader messages     | 30min |
| setup husky                   | 30min |
| open dir after actions        | ??? |
| paralellize when several files... | |
| Split PDF ?                   | 30min |

<br/>
<br/>

## Bugs
- drag n drop load nothing if one file wrong type
- stop process on app close ? (ony in dev ?)

<br/>
<br/>

## Tests
### Merge
- [ ] Merge w/o resize
- [ ] Merge w/ resize
### Compression
- [ ] Optimisaton
- [ ] Compression
- [ ] Compression Extreme
### Convert image
- [ ] Convert image w/o Resize
- [ ] Convert image w Resize
### Format
- [ ] Format to A4
### Drag n Drop
- [ ] Drag n Drop files consecutively several times

## Memory Tests (GS binary + pdfcpu lib)
| Operation | Mem before | Mem after | file size | Leaked Mem | Solution |
| ------------ | ----- | ----- | ----- | ----- | ----- |
| Start (Idle) | N/A | 28 MB | N/A | N/A | N/A |
| DragNDrop | 38 MB | >330MB | > 300 MB | 0 | N/A |
| Format to A4 | 33 MB | 33MB | 6.8 MB | 0 | N/A |
| Merge w resize | 34 | 80 | 26.4 MB | 0 | N/A |
| Merge w/o resize | 34 | 84 MB| 26.4 MB | 0 | N/A |
| Compression | 34 MB| 57 MB | 25.8MB | 23 MB | pdfcpu out in convert |
| Compression Extreme | 57 MB | 71 MB | 12.7 MB | 13 MB | pdfcpu out in convert |
| Convert image w/o Resize | 34 MB | 114 MB | 12.7 MB| 80 MB | pdfcpu out in convert |
| Convert image w Resize | 34 | 101 | 12.7MB | 67 MB | pdfcpu out in convert |

## Memory Tests V2 (only GS binary)
| Operation | Mem before | Mem after | file size | Leaked Mem | Solution |
| ------------ | ----- | ----- | ----- | ----- | ----- |
| Start (Idle) | N/A | 71 MB | N/A | N/A | N/A |
| DragNDrop | 71 MB | 74 MB |  MB | 3 MB | N/A |
| Format to A4 | 33 MB | 33MB | 6.8 MB | 0 | N/A |
| Merge w resize | 34 | 80 | 26.4 MB | 0 | N/A |
| Merge w/o resize | 74 | 77 MB| 26.4 MB | 0 | N/A |
| Compression | 34 MB| 57 MB | 25.8MB | 23 MB | pdfcpu out in convert |
| Compression Extreme | 57 MB | 71 MB | 12.7 MB | 13 MB | pdfcpu out in convert |
| Convert image w/o Resize | 34 MB | 114 MB | 12.7 MB| 80 MB | pdfcpu out in convert |
| Convert image w Resize | 34 | 101 | 12.7MB | 67 MB | pdfcpu out in convert |
