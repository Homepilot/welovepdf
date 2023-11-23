![Alt text](assets/images/logo.svg)


# WE LOVE PDF

At Homepilot, we love PDFs! And in order not to share our customers' documents to third parties, we decided to develop this simple tool for internal use.
As it relies on GhostScript which is Open Source, it felt logical to make it Open Source and freely available on the MacOS App Store.

## ℹ️ About
This application allows the user to convert images to PDF, compress PDFs sizes and merge PDF files.
As of now, it is only targetting MacOS.
As this program uses GhostScript which is under AGPL license, it must be open-sourced.
<br/>**Disclaimer : Homepilot does not take any responsiblity for any malfunctions or problems caused by the use of this software**

Note : to be self-supporting and to avoid depending on external libs or tools, this application embeds the GhostScript binary. The main point is ease-of-installation and ease-of-use.

## 🧱 Stack
This project is based on the official Wails React-TS template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## 🚀 Usage
To use the application, you need to build it first. Then, please refer to `User Manual.md` for any questions regarding usage. 

## 🛠️ Dev setup

### Building GhostScript
Currently used version : 10.02.1
Target source archive : https://github.com/ArtifexSoftware/ghostpdl-downloads/releases/download/gs10021/ghostpdl-10.02.1.tar.gz

To build the GhostScript binary into `./binary`, simply run `./build-ghostscript.sh` in a terminal window at the root level of this repository

### Installing frontend dependecies
`cd frontend && npm i`

### Run the app in watch mode
`wails dev` at the root level of this repository

You're good to go !

## 🏗️ Building
⚠️ Before building, make sure you build GhostScript as specified above.

To build a redistributable, production mode package, use `wails build`.

## 🙏 Credits
This software is built with the [Wails framework](https://wails.io/docs/introduction).
It uses : 
    - [GhostScript](https://ghostscript.readthedocs.io/en/latest/) to convert PDFs to images (used as a compression step),
    - [pdfcpu]() to convert images to PDFs and to merge PDF files
Compress icon : 
    - Image by OpenClipart-Vectors from [Pixabay](https://pixabay.com/vectors/compression-archiver-compress-149782/)
    - [License](https://pixabay.com/service/license-summary/)
