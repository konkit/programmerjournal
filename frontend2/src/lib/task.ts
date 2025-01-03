import {WallDate} from "./wall_date";

export interface Task {
    id: number;
    taskID: string;
    title: string;
    status: string;
    description: string;
    createdDate: WallDate;
    update: string;
}
