import { openLinkInBrowser } from '../../api';
import Logo from '../../assets/images/logo_light.svg';
import './AppFooter.css'

export const AppFooter = () => (
    <div id="app-footer">
        <div 
            style={{cursor: "pointer"}} 
            onClick={() => openLinkInBrowser('https://www.homepilot.fr')}
        >
            <img src={Logo} />
            <span>by Homepilot</span>
        </div>
    </div>
)