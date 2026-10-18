export enum EntryStatus {
    TaskCreated = "TaskCreated",
    TaskDone = "TaskDone",
    TaskSnoozed = "TaskSnoozed",
    TaskMigrated = "TaskMigrated",
    TaskCancelled = "TaskCancelled",
    Note = "Note",
}

export function isTask(status: EntryStatus) {
    const taskStatuses = [
        EntryStatus.TaskCreated,
        EntryStatus.TaskDone,
        EntryStatus.TaskSnoozed,
        EntryStatus.TaskMigrated,
        EntryStatus.TaskCancelled
    ]

    return taskStatuses.includes(status)
}

export function isNote(status: EntryStatus) {
    const taskStatuses = [EntryStatus.Note]
    return taskStatuses.includes(status)
}

export interface DateGroup<T> {
    date: string
    entries: T[]
}

// Groups consecutive entries with the same createdDate, keeping their order
export function groupByDate<T extends { createdDate: string }>(entries: T[]): DateGroup<T>[] {
    const groups: DateGroup<T>[] = []
    for (const entry of entries) {
        const last = groups[groups.length - 1]
        if (last && last.date === entry.createdDate) {
            last.entries.push(entry)
        } else {
            groups.push({date: entry.createdDate, entries: [entry]})
        }
    }
    return groups
}
