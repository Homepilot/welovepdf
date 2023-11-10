import { useState } from 'react';
import { MergePdfFiles } from '../../wailsjs/go/main/App';
import { FilesList } from '../components';

export const MergeFilesPage: React.FC = () => {
    const [selectedFiles, setSelectedFiles] = useState<string[]>([]);


    const mergeFiles = async () => {

        if(selectedFiles.length < 2) {
            alert('Vous devez sélectionner au moins 2 fichiers');
            return;
        }
        const result = await MergePdfFiles([...selectedFiles])
        console.log({ mergeSuccess: result })
    }
    
    return (
        <div className='container'>
            <h3>Veuillez sélectionner les fichiers à fusionner</h3>
            <FilesList onFilesSelected={setSelectedFiles} />
            {selectedFiles.length ? (
                <button onClick={mergeFiles} className="btn">Fusionner les fichiers</button>
            ): null}
        </div>
    )
}
