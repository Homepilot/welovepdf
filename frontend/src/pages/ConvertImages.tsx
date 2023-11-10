import { useState } from 'react';
import { FilesList } from '../components';

export const ConvertImagesPage: React.FC = () => {
    const [selectedFiles, setSelectedFiles] = useState<string[]>([]);


    const convertFiles = async () => {
        // const result = await MergePdfFiles(filesToConvert)
        // console.log({ mergeSuccess: result })
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

