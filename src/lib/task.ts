

export interface Task {
    id: string;
    title: string;
    done: boolean;
    created_at: string;
    updated_at: string;
    updateEntries: Record<string, UpdateEntry>;
}

export interface UpdateEntry {
    date: string;
    description: string;
    doneToday: boolean;
}
