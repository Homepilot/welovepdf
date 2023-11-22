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
        <div onClick={() => onNavigate(PageName.CONVERT_IMG)} className='box'><Image/> <div>Convertir des images</div></div>
        <div onClick={() => onNavigate(PageName.MERGE)} className='box'><FilePlus/> <div>Fusionner des PDF</div></div>
        <div onClick={() => onNavigate(PageName.COMPRESS)} className='box'><Minimize2/> <div>Comprimer des PDF</div></div>
    </div>
    </>
)
