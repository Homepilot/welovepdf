## Must have                          
| Must have 1 | |
| ---------------------| ----------|
| SIGN BUILDS                |                 ??|
| keep mem usage in check      | ??? |
<br/>
<br/>

| Known bugs |                                 |
| ---------------------| ----------|
| warning dnd multiple scrollable parents||

| Must have 2 |                                3h30 + ?? |
| ---------------------| ----------|
| hide file names in logs (personal info) replace only letters for easier debugging ||
| Backend logs - 🏆 unexpected errors 🏆
| frontend css layout v3||
| warning dnd multiple scrollable parents||
| final style touch (Abde)||
| Logger - send logs by batches ?                                | 30min |
| add tests          |                         2h30|
| add tests backend   |    2h|
| add tests frontend   |   2h30|
| test pipeline         |                      30min|
| Go file naming convention |                   15 min|

<br/>
<br/>

| Nice to have                 | 2h30 + ??? |
| --------------------- | ----------|
| Split PDF ?                                 | 30min |
| style toasts                  | 30min |
| different loader messages     | 30min |
| setup husky                   | 30min |
| font from backoffice ?        | 30min |
| buttons as component          | 30min |
| refactor layout css           | 30min |
| open dir after actions        | ??? |
| internal mode unlocked w/ PIN | |
| paralellize when several files... | |

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


#### Last reaminging leaking operations : 
  - convert image to PDF => ok but needs resizing...
    - embed & save asset file
    - ensure asset file at GS setup
  - drop file into window (search by name instead of loading file?)