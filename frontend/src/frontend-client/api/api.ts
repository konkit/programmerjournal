export * from './entry.service';
import { EntryService } from './entry.service';
export * from './note.service';
import { NoteService } from './note.service';
export * from './recurringTask.service';
import { RecurringTaskService } from './recurringTask.service';
export * from './task.service';
import { TaskService } from './task.service';
export const APIS = [EntryService, NoteService, RecurringTaskService, TaskService];
