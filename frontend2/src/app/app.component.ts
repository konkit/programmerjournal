import {Component, model, OnInit, signal} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {MatCardModule} from '@angular/material/card';
import {MatMenuModule} from '@angular/material/menu'
import {MatButtonModule} from '@angular/material/button'
import {MatIconModule} from '@angular/material/icon'
import {MatInputModule} from '@angular/material/input'
import {MatFormFieldModule} from '@angular/material/form-field'
import {AddDay, DayOfWeek, Today} from "../lib/wall_date";
import {TaskService} from "../service/task.service";
import {map, switchMap, take, tap} from "rxjs";
import {CommonModule} from '@angular/common';
import {Task} from "../lib/task";
import {TaskmenuComponent} from './components/taskmenu/taskmenu.component';
import {FormsModule} from '@angular/forms';

interface CurrentEditEntry {
  editedEntryType: string;
  entryId: number;
}

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, MatCardModule, FormsModule, CommonModule, MatButtonModule, MatMenuModule, MatIconModule, TaskmenuComponent, MatFormFieldModule, MatInputModule],
  templateUrl: './app.component.html',
  standalone: true,
  styleUrl: './app.component.scss'
})
export class AppComponent implements OnInit {
  title = 'frontend2';
  protected readonly DayOfWeek = DayOfWeek;

  todayDate = signal<string>(Today());

  taskList = signal<Task[]>([]);

  currentEdit = signal<CurrentEditEntry | null>(null);
  editedValue = model<string>("");

  constructor(private taskService: TaskService) {
  }

  ngOnInit() {
    console.log("Fetching task list")
    this.refreshTasks().subscribe()
  }

  changeDateForward() {
    this.todayDate.update((oldVal) => AddDay(oldVal, 1))
  }

  changeDateBackward() {
    this.todayDate.update((oldVal) => AddDay(oldVal, -1))
  }

  setTitleEdited(task: Task) {
    this.editedValue.set(task.title)
    this.currentEdit.set({editedEntryType: "title", entryId: task.id})
  }

  submitTitleEdit(task: Task) {
    this.taskService.setTaskTitle(task.id, this.editedValue())
        .pipe(switchMap(() => this.refreshTasks()))
        .subscribe(() => {
          this.editedValue.set("")
          this.currentEdit.set(null)
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

  private refreshTasks() {
    console.log("Refreshing tasks")
    return this.taskService.loadTaskList(this.todayDate())
      .pipe(
        tap((tasks: Task[]) => {
          console.log("New tasks:",  tasks)
          this.taskList.set(tasks)
        })
      )
  }
}
