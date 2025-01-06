import {Component, signal} from '@angular/core';
import {addMonth, ThisMonth, WallMonth} from '../../../lib/wall_date';
import {Task} from '../../../lib/task';
import {TaskService} from '../../../frontend-client';
import {tap} from 'rxjs';
import {TaskListComponent} from '../../components/task-list/task-list.component';

@Component({
  selector: 'app-month-view',
  imports: [TaskListComponent],
  templateUrl: './month-view.component.html',
  standalone: true,
  styleUrl: './month-view.component.scss'
})
export class MonthViewComponent {

  todayDate = signal<WallMonth>(ThisMonth());

  taskList = signal<Task[]>([]);

  constructor(private taskService: TaskService) {
  }

  ngOnInit() {
    console.log("Fetching task list")
    this.refreshTasks()
  }

  dateForward() {
    this.todayDate.update((oldVal) => addMonth(oldVal, 1))
    this.refreshTasks()
  }

  dateBackward() {
    this.todayDate.update((oldVal) => addMonth(oldVal, -1))
    this.refreshTasks()
  }

  refreshTasks() {
    console.log(`Refreshing tasks, today: ${this.todayDate()}`)
    return this.taskService.listTasks(this.todayDate())
      .pipe(
        tap(tasks => {
          console.log("Setting tasks")
          this.taskList.set(tasks)
        })
      )
      .subscribe()
  }
}
