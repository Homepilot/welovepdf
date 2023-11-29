import { Heart } from "react-feather";

import './AppHeader.css';
import { Button } from "../Button";

type AppHeaderProps = {
    shouldDisplayHomeBtn: boolean;
    onNavigateHome(): void;
}

export const AppHeader: React.FC<AppHeaderProps> = ({ shouldDisplayHomeBtn, onNavigateHome }) => (
    <>
        {
            shouldDisplayHomeBtn
                ? <Button extraClasses={['home-btn']} clicked={onNavigateHome}>Accueil</Button>
                : null
        }
        <div id="app-header">
            <h1>We</h1>
            <Heart size={40} color="#de2410"/> 
            <h1>PDF</h1>
        </div>
    </>
)