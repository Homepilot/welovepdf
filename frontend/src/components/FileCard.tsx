import { FileText, Image, Trash2 } from 'react-feather'
import "./FileCard.css"


type FileCardProps = {
    onDeleteCard(): void;
    fileName: string;
}

export const FileCard: React.FC<FileCardProps> = ({ fileName, onDeleteCard }) => (
    <div className="file-card">
        <span className='file-type'>{fileName.toLowerCase().endsWith('.pdf') ? <FileText /> : <Image />}</span>
        <span className='file-name'  >{fileName} </span>
        <span className='delete-file-btn' onClick={onDeleteCard} ><Trash2 color="red" /></span>
    </div>
)