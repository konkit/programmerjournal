import {Component, computed} from '@angular/core';
import {EntryListService} from '../../service/entry-list.service';
import {MatIcon} from '@angular/material/icon';
import {MatButton, MatIconButton} from '@angular/material/button';
import {getDayOfWeekFromDate, getMonthFromDate, getYearFromDate} from '../../../lib/wall_date';

@Component({
  selector: 'app-title-wrapper',
  imports: [
    MatIcon,
    MatIconButton,
    MatButton
  ],
  standalone: true,
  templateUrl: './title-wrapper.component.html',
  styleUrl: './title-wrapper.component.scss'
})
export class TitleWrapperComponent {

  todayDate = computed(() => this.entryListService.todayDate())

  currentDateString = computed<string>(() => {
    let isMonthlyDate = this.todayDate().length == 7;
    if (isMonthlyDate) {
      return `Month: ${getMonthFromDate(this.todayDate())} ${getYearFromDate(this.todayDate())}`
    } else {
      return `Day: ${getDayOfWeekFromDate(this.todayDate())}, ${this.todayDate()}`
    }
  })

  constructor(private entryListService: EntryListService) {
  }

  changeDateForward() {
    this.entryListService.dateForward().subscribe()
  }

  changeDateBackward() {
    this.entryListService.dateBackward().subscribe()
  }

  importPastTasks() {
    this.entryListService.importPastTasks().subscribe()
  }

}
