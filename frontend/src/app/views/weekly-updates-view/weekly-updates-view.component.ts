import {Component, computed, inject, OnInit, signal} from '@angular/core';
import {Entry, EntryService, WeeklySummary} from '../../../frontend-client';
import {ActivatedRoute, Router} from '@angular/router';
import {addDay, StartOfThisWeek} from '../../../lib/wall_date';
import {MatIcon} from '@angular/material/icon';
import {MatAccordion, MatExpansionModule, MatExpansionPanel} from '@angular/material/expansion';
import {MatTooltip} from '@angular/material/tooltip';
import {MarkdownPipe} from '../../components/markdown.pipe';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {MatIconButton} from '@angular/material/button';
import {JsonPipe} from '@angular/common';

let emptySummary = {
  taskSummaries: [],
  notes: [],
  tasksUpdatedThisWeek: [],
  tasksFinishedThisWeek: [],
  otherTasks: []
};

@Component({
  selector: 'app-weekly-updates-view',
  imports: [
    NavToolbarComponent,
    MarkdownPipe,
    MatIcon,
    MatIconButton,
    MatAccordion,
    MatExpansionPanel,
    MatExpansionModule,
    MatTooltip,
    JsonPipe,
  ],
  templateUrl: './weekly-updates-view.component.html',
  styleUrl: './weekly-updates-view.component.scss',
  standalone: true
})
export class WeeklyUpdatesViewComponent implements OnInit {

  private readonly entryService = inject(EntryService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router)

  weekStartDate = signal<string>(StartOfThisWeek());
  summary = signal<{ [key: string]: Entry[] | null}>({});
  currentDateString = computed<string>(() => {
    return `Updates in week ${this.weekStartDate()} - ${addDay(this.weekStartDate(), 6)}`
  })

  ngOnInit() {
    this.weekStartDate.set(this.route.snapshot.params["date"])
    this.refreshTasks()
  }

  dateForward() {
    this.weekStartDate.update((oldVal) => addDay(oldVal, 7))
    this.router.navigate(['/weekUpdates', this.weekStartDate()]);
    this.refreshTasks()
  }

  dateBackward() {
    this.weekStartDate.update((oldVal) => addDay(oldVal, -7))
    this.router.navigate(['/weekUpdates', this.weekStartDate()]);
    this.refreshTasks()
  }

  refreshTasks() {
    return this.entryService.weeklyUpdates(this.weekStartDate())
      .subscribe(summary => this.summary.set(summary))
    // return this.entryService.weeklySummary(this.weekStartDate())
    //   .subscribe(summary => this.summary.set(summary))
  }

  protected readonly Object = Object;
}
