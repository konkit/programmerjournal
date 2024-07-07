import {useEffect, useState} from "react";
import {Task} from "@/lib/task";
import TaskComponent from "@/components/task";
import {IsBeforeOrEqual, toWallDate} from "@/lib/wall_date";
import NewTaskComponent from "@/components/newTask";

interface TaskListProps {
    tasks: Task[];
    todayDate: string;
}

enum UIState {
    Idle,
    AddingTicket
}

export default function TaskList(props: TaskListProps) {
    const [tasks, setTasks] = useState<Task[]>(props.tasks);
    const [state, setState] = useState<UIState>(UIState.Idle);

    useEffect(() => {
        loadTaskList()
    }, [props.todayDate])

    const loadTaskList = async () => {
        const res = await fetch('/api/tasks/list?date=' + props.todayDate);
        const data = await res.json();
        setTasks(data.tasks);
    }

    const setTaskTitle = function(id: string, newValue: string) {
        const payload = {
            id: id,
            title: newValue,
        }
        fetch('/api/tasks/title/update', {method: "PATCH", body: JSON.stringify(payload)})
            .then(r => {
                if (r.ok) {
                    return loadTaskList()
                }
            })
    }

    const setTaskDone = function(id: string, date: string, newValue: boolean) {
        const payload = {
            id: id,
            date: date,
            done: newValue,
        }
        fetch('/api/tasks/setDone', {method: "PATCH", body: JSON.stringify(payload)})
            .then(r => {
                if (r.ok) {
                    return loadTaskList()
                }
            })
    }

    const setTaskDescription = function(id: string, date: string, newValue: string) {
        const payload = {
            id: id,
            date: date,
            description: newValue,
        }
        fetch('/api/tasks/description/update', {method: "PATCH", body: JSON.stringify(payload)})
            .then(r => {
                if (r.ok) {
                    return loadTaskList()
                }
            })
    }

    const newTask = function() {
        setState(UIState.AddingTicket)
    }

    const createTask = function(title: string, date: string) {
        const payload = {
            "title": title,
            date: date
        }

        fetch('/api/tasks/create', {method: "POST", body: JSON.stringify(payload)})
            .then(r => {
                if (r.ok) {
                    return loadTaskList()
                }
            })

    }

    const deleteTask = function(taskId: string) {
        fetch('/api/tasks/delete/' + taskId, {method: "DELETE"})
            .then(r => {
                if (r.ok) {
                    return loadTaskList()
                }
            })
    }

    return (
        <main className="flex min-h-screen flex-col items-center p-24">
            <div>
                {tasks.map((task: Task) => {
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

                {state == UIState.AddingTicket ?
                    <NewTaskComponent date={props.todayDate} createTask={createTask} />
                    :
                    <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                            onClick={newTask}>
                        Add Task
                    </button>
                }

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
