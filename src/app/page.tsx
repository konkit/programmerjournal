'use client';

import {Task, TaskID, UpdateEntry} from "@/lib/task";
import TaskList from "@/components/taskList";
import Button from "@/components/button";
import {AddDay, DayOfWeek, Today, WallDate} from "@/lib/wall_date";
import {useState} from "react";

export default function Home() {

    const [todayDate, setTodayDate] = useState(Today())

    const initTaskArr = initTasks()

    const changeDateForward = () => {
        setTodayDate((oldToday) => AddDay(oldToday, 1))
    }

    const changeDateBackward = () => {
        setTodayDate((oldToday) => AddDay(oldToday, -1));
    }

    return (
    <main className="flex min-h-screen flex-col items-center p-24">
        <div>
            <h1 className="text-3xl">
                <Button onClick={changeDateBackward}>&lt;&lt;</Button>
                {DayOfWeek(todayDate)}, {todayDate}
                <Button onClick={changeDateForward}>&gt;&gt;</Button>
            </h1>

            <TaskList tasks={initTaskArr} todayDate={todayDate}/>
        </div>
    </main>
    );
}

function initTasks() {
    const todayDate = "2024-06-06"

    const updateEntries: Record<string, UpdateEntry> = {}
    updateEntries[todayDate] = {
        date: todayDate as WallDate,
        description: "",
        doneToday: false
    }
    let taskId: TaskID = "1"
    let tasks: Task[] = [
        {
            id: taskId,
            title: "Task 1",
            done: false,
            created_at: todayDate as WallDate,
            updated_at: todayDate as WallDate,
            updateEntries: updateEntries,
        }
    ]

    return tasks
}
