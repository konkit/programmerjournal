import {Component, inject, OnInit, signal, ViewChild} from '@angular/core';
import {Task} from '../../../lib/task';
import {TaskService, TaskSummary} from '../../../frontend-client';
import {switchMap, tap} from 'rxjs';
import {RouterOutlet} from '@angular/router';
import {FormControl, FormsModule, ReactiveFormsModule} from '@angular/forms';
import {CommonModule} from '@angular/common';
import {MatCardModule} from '@angular/material/card';
import {MatMenuModule} from '@angular/material/menu'
import {MatChipsModule} from '@angular/material/chips'
import {MatButtonModule} from '@angular/material/button'
import {MatIconModule} from '@angular/material/icon'
import {MatInputModule} from '@angular/material/input'
import {MatFormFieldModule} from '@angular/material/form-field'
import {AddDay, DayOfWeek, Today} from "../../../lib/wall_date";
import {TaskmenuComponent} from '../taskmenu/taskmenu.component';
import {MatDialog,} from '@angular/material/dialog';
import {SnoozeDialogComponent} from '../snooze-dialog/snooze-dialog.component';
import {CdkDrag, CdkDragDrop, CdkDragHandle, CdkDropList} from '@angular/cdk/drag-drop';
import {MatSelectModule} from '@angular/material/select';
import {MatDrawer, MatSidenavModule} from '@angular/material/sidenav';
import {TaskSidebarComponent} from '../task-sidebar/task-sidebar.component';
import {MatToolbar} from '@angular/material/toolbar';

@Component({
  selector: 'app-task-list',
  imports: [RouterOutlet, MatSidenavModule, MatFormFieldModule, MatSelectModule, MatButtonModule, CdkDropList, CdkDrag, CdkDragHandle, MatChipsModule, MatCardModule, MatButtonModule, FormsModule, CommonModule, MatButtonModule, MatMenuModule, MatIconModule, TaskmenuComponent, MatFormFieldModule, MatInputModule, TaskSidebarComponent, ReactiveFormsModule, MatToolbar],
  templateUrl: './task-list.component.html',
  styleUrl: './task-list.component.scss',
  standalone: true,
})
export class TaskListComponent implements OnInit {
  protected readonly DayOfWeek = DayOfWeek;

  todayDate = signal<string>(Today());
  taskList = signal<Task[]>([]);

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  editedTask = signal<Task | null>(null)
  editedTaskSummary = signal<TaskSummary | null>(null)

  creatingNewTask = false

  readonly dialog = inject(MatDialog);
  createTaskFormControl = new FormControl("");

  constructor(private taskService: TaskService) {
  }

  ngOnInit() {
    console.log("Fetching task list")
    this.refreshTasks().subscribe()
  }

  changeDateForward() {
    this.todayDate.update((oldVal) => AddDay(oldVal, 1))
    this.refreshTasks().subscribe()
  }

  changeDateBackward() {
    this.todayDate.update((oldVal) => AddDay(oldVal, -1))
    this.refreshTasks().subscribe()
  }

  submitTitleEditWithValue(task: Task, e: Event) {
    let newValue = (e.target as HTMLDivElement).innerText

    if (!newValue.trim()) {
      newValue = '(empty)'
    }

    this.taskService.setTaskTitle(task.id, {title: newValue})
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe()
  }

  createNewTask() {
    // return this.taskService.createTask("", this.todayDate())
    //   .pipe(switchMap(() => this.refreshTasks()))
    //   .subscribe()

    this.creatingNewTask = true
  }

  // deleteTask(taskId: number) {
  //   return this.taskService.deleteTask(taskId)
  //     .pipe(switchMap(() => this.refreshTasks()))
  //     .subscribe()
  // }

  markTaskAsDone(task: Task) {
    this.taskService.setTaskDone(task.id, {done: true})
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe()
  }

  markTaskAsCreated(task: Task) {
    this.taskService.setTaskDone(task.id, {done: false})
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe()
  }

  snoozeTask(task: Task) {
    this.dialog.open(SnoozeDialogComponent, {width: '300px'})
      .afterClosed()
      .pipe(
        switchMap((snoozeDate) => this.taskService.snoozeTask(task.id, snoozeDate)),
        switchMap(() => this.refreshTasks())
      )
      .subscribe()
  }

  private refreshTasks() {
    console.log("Refreshing tasks")
    return this.taskService.listTasks(this.todayDate())
      .pipe(
        tap((tasks: Task[]) => {
          this.taskList.set(tasks)
        })
      )
  }

  handleDrop(e: CdkDragDrop<string[]>) {
    const id: number = this.taskList()[e.previousIndex].id
    this.taskService.changeTaskRank(id, {newRank: e.currentIndex})
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe()
  }

  importPastTasks() {
    this.taskService.importPastTasks(this.todayDate())
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe()
  }

  openUpdates(task: Task) {
    this.editedTask.set(task)

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
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe(() => {
        this.sideDrawer.close()
      })
  }

  submitNewTask() {
    let taskValue = this.createTaskFormControl.value || "";

    const payload = {
      title: taskValue,
      createdDate:  this.todayDate(),
    }

    return this.taskService.createTask(payload)
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe(() => this.creatingNewTask = false)
  }
}
