import { FileType } from './types';
import {
    OpenSaveFileDialog,
    SelectMultipleFiles
} from '../wailsjs/go/main/App';
import { 
    CompressFile,
    ConvertImageToPdf,
    MergePdfFiles,
} from '../wailsjs/go/main/PdfTools';

export async function selectMultipleFiles(fileType: FileType = FileType.PDF, selectFilesPrompt: string){
    return SelectMultipleFiles(fileType, selectFilesPrompt);
}

export async function convertFiles(filesPathes: string[]) {
    const result = await Promise.all(filesPathes.map(ConvertImageToPdf))
    console.log({ conversionSuccess: result })
    return result;
}

export async function mergeFiles (filesPathes: string[]) {
    
    if(filesPathes.length < 2) {
        console.error('Vous devez sélectionner au moins 2 fichiers');
        return false;
    }
    const targetFilePath = await OpenSaveFileDialog();
    if(!targetFilePath) return null;
    console.log('TARGET')
    const result = await MergePdfFiles(targetFilePath, [...filesPathes])
    console.log({ mergeSuccess: result })
    return result;
}

export async function compressFiles (filesPathes: string[]) {
    const result = await Promise.all(filesPathes.map(CompressFile))
    console.log({ compressionSuccess: result })
    return result;
}