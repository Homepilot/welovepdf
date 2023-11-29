import { useCallback, useEffect } from "react";

import { FilePlus, Image, Minimize2 } from "react-feather";

import { logPageVisited } from "../../api";
import ResizeA4 from '../../assets/images/resize_A4.svg'
import { AppHeader, AppFooter, Button } from "../../components";
import { PageName } from "../../types";
import "./Home.css";

type HomePageProps = {
    onNavigate(pageName: PageName): void;
}

export const HomePage: React.FC<HomePageProps> = ({ onNavigate }) => {
    const logHomePageVisited = useCallback(async () => {
        await logPageVisited(PageName.HOME)
    }, [])
    
    useEffect(() => {
        logHomePageVisited()
    }, [])

    return(
        <div id="layout">
            <AppHeader 
                shouldDisplayHomeBtn={false}
                onNavigateHome={() => onNavigate(PageName.HOME)}    
            />
            <div id="home-body">
                <div>
                    <div className='home-btn-container'>
                        <Button clicked={() => onNavigate(PageName.MERGE)}><FilePlus/> <div>Fusionner des PDF</div></Button>
                        <Button clicked={() => onNavigate(PageName.COMPRESS)}><Minimize2/> <div>Comprimer des PDF</div></Button>
                    </div>
                    <div className='home-btn-container'>
                        <Button clicked={() => onNavigate(PageName.CONVERT_IMG)}><Image/> <div>Convertir des images</div></Button>
                        <Button clicked={() => onNavigate(PageName.RESIZE)}><img src={ResizeA4} style={{width: 32}}/> <div>Formatter en A4</div></Button>
                    </div>
                </div>
            </div>
            <AppFooter/>
        </div>
    )
}