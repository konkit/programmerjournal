import Image from "next/image";
import {Task} from "@/lib/task";
import TaskComponent from "@/app/task";

export default function Home() {

    let tasks: Task[] = [
        {
            id: "1",
            title: "Task 1",
            done: false,
            created_at: "2024-06-06T00:00:00",
            updated_at: "2024-06-06T00:00:00",
            updateEntries: new Map(),
        }
    ]

  return (
    <main className="flex min-h-screen flex-col items-center p-24">
        <div>
            <h1 className="text-3xl">Monday, 2024-06-01</h1>

            <p>Tasks:</p>

            {tasks.map((task: Task) => {
                return <div key={task.id}>{TaskComponent(task)}</div>
            })}

            <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Add Task</button>
        </div>
    </main>
  );
}
