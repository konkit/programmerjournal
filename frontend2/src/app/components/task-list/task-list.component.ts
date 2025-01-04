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
import {MatToolbar} from '@angular/material/toolbar';
import {StatusButtonComponent} from '../status-button/status-button.component';
import {RouterLink} from '@angular/router';

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
    MatToolbar,
    StatusButtonComponent,
    RouterLink
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
      return `${getMonthFromDate(this.todayDate())} ${getYearFromDate(this.todayDate())}`
    } else {
      return `${getDayOfWeekFromDate(this.todayDate())}, ${this.todayDate()}`
    }
  })

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)

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
    this.creatingNewTask = true
  }

  markTaskAsDone(task: Task) {
    this.taskService.setTaskDone(task.id, {done: true})
      .subscribe(() => this.onRefreshTasks.emit())
  }

  markTaskAsCreated(task: Task) {
    this.taskService.setTaskDone(task.id, {done: false})
      .pipe(tap(() => this.onRefreshTasks.emit()))
      .subscribe()
  }

  snoozeTask(task: Task) {
    this.dialog.open(SnoozeDialogComponent, {width: '300px'})
      .afterClosed()
      .pipe(
        switchMap((snoozeDate) => {
          return this.taskService.snoozeTask(task.id, {date: snoozeDate()})
        }),
        tap(() => this.onRefreshTasks.emit())
      )
      .subscribe()
  }

  handleDrop(e: CdkDragDrop<string[]>) {
    const id: number = this.taskList()[e.previousIndex].id
    this.taskService.changeTaskRank(id, {newRank: e.currentIndex})
      .pipe(tap(() => this.onRefreshTasks.emit()))
      .subscribe()
  }

  importPastTasks() {
    this.taskService.importPastTasks(this.todayDate())
      .pipe(tap(() => this.onRefreshTasks.emit()))
      .subscribe()
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

  submitNewTask() {
    let taskValue = this.createTaskFormControl.value || "";

    const payload = {
      title: taskValue,
      createdDate: this.todayDate(),
    }

    return this.taskService.createTask(payload)
      .pipe(tap(() => this.onRefreshTasks.emit()))
      .subscribe(() => this.creatingNewTask = false)
  }
}
