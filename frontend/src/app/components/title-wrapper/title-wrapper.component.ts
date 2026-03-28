import { Component, computed, Input } from '@angular/core';
import { EntryListService } from '../../service/entry-list.service';
import { MatIcon } from '@angular/material/icon';
import { MatButton, MatIconButton } from '@angular/material/button';
import { MatBadge } from '@angular/material/badge';
import { getDayOfWeekFromDate, getMonthFromDate, getYearFromDate } from '../../../lib/wall_date';

@Component({
  selector: 'app-title-wrapper',
  imports: [
    MatIcon,
    MatIconButton,
    MatButton,
    MatBadge
  ],
  standalone: true,
  templateUrl: './title-wrapper.component.html',
  styleUrl: './title-wrapper.component.scss'
})
export class TitleWrapperComponent {

  @Input() date?: string;
  @Input() hideArrows: boolean = false;

  todayDate = computed(() => this.date || this.entryListService.todayDate())
  pendingImportsCount = computed(() => this.entryListService.pendingImportsCount())

  currentDateString = computed<string>(() => {
    let date = this.todayDate();
    let isMonthlyDate = date.length == 7;
    let isWeeklyDate = date.length == 8 && date.includes('W');

    if (isMonthlyDate) {
      return `Month: ${getMonthFromDate(date)} ${getYearFromDate(date)}`
    } else if (isWeeklyDate) {
      return `Week: ${date}`
    } else {
      return `Day: ${getDayOfWeekFromDate(date)}, ${date}`
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

}
