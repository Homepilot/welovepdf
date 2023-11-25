import { mergeFiles } from '../api/actions';
import { GenericPage } from '../components';

export const MergeFilesPage: React.FC<{onNavigateHome(): void}> = ({onNavigateHome}) =>  (
    <GenericPage 
        onNavigateHome={onNavigateHome}
        headerText='Veuillez sélectionner les fichiers à fusionner' 
        action={{
            handler: mergeFiles,
            btnLabel: 'Fusionner les fichiers',
            minFilesLength: 2
        }}  
    />
)
