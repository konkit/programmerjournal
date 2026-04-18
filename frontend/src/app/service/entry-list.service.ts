import { inject, Injectable, signal, computed } from '@angular/core';
import { Entry, EntryService, TaskService } from '../../frontend-client';
import { EMPTY, switchMap, tap } from 'rxjs';
import { addDay, addMonth, addWeek, getDateFromWeek, getWeekString, Today, toWallMonth, toWeeklyDate, isDayDate, isWeekDate, isMonthDate } from '../../lib/wall_date';
import { NoteService } from '../../frontend-client/api/note.service';
import { SnoozeMonthEntryDialogComponent } from '../components/snooze-month-dialog/snooze-month-entry-dialog.component';
import { SnoozeDayEntryDialogComponent } from '../components/snooze-day-dialog/snooze-day-entry-dialog.component';
import { MatDialog } from '@angular/material/dialog';
import {
  MigrateToDayEntryDialogComponent
} from '../components/migrate-to-day-dialog/migrate-to-day-entry-dialog.component';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { SnoozeWeekDialogComponent } from '../components/snooze-week-dialog/snooze-week-dialog.component';

@Injectable({
  providedIn: 'root'
})
export class EntryListService {

  todayDate = signal<string>(Today());

  entryList = signal<Entry[]>([]);
  refreshTrigger = signal<number>(0);
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
    this.triggerRefresh()
    this.refreshPendingImportsCount().subscribe()
    return this.entryService.listEntries(this.todayDate())
      .pipe(
        tap(entries => {
          console.log("refreshTasks - updating entryList signal")
          this.entryList.set(entries)
        })
      )
  }

  triggerRefresh() {
    this.refreshTrigger.update(n => n + 1)
  }

  refreshPendingImportsCount() {
    return this.taskService.countPastTasks(this.todayDate())
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

  createTask(taskValue: string, date?: string) {
    const payload = {
      title: taskValue,
      createdDate: date || this.todayDate(),
    }

    return this.taskService.createTask(payload)
      .pipe(
        switchMap(() => this.refreshTasks())
      )
  }

  createNote(taskValue: string, date?: string) {
    const payload = {
      title: taskValue,
      createdDate: date || this.todayDate(),
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
    } else if (isWeekEntry(entry)) {
      return this.dialog.open(SnoozeWeekDialogComponent, { width: '300px' })
        .afterClosed()
        .pipe(
          switchMap((snoozeWeekDate) => {
            return this.taskService.snoozeTask(entry.id, { date: snoozeWeekDate() })
          }),
          switchMap(() => this.refreshTasks())
        )
    } else {
      console.error("Unrecognized entry date type")
      return EMPTY
    }
  }

  migrateToMonthly(entry: Entry) {
    let monthlyDate: string;
    if (isWeekEntry(entry)) {
      // Parse 2024-W01
      let year = parseInt(entry.createdDate.substring(0, 4))
      let week = parseInt(entry.createdDate.substring(6, 8))
      let date = getDateFromWeek(year, week);
      monthlyDate = toWallMonth(date);
    } else {
      //TODO: Temporary migrate to the same month. Add monthly datepicker for a final solution
      monthlyDate = entry.createdDate.substring(0, 7)
    }

    return this.taskService.migrateTaskToMonthlyLog(entry.id, { date: monthlyDate })
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

  migrateToDailyDirect(entryId: number, date: string) {
    return this.taskService.migrateTaskToDailyLog(entryId, { date: date })
      .pipe(switchMap(() => this.refreshTasks()))
  }

  migrateToWeekly(entry: Entry) {
    //TODO: Temporary migrate to the same month. Add montly datepicker for a final solution
    let weeklyDate = toWeeklyDate(entry.createdDate)
    return this.taskService.migrateTaskToWeeklyLog(entry.id, { date: weeklyDate })
      .pipe(switchMap(() => this.refreshTasks()))
  }

  deleteEntry(entryId: number) {
    return this.entryService.deleteEntry(entryId)
      .pipe(switchMap(() => this.refreshTasks()))
  }

  importPastTasks() {
    return this.taskService.importPastTasks(this.todayDate())
      .pipe(
        switchMap(() => this.refreshTasks()),
        tap(() => this._snackBar.open("Past tasks migrated")),
        switchMap(() => this.refreshPendingImportsCount())
      )
  }

  handleDrop(targetIndex: number, currentRank: number) {
    const id: number = this.entryList().find(entry => entry.rank == currentRank)!.id

    return this.entryService.changeRank(id, { newRank: targetIndex })
      .pipe(
        switchMap(() => this.refreshTasks())
      )
  }

  getTaskSummary(taskId: number) {
    return this.taskService.getTaskSummary(taskId)
  }

  getWeeklyTasks() {
    const weeklyDateString = computed(() => toWeeklyDate(this.todayDate()))
    return this.entryService.listEntries(weeklyDateString())
  }

  migrateWeeklyTaskToToday(entryId: number) {
    return this.taskService.migrateTaskToDailyLog(entryId, { date: this.todayDate() }).pipe(
      switchMap(() => this.refreshTasks())
    )
  }

}

function isMonthEntry(entry: Entry): boolean {
  return isMonthDate(entry.createdDate);
}

function isDayEntry(entry: Entry): boolean {
  return isDayDate(entry.createdDate)
}

function isWeekEntry(entry: Entry): boolean {
  return isWeekDate(entry.createdDate)
}

function dateToString(date: Date): string {
  return `${date.getFullYear()}-${('0' + (date.getMonth() + 1)).slice(-2)}-${('0' + date.getDate()).slice(-2)}`
}
