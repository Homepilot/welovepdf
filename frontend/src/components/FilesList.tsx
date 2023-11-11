import { useState } from 'react';
import { DragDropContext, Droppable, Draggable, OnDragEndResponder } from "react-beautiful-dnd";
import { FileCard } from './FileCard';
import './FilesList.css';
import { FileType } from '../types';
import { selectMultipleFiles } from '../actions';

export const FilesList: React.FC<React.PropsWithChildren<{onSelectionUpdated(filePathes: string[]): void, filesType?: FileType}>> = ({ onSelectionUpdated, filesType, children }) => {
    const [selectedFiles, setSelectedFiles] = useState<{ id: string}[]>([]);

    const selectFiles = async () => {
        const files = await selectMultipleFiles(filesType);
        const newSelection = Array.from(new Set([...selectedFiles.map(({id}) => id), ...files]));
        const selectionWithIds = newSelection.map(id => ({id}))
        setSelectedFiles(selectionWithIds);
        onSelectionUpdated(newSelection);
    }

    const removeFileFromList = (fileId: string) => {
        const newSelectionWithIds = selectedFiles.filter(({id}) => id !== fileId);
        setSelectedFiles(newSelectionWithIds);
        onSelectionUpdated(newSelectionWithIds.map(({id}) => id))
    } 

    const onDragEnd: OnDragEndResponder = (result) => {
        console.log('ON DRAG END')
        // dropped outside the list
        if (!result.destination) {
          return;
        }
    
        const reorderedItems = reorder(
          selectedFiles,
          result.source.index,
          result.destination.index
        );
    
        setSelectedFiles(reorderedItems)
      }

    return (
        
        <div className='files-list'>
            <div className='btn-container'>
                <button onClick={selectFiles} className="btn">Choisir des fichiers</button>
                {children}
            </div>
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
                                        <FileCard fileName={item.id} onDeleteCard={() => removeFileFromList(item.id)} />
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
};

const grid = 8;

// @ts-expect-error unknown object type
function getItemStyle (isDragging: boolean, draggableStyle) {
    return {
        // some basic styles to make the items look a bit nicer
        userSelect: "none",
        padding: grid * 2,
        margin: `0 0 ${grid}px 0`,
        
        // change background colour if dragging
        background: isDragging ? "lightgreen" : "grey",

        // styles we need to apply on draggables
        ...draggableStyle
    }
};

function getListStyle(isDraggingOver: boolean){
    return {
        background: isDraggingOver ? "lightblue" : "lightgrey",
        padding: grid,
        width: '99%',
        margin: '2rem auto',
    }
};
