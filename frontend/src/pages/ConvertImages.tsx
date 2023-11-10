import { useState } from 'react';
import { FilesList } from '../components';
import { ConvertImageToPdf } from '../../wailsjs/go/main/App';

export const ConvertImagesPage: React.FC = () => {
    const [selectedFiles, setSelectedFiles] = useState<string[]>([]);


    const convertFiles = async () => {
        const result = await Promise.all(selectedFiles.map(ConvertImageToPdf))
        console.log({ conversionSuccess: result })
    }
    
    return (
        <div className='container'>
            <div>
                <h3>Veuillez sélectionner les fichiers à convertir</h3>
            </div>
            <FilesList onFilesSelected={setSelectedFiles} />
            {selectedFiles.length ? (
                <button onClick={convertFiles} className="btn">Convertir les fichiers</button>
            ): null}
        </div>
    )
}

