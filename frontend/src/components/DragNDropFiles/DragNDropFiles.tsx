import { FileUploader } from "react-drag-drop-files";
import toast from "react-hot-toast";

import {  FileType } from "../../types";

import "./DragNDropFiles.css"
import { ALLOWED_IMAGE_EXTENTIONS } from "./constants";
import { usePreventDropzoneClick } from "./hooks";

const isPdf = (fileName: string) => fileName.toLowerCase().endsWith('.pdf')
const isImage = (fileName: string) => {
  const splitted = fileName.split('.')
  if(!splitted.length) return false;
  return ALLOWED_IMAGE_EXTENTIONS.includes(splitted[splitted.length -1].toUpperCase())
}

export const DragDrop: React.FC<React.PropsWithChildren & { filesType: FileType, onFilesDropped(fileNames:File[]): Promise<void> }> = ({
  children,  
  filesType,
  onFilesDropped,
}) => {
    const dropzoneClassName = 'dropzone';
    usePreventDropzoneClick([dropzoneClassName])


    const handleFilesDropped = async (files: Iterable<File>) => {
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
        classes={[dropzoneClassName]}
    >
      {children}
    </FileUploader>
  );
}

export default DragDrop;
