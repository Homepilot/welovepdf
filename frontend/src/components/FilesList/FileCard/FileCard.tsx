import { FileText, Image, Trash2 } from 'react-feather'
import "./FileCard.css"


type FileCardProps = {
    onDeleteCard(): void;
    fileName: string;
    filePath: string;
}

export const FileCard: React.FC<FileCardProps> = ({ fileName, filePath, onDeleteCard }) => (
    <div className="file-card">
        <div className="file-info-container">
            <div className="file-info">
                <div className='file-name'>
                    {fileName.toLowerCase().endsWith('.pdf') ? <FileText /> : <Image />}
                    {fileName}
                </div>
                <div className='file-path'>{formatPathString(filePath)} </div>
            </div>
        </div>
        <div className='delete-file-btn' onClick={onDeleteCard} ><Trash2 className='delete-icon'/></div>
    </div>
)

function formatPathString(pathString: string) {
    const formattedPathString = pathString.startsWith('/') ? pathString.substring(1) :  pathString;
    const splittedBySlash = formattedPathString.split('/');
    return splittedBySlash.slice(0, splittedBySlash.length - 1).join(' > ');
}