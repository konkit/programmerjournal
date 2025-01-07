import {Component, signal} from '@angular/core';
import {addMonth, ThisMonth, WallMonth} from '../../../lib/wall_date';
import {Entry, EntryService, TaskService} from '../../../frontend-client';
import {tap} from 'rxjs';
import {EntryListComponent} from '../../components/entry-list/entry-list.component';

@Component({
  selector: 'app-month-view',
  imports: [EntryListComponent],
  templateUrl: './month-view.component.html',
  standalone: true,
  styleUrl: './month-view.component.scss'
})
export class MonthViewComponent {

  todayDate = signal<WallMonth>(ThisMonth());

  entryList = signal<Entry[]>([]);

  constructor(private taskService: TaskService,
              private entryService: EntryService) {
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
    return this.entryService.listEntries(this.todayDate())
      .pipe(
        tap(entries => {
          console.log("Setting tasks")
          this.entryList.set(entries)
        })
      )
      .subscribe()
  }
}
