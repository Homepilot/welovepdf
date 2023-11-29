import { logHomepilotLinkClicked, openLinkInBrowser } from '../../api';
import Logo from '../../assets/images/logo_light.svg';
import './AppFooter.css'

export const AppFooter = () => (
    <div id="app-footer">
        <div 
            className='app-footer-container'
            style={{cursor: "pointer"}} 
            onClick={async () => {
                openLinkInBrowser('https://www.homepilot.fr');
                await logHomepilotLinkClicked();
            }}
        >
            <img id="footer-logo" src={Logo} />
            <span>by Homepilot</span>
        </div>
    </div>
)