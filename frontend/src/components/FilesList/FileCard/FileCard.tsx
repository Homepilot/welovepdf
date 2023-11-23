import { FileText, Image, Trash2 } from 'react-feather'
import "./FileCard.css"


type FileCardProps = {
    onDeleteCard(): void;
    fileName: string;
}

export const FileCard: React.FC<FileCardProps> = ({ fileName, onDeleteCard }) => (
    <div className="file-card">
        <div className='file-info'>
        <span className='file-type'>{fileName.toLowerCase().endsWith('.pdf') ? <FileText /> : <Image />}</span>
        <span className='file-name'  >{formatPathString(fileName)} </span>
        </div>
        <span className='delete-file-btn' onClick={onDeleteCard} ><Trash2 className='delete-icon'/></span>
    </div>
)

function formatPathString(pathString: string) {
    let formattedPathString = pathString;
    if(formattedPathString.startsWith('/')) formattedPathString = formattedPathString.substring(1);
    return formattedPathString.split('/').join(' > ')
}