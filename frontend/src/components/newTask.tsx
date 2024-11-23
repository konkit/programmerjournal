'use client';

import {useState} from "react";
import Button from "@/components/button";

export interface NewTaskComponentProps {
    date: string;
    createTask: (title: string, date: string) => void;
}

export default function NewTaskComponent(props: NewTaskComponentProps) {
    const [title, setTitle] = useState("");

    const createTask = () => {
        props.createTask(title, props.date)
    }

    return (
        <div className="p-4 flex flex-row">
            <div className="flex flex-col">
                <input onChange={(e) => setTitle(e.target.value)}
                       value={title}
                       autoFocus />
            </div>

            <div>
                <Button onClick={createTask}>Submit</Button>
            </div>
        </div>
    )
}


