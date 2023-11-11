import { useState } from 'react';
import { FilesList } from '../components';

type GenericPageProps = {
    headerText: string;
    action: {
        btnLabel: string;
        handler(filesToHandle: string[]): Promise<void>;
        minFilesLength: number;
    }
}

export const GenericPage: React.FC<GenericPageProps> = ({
    headerText,
    action
}) => {
    const [selectedFiles, setSelectedFiles] = useState<string[]>([]);

    return (
        <div className='container'>
            <div>
                <h3>{headerText}</h3>
            </div>
            <FilesList onSelectionUpdated={setSelectedFiles} />
            <button
                onClick={() => action.handler(selectedFiles)}
                disabled={selectedFiles.length < action.minFilesLength} 
                className="btn"
            >
                { action.btnLabel}
            </button>
        </div>
    )
}

