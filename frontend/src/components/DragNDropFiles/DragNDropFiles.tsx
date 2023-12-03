import { useCallback } from "react";

import { FileUploader } from "react-drag-drop-files";
import toast from "react-hot-toast";

import { findFilePathByName } from "../../api";
import {  FileInfo, FileType } from "../../types";

import "./DragNDropFiles.css"
import { ALLOWED_IMAGE_EXTENTIONS } from "./constants";
import { usePreventDropzoneClick } from "./hooks";

type DragDropProps = {
   filesType: FileType;
   onFilesDropped(fileNames:FileInfo[]): void;
   setIsLoading(isLoading: boolean): void;
}

export const DragDrop: React.FC<React.PropsWithChildren & DragDropProps> = ({
  children,  
  filesType,
  onFilesDropped,
  setIsLoading,
}) => {
    const dropzoneClassName = 'dropzone';
    usePreventDropzoneClick([dropzoneClassName])

  const searchFilesInFS = useCallback(async (files: File[]): Promise<{successes: FileInfo[], failures: string[] }> => {
    try {
        const results = await Promise.all(files.map(fileInfo => findFilePathByName(fileInfo.name, fileInfo.size, fileInfo.lastModified)))
        return results.reduce(
            (acc, path, index) => path ? 
            ({ ...acc, successes: [...acc.successes, { id: path, name: files[index].name }]}) : 
            ({ ...acc, failures: [...acc.failures, files[index].name] })
        , {failures: [], successes: []} as {failures: string[], successes: FileInfo[]})
      } catch (error) {
        console.error(error)
        toast.error("Erreur lors de l'ajout des fichiers");
      }
      return { successes: [], failures: [] }      
    }, [findFilePathByName])
    
    
    const handleFilesDropped = useCallback(async (files: Iterable<File>) => {
      setIsLoading(true);
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

      if(!validFiles.length){
        setIsLoading(false);
        invalidFiles.forEach((file) => toast.error(`Le fichier "${file.name}" est de type invalide`));
        return        
      }

      const { successes, failures } = await searchFilesInFS(validFiles);
      setIsLoading(false);

      if (successes.length){
        onFilesDropped(successes)
        toast.success(`${successes.length} fichier(s) ajouté(s)`)
      }
      if(invalidFiles.length){
        invalidFiles.forEach(fileName => toast.error(`"${fileName}" : type de fichier non accepté`))
      }
      if(failures.length){
        failures.forEach(fileName => toast.error(`"${fileName}" : impossible de trouver le fichier, veuillez essayer de l'ajouter manuellement`))
      }
  }, [setIsLoading, onFilesDropped, searchFilesInFS])

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

function isPdf (fileName: string){
  return  fileName.toLowerCase().endsWith('.pdf')
} 

function isImage(fileName: string) {
  const splitted = fileName.split('.')
  if(!splitted.length) return false;
  return ALLOWED_IMAGE_EXTENTIONS.includes(splitted[splitted.length -1].toUpperCase())
}

