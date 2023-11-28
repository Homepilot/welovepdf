import { convertFiles } from '../api/actions';
import { GenericPage } from '../components';
import { FileType, PageName } from '../types';

export const ConvertImagesPage: React.FC<{onNavigateHome(): void}> = ({onNavigateHome}) => (
    <GenericPage
        pageName={PageName.CONVERT_IMG}
        inputFilesType={FileType.IMAGE}
        onNavigateHome={onNavigateHome} 
        headerText='Veuillez sélectionner les fichiers à convertir' 
        action={{
            handler: convertFiles,
            btnLabel: 'Convertir les fichiers',
            minFilesLength: 1
        }}
    />
)

