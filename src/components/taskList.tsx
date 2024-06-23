import {useState} from "react";
import {Task} from "@/lib/task";
import TaskComponent from "@/components/task";
import {IsBeforeOrEqual} from "@/lib/wall_date";

interface TaskListProps {
    tasks: Task[];
    todayDate: string;
}

export default function TaskList(props: TaskListProps) {
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

    const newTaskPrompt = function() {
        setTasks((oldTasks) => {
            let newTask: Task = {
                id: crypto.randomUUID(),
                title: "New Task",
                done: false,
                created_at: "2024-06-06T00:00:00",
                updated_at: "2024-06-06T00:00:00",
                updateEntries: {},
            }
            return [...oldTasks, newTask]
        })
    }

    const deleteTask = function(taskId: string) {
        setTasks((oldTasks) => {
            return [...oldTasks.filter((t) => t.id !== taskId)]
        })
    }

    return (
        <main className="flex min-h-screen flex-col items-center p-24">
            <div>
                {tasks.filter(t => IsBeforeOrEqual(t.created_at, props.todayDate)).map((task: Task) => {
                    return <div key={task.id}>
                        <TaskComponent task={task}
                                       date={props.todayDate}
                                       setTaskTitle={setTaskTitle}
                                       setTaskDone={setTaskDone}
                                       setTaskDescription={setTaskDescription}
                                       deleteTask={deleteTask}
                        />
                    </div>
                })}

                <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                        onClick={newTaskPrompt}>
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
