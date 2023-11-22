import { compressFilesExtreme } from '../actions';
import { GenericPage } from '../components';

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