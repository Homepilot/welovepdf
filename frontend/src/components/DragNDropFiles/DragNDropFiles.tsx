import { useEffect } from "react";

import { FileUploader } from "react-drag-drop-files";

import DragNDropIcon from '../../assets/images/drag_n_drop.svg'
import { FileType } from "../../types";
import "./DragNDropFiles.css"

const imgFileExtensions = ["JPG", "JPEG", "PNG", "GIF"];

export const DragDrop: React.FC<{ filesType: FileType, onFilesDropped(file:File[]): void }> = ({
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

  return (
    <FileUploader 
        handleChange={onFilesDropped}
        hoverTitle="Glisser des fichiers dans la f+enêtre pour les ajouter à la liste" 
        multiple={true} 
        name="file" 
        label="Vous pouvez aussi glisser des fichiers ici"
        maxSize={50}
        types={allowedExtensions} 
        classes={["dropzone"]}
    >
      <div className="drop-container">
        <img className="drag-n-drop-icon" src={DragNDropIcon} />
        <span>Vous pouvez aussi glisser des fichiers ici</span>
      </div>
    </FileUploader>
  );
}

export default DragDrop;

function preventClickEventDefault(ev: MouseEvent){
  ev.preventDefault();
  ev.stopPropagation()
}