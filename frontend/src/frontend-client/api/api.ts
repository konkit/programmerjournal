export * from './entry.service';
import { EntryService } from './entry.service';
export * from './recurringTask.service';
import { RecurringTaskService } from './recurringTask.service';
export * from './task.service';
import { TaskService } from './task.service';
export const APIS = [EntryService, RecurringTaskService, TaskService];
