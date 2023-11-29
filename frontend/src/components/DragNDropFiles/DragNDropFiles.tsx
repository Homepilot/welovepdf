import { useEffect } from "react";

import { FileUploader } from "react-drag-drop-files";
import toast from "react-hot-toast";

// import DragNDropIcon from '../../assets/images/drag_n_drop.svg'
import {  FileType } from "../../types";
import "./DragNDropFiles.css"

const imgFileExtensions = ["JPG", "JPEG", "PNG", "GIF"];

const isPdf = (fileName: string) => fileName.toLowerCase().endsWith('.pdf')
const isImage = (fileName: string) => {
  const splitted = fileName.split('.')
  if(!splitted.length) return false;
  return imgFileExtensions.includes(splitted[splitted.length -1].toUpperCase())
}

export const DragDrop: React.FC<React.PropsWithChildren & { filesType: FileType, onFilesDropped(fileNames:File[]): Promise<void> }> = ({
  children,  
  filesType,
  onFilesDropped,
}) => {
    const allowedExtensions = filesType === FileType.IMAGE ? imgFileExtensions : ['PDF']

    useEffect(() => {
      const dropzoneElement = window.document.getElementsByClassName('dropzone')[0] as undefined | HTMLLabelElement;
      if(!dropzoneElement){
        console.error('Error finding input element')
        return
      }
      const inputElement = Array.from(dropzoneElement.children).find((elem => elem.tagName.toLowerCase() === 'input')) as undefined | HTMLInputElement
      if(!inputElement){
        console.error('Error finding input element')
        return
      }

      inputElement.onclick = preventClickEventDefault;
      dropzoneElement.onclick = preventClickEventDefault;
      console.log('click events successfully prevented on dropzone')
    }, []) 

    const handleFilesDropped = async (files: Iterable<File>) => {
      console.log({files})
      const {invalidFiles, validFiles} = Array.from(files).reduce((acc, file) => {
        const isValid = filesType === FileType.PDF ? isPdf(file.name) : isImage(file.name)
        return isValid ? ({
          ...acc,
          validFiles: [...acc.validFiles, file]
        }) : ({
          ...acc,
          invalidFiles: [...acc.invalidFiles, file]
        })
    }, { validFiles: [], invalidFiles: [] } as { validFiles: File[], invalidFiles: File[] })

    if(invalidFiles.length){
      invalidFiles.forEach((file) => toast.error(`Le fichier "${file.name}" est de type invalide`));
    }
    if(!validFiles.length){
      return        
    }
    await onFilesDropped(validFiles)
  }

  return (
    <FileUploader 
        handleChange={handleFilesDropped}
        hoverTitle="Glisser des fichiers dans la f+enêtre pour les ajouter à la liste" 
        multiple={true} 
        name="file" 
        label="Vous pouvez aussi glisser des fichiers ici"
        maxSize={50}
        types={allowedExtensions} 
        classes={["dropzone"]}
    >
      {/* <div className="drop-container">
        <img className="drag-n-drop-icon" src={DragNDropIcon} />
        <span>Vous pouvez aussi glisser des fichiers ici</span>
      </div> */}
      {children}
    </FileUploader>
  );
}

export default DragDrop;

function preventClickEventDefault(ev: MouseEvent){
  ev.preventDefault();
  ev.stopPropagation()
}