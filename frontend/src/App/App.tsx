import {useMemo, useState} from 'react';

import { Toaster } from 'react-hot-toast';

import './App.css';
import { CompressFilesPage, ConvertImagesPage, HomePage, MergeFilesPage, ResizeFilesPage } from '../pages';
import { PageName } from '../types';


export function App() {
    const [currentPage, setCurrentPage] = useState<PageName>(PageName.HOME);

    function onNavigate(pageName: PageName){
        setCurrentPage(pageName);
    }

    const pageComponent = useMemo(() => {
        switch(currentPage){
            case PageName.COMPRESS: return <CompressFilesPage onNavigateHome={() => onNavigate(PageName.HOME)} />;
            case PageName.CONVERT_IMG: return <ConvertImagesPage onNavigateHome={() => onNavigate(PageName.HOME)} />;
            case PageName.MERGE: return <MergeFilesPage onNavigateHome={() => onNavigate(PageName.HOME)} />;
            case PageName.RESIZE: return <ResizeFilesPage onNavigateHome={() => onNavigate(PageName.HOME)} />;
            default: return <HomePage onNavigate={onNavigate}/>;
        }
    }, [currentPage]);

    return (
        <>
            { pageComponent}
            <Toaster position="bottom-right"/>
         </>
    )
}
