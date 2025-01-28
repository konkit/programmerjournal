import {inject, Injectable, signal} from '@angular/core';
import {Entry, EntryService, TaskService} from '../../frontend-client';
import {EMPTY, empty, Observable, switchMap, tap} from 'rxjs';
import {Today} from '../../lib/wall_date';
import {NoteService} from '../../frontend-client/api/note.service';
import {SnoozeMonthEntryDialogComponent} from '../components/snooze-month-dialog/snooze-month-entry-dialog.component';
import {SnoozeDayEntryDialogComponent} from '../components/snooze-day-dialog/snooze-day-entry-dialog.component';
import {MatDialog} from '@angular/material/dialog';
import {
  MigrateToDayEntryDialogComponent
} from '../components/migrate-to-day-dialog/migrate-to-day-entry-dialog.component';
import {MatSnackBar} from '@angular/material/snack-bar';
import {CdkDragDrop} from '@angular/cdk/drag-drop';

@Injectable({
  providedIn: 'root'
})
export class EntryListService {

  todayDate = signal<string>(Today());

  entryList = signal<Entry[]>([]);

  readonly dialog = inject(MatDialog);
  private _snackBar = inject(MatSnackBar);

  constructor(private entryService: EntryService,
              private noteService: NoteService,
              private taskService: TaskService) {
  }

  refreshTasks() {
    console.log("Refreshing tasks")
    return this.entryService.listEntries(this.todayDate())
      .pipe(
        tap(entries => {
          console.log("Setting tasks")
          this.entryList.set(entries)
        })
      )
  }

  setTitle(newValue: string, entry: Entry) {
    if (!newValue.trim()) {
      newValue = '(empty)'
    }

    this.entryService.setTitle(entry.id, {title: newValue})
      .pipe(tap(() => this.refreshTasks()))
      .subscribe()
  }

  createTask(taskValue: string) {
    const payload = {
      title: taskValue,
      createdDate: this.todayDate(),
    }

    return this.taskService.createTask(payload)
      .pipe(
        switchMap(() => this.refreshTasks())
      )
  }

  createNote(taskValue: string) {
    const payload = {
      title: taskValue,
      createdDate: this.todayDate(),
    }

    return this.noteService.createNote(payload)
      .pipe(switchMap(() => this.refreshTasks()))
  }

  markTaskAsDone(entryId: number) {
    return this.taskService.setTaskDone(entryId, {done: true})
      .pipe(switchMap(() => this.refreshTasks()))
  }

  markTaskAsCreated(entryId: number) {
    return this.taskService.setTaskDone(entryId, {done: false})
      .pipe(switchMap(() => this.refreshTasks()))
  }

  snoozeTask(entry: Entry) {
    // Add different modal depending on day or month

    if (isMonthEntry(entry)) {
      return this.dialog.open(SnoozeMonthEntryDialogComponent, {width: '300px'})
        .afterClosed()
        .pipe(
          switchMap((snoozeDate) => {
            return this.taskService.snoozeTask(entry.id, {date: snoozeDate()})
          }),
          switchMap(() => this.refreshTasks())
        )
    } else if (isDayEntry(entry)) {
     return this.dialog.open(SnoozeDayEntryDialogComponent, {width: '300px'})
        .afterClosed()
        .pipe(
          switchMap((snoozeDate) => {
            return this.taskService.snoozeTask(entry.id, {date: dateToString(snoozeDate())})
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
    return this.taskService.migrateTaskToMonthlyLog(entry.id, {date: monthlyDate})
      .pipe(tap(() => this.refreshTasks()))
  }

  migrateToDaily(entry: Entry) {
    return this.dialog.open(MigrateToDayEntryDialogComponent, {width: '300px'})
      .afterClosed()
      .pipe(
        switchMap((snoozeDate) => {
          return this.taskService.migrateTaskToDailyLog(entry.id, {date: dateToString(snoozeDate())})
        }),
        switchMap(() => this.refreshTasks())
      )
  }

  deleteEntry(entryId: number) {
    return this.entryService.deleteEntry(entryId)
      .pipe(tap(() => this.refreshTasks()))

  }

  importPastTasks() {
    return this.taskService.importPastTasks(this.todayDate())
      .pipe(
        switchMap(() => this.refreshTasks()),
        tap(() => this._snackBar.open("Past tasks migrated"))
      )
  }

  handleDrop(targetIndex: number, currentRank: number) {
    const id: number = this.entryList().find(entry => entry.rank == currentRank)!.id

    return this.entryService.changeRank(id, {newRank: targetIndex})
      .pipe(
        switchMap(() => this.refreshTasks())
      )
  }

  handleDropToPriority(priorityIndex: number, targetRank: number) {
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

    const id: number = this.entryList().find(entry => entry.rank == targetRank)!.id

    return this.entryService.changeRank(id, {newRank: newRank})
      .pipe(switchMap(() => this.refreshTasks()))
  }

  getTaskSummary(taskId: number) {
    return this.taskService.getTaskSummary(taskId)
  }
}

function isMonthEntry(entry: Entry): boolean {
  return entry.createdDate.length === 7; // 2024-12
}

function isDayEntry(entry: Entry): boolean {
  return entry.createdDate.length === 10; // 2024-12-12
}

function dateToString(date: Date): string {
  return `${date.getFullYear()}-${('0' + (date.getMonth() + 1)).slice(-2)}-${('0' + date.getDate()).slice(-2)}`
}
