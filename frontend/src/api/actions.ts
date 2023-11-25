import {
    OpenSaveFileDialog,
    PromptUserSelect,
    SelectMultipleFiles
} from '../../wailsjs/go/models/App';
import {
    ConvertImageToPdf,
    CompressFile,
    MergePdfFiles,
    OptimizePdfFile,
    ResizePdfFileToA4,
} from '../../wailsjs/go/models/PdfHandler';
import {
    BrowserOpenURL
} from '../../wailsjs/runtime/runtime';
import { CompressionMode, FileType } from '../types';

export async function selectMultipleFiles(fileType: FileType = FileType.PDF, selectFilesPrompt: string){
    return SelectMultipleFiles(fileType, selectFilesPrompt);
}

export async function resizeToA4(filesPathes: string[]) {
    const shouldResize = await chooseShouldResize();
    if(shouldResize === null) return null

    const result = await Promise.all(filesPathes.map(path => ResizePdfFileToA4(path)));
    console.log({ conversionSuccess: result })
    return result;
}

export async function convertFiles(filesPathes: string[]) {
    const shouldResize = await chooseShouldResize();
    if(shouldResize === null) return null

    const result = await Promise.all(filesPathes.map(path => ConvertImageToPdf(path, shouldResize)));
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

    const chosenCompressionMode = await chooseCompressionMode() as CompressionMode | '';

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

export function chooseCompressionMode(){
    return PromptUserSelect({
        Title:        "Mode de compression",
		Message:      "Choississez un mode de compression",
		Buttons:      ["Optimisation", "Compression", "Compression extrême"],
        Icon:         "compress",
    })
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

export function createTempFilesFromUpload(files: File[]){

    return files.map(file => ({name: file.name, path: file.name  }))
}