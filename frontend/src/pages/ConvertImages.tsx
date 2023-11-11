import { GenericPage } from '../components';
import { convertFiles } from '../actions';
import { FileType } from '../types';

export const ConvertImagesPage: React.FC = () => (
    <GenericPage 
        headerText='Veuillez sélectionner les fichiers à convertir' 
        action={{
            handler: convertFiles,
            btnLabel: 'Convertir les fichiers',
            minFilesLength: 1
        }}
        filesType={FileType.IMAGE}
    />
)

