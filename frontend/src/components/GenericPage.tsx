import { useState } from 'react';
import { FilesList } from '../components';
import { FileType } from '../types';
import toast from 'react-hot-toast';
import { selectMultipleFiles } from '../actions';

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

    const removeFileFromList = (fileId: string) => {
        const newSelectionWithIds = selectedFiles.filter(({id}) => id !== fileId);
        setSelectedFiles(newSelectionWithIds);
    } 

    const selectFiles = async () => {
        const files = await selectMultipleFiles(filesType, selectFilesPrompt ?? headerText);
        const newSelection = Array.from(new Set([...selectedFiles.map(({id}) => id), ...files]));
        const selectionWithIds = newSelection.map(id => ({id}))
        setSelectedFiles(selectionWithIds);
    }

    const emptyList = () => {
        setSelectedFiles([]);
    }


    async function runHandler(){
        
       const result = await action.handler(selectedFiles.map(({id}) => id));

       if(result === null) return;
       
       if(!Array.isArray(result)){
        if(result) {
            toast.success('Opération réussie');
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
            emptyList();
            return;
        }
        
        toast.success(`L'opération a réussi pour ${success} fichiers`);
        toast.error(`L'opération a échoué pour ${failures} fichiers`);
        
        emptyList();
    }

    return (
        <div className='container'>
            <div>
                <h3>{headerText}</h3>
            </div>
            <div className='btn-container'>
                <button disabled={!selectedFiles.length} onClick={() => setSelectedFiles([])} className="btn">Vider la liste</button>
                <button onClick={selectFiles} className="btn">Choisir des fichiers</button>
                <button
                    onClick={runHandler}
                    disabled={selectedFiles.length < action.minFilesLength} 
                    className="btn"
                >
                    { action.btnLabel}
                </button>
            </div>
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

