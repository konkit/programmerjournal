import {Component, computed} from '@angular/core';
import {EntryListService} from '../../service/entry-list.service';
import {MatIcon} from '@angular/material/icon';
import {MatButton, MatIconButton} from '@angular/material/button';
import {MatTooltip} from '@angular/material/tooltip';
import {
  formatDayTitle,
  getDateFromWeek,
  getMonthFromDate,
  getWeekNumber,
  getYearFromDate,
  isMonthDate,
  isWeekDate,
  ThisMonth,
  ThisWeek,
  Today
} from '../../../lib/wall_date';

type LogKind = 'day' | 'week' | 'month';

const SHORT_MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]

@Component({
  selector: 'app-title-wrapper',
  imports: [
    MatIcon,
    MatIconButton,
    MatButton,
    MatTooltip
  ],
  standalone: true,
  templateUrl: './title-wrapper.component.html',
  styleUrl: './title-wrapper.component.scss'
})
export class TitleWrapperComponent {

  date = computed(() => this.entryListService.todayDate())
  pendingImportsCount = computed(() => this.entryListService.pendingImportsCount())

  kind = computed<LogKind>(() => {
    const date = this.date()
    if (isMonthDate(date)) return 'month'
    if (isWeekDate(date)) return 'week'
    return 'day'
  })

  title = computed<string>(() => {
    const date = this.date()
    switch (this.kind()) {
      case 'month':
        return `${getMonthFromDate(date)} ${getYearFromDate(date)}`
      case 'week':
        return `Week ${getWeekNumber(date)}`
      default:
        return formatDayTitle(date)
    }
  })

  // "19 – 25 Oct 2026" for weeks, nothing for days and months
  subtitle = computed<string>(() => {
    if (this.kind() !== 'week') return ''
    const date = this.date()
    const monday = getDateFromWeek(getYearFromDate(date), getWeekNumber(date))
    const sunday = new Date(monday.getFullYear(), monday.getMonth(), monday.getDate() + 6)
    return formatRange(monday, sunday)
  })

  isCurrent = computed(() => {
    const current = {day: Today(), week: ThisWeek(), month: ThisMonth()}[this.kind()]
    return this.date() === current
  })

  currentLabel = computed(() => ({day: 'Today', week: 'This week', month: 'This month'})[this.kind()])
  previousLabel = computed(() => `Previous ${this.kind()}`)
  nextLabel = computed(() => `Next ${this.kind()}`)

  constructor(private entryListService: EntryListService) {
  }

  changeDateForward() {
    this.entryListService.dateForward().subscribe()
  }

  changeDateBackward() {
    this.entryListService.dateBackward().subscribe()
  }

  goToCurrent() {
    this.entryListService.goToCurrent().subscribe()
  }

  importPastTasks() {
    this.entryListService.importPastTasks().subscribe()
  }
}

function formatRange(start: Date, end: Date): string {
  const endPart = `${end.getDate()} ${SHORT_MONTHS[end.getMonth()]} ${end.getFullYear()}`
  if (start.getFullYear() !== end.getFullYear()) {
    return `${start.getDate()} ${SHORT_MONTHS[start.getMonth()]} ${start.getFullYear()} – ${endPart}`
  }
  if (start.getMonth() !== end.getMonth()) {
    return `${start.getDate()} ${SHORT_MONTHS[start.getMonth()]} – ${endPart}`
  }
  return `${start.getDate()} – ${endPart}`
}
