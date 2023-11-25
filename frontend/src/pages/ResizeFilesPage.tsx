import { resizeToA4 } from '../api/actions';
import { GenericPage } from '../components';

export const ResizeFilesPage: React.FC<{onNavigateHome(): void}> = ({onNavigateHome}) => (
    <GenericPage 
        onNavigateHome={onNavigateHome}
        headerText='Veuillez sélectionner les fichiers à comprimer' 
        action={{
            handler: resizeToA4,
            btnLabel: 'Comprimer les fichiers',
            minFilesLength: 1
        }}  
    />
)