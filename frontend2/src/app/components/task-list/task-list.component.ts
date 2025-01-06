import {Component, computed, inject, input, output, signal, ViewChild} from '@angular/core';
import {Task} from '../../../lib/task';
import {TaskService, TaskSummary} from '../../../frontend-client';
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
import {SnoozeDialogComponent} from '../snooze-dialog/snooze-dialog.component';
import {CdkDrag, CdkDragDrop, CdkDragHandle, CdkDropList} from '@angular/cdk/drag-drop';
import {MatSelectModule} from '@angular/material/select';
import {MatDrawer, MatSidenavModule} from '@angular/material/sidenav';
import {TaskSidebarComponent} from '../task-sidebar/task-sidebar.component';
import {StatusButtonComponent} from '../status-button/status-button.component';
import {NavToolbarComponent} from '../nav-toolbar/nav-toolbar.component';
import {MatSnackBar} from '@angular/material/snack-bar';

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
    TaskSidebarComponent,
    ReactiveFormsModule,
    StatusButtonComponent,
    NavToolbarComponent
  ],
  selector: 'app-task-list',
  standalone: true,
  styleUrl: './task-list.component.scss',
  templateUrl: './task-list.component.html',
})
export class TaskListComponent {
  todayDate = input("")
  taskList = input<Task[]>([])

  dateForward = output<void>()
  dateBackward = output<void>()
  onRefreshTasks = output<void>()

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

  creatingNewTask = false

  readonly dialog = inject(MatDialog);
  createTaskFormControl = new FormControl("");

  constructor(private taskService: TaskService) {}

  changeDateForward() {
    this.dateForward.emit()
  }

  changeDateBackward() {
    this.dateBackward.emit()
  }

  submitTitleEditWithValue(task: Task, e: Event) {
    let newValue = (e.target as HTMLDivElement).innerText

    if (!newValue.trim()) {
      newValue = '(empty)'
    }

    this.taskService.setTaskTitle(task.id, {title: newValue})
      .pipe(tap(() => this.onRefreshTasks.emit()))
      .subscribe()
  }

  setCreatingNewTaskState() {
    this.createTaskFormControl.setValue("")
    this.creatingNewTask = true
  }

  submitNewTask() {
    let taskValue = this.createTaskFormControl.value || "";

    const payload = {
      title: taskValue,
      createdDate: this.todayDate(),
    }

    return this.taskService.createTask(payload)
      .subscribe(() => {
        this.onRefreshTasks.emit()
        this.creatingNewTask = false
      })
  }

  markTaskAsDone(task: Task) {
    this.taskService.setTaskDone(task.id, {done: true})
      .subscribe(() => this.onRefreshTasks.emit())
  }

  markTaskAsCreated(task: Task) {
    this.taskService.setTaskDone(task.id, {done: false})
      .subscribe(() => this.onRefreshTasks.emit())
  }

  snoozeTask(task: Task) {
    this.dialog.open(SnoozeDialogComponent, {width: '300px'})
      .afterClosed()
      .pipe(
        switchMap((snoozeDate) => {
          return this.taskService.snoozeTask(task.id, {date: snoozeDate()})
        }),
      )
      .subscribe(() => {
        this.onRefreshTasks.emit()
      })
  }

  handleDrop(e: CdkDragDrop<string[]>) {
    const id: number = this.taskList()[e.previousIndex].id
    this.taskService.changeTaskRank(id, {newRank: e.currentIndex})
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

  openUpdates(task: Task) {
    this.taskService.getTaskSummary(task.id)
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
    return this.taskService.deleteTask(taskId)
      .pipe(tap(() => this.onRefreshTasks.emit()))
      .subscribe(() => {
        this.sideDrawer.close()
      })
  }

  migrateToMonthly(task: Task) {
    //TODO: Temporary migrate to the same month. Add montly datepicker for a final solution
    let monthlyDate = task.createdDate.substring(0, 7)
    this.taskService.migrateToMonthlyLog(task.id, {date: monthlyDate})
      .pipe(tap(() => this.onRefreshTasks.emit()))
      .subscribe()
  }
}
