import {WallDate} from "@/lib/wall_date";

export interface Task {
    id: number;
    taskID: string;
    title: string;
    status: string;
    createdDate: WallDate;
    update: string;
}
