import { GenericPage } from '../components';
import { compressFilesExtreme } from '../actions';

export const CompressFilesPage: React.FC = () => (
    <GenericPage 
        headerText='Veuillez sélectionner les fichiers à comprimer' 
        action={{
            handler: compressFilesExtreme,
            btnLabel: 'Comprimer les fichiers',
            minFilesLength: 1
        }}  
    />
)