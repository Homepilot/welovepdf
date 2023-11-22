import { FilePlus, Image, Minimize2 } from "react-feather";

import { PageName } from "../types";
import "./Home.css";

type HomePageProps = {
    onNavigate(pageName: PageName): void;
}

export const HomePage: React.FC<HomePageProps> = ({ onNavigate }) => (
    <>
    {/* <h3>Que souhaitez vous faire ?</h3> */}
    <div className='home-container'>
        <div onClick={() => onNavigate(PageName.CONVERT_IMG)} className='hp-btn'><Image/> <div>Convertir des images</div></div>
        <div onClick={() => onNavigate(PageName.MERGE)} className='hp-btn'><FilePlus/> <div>Fusionner des PDF</div></div>
        <div onClick={() => onNavigate(PageName.COMPRESS)} className='hp-btn'><Minimize2/> <div>Comprimer des PDF</div></div>
    </div>
    </>
)
