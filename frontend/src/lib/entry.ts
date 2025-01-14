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
