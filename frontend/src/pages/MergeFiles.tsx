import { GenericPage } from '../components';
import { mergeFiles } from '../actions';

export const MergeFilesPage: React.FC = () =>  (
    <GenericPage 
        headerText='Veuillez sélectionner les fichiers à fusionner' 
        action={{
            handler: mergeFiles,
            btnLabel: 'Fusionner les fichiers',
            minFilesLength: 2
        }}  
    />
)
