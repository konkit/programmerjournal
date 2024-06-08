'use client';

import {Task} from "@/lib/task";
import {useState} from "react";

export interface TaskComponentProps {
    task: Task;
    date: string;
    setTaskTitle: (id: string, newValue: string) => void
    setTaskDone: (id: string, newValue: boolean) => void
    setTaskDescription: (id: string, newValue: string) => void
}

export default function TaskComponent(props: TaskComponentProps) {
    const [descriptionEdited, setDescriptionEdited] = useState(false);

    const initialDescription = getTaskDescription(props.task, props.date)
    const [description, setDescription] = useState(initialDescription);

    const toggleDescriptionEdit = () => setDescriptionEdited(!descriptionEdited)

    const finishDescriptionEdit = () => {
        props.setTaskDescription(props.task.id, description)
        toggleDescriptionEdit()
    }

    let descriptionComponent
    if (!descriptionEdited) {
        descriptionComponent = <div onClick={toggleDescriptionEdit}>{description}</div>
    } else {
        descriptionComponent = <textarea onBlur={finishDescriptionEdit}
                                         autoFocus
                                         className="text-black"
                                         onChange={(e) => setDescription(e.target.value)}
                                         value={description} />
    }

    return (
        <div className="p-4 flex flex-row">
            <div>
                <input type="checkbox" className="mx-4" checked={props.task.done} onChange={(e) => props.setTaskDone(props.task.id, !props.task.done)} />
            </div>

            <div className="flex flex-col">
                <span>{props.task.title}</span>
                {descriptionComponent}
            </div>
        </div>
    )
}

function getTaskDescription(task: Task, date: string) {
    return "task description asdf"
}


