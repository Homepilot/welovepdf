import { PageName } from "../App";

type HomePageProps = {
    onNavigate(pageName: PageName): void;
}

export const HomePage: React.FC<HomePageProps> = ({ onNavigate }) => (
    <>
    <h3>Que souhaitez vous faire ?</h3>
    <div className='flex-container'>
        <div onClick={() => onNavigate(PageName.CONVERT_IMG)} className='box'>Convertir des images en PDF</div>
        <div onClick={() => onNavigate(PageName.MERGE)} className='box'>Fusionner des PDF</div>
        <div onClick={() => onNavigate(PageName.COMPRESS)} className='box'>Comprimer un PDF</div>
        <div onClick={() => onNavigate(PageName.SPLIT)} className='box'>Diviser un PDF en plusieurs</div>
    </div>
    </>
)
