// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {utils} from '../models';
import {embed} from '../models';
import {models} from '../models';

export function CompressFile(arg1:string,arg2:number):Promise<boolean>;

export function ConvertImageToPdf(arg1:string,arg2:boolean):Promise<boolean>;

export function CreateTempFilesFromUpload(arg1:Array<number>):Promise<string>;

export function Init(arg1:utils.CustomLogger,arg2:embed.FS):Promise<models.PdfService>;

export function MergePdfFiles(arg1:string,arg2:Array<string>,arg3:boolean):Promise<boolean>;

export function ResizePdfFileToA4(arg1:string):Promise<boolean>;
