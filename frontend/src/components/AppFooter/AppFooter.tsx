import { logHomepilotLinkClicked, openLinkInBrowser } from '../../api';
import Logo from '../../assets/images/logo_light.svg';
import './AppFooter.css'

export const AppFooter = () => (
    <div id="app-footer">
        <div 
            style={{cursor: "pointer"}} 
            onClick={async () => {
                openLinkInBrowser('https://www.homepilot.fr');
                await logHomepilotLinkClicked();
            }}
        >
            <img src={Logo} />
            <span>by Homepilot</span>
        </div>
    </div>
)