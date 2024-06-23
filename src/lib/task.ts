import {WallDate} from "@/lib/wall_date";

export type TaskID = string;

export interface Task {
    id: TaskID;
    title: string;
    done: boolean;
    created_at: WallDate;
    updated_at: WallDate;
    snoozed_until?: WallDate;
    updateEntries: Record<WallDate, UpdateEntry>;
}

export interface UpdateEntry {
    date: WallDate;
    description: string;
    doneToday: boolean;
}
