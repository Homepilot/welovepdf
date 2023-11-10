
import { useState } from 'react';
import { SelectMultipleFiles } from '../../wailsjs/go/main/App';

export const FilesList: React.FC<{onFilesSelected(filePathes: string[]): void}> = ({ onFilesSelected }) => {
    const [selectedFiles, setSelectedFiles] = useState<string[]>([]);

    const selectFiles = async () => {
        const files = await SelectMultipleFiles();
        const newSelection = Array.from(new Set([...selectedFiles, ...files]));
        setSelectedFiles(newSelection);
        onFilesSelected(newSelection);
    }

    return (
        <div className='files-list'>
            <div>
                <button onClick={selectFiles} className="btn">Choisir des fichiers</button>
            </div>
            {selectedFiles.length ? selectedFiles.map((path) => <h3 key={path} >{path}</h3>) : <h3>Aucun fichier sélectionné</h3>}
        </div>
    )
}