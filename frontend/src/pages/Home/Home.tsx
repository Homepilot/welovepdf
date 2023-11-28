import { useCallback, useEffect } from "react";

import { FilePlus, Image, Minimize2 } from "react-feather";

import { logPageVisited } from "../../api";
import ResizeA4 from '../../assets/images/resize_A4.svg'
import { AppHeader, AppFooter } from "../../components";
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
            <div id="home-header">
                <AppHeader 
                    shouldDisplayHomeBtn={false}
                    onNavigateHome={() => onNavigate(PageName.HOME)}    
                    />
            </div>
            <div id="home-container">
                <div className='home-container'>
                    <div onClick={() => onNavigate(PageName.MERGE)} className='hp-btn'><FilePlus/> <div>Fusionner des PDF</div></div>
                    <div onClick={() => onNavigate(PageName.COMPRESS)} className='hp-btn'><Minimize2/> <div>Comprimer des PDF</div></div>
                </div>
                <div className='home-container'>
                    <div onClick={() => onNavigate(PageName.CONVERT_IMG)} className='hp-btn'><Image/> <div>Convertir des images</div></div>
                    <div onClick={() => onNavigate(PageName.RESIZE)} className='hp-btn'><img src={ResizeA4} style={{width: 32}}/> <div>Formatter en A4</div></div>
                </div>
            </div>
            <AppFooter/>
        </div>
    )
}