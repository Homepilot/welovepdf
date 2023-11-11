import { useState } from 'react';
import { FilesList } from '../components';
import { FileType } from '../types';
import toast from 'react-hot-toast';

type GenericPageProps = {
    headerText: string;
    filesType?: FileType;
    action: {
        btnLabel: string;
        handler(filesToHandle: string[]): Promise<boolean | boolean[]>;
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
    const [selectedFiles, setSelectedFiles] = useState<string[]>([]);

    async function runHandler(){
        
       const result = await action.handler(selectedFiles)
       
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
                return;
            }

            if(success === 0) {
                toast.error("L'opération a échoué pour tous les fichiers");
                return;
            }

            toast.success(`L'opération a réussi pour ${success} fichiers`);
            toast.error(`L'opération a échoué pour ${success} fichiers`);
    }

    return (
        <div className='container'>
            <div>
                <h3>{headerText}</h3>
            </div>
            <FilesList 
                filesType={filesType} 
                onSelectionUpdated={setSelectedFiles} 
                selectFilesPrompt={selectFilesPrompt || headerText}
            >
                <button
                    onClick={runHandler}
                    disabled={selectedFiles.length < action.minFilesLength} 
                    className="btn"
                >
                    { action.btnLabel}
                </button>
            </FilesList>
        </div>
    )
}

