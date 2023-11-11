import { PageName } from "../App";

type HomePageProps = {
    onNavigate(pageName: PageName): void;
}

export const HomePage: React.FC<HomePageProps> = ({ onNavigate }) => (
    <>
    {/* <h3>Que souhaitez vous faire ?</h3> */}
    <div className='flex-container'>
        <div onClick={() => onNavigate(PageName.CONVERT_IMG)} className='box'>Convertir des images en PDF</div>
        <div onClick={() => onNavigate(PageName.MERGE)} className='box'>Fusionner des PDF</div>
        <div onClick={() => onNavigate(PageName.COMPRESS)} className='box'>Comprimer des PDF</div>
    </div>
    </>
)
