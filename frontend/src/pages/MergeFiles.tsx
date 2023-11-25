import { mergeFiles } from '../api/actions';
import { GenericPage } from '../components';
import { FileType } from '../types';

export const MergeFilesPage: React.FC<{onNavigateHome(): void}> = ({onNavigateHome}) =>  (
    <GenericPage 
        inputFilesType={FileType.PDF}
        onNavigateHome={onNavigateHome}
        headerText='Veuillez sélectionner les fichiers à fusionner' 
        action={{
            handler: mergeFiles,
            btnLabel: 'Fusionner les fichiers',
            minFilesLength: 2
        }}  
    />
)
