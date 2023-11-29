import {
    FindFileInUserDir,
    OpenSaveFileDialog,
    PromptUserSelect,
    SelectMultipleFiles
} from '../../wailsjs/go/models/App';
import {
    ConvertImageToPdf,
    CompressFile,
    MergePdfFiles,
    ResizePdfFileToA4,
    CreateTempFilesFromUpload,
} from '../../wailsjs/go/models/PdfService';
import {
    BrowserOpenURL
} from '../../wailsjs/runtime/runtime';
import { CompressionMode, FileInfo, FileType, PageName } from '../types';

import { logOperationCanceledByUser } from './logger';

export async function selectMultipleFiles(fileType: FileType = FileType.PDF, selectFilesPrompt: string){
    return SelectMultipleFiles(fileType, selectFilesPrompt);
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export async function resizeToA4(filesPathes: string[], _batchId: string = "unknown_batch") {
    const result = await Promise.all(filesPathes.map(path => ResizePdfFileToA4(path)));
    console.log({ conversionSuccess: result })
    return result;
}

export async function convertFiles(filesPathes: string[], batchId: string = "unknown_batch") {
    const shouldResize = await chooseShouldResize();
    console.log({shouldResize})
    if(shouldResize === null) {
        logOperationCanceledByUser(PageName.CONVERT_IMG, batchId)
        return null
    }

    const result = await Promise.all(filesPathes.map(path => ConvertImageToPdf(path, shouldResize)));
    console.log({ conversionSuccess: result })
    return result;
}

export async function mergeFiles (filesPathes: string[], batchId: string = "unknown_batch") {
    
    if(filesPathes.length < 2) {
        console.error('Vous devez sélectionner au moins 2 fichiers');
        return false;
    }
    const targetFilePath = await OpenSaveFileDialog();
    console.log({ targetFilePath })
    if(!targetFilePath || targetFilePath === '.pdf') {
        logOperationCanceledByUser(PageName.MERGE, batchId)
        return null
    }
    const shouldResize = await chooseShouldResize();
    if(shouldResize === null) {
        logOperationCanceledByUser(PageName.MERGE, batchId)
        return null
    }

    const result = await MergePdfFiles(targetFilePath, [...filesPathes], shouldResize)
    console.log({ mergeSuccess: result })
    return result;
}

export async function compressFiles (filesPathes: string[], batchId: string = "unknown_batch"): Promise<boolean[] | null> {
    const resultsArray = [];

    const chosenCompressionMode = await chooseCompressionMode() as CompressionMode | '';

    if(!chosenCompressionMode){
        logOperationCanceledByUser(PageName.COMPRESS, batchId)
        return null
    }
    const targetImageQuality = chosenCompressionMode === CompressionMode.EXTREME ? 10 : 20;

    for (const file of filesPathes){
        const result = await CompressFile(file, targetImageQuality)
        resultsArray.push(result)
    }

    console.log({ compressionSuccess: resultsArray })
    return resultsArray;
}

export async function chooseCompressionMode(): Promise<CompressionMode | null>{
    const result = await PromptUserSelect({
        Title:        "Mode de compression",
		Message:      "Choississez un mode de compression",
		Buttons:      [ "Compression", "Compression extrême"],
        Icon:         "compress",
    }) as CompressionMode | '';

    if(result === "") return null;
    return result;
}

export async function chooseShouldResize(): Promise<boolean | null>{
    const result = await PromptUserSelect({
        Title:        "Formattage A4",
		Message:      "Souhaitez convertir le fichier au format A4?",
		Buttons:      ["Oui", "Non"],
        Icon:         "resizeA4",
    })

    if(result === "") return null;
    return result === 'Oui'
}

export function openLinkInBrowser(url: string){
    return BrowserOpenURL(url)
}


export async function createTempFilesFromUpload(files: File[]): Promise<FileInfo[]> {
    const filesAsArrBuff = await Promise.all(files.map(file => file.arrayBuffer()));
    const filesAsUint8Arr = filesAsArrBuff.map(arrBuff => new Uint8Array(arrBuff));
    const result = await Promise.all(filesAsUint8Arr.map<Promise<string>>(file => CreateTempFilesFromUpload([...file])))

    const newFileInfos = result.reduce<FileInfo[]>((fileInfos, filePath, i) => {
        if(!filePath) return fileInfos;
        
        return [...fileInfos, {
            id: filePath,
            name: files[i].name
        }]
    }, [] as FileInfo[])
    
    console.log({newFileInfos})
    return newFileInfos
}

export async function findFilePathByName(fileName: string, size: number, lastModifiedAt: number): Promise<string | null> {
    const matchingFile = await FindFileInUserDir(fileName, size, lastModifiedAt);

    return matchingFile || null
}
// type FnToRun<Args extends unknown[], Return = boolean> = (...args: Args) => Return | Promise<Return>;

// async function noFail<T extends unknown[], R>(fnToRun: FnToRun<T, R>, errorValue: R, ...args: T): Promise<R>{
//     try {
//         let result = fnToRun(...args);
//         if(result instanceof Promise) result = await result;
//         return result; 
//     } catch (error) {
//         console.error(`Error running fn "${fnToRun.name}"`, args)
//         return errorValue
//     }
// }