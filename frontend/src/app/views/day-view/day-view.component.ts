import {Component, OnInit} from '@angular/core';
import {AddDay, Today} from '../../../lib/wall_date';
import {EntryListComponent} from '../../components/entry-list/entry-list.component';
import {EntryListService} from '../../service/entry-list.service';

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

  constructor(private entryListService: EntryListService) {
  }

  ngOnInit() {
    this.entryListService.todayDate.set(Today())
    this.entryListService.refreshTasks().subscribe()
  }

  dateForward() {
    this.entryListService.todayDate.update((oldVal) => AddDay(oldVal, 1))
    this.entryListService.refreshTasks().subscribe()
  }

  dateBackward() {
    this.entryListService.todayDate.update((oldVal) => AddDay(oldVal, -1))
    this.entryListService.refreshTasks().subscribe()
  }
}
