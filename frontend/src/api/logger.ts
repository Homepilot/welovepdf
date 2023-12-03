import toast from 'react-hot-toast';

import * as Logger from '../../wailsjs/go/models/FrontendLogger';
import {PageName} from '../types'

export async function logHomepilotLinkClicked(){
    await Logger.Info('Homepilot link clicked', "")
}

export async function logPageVisited(pageName: PageName){
    await Logger.Info(`${pageName} page visited`, "")
}

export async function logOperationStarted(pageName: PageName, batchId: string){
    await Logger.Info(`${pageName} : operation started`, JSON.stringify({ batchId}))
}

export async function logOperationCanceledByUser(pageName: PageName, batchId: string){
    await Logger.Info(`${pageName} : operation canceled by user`, JSON.stringify({ batchId}))
}

export async function notifyAndLogOperationsResult(operationName: string, batchId: string, {successes, failures}: { successes: number, failures: number }){
    if(!successes && !failures){
        console.error('No results to log', { successes, failures })
        await Logger.Warn(`${operationName} : no results to log`, JSON.stringify({ batchId }))
        return
    }

    if(!failures){
        toast.success('Opération réussie pour tous les fichiers');
        await Logger.Info(
            `${operationName} : operation succeeded for all files`,
            JSON.stringify({ successes, failures, batchId}))
        return
    }

    if(!successes){
        toast.error("L'opération a échoué pour tous les fichiers");
        await Logger.Error(
            `${operationName} : operation failed for all files`,
            JSON.stringify({successes, failures, batchId}))
            return
        }
        
        toast.success(`L'opération a réussi pour ${successes} fichiers`);
        toast.error(`L'opération a échoué pour ${failures} fichiers`);
        await Logger.Warn(
            `${operationName} : operation failed for some files`,
            JSON.stringify({successes, failures, batchId}))
    return   
}
