# README
This application allows the user to convert images to PDF, compress PDFs sizes and merge PDF files.
As of now, it is only targetting MacOS.
This program using GhostScript which is under AGPL license, it must be open-sourced.
<br/>**Disclaimer : Homepilot does not take any responsiblity for any malfunctions or problems caused by the use of this software**

## About

This is the official Wails React-TS template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Dev setup

### Building GhostScript
Currently used version: 10.02.1
Target source archive : https://github.com/ArtifexSoftware/ghostpdl-downloads/releases/download/gs10021/ghostpdl-10.02.1.tar.gz

1. Extract the source code archive at the root directory of this repo.
2. From inside the extracted build directory, run in a terminal : bash```
    ./configure
    make O=binary
   ```
3. Remove ./ghostpdf-XX.XX.X.tar.gz

You're good to go !

## Building
! Before building, make sure you build GhostScript as specified above
To build a redistributable, production mode package, use `wails build`.

## Credits & License
This software uses ghostscript to convert pdfs to images.
As GhostScript is under the AGPL license, this software also must be open sourceù

# Ghostscript build command
/opt/homebrew/bin/gcc-13
./configure CC="/opt/homebrew/bin/gcc-13 -arch i386 -arch x86_64 -arch ppc" CPP="/opt/homebrew/bin/gcc-13 -E"
./configure CC="/opt/homebrew/bin/gcc-13 -arch i386 -arch x86_64 -arch ppc" CPP="/opt/homebrew/bin/gcc-13 -E"