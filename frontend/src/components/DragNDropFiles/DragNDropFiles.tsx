import { useEffect } from "react";

import { FileUploader } from "react-drag-drop-files";

import { FileType } from "../../types";
import "./DragNDropFiles.css"

const imgFileExtensions = ["JPG", "JPEG", "PNG", "GIF"];

export const DragDrop: React.FC<React.PropsWithChildren & { filesType: FileType, onFilesDropped(file:File[]): void }> = ({
    filesType,
    onFilesDropped,
    children,
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

  return (
    <FileUploader 
        handleChange={onFilesDropped} 
        onClick={() => {}}
        hoverTitle="Glisser des fichiers dans la denêtre pour les ajouter à la liste" 
        multiple={true} 
        name="file" 
        maxSize={50}
        types={allowedExtensions} 
        classes={["dropzone"]}
    >
      {children}
      </FileUploader>
  );
}

export default DragDrop;

function preventClickEventDefault(ev: MouseEvent){
  ev.preventDefault();
  ev.stopPropagation()
}