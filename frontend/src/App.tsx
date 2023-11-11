import {useMemo, useState} from 'react';
import { Toaster } from 'react-hot-toast';
import './App.css';
import { CompressFilesPage, ConvertImagesPage, HomePage, MergeFilesPage } from './pages';
import { AppHeader } from './components';

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
            <div id="app-container">
                <AppHeader 
                    shouldDisplayHomeBtn={currentPage !== PageName.HOME}
                    onNavigateHome={() => onNavigate(PageName.HOME)}    
                    />
                { pageComponent}
            </div>
            <Toaster position="bottom-right"/>
         </div>
    )
}

export default App
