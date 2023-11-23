import {useMemo, useState} from 'react';

import { Toaster } from 'react-hot-toast';

import './App.css';
import Logo from './assets/images/logo.svg';
import { AppHeader } from './components';
import { CompressFilesPage, ConvertImagesPage, HomePage, MergeFilesPage } from './pages';
import { PageName } from './types';


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
            <div id="app-header">
                <AppHeader 
                    shouldDisplayHomeBtn={currentPage !== PageName.HOME}
                    onNavigateHome={() => onNavigate(PageName.HOME)}    
                    />
            </div>
            <div id="app-container">
                { pageComponent}
            </div>
            <div id="app-footer">
                <div><img src={Logo} /><span>by Homepilot</span> </div>
            </div>
            <Toaster position="bottom-right"/>
         </div>
    )
}

export default App
