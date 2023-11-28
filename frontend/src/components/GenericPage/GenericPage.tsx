import { useCallback, useEffect, useState } from 'react';

import { v4 as uuidv4 } from 'uuid';

import { createTempFilesFromUpload, 
    notifyAndLogOperationsResult, 
    logPageVisited,
    selectMultipleFiles,
    logOperationStarted, } from '../../api';
import { FileInfo, FileType, PageName } from '../../types';
import { AppFooter } from '../AppFooter';
import { AppHeader } from '../AppHeader';
import { Backdrop } from '../Backdrop';
import { DragDrop } from '../DragNDropFiles';
import { FilesList } from '../FilesList';
import './GenericPage.css';

type GenericPageProps = {
    pageName: PageName
    headerText: string;
    inputFilesType: FileType;
    action: {
        btnLabel: string;
        handler(filesToHandle: string[]): Promise<boolean | boolean[] | null>;
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
    
    const logHomePageVisited = useCallback(async () => {
        await logPageVisited(pageName)
    }, [])
    
    useEffect(() => {
        logHomePageVisited()
    }, [])
    // TODO : should use useCallback here ?
    const removeFileFromList = (fileId: string) => {
        const newSelectionWithIds = selectedFiles.filter(({id}) => id !== fileId);
        setSelectedFiles(newSelectionWithIds);
    } 


    // TODO : should use useCallback here ?
    const selectFiles = async () => {
        setIsLoading(true);
        const files = await selectMultipleFiles(inputFilesType, selectFilesPrompt ?? headerText);
        addFilesToSelectionList(files.map((filePath: string) => ({ name: filePath, id: filePath })))
        setIsLoading(false);
    }
    
    // TODO : should use useCallback here ?
    const emptyList = () => {
        setSelectedFiles([]);
    }
    
    const addFilesToSelectionList = (files: FileInfo[]) => {
        const newSelectionMap = [...selectedFiles, ...files].reduce<Map<string, FileInfo>>((map, fileInfo) => {
            if(!map.has(fileInfo.id)){
                map.set(fileInfo.id, fileInfo)
            } 
            return map;
        }, new Map())

        const updatedSelection = [...newSelectionMap].map(([, fileInfo]) => fileInfo)
        setSelectedFiles(updatedSelection);
    }

    // TODO : should use useCallback here ?
    // TODO Should split

    async function runHandler(){
        setIsLoading(true);
        const batchId = uuidv4();
        logOperationStarted(pageName, batchId);

        const includedFiles = [...selectedFiles];
        const result = await action.handler(includedFiles.map(({id}) => id));
        
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
            emptyList();
            return;
        }
        setSelectedFiles(includedFiles.filter((_, index) => result[index]));
    }

    async function handleFileDrop(files: File[]){
        setIsLoading(true)
        console.log(files)
        const filesArray = Array.from(files)
        const newFileInfos = await createTempFilesFromUpload(filesArray)
        addFilesToSelectionList(newFileInfos)
        setIsLoading(false)
    }


    return (
        <div id="generic-layout">
                <Backdrop isVisible={isLoading} />
                    <div id="generic-page-header">
                        <AppHeader 
                            shouldDisplayHomeBtn={true}
                            onNavigateHome={onNavigateHome}    
                            />
                    </div>
                    <div id="generic-page-container">
                        <div id="page-header">
                            <div id="page-header-text">
                                <h3>{headerText}</h3>
                            </div>
                            <div id='btn-container'>
                                <span onClick={() => setSelectedFiles([])} className={selectedFiles.length ? 'action-btn' : 'action-btn-disabled'}>Vider la liste</span>
                                <span onClick={selectFiles} className="action-btn">{`${selectedFiles.length ? 'Ajouter' : 'Choisir'} des fichiers`}</span>
                                <span
                                    onClick={runHandler}
                                    className={selectedFiles.length >= action.minFilesLength ? 'action-btn' : 'action-btn-disabled'}
                                >
                                    { action.btnLabel}
                                </span>
                            </div>
                        </div>
                            <FilesList 
                                selectedFiles={selectedFiles}
                                onRemoveFileFromList={removeFileFromList}
                                filesType={inputFilesType} 
                                onSelectionUpdated={addFilesToSelectionList} 
                                selectFilesPrompt={selectFilesPrompt || headerText}
                            />
                    </div>
                    <DragDrop filesType={inputFilesType} onFilesDropped={handleFileDrop} />
                    <AppFooter/>
                </div>
    )
}

