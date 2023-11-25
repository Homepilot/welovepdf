import { convertFiles } from '../api/actions';
import { GenericPage } from '../components';
import { FileType } from '../types';

export const ConvertImagesPage: React.FC<{onNavigateHome(): void}> = ({onNavigateHome}) => (
    <GenericPage
        onNavigateHome={onNavigateHome} 
        headerText='Veuillez sélectionner les fichiers à convertir' 
        action={{
            handler: convertFiles,
            btnLabel: 'Convertir les fichiers',
            minFilesLength: 1
        }}
        filesType={FileType.IMAGE}
    />
)

