import DragNDropIcon from '../../../assets/images/drag_n_drop.svg';
import './EmptyList.css';

export const EmptyList: React.FC = () => (
    <div className="empty-list">
        <h3 className="disclaimer-title">Aucun fichier sélectionné</h3>
        <img className="drag-n-drop-icon" src={DragNDropIcon} />
        <span>Vous pouvez aussi glisser des fichiers ici</span>
    </div>
)