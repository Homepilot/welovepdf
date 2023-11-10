import { useState } from 'react';
import { CompressPdfFile } from '../../wailsjs/go/main/App';
import { FilesList } from '../components';

export const CompressFilesPage: React.FC = () => {
    const [selectedFiles, setSelectedFiles] = useState<string[]>([]);


    const compressFiles = async () => {
        const result = await Promise.all(selectedFiles.map(CompressPdfFile))
        console.log({ compressionSuccess: result })
    }
    
    return (
        <div className='container'>
            <h3>Veuillez sélectionner les fichiers à comprimer</h3>
            <FilesList onFilesSelected={setSelectedFiles} />
            {selectedFiles.length ? (
                <button disabled={!selectedFiles.length} onClick={compressFiles} className="btn">Comprimer les fichiers</button>
            ): null}
        </div>
    )
}
