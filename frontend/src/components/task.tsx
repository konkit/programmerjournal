'use client';

import {Task} from "@/lib/task";
import {useEffect, useState} from "react";
import Button from "@/components/button";

export interface TaskComponentProps {
    task: Task;
    date: string;
    setTaskTitle: (id: number, newValue: string) => void
    setTaskDone: (id: number, date: string, task: Task) => void
    setTaskDescription: (id: number, date: string, newValue: string) => void
    snoozeTask: (id: number) => void
    deleteTask: (id: number) => void
}

export default function TaskComponent(props: TaskComponentProps) {
    const [titleEdited, setTitleEdited] = useState(false);
    const [descriptionEdited, setDescriptionEdited] = useState(false);

    const [title, setTitle] = useState(props.task.title);
    const [description, setDescription] = useState(props.task.update);

    useEffect(() => {
        setTitle(props.task.title)
        setDescription(props.task.update)
    }, [props.task])

    const toggleTitleEdit = () => setTitleEdited(!titleEdited)
    const toggleDescriptionEdit = () => setDescriptionEdited(!descriptionEdited)

    const finishDescriptionEdit = () => {
        props.setTaskDescription(props.task.id, props.date, description)
        toggleDescriptionEdit()
    }

    const finishTitleEdit = () => {
        props.setTaskTitle(props.task.id, title)
        toggleTitleEdit()
    }

    const deleteTask = () => {
        props.deleteTask(props.task.id)
    }

    let descriptionComponent
    if (!descriptionEdited) {
        descriptionComponent = <div onClick={toggleDescriptionEdit}>{description || "Enter update for today" }</div>
    } else {
        descriptionComponent = <textarea onBlur={finishDescriptionEdit}
                                         placeholder="Enter update for today"
                                         autoFocus
                                         onChange={(e) => setDescription(e.target.value)}
                                         value={description} />
    }

    let titleComponent
    if (!titleEdited) {
        titleComponent = <div onClick={toggleTitleEdit}>{props.task.title}</div>
    } else {
        titleComponent = <input onChange={(e) => setTitle(e.target.value)}
                                value={title}
                                autoFocus
                                onBlur={finishTitleEdit} />
    }

    return (
        <div className="p-4 flex flex-row">
            <div>
                <input type="checkbox" className="mx-4"
                       checked={props.task.status == "Done"}
                       onChange={(e) => props.setTaskDone(props.task.id, props.date, props.task)}
                />
            </div>

            <div className="flex flex-col">
                {titleComponent}
                {descriptionComponent}
            </div>

            <div>
                <Button onClick={() => props.snoozeTask(props.task.id)}>Snooze</Button>
                <Button onClick={deleteTask}>Delete</Button>
                {/*<Button>Mark as Done</Button>*/}
                {/*<Button>Mark as Irrelevant</Button>*/}
            </div>
        </div>
    )
}
