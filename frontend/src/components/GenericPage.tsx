import { useState } from 'react';
import { FilesList } from '../components';
import { FileType } from '../types';

type GenericPageProps = {
    headerText: string;
    filesType?: FileType;
    action: {
        btnLabel: string;
        handler(filesToHandle: string[]): Promise<void>;
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
                    onClick={() => action.handler(selectedFiles)}
                    disabled={selectedFiles.length < action.minFilesLength} 
                    className="btn"
                >
                    { action.btnLabel}
                </button>
            </FilesList>
        </div>
    )
}

