import { useCallback } from "react";

import { DragDropContext, Droppable, Draggable, OnDragEndResponder } from "react-beautiful-dnd";

import { FileInfo, FileType } from '../../types';

import { EmptyList } from "./EmptyList";
import { FileCard } from './FileCard';

import './FilesList.css';

type FilesListProps = {
    onSelectionReordered(newSelection: FileInfo[]): void;
    onRemoveFileFromList(fileId: string): void;
    selectedFiles: FileInfo[];
    filesType?: FileType;
    selectFilesPrompt: string;
}

export const FilesList: React.FC<React.PropsWithChildren<FilesListProps>> = ({ onSelectionReordered, onRemoveFileFromList, selectedFiles }) => {

    const onDragEnd: OnDragEndResponder = useCallback((result) => {
        if (!result.destination) {
          return;
        }
    
        const reorderedItems = reorder(
          selectedFiles,
          result.source.index,
          result.destination.index
        );
        onSelectionReordered(reorderedItems)
    }, [onSelectionReordered, selectedFiles])

    return (
        <div id='files-list'>
            {
                !selectedFiles.length ?  
                <EmptyList/> : 
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
                                        <FileCard fileName={item.name} onDeleteCard={() => onRemoveFileFromList(item.id)} />
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
function reorder (list: FileInfo[], startIndex: number, endIndex: number) {
  const result = Array.from(list);
  const [removed] = result.splice(startIndex, 1);
  result.splice(endIndex, 0, removed);

  return result;
}


// @ts-expect-error unknown object type
function getItemStyle (isDragging: boolean, draggableStyle) {
    const grid = 8;
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
