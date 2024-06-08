import {useState} from "react";
import {Task} from "@/lib/task";
import TaskComponent from "@/app/task";

interface TaskListProps {
    tasks: Task[];
    todayDate: string;
}

export default function TaskList(props: TaskListProps) {
    console.log("Render TaskList")

    const [tasks, setTasks] = useState<Task[]>(props.tasks);

    const setTaskTitle = function(id: string, newValue: string) {
        setTasks((oldTasks) => {
            return oldTasks.map((t) => {
                if (t.id !== id) {
                    return t
                }

                const newTask = cloneTask(t);
                newTask.title = newValue;
                return newTask;
            })
        })
    }

    const setTaskDone = function(id: string, newValue: boolean) {
        setTasks((oldTasks) => {
            return oldTasks.map((t) => {
                if (t.id !== id) {
                    return t
                }

                const newTask = cloneTask(t);
                newTask.done = newValue;
                return newTask;
            })
        })
    }

    const setTaskDescription = function(id: string, newValue: string) {
        setTasks((oldTasks) => {
            return oldTasks.map((t) => {
                if (t.id !== id) {
                    return t
                }

                const newTask = cloneTask(t);
                if (!newTask.updateEntries[props.todayDate]) {
                    newTask.updateEntries[props.todayDate] = {}
                }
                newTask.updateEntries[props.todayDate]!.description = newValue
                return newTask;
            })
        })
    }


    return (
        <main className="flex min-h-screen flex-col items-center p-24">
            <div>
                {JSON.stringify(tasks)}

                <h1 className="text-3xl">Monday, 2024-06-01</h1>

                <p>Tasks:</p>

                {tasks.map((task: Task) => {
                    return <div key={task.id}>
                        <TaskComponent task={task}
                                       date={props.todayDate}
                                       setTaskTitle={setTaskTitle}
                                       setTaskDone={setTaskDone}
                                       setTaskDescription={setTaskDescription} />
                    </div>
                })}

                <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
                    Add Task
                </button>
            </div>
        </main>
    );
}

function cloneTask(task: Task) {
    let result = JSON.parse(JSON.stringify(task))
    if (result.updateEntries == null) {
        result.updateEntries = {}
    }
    return result
}
