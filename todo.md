- -------------------------------------------------------------
    Must Have                                           
- -------------------------------------------------------------
- [ ] remove auto resize on convert
- [ ] Test img resize
- [ ] Test img rotate
- -------------------------------------------------------------
                                                        3h
- -------------------------------------------------------------
- [ ] FileCard components, alignment & overflow         30m
- [ ] batch Ids in logs                                 45m
- [ ] disclaimer about allowed image file types         15m
- [ ] drag n drop load nothing if one file wrong type
- [ ] Toast if wrong file type                          15m
- [ ] 1 line comment on all function (Golang)           30m
- [ ] check all hooks' dependencies                     30m
- [ ] new logo by iad                                   10m
---------------------------------------------------------------
- [ ] check new mem usage       
- [ ] SIGN BUILDS
- [ ] DOCS
<br/>

---------------------------------------------------------------
    Improvements                                        3h30
---------------------------------------------------------------
- [ ]  Logger - send logs by batches ?                  30m 
- [ ] improved logger (use slog only & everywhere)      1h30 
- [ ] hide file names in logs (personal info)           30m
        replace only letters for easier debugging 
- [ ] Backend logs - 🏆 unexpected errors 🏆
- [ ] open dir after actions 
<br/>

---------------------------------------------------------------
    Automated Tests                                     9h
---------------------------------------------------------------
- [ ] add tests backend                                 6h
- [ ] add tests frontend                                2h30 
- [ ] lint + test pipeline                              30min 
<br/>

---------------------------------------------------------------
   Nice to have                                         3h30
---------------------------------------------------------------
- [ ] cache gs binary in build pipeline                 30min 
- [ ] cache node modules & go in test/build pipelines   30min 
- [ ] format frontend code                              30min 
- [ ] different loader messages                         30min 
- [ ] setup husky                                       30min 
- [ ] paralellize when several files... 
- [ ] Split PDF                                         1h 
<br/>

---------------------------------------------------------------
   BUGS                                         3h30
---------------------------------------------------------------
- [ ] drag n drop load nothing if one file wrong type
- [ ] processd files dont disapeear from list (1 file OR several)
- [ ] stop process on app close ? (ony in dev ?)

<br/>
<br/>


## Tests
| [ ] To test | ✅ Pass  | ❌ Fail |
| ------------ | -------- | ------- |

### Backend Tests
### Merge
- ✅ Merge w/o resize
- ✅ Merge w/ resize
### Compression
- ✅ Compression
- ✅ Compression Extreme
### Convert image
- ✅ Convert JPEG image to PDF
- ✅ Convert PNG/TIFF/WEBP image to PDF
### Format
- ✅ Format to A4
### Drag n Drop
- ✅ Drag n Drop files consecutively several times

### Frontend Tests (additional functionalities to test)
- ✅ reorder items in list
- ✅ empty list
- ✅ add files



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
| Convert JPEG image w Resize | 34 | 101 | 12.7MB | 67 MB | pdfcpu out in convert |
| Convert JPEG image w/o Resize | 34 MB | 114 MB | 12.7 MB| 80 MB | pdfcpu out in convert |
| Convert  image w/o Resize | 34 MB | 114 MB | 12.7 MB| 80 MB | pdfcpu out in convert |
| Convert  image w Resize | 34 | 101 | 12.7MB | 67 MB | pdfcpu out in convert |
| Convert PNG image w/o Resize | 34 MB | 114 MB | 12.7 MB| 80 MB | pdfcpu out in convert |
| Convert PNG image w Resize | 34 | 101 | 12.7MB | 67 MB | pdfcpu out in convert |

## Memory Tests V2 (only GS binary)
| Operation | Mem before | Mem after | file size | Leaked Mem | Solution |
| ------------ | ----- | ----- | ----- | ----- | ----- |
| Start (Idle) | N/A | 71 MB | N/A | N/A | N/A |
| DragNDrop | 71 MB | 74 MB |  MB | 3 MB | N/A |
| Format to A4 |71 MB | 73 MB | 1 MB | 2MB | N/A |
| Merge w resize | 34 | 80 | 26.4 MB | 0 | N/A |
| Merge w/o resize | 74 | 77 MB| 26.4 MB | 0 | N/A |
| Compression | 34 MB| 57 MB | 25.8MB | 23 MB | pdfcpu out in convert |
| Compression Extreme | 57 MB | 71 MB | 12.7 MB | 13 MB | pdfcpu out in convert |
| Convert image w/o Resize | 34 MB | 114 MB | 12.7 MB| 80 MB | pdfcpu out in convert |
| Convert image w Resize | 34 | 101 | 12.7MB | 67 MB | pdfcpu out in convert |
