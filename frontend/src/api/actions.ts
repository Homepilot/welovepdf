import {
    CompressFile,
    MergePdfFiles,
    ResizePdfFileToA4,
    ConvertImageToPdf,
    RemoveFile,
} from '../../wailsjs/go/models/PdfService';
import {
    OpenOutputDir,
    SearchFileInUserDir,
    OpenSaveFileDialog,
    PromptUserSelect,
    SelectMultipleFiles
} from '../../wailsjs/go/models/UserPrompter';
import {
    BrowserOpenURL
} from '../../wailsjs/runtime/runtime';
import { CompressionMode, FileType, PageName } from '../types';

import { logOperationCanceledByUser } from './logger';

export async function selectMultipleFiles(fileType: FileType = FileType.PDF, selectFilesPrompt: string){
    return SelectMultipleFiles(fileType, selectFilesPrompt);
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export async function resizeToA4(filesPathes: string[], _batchId: string = "unknown_batch") {
    return Promise.all(filesPathes.map(path => ResizePdfFileToA4(path)));
}

export async function convertImagesToPdf(
        filesPathes: string[], 
    ) {
    return Promise.all(filesPathes.map(path => ConvertImageToPdf(path)));
}

export async function mergeFiles (filesPathes: string[], batchId: string = "unknown_batch") {
    
    if(filesPathes.length < 2) {
        console.error('Vous devez sélectionner au moins 2 fichiers');
        return false;
    }
    const targetFilePath = await OpenSaveFileDialog();
    if(!targetFilePath || targetFilePath === '.pdf') {
        logOperationCanceledByUser(PageName.MERGE, batchId)
        return null
    }
    const shouldResize = await chooseShouldResize();
    if(shouldResize === null) {
        logOperationCanceledByUser(PageName.MERGE, batchId)
        return null
    }

    const mergeResult = await MergePdfFiles(targetFilePath, [...filesPathes]);
    if(!shouldResize) return mergeResult
    if(!mergeResult) return mergeResult

    const [resizeResult] = await resizeToA4([targetFilePath])
    await RemoveFile(targetFilePath)
    return resizeResult
}

export async function compressFiles (filesPathes: string[], batchId: string = "unknown_batch"): Promise<boolean[] | null> {
    const resultsArray = [];

    const chosenCompressionMode = await chooseCompressionMode() as CompressionMode | '';

    if(!chosenCompressionMode){
        logOperationCanceledByUser(PageName.COMPRESS, batchId)
        return null
    }
    const targetImageQuality = chosenCompressionMode === CompressionMode.EXTREME ? 10 : 30;

    for (const file of filesPathes){
        const result = await CompressFile(file, targetImageQuality)
        resultsArray.push(result)
    }
    return resultsArray;
}

export async function chooseCompressionMode(): Promise<CompressionMode | null>{
    const result = await PromptUserSelect({
        Title:        "Mode de compression",
		Message:      "Veuillez choisir un mode de compression",
		Buttons:      [ "Compression", "Compression extrême"],
        Icon:         "compress",
    }) as CompressionMode | '';

    if(result === "") return null;
    return result;
}

export async function chooseShouldResize(): Promise<boolean | null>{
    const result = await PromptUserSelect({
        Title:        "Formattage A4",
		Message:      "Souhaitez-vous convertir le fichier au format A4?",
		Buttons:      ["Oui", "Non"],
        Icon:         "resizeA4",
    })

    if(result === "") return null;
    return result === 'Oui'
}

export function openLinkInBrowser(url: string){
    return BrowserOpenURL(url)
}

export async function findFilePathByName(fileName: string, size: number, lastModifiedAt: number): Promise<string | null> {
    const matchingFile = await SearchFileInUserDir(fileName, size, lastModifiedAt);

    return matchingFile || null
}

export async function openOutputDirInFinder(){
    await OpenOutputDir()
}