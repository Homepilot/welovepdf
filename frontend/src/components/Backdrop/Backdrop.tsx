
import './Backdrop.css';

export const Backdrop: React.FC<{ isVisible: boolean }> = ({ isVisible }) => !isVisible ? null : (
    <div className="backdrop">
        <div className="backdrop-container">

        <span className="loader"></span>
            <span className='disclaimer' >Opération en cours, veuilez patienter</span>
        </div>
    </div>
)