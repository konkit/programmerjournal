import {Component, computed, inject, input, output, signal, ViewChild} from '@angular/core';
import {Entry, EntryService, TaskService, TaskSummary} from '../../../frontend-client';
import {switchMap, tap} from 'rxjs';
import {FormControl, FormsModule, ReactiveFormsModule} from '@angular/forms';
import {CommonModule} from '@angular/common';
import {MatCardModule} from '@angular/material/card';
import {MatMenuModule} from '@angular/material/menu'
import {MatChipsModule} from '@angular/material/chips'
import {MatButtonModule} from '@angular/material/button'
import {MatIconModule} from '@angular/material/icon'
import {MatInputModule} from '@angular/material/input'
import {MatFormFieldModule} from '@angular/material/form-field'
import {getDayOfWeekFromDate, getMonthFromDate, getYearFromDate} from "../../../lib/wall_date";
import {MatDialog,} from '@angular/material/dialog';
import {CdkDrag, CdkDragDrop, CdkDragHandle, CdkDropList} from '@angular/cdk/drag-drop';
import {MatSelectModule} from '@angular/material/select';
import {MatDrawer, MatSidenavModule} from '@angular/material/sidenav';
import {EntrySidebarComponent} from '../entry-sidebar/entry-sidebar.component';
import {StatusButtonComponent} from '../status-button/status-button.component';
import {NavToolbarComponent} from '../nav-toolbar/nav-toolbar.component';
import {MatSnackBar} from '@angular/material/snack-bar';
import {MatList, MatListItem} from '@angular/material/list';
import {NoteService} from '../../../frontend-client/api/note.service';
import {SnoozeMonthEntryDialogComponent} from '../snooze-month-dialog/snooze-month-entry-dialog.component';
import {SnoozeDayEntryDialogComponent} from '../snooze-day-dialog/snooze-day-entry-dialog.component';
import {MigrateToDayEntryDialogComponent} from '../migrate-to-day-dialog/migrate-to-day-entry-dialog.component';

export enum EditorStateEnum {
  IDLE,
  EDITING_NEW_TASK,
  EDITING_NEW_NOTE,
}

@Component({
  imports: [
    MatSidenavModule,
    MatFormFieldModule,
    MatSelectModule,
    MatButtonModule,
    CdkDropList,
    CdkDrag,
    CdkDragHandle,
    MatChipsModule,
    MatCardModule,
    MatButtonModule,
    FormsModule,
    CommonModule,
    MatButtonModule,
    MatMenuModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    EntrySidebarComponent,
    ReactiveFormsModule,
    StatusButtonComponent,
    NavToolbarComponent,
    MatList,
    MatListItem
  ],
  selector: 'app-entry-list',
  standalone: true,
  styleUrl: './entry-list.component.scss',
  templateUrl: './entry-list.component.html',
})
export class EntryListComponent {
  todayDate = input("")
  entryList = input<Entry[]>([])

  dateForward = output<void>()
  dateBackward = output<void>()
  onRefreshTasks = output<void>()

  entryPriority1 = computed(() => this.entryList().find(e => e.rank === -3))
  entryPriority2 = computed(() => this.entryList().find(e => e.rank === -2))
  entryPriority3 = computed(() => this.entryList().find(e => e.rank === -1))

  nonPriorityCount = computed(() => this.entryList().filter(x => x.rank >= 0).length)

  currentDateString = computed<string>(() => {
    let isMonthlyDate = this.todayDate().length == 7;
    if (isMonthlyDate) {
      return `Month: ${getMonthFromDate(this.todayDate())} ${getYearFromDate(this.todayDate())}`
    } else {
      return `Day: ${getDayOfWeekFromDate(this.todayDate())}, ${this.todayDate()}`
    }
  })

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)

  private _snackBar = inject(MatSnackBar);

  EditorStateEnum = EditorStateEnum;
  editorState = EditorStateEnum.IDLE

  readonly dialog = inject(MatDialog);
  newEntryFormControl = new FormControl("");

  constructor(private taskService: TaskService,
              private entryService: EntryService,
              private noteService: NoteService) {}

  changeDateForward() {
    this.dateForward.emit()
  }

  changeDateBackward() {
    this.dateBackward.emit()
  }

  submitTitleEditWithValue(entry: Entry, e: Event) {
    let newValue = (e.target as HTMLDivElement).innerText

    if (!newValue.trim()) {
      newValue = '(empty)'
    }

    this.entryService.setTitle(entry.id, {title: newValue})
      .pipe(tap(() => this.onRefreshTasks.emit()))
      .subscribe()
  }

  setCreatingNewTaskState() {
    this.newEntryFormControl.setValue("")
    this.editorState = EditorStateEnum.EDITING_NEW_TASK;
  }

  setCreatingNewNoteState() {
    this.newEntryFormControl.setValue("")
    this.editorState = EditorStateEnum.EDITING_NEW_NOTE;
  }

  submitNewTask() {
    let taskValue = this.newEntryFormControl.value || "";

    const payload = {
      title: taskValue,
      createdDate: this.todayDate(),
    }

    return this.taskService.createTask(payload)
      .subscribe(() => {
        this.onRefreshTasks.emit()
        this.editorState = EditorStateEnum.IDLE;
      })
  }

  submitNewNote() {
    let taskValue = this.newEntryFormControl.value || "";

    const payload = {
      title: taskValue,
      createdDate: this.todayDate(),
    }

    return this.noteService.createNote(payload)
      .subscribe(() => {
        this.onRefreshTasks.emit()
        this.editorState = EditorStateEnum.IDLE;
      })
  }

  cancelEdit() {
    this.newEntryFormControl.setValue("")
    this.editorState = EditorStateEnum.IDLE;
  }

  markTaskAsDone(entry: Entry) {
    this.taskService.setTaskDone(entry.id, {done: true})
      .subscribe(() => this.onRefreshTasks.emit())
  }

  markTaskAsCreated(entry: Entry) {
    this.taskService.setTaskDone(entry.id, {done: false})
      .subscribe(() => this.onRefreshTasks.emit())
  }

  snoozeTask(entry: Entry) {
    // Add different modal depending on day or month

    if (isMonthEntry(entry)) {
      this.dialog.open(SnoozeMonthEntryDialogComponent, {width: '300px'})
        .afterClosed()
        .pipe(
          switchMap((snoozeDate) => {
            return this.taskService.snoozeTask(entry.id, {date: snoozeDate()})
          }),
        )
        .subscribe(() => {
          this.onRefreshTasks.emit()
        })
    } else if (isDayEntry(entry)) {
      this.dialog.open(SnoozeDayEntryDialogComponent, {width: '300px'})
        .afterClosed()
        .pipe(
          switchMap((snoozeDate) => {
            return this.taskService.snoozeTask(entry.id, {date: dateToString(snoozeDate())})
          }),
        )
        .subscribe(() => {
          this.onRefreshTasks.emit()
        })
    } else {
      console.error("Unrecognized entry date type")
    }
  }

  handleDrop(e: CdkDragDrop<string[]>) {
    console.log("handleDrop", "e", e, "e.item.data", e.item.data)

    const id: number = this.entryList().find(entry => entry.rank == e.item.data)!.id

    console.log("id: ", id)

    this.entryService.changeRank(id, {newRank: e.currentIndex})
      .subscribe(() => {
        this.onRefreshTasks.emit()
      })
  }

  handleDropToPriority(e: CdkDragDrop<string[]>, priorityIndex: number) {
    console.log("handleDropToPriority", "e", e, "priorityIndex", priorityIndex)

    let newRank: number
    if (priorityIndex == 1) {
      newRank = -3
    } else if (priorityIndex == 2) {
      newRank = -2
    } else if (priorityIndex == 3) {
      newRank = -1
    } else {
      console.error("Invalid priority index: ", priorityIndex)
      return
    }

    const id: number = this.entryList().find(entry => entry.rank == e.item.data)!.id

    this.entryService.changeRank(id, {newRank: newRank})
      .subscribe(() => {
        this.onRefreshTasks.emit()
      })
  }

  importPastTasks() {
    this.taskService.importPastTasks(this.todayDate())
      .subscribe(() => {
        this.onRefreshTasks.emit()
        this._snackBar.open("Past tasks migrated")
      })
  }

  openUpdates(entry: Entry) {
    this.taskService.getTaskSummary(entry.id)
      .subscribe((ts) => {
        console.log("openUpdates")
        this.editedTaskSummary.set(ts)
        this.sideDrawer.open()
      })
  }

  reloadTaskSummary(taskId: number) {
    this.taskService.getTaskSummary(taskId)
      .subscribe((ts) => {
        console.log("reloadTaskSummary")
        this.editedTaskSummary.set(ts)
      })
  }

  deleteTaskFromSidebar(taskId: number) {
    return this.entryService.deleteEntry(taskId)
      .pipe(tap(() => this.onRefreshTasks.emit()))
      .subscribe(() => {
        this.sideDrawer.close()
      })
  }

  migrateToMonthly(entry: Entry) {
    //TODO: Temporary migrate to the same month. Add montly datepicker for a final solution
    let monthlyDate = entry.createdDate.substring(0, 7)
    this.taskService.migrateTaskToMonthlyLog(entry.id, {date: monthlyDate})
      .pipe(tap(() => this.onRefreshTasks.emit()))
      .subscribe()
  }

  migrateToDaily(entry: Entry) {
    this.dialog.open(MigrateToDayEntryDialogComponent, {width: '300px'})
      .afterClosed()
      .pipe(
        switchMap((snoozeDate) => {
          return this.taskService.migrateTaskToDailyLog(entry.id, {date: dateToString(snoozeDate())})
        }),
      )
      .subscribe(() => {
        this.onRefreshTasks.emit()
      })
  }
}

function dateToString(date: Date): string {
  return `${date.getFullYear()}-${('0' + (date.getMonth() + 1)).slice(-2)}-${('0' + date.getDate()).slice(-2)}`
}

function isMonthEntry(entry: Entry): boolean {
  return entry.createdDate.length === 7; // 2024-12
}

function isDayEntry(entry: Entry): boolean {
  return entry.createdDate.length === 10; // 2024-12-12
}
