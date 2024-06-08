'use client';

import {Task, UpdateEntry} from "@/lib/task";
import TaskList from "@/app/taskList";

export default function Home() {

    let todayDate = "2024-06-06";

    const initTaskArr = initTasks(todayDate)

    return (
    <main className="flex min-h-screen flex-col items-center p-24">
        <div>
            <TaskList tasks={initTaskArr}  />
        </div>
    </main>
  );
}

function initTasks(todayDate: string) {
    const updateEntries: Record<string, UpdateEntry> = {}
    updateEntries[todayDate] = {
        date: todayDate,
        description: "",
        doneToday: false
    }
    let tasks: Task[] = [
        {
            id: "1",
            title: "Task 1",
            done: false,
            created_at: "2024-06-06T00:00:00",
            updated_at: "2024-06-06T00:00:00",
            updateEntries: updateEntries,
        }
    ]

    return tasks
}
