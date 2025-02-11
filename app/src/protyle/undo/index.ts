import {onTransaction, transaction} from "../wysiwyg/transaction";
import {preventScroll} from "../scroll/preventScroll";
import {Constants} from "../../constants";
import {hideElements} from "../ui/hideElements";
import {scrollCenter} from "../../util/highlightById";

interface IOperations {
    doOperations: IOperation[],
    undoOperations: IOperation[]
}

class Undo {
    private hasUndo = false;
    private redoStack: IOperations[];
    private undoStack: IOperations[];

    constructor() {
        this.redoStack = [];
        this.undoStack = [];
    }

    public undo(protyle: IProtyle) {
        if (protyle.disabled) {
            return;
        }
        if (this.undoStack.length === 0) {
            return;
        }
        const state = this.undoStack.pop();
        this.render(protyle, state, false);
        this.hasUndo = true;
        this.redoStack.push(state);
    }


    public redo(protyle: IProtyle) {
        if (protyle.disabled) {
            return;
        }
        if (this.redoStack.length === 0) {
            return;
        }
        const state = this.redoStack.pop();
        this.render(protyle, state, true);
        this.undoStack.push(state);
    }

    private render(protyle: IProtyle, state: IOperations, redo: boolean) {
        hideElements(["hint"], protyle);
        protyle.wysiwyg.lastHTMLs = {};
        if (!redo) {
            state.undoOperations.forEach(item => {
               onTransaction(protyle, item, true);
            });
            transaction(protyle, state.undoOperations);
        } else {
            state.doOperations.forEach(item => {
                onTransaction(protyle, item, true);
            });
            transaction(protyle, state.doOperations);
        }
        preventScroll(protyle);
        scrollCenter(protyle);
    }

    public replace(doOperations: IOperation[]) {
        this.undoStack[this.undoStack.length - 1].doOperations = doOperations;
    }

    public add( doOperations: IOperation[], undoOperations: IOperation[]) {
        this.undoStack.push({undoOperations, doOperations});
        if (this.undoStack.length > Constants.SIZE_UNDO) {
            this.undoStack.shift();
        }
        if (this.hasUndo) {
            this.redoStack = [];
            this.hasUndo = false;
        }
    }

    public clear() {
        this.undoStack = [];
        this.redoStack = [];
    }
}

export {Undo};
