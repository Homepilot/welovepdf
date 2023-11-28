import { resizeToA4 } from '../api/actions';
import { GenericPage } from '../components';
import { FileType, PageName } from '../types';

export const ResizeFilesPage: React.FC<{onNavigateHome(): void}> = ({onNavigateHome}) => (
    <GenericPage 
        pageName={PageName.RESIZE}
        inputFilesType={FileType.PDF}
        onNavigateHome={onNavigateHome}
        headerText='Veuillez sélectionner les fichiers à formatter' 
        action={{
            handler: resizeToA4,
            btnLabel: 'Formatter les fichiers',
            minFilesLength: 1
        }}  
    />
)