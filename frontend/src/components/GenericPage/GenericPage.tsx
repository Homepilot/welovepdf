import { useCallback, useEffect, useState } from 'react';

import toast from 'react-hot-toast';

import { createTempFilesFromUpload, 
    logPageVisited,
    selectMultipleFiles } from '../../api';
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
        const includedFiles = [...selectedFiles];
        const result = await action.handler(includedFiles.map(({id}) => id));
        
        if(!result) {
            setIsLoading(false);
            return;
        }
        
        if(!Array.isArray(result)){
            if(result) {
                toast.success('Opération réussie');
                setIsLoading(false);
                emptyList();
                return;
            }
            toast.error("L'opération a échoué");
            setIsLoading(false);
            
            return;
        }
        
        const { success, failures } = result.reduce<{success: number, failures: number}>(
            (acc, operationResult) => operationResult 
            ? { success: acc.success + 1, failures: acc.failures } 
            : { success: acc.success, failures: acc.failures + 1 }, 
            { success: 0, failures: 0 })
            
            if(failures === 0) {
                toast.success('Opération réussie pour tous les fichiers');
                emptyList();
                setIsLoading(false);
                return;
            }
            
            if(success === 0) {
                toast.error("L'opération a échoué pour tous les fichiers");
                setIsLoading(false);
                return;
            }
            
            toast.success(`L'opération a réussi pour ${success} fichiers`);
            toast.error(`L'opération a échoué pour ${failures} fichiers`);
            
            setSelectedFiles(includedFiles.filter((_, index) => result[index]));
            setIsLoading(false);
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

