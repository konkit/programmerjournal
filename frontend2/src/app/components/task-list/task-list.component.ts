import {Component, inject, OnInit, signal} from '@angular/core';
import {Task} from '../../../lib/task';
import {TaskService} from '../../../service/task.service';
import {switchMap, tap} from 'rxjs';
import {RouterOutlet} from '@angular/router';
import {FormsModule} from '@angular/forms';
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
import {CdkDragDrop, CdkDropList, CdkDrag, CdkDragHandle, moveItemInArray} from '@angular/cdk/drag-drop';

@Component({
  selector: 'app-task-list',
  imports: [RouterOutlet, CdkDropList, CdkDrag, CdkDragHandle, MatChipsModule, MatCardModule, MatButtonModule, FormsModule, CommonModule, MatButtonModule, MatMenuModule, MatIconModule, TaskmenuComponent, MatFormFieldModule, MatInputModule],
  templateUrl: './task-list.component.html',
  styleUrl: './task-list.component.scss',
  standalone: true,
})
export class TaskListComponent implements OnInit {
  protected readonly DayOfWeek = DayOfWeek;

  todayDate = signal<string>(Today());
  taskList = signal<Task[]>([]);

  readonly dialog = inject(MatDialog);

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

    this.taskService.setTaskTitle(task.id, newValue)
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe()
  }

  submitTaskUpdate(task: Task, e: Event) {
    let newValue = (e.target as HTMLSpanElement).innerText

    this.taskService.setTaskUpdate(task.id, newValue)
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe()
  }

  createNewTask() {
    return this.taskService.createTask("", this.todayDate())
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe()
  }

  deleteTask(task: Task) {
    return this.taskService.deleteTask(task.id)
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe()
  }

  markTaskAsDone(task: Task) {
    this.taskService.setTaskDone(task.id, task)
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
    return this.taskService.loadTaskList(this.todayDate())
      .pipe(
        tap((tasks: Task[]) => {
          console.log("New tasks:", tasks)
          this.taskList.set(tasks)
        })
      )
  }

  handleDrop(e: CdkDragDrop<string[]>) {
    const id: number = this.taskList()[e.previousIndex].id
    this.taskService.handleDrop(id, e.currentIndex)
      .pipe(
        switchMap(() => this.refreshTasks())
      )
      .subscribe()
  }
}
