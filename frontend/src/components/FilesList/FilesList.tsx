import { DragDropContext, Droppable, Draggable, OnDragEndResponder } from "react-beautiful-dnd";

import { FileType } from '../../types';

import { FileCard } from './FileCard';

import './FilesList.css';

type FileInfo = {
    id: string;
}

type FilesListProps = {
    onSelectionUpdated(newSelection: FileInfo[]): void;
    onRemoveFileFromList(fileId: string): void;
    selectedFiles: FileInfo[];
    filesType?: FileType;
    selectFilesPrompt: string;
}

export const FilesList: React.FC<React.PropsWithChildren<FilesListProps>> = ({ onSelectionUpdated, onRemoveFileFromList, selectedFiles }) => {

    const onDragEnd: OnDragEndResponder = (result) => {
        if (!result.destination) {
          return;
        }
    
        const reorderedItems = reorder(
          selectedFiles,
          result.source.index,
          result.destination.index
        );
    
        onSelectionUpdated(reorderedItems)
      }

    return (
        
        <div className='files-list'>

            {
                !selectedFiles.length ?  
                <h3>Aucun fichier sélectionné</h3> : 
                (
                    <DragDropContext onDragEnd={onDragEnd}>
                        <Droppable droppableId="droppable">
                        {(provided, snapshot) => (
                            <div
                            {...provided.droppableProps}
                            ref={provided.innerRef}
                            style={getListStyle(snapshot.isDraggingOver)}
                            >
                            {selectedFiles.map((item, index) => (
                                <Draggable key={item.id} draggableId={item.id} index={index}>
                                {(provided, snapshot) => (
                                    <div
                                        ref={provided.innerRef}
                                        {...provided.draggableProps}
                                        {...provided.dragHandleProps}
                                        style={getItemStyle(
                                            snapshot.isDragging,
                                            provided.draggableProps.style
                                        )}
                                    >
                                        <FileCard fileName={item.id} onDeleteCard={() => onRemoveFileFromList(item.id)} />
                                    </div>
                                )}
                                </Draggable>
                            ))}
                            {provided.placeholder}
                            </div>
                        )}
                        </Droppable>
                    </DragDropContext>
                )
            }
        </div>
    )
}

// a little function to help us with reordering the result
function reorder (list: {id: string}[], startIndex: number, endIndex: number) {
    console.log('REORDER')
  const result = Array.from(list);
  const [removed] = result.splice(startIndex, 1);
  result.splice(endIndex, 0, removed);

  return result;
}

const grid = 8;

// @ts-expect-error unknown object type
function getItemStyle (isDragging: boolean, draggableStyle) {
    return {
        // some basic styles to make the items look a bit nicer
        userSelect: "none",
        padding: grid * 2,
        margin: `0 0 ${grid}px 0`,
        
        // change background colour if dragging
        background: isDragging ? "#9c9c9c" : "#282824",
        border: "solid 1px #dbdbda",
        borderRadius: '2px',
        // styles we need to apply on draggables
        ...draggableStyle
    }
}

function getListStyle(isDraggingOver: boolean){
    return {
        background: isDraggingOver ? "#484848" : "#282824",
        padding: '0.5rem 0',
        width: '99%',
        margin: '0 auto',
    }
}