import {Component, computed, signal} from '@angular/core';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {
  AddDay, addMonth,
  getDayOfWeekFromDate,
  getMonthFromDate,
  getYearFromDate,
  StartOfThisWeek,
  Today
} from '../../../lib/wall_date';
import {Task} from '../../../lib/task';
import {TaskService, TaskWeeklySummary} from '../../../frontend-client';
import {tap} from 'rxjs';
import {JsonPipe} from '@angular/common';
import {MarkdownPipe} from '../../components/markdown.pipe';
import {MatIcon} from '@angular/material/icon';
import {MatIconButton} from '@angular/material/button';
import {MatAccordion, MatExpansionModule, MatExpansionPanel} from '@angular/material/expansion';

@Component({
  selector: 'app-weekly-summary-view',
  imports: [
    NavToolbarComponent,
    MarkdownPipe,
    MatIcon,
    MatIconButton,
    MatAccordion,
    MatExpansionPanel,
    MatExpansionModule,
],
  templateUrl: './weekly-summary-view.component.html',
  standalone: true,
  styleUrl: './weekly-summary-view.component.scss'
})
export class WeeklySummaryViewComponent {

  weekStartDate = signal<string>(StartOfThisWeek());

  taskSummaryList = signal<TaskWeeklySummary[]>([]);

  currentDateString = computed<string>(() => {
    return `Week ${this.weekStartDate()} - ${AddDay(this.weekStartDate(), 6)}`
  })

  constructor(private taskService: TaskService) {
  }

  ngOnInit() {
    this.refreshTasks()
  }

  dateForward() {
    this.weekStartDate.update((oldVal) => AddDay(oldVal, 7))
    this.refreshTasks()
  }

  dateBackward() {
    this.weekStartDate.update((oldVal) => AddDay(oldVal, -7))
    this.refreshTasks()
  }

  refreshTasks() {
    return this.taskService.weeklySummary(this.weekStartDate())
      .subscribe(tasks => this.taskSummaryList.set(tasks))
  }

}
