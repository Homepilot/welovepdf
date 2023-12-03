import { useCallback, useEffect, useState } from 'react';

import { v4 as uuidv4 } from 'uuid';

import {  
    notifyAndLogOperationsResult, 
    logPageVisited,
    selectMultipleFiles,
    logOperationStarted,
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
    }, [logPageVisited, pageName])
    
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
        addFilesToSelectionList(files.map((filePath: string) => ({ name: getFileNameFromPath(filePath), id: filePath })))
        setIsLoading(false);
    }, [setIsLoading, selectMultipleFiles, inputFilesType, selectFilesPrompt, headerText, addFilesToSelectionList, getFileNameFromPath])

    const runHandler = useCallback(async () => {
        setIsLoading(true);
        const batchId = uuidv4();
        logOperationStarted(pageName, batchId);

        const includedFiles = [...selectedFiles];
        const result = await action.handler(includedFiles.map(({id}) => id), batchId);
        
        console.log({result})
        if(result == null || !['boolean', 'object'].includes(typeof result)) {
            setIsLoading(false);
            return;
        }

        if(!Array.isArray(result)){
            const operationResult = { successes: result ? 1 : 0, failures: result ? 0 : 1 }
            if(result){
                setSelectedFiles([])
            }
            setIsLoading(false);
            notifyAndLogOperationsResult(pageName, batchId, operationResult);
            return
        }

        const { successes, failures } = result.reduce<{successes: number, failures: number}>(
                (acc, operationResult) => operationResult 
                ? { successes: acc.successes + 1, failures: acc.failures } 
                : { successes: acc.successes, failures: acc.failures + 1 }, 
            { successes: 0, failures: 0 })

        setIsLoading(false);
        notifyAndLogOperationsResult(pageName, batchId, { successes, failures })
        setSelectedFiles(includedFiles.filter((_, index) => !result[index]));
    }, [setIsLoading, logOperationStarted, action, selectedFiles, setSelectedFiles, notifyAndLogOperationsResult])

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
                <DragDrop filesType={inputFilesType} onFilesDropped={addFilesToSelectionList} setIsLoading={setIsLoading}>
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


function getFileNameFromPath(pathString: string) {
    const splitted = pathString.split('/')
    return splitted[splitted.length -1]
}