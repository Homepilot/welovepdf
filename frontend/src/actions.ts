import {
    OpenSaveFileDialog,
    SelectMultipleFiles
} from '../wailsjs/go/main/App';
import {
    ConvertImageToPdf,
    CompressFileExtreme,
    MergePdfFiles,
    OptimizePdfFile
} from '../wailsjs/go/main/PdfUtils';

import { FileType } from './types';

export async function selectMultipleFiles(fileType: FileType = FileType.PDF, selectFilesPrompt: string){
    return SelectMultipleFiles(fileType, selectFilesPrompt);
}

export async function convertFiles(filesPathes: string[]) {
    const result = await Promise.all(filesPathes.map((ConvertImageToPdf as (index: string) => Promise<boolean>)))
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

    const result = await MergePdfFiles(targetFilePath, [...filesPathes])
    console.log({ mergeSuccess: result })
    return result;
}

export async function optimizeFiles (filesPathes: string[]) {
    const result = await Promise.all(filesPathes.map(OptimizePdfFile))
    console.log({ optimizationSuccess: result })
    return result;
}

export async function compressFilesExtreme (filesPathes: string[]) {
    const resultsArray = [];
    for (const file of filesPathes){
        const result = await CompressFileExtreme(file)
        resultsArray.push(result)

    }
    console.log({ compressionSuccess: resultsArray })
    return resultsArray;
}