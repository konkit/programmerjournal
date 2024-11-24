import {Component, inject, OnInit, signal} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {MatCardModule} from '@angular/material/card';
import {MatMenuModule} from '@angular/material/menu'
import {MatButtonModule} from '@angular/material/button'
import {MatIconModule} from '@angular/material/icon'
import {MatInputModule} from '@angular/material/input'
import {MatFormFieldModule} from '@angular/material/form-field'
import {AddDay, DayOfWeek, Today} from "../lib/wall_date";
import {TaskService} from "../service/task.service";
import {switchMap, tap} from "rxjs";
import {CommonModule} from '@angular/common';
import {Task} from "../lib/task";
import {TaskmenuComponent} from './components/taskmenu/taskmenu.component';
import {FormsModule} from '@angular/forms';
import {
  MatDialog,
  MatDialogActions,
  MatDialogClose,
  MatDialogContent,
  MatDialogRef,
  MatDialogTitle,
} from '@angular/material/dialog';
import {SnoozeDialogComponent} from './components/snooze-dialog/snooze-dialog.component';
import {MatNativeDateModule} from '@angular/material/core';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, MatCardModule, MatButtonModule, FormsModule, CommonModule, MatButtonModule, MatMenuModule, MatIconModule, TaskmenuComponent, MatFormFieldModule, MatInputModule],
  templateUrl: './app.component.html',
  standalone: true,
  styleUrl: './app.component.scss'
})
export class AppComponent implements OnInit {
  title = 'Programmer Journal';
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

    this.taskService.setTaskTitle(task.id, newValue)
      .pipe(switchMap(() => this.refreshTasks()))
      .subscribe(() => {
      })
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

  markTaskAsCancelled(task: Task) {
    this.taskService.setTaskCancelled(task.id, task)
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
}
