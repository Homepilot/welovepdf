import DragNDropIcon from '../../assets/images/drag_n_drop.svg'

export const EmptyList: React.FC = () => (
    <div className="no-files-disclaimer">
    <h3 className="disclaimer-title">Aucun fichier sélectionné</h3>
    <img src={DragNDropIcon} className='drag-n-drop-icon' />
    <p className="disclaimer-subtitle">Vous pouvez aussi glisser des fichiers ici</p>
</div>
)