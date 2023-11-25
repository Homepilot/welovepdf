import {
    ChooseCompressionMode,
    OpenSaveFileDialog,
    SelectMultipleFiles
} from '../../wailsjs/go/models/App';
import {
    ConvertImageToPdf,
    CompressFile,
    MergePdfFiles,
    OptimizePdfFile
} from '../../wailsjs/go/models/PdfHandler';
import { CompressionMode, FileType } from '../types';

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

export async function compressFiles (filesPathes: string[]): Promise<boolean[] | null> {
    const resultsArray = [];

    const chosenCompressionMode = await ChooseCompressionMode() as CompressionMode | '';

    if(!chosenCompressionMode) return null;

    if(chosenCompressionMode === CompressionMode.OPTIMIZE) return optimizeFiles(filesPathes);

    const targetImageQuality = CompressionMode.EXTREME ? 10 : 20;

    for (const file of filesPathes){
        const result = await CompressFile(file, targetImageQuality)
        resultsArray.push(result)
    }

    console.log({ compressionSuccess: resultsArray })
    return resultsArray;
}