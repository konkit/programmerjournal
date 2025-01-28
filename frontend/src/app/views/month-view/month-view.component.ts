import {Component} from '@angular/core';
import {addMonth, ThisMonth} from '../../../lib/wall_date';
import {EntryListComponent} from '../../components/entry-list/entry-list.component';
import {EntryListService} from '../../service/entry-list.service';

@Component({
  selector: 'app-month-view',
  imports: [EntryListComponent],
  templateUrl: './month-view.component.html',
  standalone: true,
  styleUrl: './month-view.component.scss'
})
export class MonthViewComponent {

  constructor(private entryListService: EntryListService) {
  }

  ngOnInit() {
    this.entryListService.todayDate.set(ThisMonth())
    this.entryListService.refreshTasks().subscribe()
  }

  dateForward() {
    this.entryListService.todayDate.update((oldVal) => addMonth(oldVal, 1))
    this.entryListService.refreshTasks().subscribe()
  }

  dateBackward() {
    this.entryListService.todayDate.update((oldVal) => addMonth(oldVal, -1))
    this.entryListService.refreshTasks().subscribe()
  }
}
