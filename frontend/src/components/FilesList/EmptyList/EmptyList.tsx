import DragNDropIcon from '../../../assets/images/drag_n_drop.svg';
import { FileType } from '../../../types';
import { ALLOWED_IMAGE_EXTENTIONS } from '../../DragNDropFiles/constants';
import './EmptyList.css';

export const EmptyList: React.FC<{filesType: FileType}> = ({ filesType }) => (
    <div className="empty-list">
        <h3 className="disclaimer-title">Aucun fichier sélectionné</h3>
        <img className="drag-n-drop-icon" src={DragNDropIcon} />
        <span>Vous pouvez aussi glisser des fichiers ici</span>
        {filesType === FileType.PDF ? null : (
            <span className="files-type"><i>{`(fichiers acceptés: ${ALLOWED_IMAGE_EXTENTIONS.map(ext => ext.toLowerCase()).join(', ')})`}</i></span>
        )}
    </div>
)