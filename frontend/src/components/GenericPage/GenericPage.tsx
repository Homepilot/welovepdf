import { useCallback, useEffect, useState } from 'react';

import toast from 'react-hot-toast';
import { v4 as uuidv4 } from 'uuid';

import {  
    notifyAndLogOperationsResult, 
    logPageVisited,
    selectMultipleFiles,
    logOperationStarted,
    findFilePathByName,
} from '../../api';
import { FileInfo, FileType, PageName } from '../../types';
import { AppFooter } from '../AppFooter';
import { AppHeader } from '../AppHeader';
import { Backdrop } from '../Backdrop';
import { DragDrop } from '../DragNDropFiles';
import { FilesList } from '../FilesList';
import './GenericPage.css';
import { PageHeader } from '../PageHeader';

type GenericPageProps = {
    pageName: PageName
    headerText: string;
    inputFilesType: FileType;
    action: {
        btnLabel: string;
        handler(filesToHandle: string[], batchId: string): Promise<boolean | boolean[] | null>;
        minFilesLength: number;
    };
    selectFilesPrompt?: string;
    onNavigateHome(): void;
}

export const GenericPage: React.FC<GenericPageProps> = ({
    headerText,
    action,
    inputFilesType,
    selectFilesPrompt,
    onNavigateHome,
    pageName
}) => {
    const [isLoading, setIsLoading] = useState(false);
    const [selectedFiles, setSelectedFiles] = useState<FileInfo[]>([]);
    
    const logGenericPageVisited = useCallback(async () => {
        await logPageVisited(pageName)
    }, [pageName])
    
    useEffect(() => {
        logGenericPageVisited()
    }, [])

    const removeFileFromList = useCallback((fileId: string) => {
        const newSelectionWithIds = selectedFiles.filter(({id}) => id !== fileId);
        setSelectedFiles(newSelectionWithIds);
    }, [selectedFiles, setSelectedFiles]) 

    const addFilesToSelectionList = useCallback((files: FileInfo[]) => {
        const newSelectionMap = [...selectedFiles, ...files].reduce<Map<string, FileInfo>>((map, fileInfo) => {
            if(!map.has(fileInfo.id)){
                map.set(fileInfo.id, fileInfo)
            } 
            return map;
        }, new Map())

        const updatedSelection = [...newSelectionMap].map(([, fileInfo]) => fileInfo)
        setSelectedFiles(updatedSelection);
    }, [selectedFiles, setSelectedFiles])

    const selectFiles = useCallback(async () => {
        setIsLoading(true);
        const files = await selectMultipleFiles(inputFilesType, selectFilesPrompt ?? headerText);
        addFilesToSelectionList(files.map((filePath: string) => ({ name: filePath, id: filePath })))
        setIsLoading(false);
    }, [setIsLoading, inputFilesType, selectFilesPrompt, headerText, addFilesToSelectionList])

    const runHandler = useCallback(async () => {
        setIsLoading(true);
        const batchId = uuidv4();
        logOperationStarted(pageName, batchId);

        const includedFiles = [...selectedFiles];
        const result = await action.handler(includedFiles.map(({id}) => id), batchId);
        
        if(!['boolean', 'object'].includes(typeof result)) {
            setIsLoading(false);
            return;
        }

        const { successes, failures } = 
            !Array.isArray(result)
            ? { successes: result ? 1 : 0, failures: result ? 0 : 1 }
            : result.reduce<{successes: number, failures: number}>(
                (acc, operationResult) => operationResult 
                ? { successes: acc.successes + 1, failures: acc.failures } 
                : { successes: acc.successes, failures: acc.failures + 1 }, 
            { successes: 0, failures: 0 })


        setIsLoading(false);
        notifyAndLogOperationsResult(pageName, batchId, { successes, failures })
        if(!Array.isArray(result)){
            if(result){
                setSelectedFiles([]);
            }
            return;
        }
        setSelectedFiles(includedFiles.filter((_, index) => result[index]));
    }, [setIsLoading, action, selectedFiles, setSelectedFiles])

    const handleFilesDropped = useCallback(async (fileNames: File[]) => {
        setIsLoading(true)
        let failuresNames: string[] = [];
        try {
            const results = await Promise.all(fileNames.map(fileInfo => findFilePathByName(fileInfo.name, fileInfo.size, fileInfo.lastModified)))
            const { failures, successes } = results.reduce(
                (acc, path, index) => path ? 
                ({ ...acc, successes: [...acc.successes, { id: path, name: fileNames[index].name }]}) : 
                ({ ...acc, failures: [...acc.failures, fileNames[index].name] })
            , {failures: [], successes: []} as {failures: string[], successes: FileInfo[]})
            failuresNames = failures;
            addFilesToSelectionList(successes)
        } catch (error) {
            console.error(error)
            toast.error("Erreur lors de l'ajout des fichiers");
        }
        setIsLoading(false)
        if (!failuresNames.length) return
        failuresNames.forEach((fileName) => toast.error(`${fileName}: erreur lors de l'import du fichier, essayez de l'ajouter manuellement`))
    }, [setIsLoading, pageName, selectedFiles, addFilesToSelectionList, ])

    return (
        <>
            <Backdrop isVisible={isLoading} />
            <div id="generic-layout">
                <AppHeader 
                    shouldDisplayHomeBtn={true}
                    onNavigateHome={onNavigateHome}    
                />
                <PageHeader
                    headerText={headerText}
                    actionLabel={action.btnLabel}
                    isSelectionEmpty={!selectedFiles.length}
                    isActionDisabled={selectedFiles.length < action.minFilesLength}
                    onEmptyList={() => setSelectedFiles([])}
                    onSelectFiles={selectFiles}
                    onRunAction={runHandler}
                />
                <DragDrop filesType={inputFilesType} onFilesDropped={handleFilesDropped} >
                    <FilesList 
                        selectedFiles={selectedFiles}
                        onRemoveFileFromList={removeFileFromList}
                        filesType={inputFilesType} 
                        onSelectionReordered={setSelectedFiles} 
                        selectFilesPrompt={selectFilesPrompt || headerText}
                    />
                </DragDrop>
                <AppFooter/>
            </div>
        </>
    )
}

