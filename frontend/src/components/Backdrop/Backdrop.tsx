
import Logo from '../../assets/images/logo_light.svg';
import './Backdrop.css';

export const Backdrop: React.FC<{ isVisible: boolean }> = ({ isVisible }) => !isVisible ? null : (
    <div className="backdrop">
        <img src={Logo} />
    </div>
)