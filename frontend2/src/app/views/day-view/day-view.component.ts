import {Component, OnInit, signal} from '@angular/core';
import {AddDay, Today} from '../../../lib/wall_date';
import {tap} from 'rxjs';
import {Entry, EntryService, TaskService} from '../../../frontend-client';
import {EntryListComponent} from '../../components/entry-list/entry-list.component';

@Component({
  selector: 'app-day-view',
  imports: [
    EntryListComponent
  ],
  templateUrl: './day-view.component.html',
  standalone: true,
  styleUrl: './day-view.component.scss'
})
export class DayViewComponent implements OnInit {

  todayDate = signal<string>(Today());

  entryList = signal<Entry[]>([]);

  constructor(private taskService: TaskService,
              private entryService: EntryService) {
  }

  ngOnInit() {
    console.log("Fetching task list")
    this.refreshTasks()
  }

  dateForward() {
    this.todayDate.update((oldVal) => AddDay(oldVal, 1))
    this.refreshTasks()
  }

  dateBackward() {
    this.todayDate.update((oldVal) => AddDay(oldVal, -1))
    this.refreshTasks()
  }

  refreshTasks() {
    console.log("Refreshing tasks")
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
