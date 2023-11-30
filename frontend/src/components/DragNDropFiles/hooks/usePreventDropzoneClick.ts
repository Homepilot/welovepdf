import { useEffect } from "react";


/**
 * This hooks prevents the Dropzone to open the OpenFileDialog when clicking on children elements
 * (in our case : the files list and the file cards + delete button)
 * It is done by finding the input tag among the dropzone's children.
 * The dropzone  which is found by className.
 * !! If there are several dropzones in the same page, every one should have an identifying className
 * This hook is called w/ an array to cover such cases
 * @param dropzonesClassNames the className added to the label element wrapping the dropzone input
 */
export const usePreventDropzoneClick = (dropzonesClassNames: string[]) => {
    useEffect(() => {
        dropzonesClassNames.forEach(preventClickEventDefaultChildInput)
    }, []) 
}

function preventClickEventDefaultChildInput(parentClassName: string){
    const parentElement = 
        window.document
            .getElementsByClassName(parentClassName)[0] as undefined | HTMLLabelElement;
    if(!parentElement){
      console.error(`Error finding dropzone label element w/ className="${parentClassName}"`)
      return
    }
    
    const inputElement = 
    Array.from(parentElement.children)
    .find((elem => elem.tagName.toLowerCase() === 'input')) as  HTMLInputElement | undefined;
    if(!inputElement){
        console.error(`Error finding dropzone input element in parent w/ className="${parentClassName}"`)
      return
    }
    inputElement.onclick = (ev) => { ev.preventDefault() };
}