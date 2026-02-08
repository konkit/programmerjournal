import {inject, Injectable, signal} from '@angular/core';
import {Entry, EntryService, TaskService} from '../../frontend-client';
import {EMPTY, forkJoin, switchMap, tap} from 'rxjs';
import {addDay, addMonth, addWeek, ThisWeek, Today} from '../../lib/wall_date';
import {NoteService} from '../../frontend-client/api/note.service';
import {SnoozeMonthEntryDialogComponent} from '../components/snooze-month-dialog/snooze-month-entry-dialog.component';
import {SnoozeDayEntryDialogComponent} from '../components/snooze-day-dialog/snooze-day-entry-dialog.component';
import {MatDialog} from '@angular/material/dialog';
import {
  MigrateToDayEntryDialogComponent
} from '../components/migrate-to-day-dialog/migrate-to-day-entry-dialog.component';
import {MatSnackBar} from '@angular/material/snack-bar';
import {Router} from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class WeeklyEntryListService {

  // Used as the "View Date" (Week Date)
  todayDate = signal<string>(ThisWeek());

  // Used for the "Today" column
  realTodayDate = signal<string>(Today());

  dailyEntryList = signal<Entry[]>([]);
  weeklyEntryList = signal<Entry[]>([]);
  pendingImportsCount = signal<number>(0);

  readonly dialog = inject(MatDialog);
  private _snackBar = inject(MatSnackBar);

  constructor(private entryService: EntryService,
    private noteService: NoteService,
    private router: Router,
    private taskService: TaskService) {
  }

  dateForward() {
    if (isDayDate(this.todayDate())) {
      this.todayDate.update((oldVal) => addDay(oldVal, 1))
      this.router.navigate(['/day', this.todayDate()]);
    } else if (isMonthDate(this.todayDate())) {
      this.todayDate.update((oldVal) => addMonth(oldVal, 1))
      this.router.navigate(['/month', this.todayDate()]);
    } else if (isWeekDate(this.todayDate())) {
      this.todayDate.update((oldVal) => addWeek(oldVal, 1))
      this.router.navigate(['/week', this.todayDate()]);
    } else {
      console.error("Unrecognized date format", this.todayDate())
    }
    return this.refreshTasks()
  }

  dateBackward() {
    if (isDayDate(this.todayDate())) {
      this.todayDate.update((oldVal) => addDay(oldVal, -1))
      this.router.navigate(['/day', this.todayDate()]);
    } else if (isMonthDate(this.todayDate())) {
      this.todayDate.update((oldVal) => addMonth(oldVal, -1))
      this.router.navigate(['/month', this.todayDate()]);
    } else if (isWeekDate(this.todayDate())) {
      this.todayDate.update((oldVal) => addWeek(oldVal, -1))
      this.router.navigate(['/week', this.todayDate()]);
    } else {
      console.error("Unrecognized date format", this.todayDate())
    }
    return this.refreshTasks()
  }

  refreshTasks() {
    console.log("refreshTasks")
    this.refreshPendingImportsCount().subscribe()

    const daily$ = this.entryService.listEntries(this.realTodayDate())
      .pipe(
        tap(entries => {
          this.dailyEntryList.set(entries)
        })
      );

    const weekly$ = this.entryService.listEntries(this.todayDate())
      .pipe(
        tap(entries => {
          this.weeklyEntryList.set(entries)
        })
      );

    return forkJoin([daily$, weekly$]);
  }

  refreshPendingImportsCount() {
    return this.taskService.countPastTasks(this.realTodayDate())
      .pipe(
        tap(resp => {
          this.pendingImportsCount.set(resp.count)
        })
      )
  }

  setTitle(newValue: string, entry: Entry) {
    if (!newValue.trim()) {
      newValue = '(empty)'
    }

    return this.entryService.setTitle(entry.id, { title: newValue })
      .pipe(switchMap(() => this.refreshTasks()))
  }

  setTaskUpdate(newValue: string, entry: Entry) {
    return this.taskService.setTaskUpdate(entry.id, { update: newValue })
      .pipe(switchMap(() => this.refreshTasks()))
  }

  createTask(taskValue: string) {
    // Default to daily task if not specified
    return this.createDailyTask(taskValue);
  }

  createNote(taskValue: string) {
    // Default to daily note if not specified
    return this.createDailyNote(taskValue);
  }

  createDailyTask(taskValue: string) {
    const payload = {
      title: taskValue,
      createdDate: this.realTodayDate(),
    }

    return this.taskService.createTask(payload)
      .pipe(
        switchMap(() => this.refreshTasks())
      )
  }

  createWeeklyTask(taskValue: string) {
    const payload = {
      title: taskValue,
      createdDate: this.todayDate(),
    }

    return this.taskService.createTask(payload)
      .pipe(
        switchMap(() => this.refreshTasks())
      )
  }

  createDailyNote(taskValue: string) {
    const payload = {
      title: taskValue,
      createdDate: this.realTodayDate(),
    }

    return this.noteService.createNote(payload)
      .pipe(switchMap(() => this.refreshTasks()))
  }

  createWeeklyNote(taskValue: string) {
    const payload = {
      title: taskValue,
      createdDate: this.todayDate(),
    }

    return this.noteService.createNote(payload)
      .pipe(switchMap(() => this.refreshTasks()))
  }

  markTaskAsDone(entryId: number) {
    return this.taskService.setTaskDone(entryId, { done: true })
      .pipe(switchMap(() => this.refreshTasks()))
  }

  markTaskCancelled(entryId: number) {
    return this.taskService.cancelTask(entryId)
      .pipe(switchMap(() => this.refreshTasks()))
  }

  markTaskAsCreated(entryId: number) {
    return this.taskService.setTaskDone(entryId, { done: false })
      .pipe(switchMap(() => this.refreshTasks()))
  }

  snoozeTask(entry: Entry) {
    // Add different modal depending on day or month

    if (isMonthEntry(entry)) {
      return this.dialog.open(SnoozeMonthEntryDialogComponent, { width: '300px' })
        .afterClosed()
        .pipe(
          switchMap((snoozeDate) => {
            return this.taskService.snoozeTask(entry.id, { date: snoozeDate() })
          }),
          switchMap(() => this.refreshTasks())
        )
    } else if (isDayEntry(entry)) {
      return this.dialog.open(SnoozeDayEntryDialogComponent, { width: '300px' })
        .afterClosed()
        .pipe(
          switchMap((snoozeDate) => {
            return this.taskService.snoozeTask(entry.id, { date: dateToString(snoozeDate()) })
          }),
          switchMap(() => this.refreshTasks())
        )
    } else {
      console.error("Unrecognized entry date type")
      return EMPTY
    }
  }

  migrateToMonthly(entry: Entry) {
    //TODO: Temporary migrate to the same month. Add montly datepicker for a final solution
    let monthlyDate = entry.createdDate.substring(0, 7)
    return this.taskService.migrateTaskToMonthlyLog(entry.id, { date: monthlyDate })
      .pipe(switchMap(() => this.refreshTasks()))
  }

  migrateToWeekly(entry: Entry) {
    //TODO: Temporary migrate to the current week. Add weekly datepicker for a final solution
    let weeklyDate = ThisWeek()
    return this.taskService.migrateTaskToWeeklyLog(entry.id, { date: weeklyDate })
      .pipe(switchMap(() => this.refreshTasks()))
  }

  migrateToDaily(entry: Entry) {
    return this.dialog.open(MigrateToDayEntryDialogComponent, { width: '300px' })
      .afterClosed()
      .pipe(
        switchMap((snoozeDate) => {
          return this.taskService.migrateTaskToDailyLog(entry.id, { date: dateToString(snoozeDate()) })
        }),
        switchMap(() => this.refreshTasks())
      )
  }

  deleteEntry(entryId: number) {
    return this.entryService.deleteEntry(entryId)
      .pipe(switchMap(() => this.refreshTasks()))
  }

  importPastTasks() {
    return this.taskService.importPastTasks(this.realTodayDate())
      .pipe(
        switchMap(() => this.refreshTasks()),
        tap(() => this._snackBar.open("Past tasks migrated")),
        switchMap(() => this.refreshPendingImportsCount())
      )
  }

  handleDropDaily(targetIndex: number, currentRank: number) {
    const id: number = this.dailyEntryList().find(entry => entry.rank == currentRank)!.id

    return this.entryService.changeRank(id, { newRank: targetIndex })
      .pipe(
        switchMap(() => this.refreshTasks())
      )
  }

  handleDropWeekly(targetIndex: number, currentRank: number) {
    const id: number = this.weeklyEntryList().find(entry => entry.rank == currentRank)!.id

    return this.entryService.changeRank(id, { newRank: targetIndex })
      .pipe(
        switchMap(() => this.refreshTasks())
      )
  }

  handleDropToPriorityDaily(priorityIndex: number, targetRank: number) {
    let newRank: number
    if (priorityIndex == 0) {
      newRank = -3
    } else if (priorityIndex == 1) {
      newRank = -2
    } else if (priorityIndex == 2) {
      newRank = -1
    } else {
      console.error("Invalid priority index: ", priorityIndex)
      return EMPTY
    }

    const id: number = this.dailyEntryList().find(entry => entry.rank == targetRank)!.id

    return this.entryService.changeRank(id, { newRank: newRank })
      .pipe(switchMap(() => this.refreshTasks()))
  }

  handleDropToPriorityWeekly(priorityIndex: number, targetRank: number) {
    let newRank: number
    if (priorityIndex == 0) {
      newRank = -3
    } else if (priorityIndex == 1) {
      newRank = -2
    } else if (priorityIndex == 2) {
      newRank = -1
    } else {
      console.error("Invalid priority index: ", priorityIndex)
      return EMPTY
    }

    const id: number = this.weeklyEntryList().find(entry => entry.rank == targetRank)!.id

    return this.entryService.changeRank(id, { newRank: newRank })
      .pipe(switchMap(() => this.refreshTasks()))
  }

  getTaskSummary(taskId: number) {
    return this.taskService.getTaskSummary(taskId)
  }
}

function isMonthEntry(entry: Entry): boolean {
  return isMonthDate(entry.createdDate);
}

function isMonthDate(date: string): boolean {
  return date.length === 7; // 2024-12
}

function isDayEntry(entry: Entry): boolean {
  return isDayDate(entry.createdDate)
}

function isDayDate(date: string): boolean {
  return date.length === 10; // 2024-12-12
}

function isWeekDate(date: string): boolean {
  return date.length === 8 && date.includes('W'); // 2024-W01
}

function dateToString(date: Date): string {
  return `${date.getFullYear()}-${('0' + (date.getMonth() + 1)).slice(-2)}-${('0' + date.getDate()).slice(-2)}`
}
