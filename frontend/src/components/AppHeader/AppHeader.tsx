import { Heart } from "react-feather";
import './AppHeader.css';

type AppHeaderProps = {
    shouldDisplayHomeBtn: boolean;
    onNavigateHome(): void;
}

export const AppHeader: React.FC<AppHeaderProps> = ({ shouldDisplayHomeBtn, onNavigateHome }) => (
    <div className="app-header">
        {
            shouldDisplayHomeBtn
                ? <span className="hp-btn home-btn" onClick={onNavigateHome}>Accueil</span> 
                : null
        }
        <div className="header-container">
            <h1>We</h1>
            <Heart size={40} color="#de2410"/> 
            <h1>PDF</h1>
        </div>
    </div>
)