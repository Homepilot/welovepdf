import { useState } from 'react';

import toast from 'react-hot-toast';

import { selectMultipleFiles } from '../../api/actions';
import { FileType } from '../../types';
import { Backdrop } from '../Backdrop';
import { FilesList } from '../FilesList';

import './GenericPage.css';

type GenericPageProps = {
    headerText: string;
    filesType?: FileType;
    action: {
        btnLabel: string;
        handler(filesToHandle: string[]): Promise<boolean | boolean[] | null>;
        minFilesLength: number;
    };
    selectFilesPrompt?: string;
}

export const GenericPage: React.FC<GenericPageProps> = ({
    headerText,
    action,
    filesType,
    selectFilesPrompt,
}) => {
    const [isLoading, setIsLoading] = useState(false);
    const [selectedFiles, setSelectedFiles] = useState<{ id: string}[]>([]);

    // TODO : should use useCallback here ?
    const removeFileFromList = (fileId: string) => {
        const newSelectionWithIds = selectedFiles.filter(({id}) => id !== fileId);
        setSelectedFiles(newSelectionWithIds);
    } 


    // TODO : should use useCallback here ?
    const selectFiles = async () => {
        setIsLoading(true);
        const files = await selectMultipleFiles(filesType, selectFilesPrompt ?? headerText);
        const newSelection = Array.from(new Set([...selectedFiles.map(({id}) => id), ...files]));
        const selectionWithIds = newSelection.map(id => ({id}))
        setSelectedFiles(selectionWithIds);
        setIsLoading(false);
    }
    
    // TODO : should use useCallback here ?
    const emptyList = () => {
        setSelectedFiles([]);
    }
    
    // TODO : should use useCallback here ?
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

    return (
        <>
            <Backdrop isVisible={isLoading} />
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
            <div id="page-body">
                <FilesList 
                    selectedFiles={selectedFiles}
                    onRemoveFileFromList={removeFileFromList}
                    filesType={filesType} 
                    onSelectionUpdated={setSelectedFiles} 
                    selectFilesPrompt={selectFilesPrompt || headerText}
                />
            </div>
        </>
    )
}

