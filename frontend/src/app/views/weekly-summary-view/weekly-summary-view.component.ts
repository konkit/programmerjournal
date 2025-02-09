import {Component, computed, inject, OnInit, signal} from '@angular/core';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {addDay, StartOfThisWeek} from '../../../lib/wall_date';
import {EntryService, TaskService, WeeklySummary} from '../../../frontend-client';
import {MarkdownPipe} from '../../components/markdown.pipe';
import {MatIcon} from '@angular/material/icon';
import {MatIconButton} from '@angular/material/button';
import {MatAccordion, MatExpansionModule, MatExpansionPanel} from '@angular/material/expansion';
import {ActivatedRoute} from '@angular/router';

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
export class WeeklySummaryViewComponent implements OnInit {

  weekStartDate = signal<string>(StartOfThisWeek());

  summary = signal<WeeklySummary>({taskSummaries: [], notes: []});

  currentDateString = computed<string>(() => {
    return `Summary of the week ${this.weekStartDate()} - ${addDay(this.weekStartDate(), 6)}`
  })

  private readonly route = inject(ActivatedRoute);

  constructor(private taskService: TaskService, private entryService: EntryService) {
  }

  ngOnInit() {
    this.weekStartDate.set(this.route.snapshot.params["date"])
    this.refreshTasks()
  }

  dateForward() {
    this.weekStartDate.update((oldVal) => addDay(oldVal, 7))
    this.refreshTasks()
  }

  dateBackward() {
    this.weekStartDate.update((oldVal) => addDay(oldVal, -7))
    this.refreshTasks()
  }

  refreshTasks() {
    return this.entryService.weeklySummary(this.weekStartDate())
      .subscribe(summary => this.summary.set(summary))
  }

}
