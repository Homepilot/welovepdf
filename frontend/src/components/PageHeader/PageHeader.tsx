import { Button } from '../Button';
import './PageHeader.css'

type PageHeaderProps = {
    headerText: string;
    actionLabel: string;
    isSelectionEmpty: boolean;
    isActionDisabled: boolean;
    onEmptyList(): void;
    onSelectFiles(): void;
    onRunAction(): void;
}

export const PageHeader: React.FC<PageHeaderProps> = ({
    actionLabel,
    headerText,
    isSelectionEmpty,
    isActionDisabled,
    onEmptyList,
    onRunAction,
    onSelectFiles,
}) => (
    <div id="page-header">
        <div id="page-header-text">
            <h3>{headerText}</h3>
        </div>
        <div id='btn-container'>
            <Button
                disabled={isSelectionEmpty}
                clicked={onEmptyList}
            >
                Vider la liste
            </Button>
            <Button
                clicked={onSelectFiles}
            >
                {`${isSelectionEmpty ? 'Choisir' : 'Ajouter'} des fichiers`}
            </Button>
            <Button
                disabled={isActionDisabled}
                clicked={onRunAction}
            >
                {actionLabel}
            </Button>
        </div>
    </div>
)