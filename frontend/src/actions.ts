import { FileType } from './types';
import { 
    CompressPdfFile,
    ConvertImageToPdf,
    MergePdfFiles,
    SelectMultipleFiles
} from '../wailsjs/go/main/App';

export async function selectMultipleFiles(fileType: FileType = FileType.PDF, selectFilesPrompt: string){
    return SelectMultipleFiles(fileType, selectFilesPrompt);
}

export async function convertFiles(filesPathes: string[]) {
    const result = await Promise.all(filesPathes.map(ConvertImageToPdf))
    console.log({ conversionSuccess: result })
}

export async function mergeFiles (filesPathes: string[]) {

    if(filesPathes.length < 2) {
        console.error('Vous devez sélectionner au moins 2 fichiers');
        return;
    }
    const result = await MergePdfFiles([...filesPathes])
    console.log({ mergeSuccess: result })
}

export async function compressFiles (filesPathes: string[]) {
    const result = await Promise.all(filesPathes.map(CompressPdfFile))
    console.log({ compressionSuccess: result })
}