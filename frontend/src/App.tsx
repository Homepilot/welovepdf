import {useMemo, useState} from 'react';
import './App.css';
import { CompressFilesPage, ConvertImagesPage, HomePage, MergeFilesPage } from './pages';

export enum PageName {
    HOME = 'HOME',
    MERGE = 'MERGE',
    CONVERT_IMG = 'CONVERT_IMG',
    COMPRESS = 'COMPRESS',
}

function App() {
    const [currentPage, setCurrentPage] = useState<PageName>(PageName.HOME);

    function onNavigate(pageName: PageName){
        setCurrentPage(pageName);
    }

    const pageComponent = useMemo(() => {
        switch(currentPage){
            case PageName.COMPRESS: return <CompressFilesPage />;
            case PageName.CONVERT_IMG: return <ConvertImagesPage />;
            case PageName.MERGE: return <MergeFilesPage />;
            default: return <HomePage onNavigate={onNavigate}/>;
        }
    }, [currentPage]);

    return (
        <div id="app">
            { currentPage !== PageName.HOME && <span className="box home-btn" onClick={() => onNavigate(PageName.HOME)} >Accueil</span> }
            <h1>Homepilot ❤ PDF</h1>
            { pageComponent}
         </div>
    )
}

export default App
