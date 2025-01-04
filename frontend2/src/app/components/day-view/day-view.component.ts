import {Component, OnInit, signal} from '@angular/core';
import {TaskListComponent} from '../task-list/task-list.component';
import {AddDay, Today} from '../../../lib/wall_date';
import {Task} from '../../../lib/task';
import {tap} from 'rxjs';
import {TaskService} from '../../../frontend-client';

@Component({
  selector: 'app-day-view',
  imports: [
    TaskListComponent
  ],
  templateUrl: './day-view.component.html',
  standalone: true,
  styleUrl: './day-view.component.scss'
})
export class DayViewComponent implements OnInit {

  todayDate = signal<string>(Today());

  taskList = signal<Task[]>([]);

  constructor(private taskService: TaskService) {
  }

  ngOnInit() {
    console.log("Fetching task list")
    this.refreshTasks().subscribe()
  }

  dateForward() {
    this.todayDate.update((oldVal) => AddDay(oldVal, 1))
    this.refreshTasks().subscribe()
  }

  dateBackward() {
    this.todayDate.update((oldVal) => AddDay(oldVal, -1))
    this.refreshTasks().subscribe()
  }

  refreshTasks() {
    console.log("Refreshing tasks")
    return this.taskService.listTasks(this.todayDate())
      .pipe(
        tap(tasks => {
          console.log("Setting tasks")
          this.taskList.set(tasks)
        })
      )
  }

}
