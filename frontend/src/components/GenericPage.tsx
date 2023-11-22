import { useState } from 'react';

import toast from 'react-hot-toast';

import { selectMultipleFiles } from '../actions';
import { FilesList } from '../components';
import { FileType } from '../types';

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
    const [selectedFiles, setSelectedFiles] = useState<{ id: string}[]>([]);

    // TODO : should use useCallback here ?
    const removeFileFromList = (fileId: string) => {
        const newSelectionWithIds = selectedFiles.filter(({id}) => id !== fileId);
        setSelectedFiles(newSelectionWithIds);
    } 


    // TODO : should use useCallback here ?
    const selectFiles = async () => {
        const files = await selectMultipleFiles(filesType, selectFilesPrompt ?? headerText);
        const newSelection = Array.from(new Set([...selectedFiles.map(({id}) => id), ...files]));
        const selectionWithIds = newSelection.map(id => ({id}))
        setSelectedFiles(selectionWithIds);
    }

    // TODO : should use useCallback here ?
    const emptyList = () => {
        setSelectedFiles([]);
    }

    // TODO : should use useCallback here ?
    async function runHandler(){
        const includedFiles = [...selectedFiles];
       const result = await action.handler(includedFiles.map(({id}) => id));

       if(result === null) return;
       
       if(!Array.isArray(result)){
        if(result) {
            toast.success('Opération réussie');
            emptyList();
            return;
        }
        toast.error("L'opération a échoué");
        
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
            return;
        }
        
        if(success === 0) {
            toast.error("L'opération a échoué pour tous les fichiers");
            return;
        }
        
        toast.success(`L'opération a réussi pour ${success} fichiers`);
        toast.error(`L'opération a échoué pour ${failures} fichiers`);
        
        setSelectedFiles(includedFiles.filter((_, index) => result[index]));
    }

    return (
        <div className='container'>
            <div>
                <h3>{headerText}</h3>
            </div>
            <div className='btn-container'>
                <span onClick={() => setSelectedFiles([])} className={selectedFiles.length ? 'hp-btn' : 'hp-btn-disabled'}>Vider la liste</span>
                <span onClick={selectFiles} className="hp-btn">{`${selectedFiles.length ? 'Ajouter' : 'Choisir'} des fichiers`}</span>
                <span
                    onClick={runHandler}
                    className={selectedFiles.length >= action.minFilesLength ? 'hp-btn' : 'hp-btn-disabled'}
                >
                    { action.btnLabel}
                </span>
            </div>
            <hr/>
            <FilesList 
                selectedFiles={selectedFiles}
                onRemoveFileFromList={removeFileFromList}
                filesType={filesType} 
                onSelectionUpdated={setSelectedFiles} 
                selectFilesPrompt={selectFilesPrompt || headerText}
            >
            </FilesList>
        </div>
    )
}

