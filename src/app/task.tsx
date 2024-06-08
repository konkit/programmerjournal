import {Task} from "@/lib/task";


export default function TaskComponent(task: Task) {
    return (
        <div>
            <input type="checkbox" />
            <span>{task.title}</span>
        </div>
    )
}
