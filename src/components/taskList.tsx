import {useEffect, useState} from "react";
import {Task} from "@/lib/task";
import TaskComponent from "@/components/task";
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

    const [snoozeVisible, setSnoozeVisible] = useState(false);
    const [snoozedTaskId, setSnoozedTaskId] = useState<number | null>(null)
    const [snoozedTaskDate, setSnoozedTaskDate] = useState<string>("")

    useEffect(() => {
        loadTaskList()
    }, [props.todayDate])

    const loadTaskList = async () => {
        const res = await fetch('/api/tasks/list/' + props.todayDate);
        const data = await res.json();
        setTasks(data);
    }

    const setTaskTitle = function(id: number, newValue: string) {
        const payload = {
            title: newValue,
        }
        fetch(`/api/tasks/${id}/setTitle`, {method: "PATCH", body: JSON.stringify(payload)})
            .then(r => {
                if (r.ok) {
                    return loadTaskList()
                }
            })
    }

    const setTaskDone = function(id: number, date: string, task: Task) {
        let currentValue = task.status == "Done"
        let newValue = !currentValue

        const payload = {
            done: newValue,
        }
        fetch(`/api/tasks/${id}/setDone`, {method: "PATCH", body: JSON.stringify(payload)})
            .then(r => {
                if (r.ok) {
                    return loadTaskList()
                }
            })
    }

    const setTaskDescription = function(id: number, date: string, newValue: string) {
        const payload = {
            update: newValue,
        }
        fetch(`/api/tasks/${id}/setUpdate`, {method: "PATCH", body: JSON.stringify(payload)})
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
            "createdDate": date,
        }

        fetch('/api/tasks/create', {method: "POST", body: JSON.stringify(payload)})
            .then(r => {
                if (r.ok) {
                    setState(UIState.Idle)
                    return loadTaskList()
                }
            })
    }

    const deleteTask = function(taskId: number) {
        fetch(`/api/tasks/${taskId}/delete/`, {method: "DELETE"})
            .then(r => {
                if (r.ok) {
                    return loadTaskList()
                }
            })
    }

    const showSnoozeModal = (taskID: number) => {
        setSnoozeVisible(true)
        setSnoozedTaskId(taskID)
    }

    const submitSnooze = (taskId: number | null, date: string) => {
        if (taskId == null) {
            console.error("Cannot snooze task, taskID is null")
            setSnoozeVisible(false)
            return;
        }
        const payload = {
            date: date,
        }
        fetch(`/api/tasks/${taskId}/snooze`, {method: "PATCH", body: JSON.stringify(payload)})
            .then(r => {
                setSnoozeVisible(false)
                if (r.ok) {
                    return loadTaskList()
                }
            })
    }

    const cancelSnooze = () => {
        setSnoozeVisible(false)
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
                                       snoozeTask={showSnoozeModal}
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

            {snoozeVisible ? <div>
                <div id="default-modal" tabIndex={-1} className="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full">
                    <div className="relative p-4 w-full max-w-2xl max-h-full">
                        <div className="relative bg-white rounded-lg shadow dark:bg-gray-700">
                            <div className="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:border-gray-600">
                                <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                                    Snooze
                                </h3>
                            </div>
                            <div className="p-4 md:p-5 space-y-4">
                                <input placeholder="2024-04-04" onChange={(e) => setSnoozedTaskDate(e.target.value)}/>
                            </div>
                            <div className="flex items-center p-4 md:p-5 border-t border-gray-200 rounded-b dark:border-gray-600">
                                <button onClick={cancelSnooze} className="py-2.5 px-5 ms-3 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700">
                                    Cancel
                                </button>
                                <button onClick={() => submitSnooze(snoozedTaskId, snoozedTaskDate)} type="button" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
                                    Snooze
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

            </div> : <div></div>
            }
        </main>
    );
}
