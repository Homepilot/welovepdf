import { compressFiles } from '../api/actions';
import { GenericPage } from '../components';
import { FileType } from '../types';

export const CompressFilesPage: React.FC<{onNavigateHome(): void}> = ({onNavigateHome}) => (
    <GenericPage 
        inputFilesType={FileType.PDF}
        onNavigateHome={onNavigateHome}
        headerText='Veuillez sélectionner les fichiers à comprimer' 
        action={{
            handler: compressFiles,
            btnLabel: 'Comprimer les fichiers',
            minFilesLength: 1
        }}  
    />
)