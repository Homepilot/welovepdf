import { resizeToA4 } from '../api/actions';
import { GenericPage } from '../components';

export const ResizeFilesPage: React.FC = () => (
    <GenericPage 
        headerText='Veuillez sélectionner les fichiers à comprimer' 
        action={{
            handler: resizeToA4,
            btnLabel: 'Comprimer les fichiers',
            minFilesLength: 1
        }}  
    />
)