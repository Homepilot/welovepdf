import { Heart } from "react-feather";
import './AppHeader.css';

type AppHeaderProps = {
    shouldDisplayHomeBtn: boolean;
    onNavigateHome(): void;
}

export const AppHeader: React.FC<AppHeaderProps> = ({ shouldDisplayHomeBtn, onNavigateHome }) => (
    <>
        {
            shouldDisplayHomeBtn
                ? <span id="home-btn" className="hp-btn" onClick={onNavigateHome}>Accueil</span> 
                : null
        }
        <div id="app-header">
            <h1>We</h1>
            <Heart size={40} color="#de2410"/> 
            <h1>PDF</h1>
        </div>
    </>
)