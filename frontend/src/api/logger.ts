import {
    Info,
} from '../../wailsjs/go/models/Logger';
import {PageName} from '../types'

export async function logPageVisited(pageName: PageName){
    Info(`${pageName} page visited`, new Map<string, unknown>())
}