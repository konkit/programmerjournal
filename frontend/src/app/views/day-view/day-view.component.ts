import { CdkDragDrop, CdkDropList } from '@angular/cdk/drag-drop';
import { NgTemplateOutlet } from '@angular/common';
import { Component, computed, effect, inject, OnInit, signal, ViewChild } from '@angular/core';
import { MatIconButton } from '@angular/material/button';
import { MatIcon } from '@angular/material/icon';
import { MatMenu, MatMenuItem, MatMenuTrigger } from '@angular/material/menu';
import { MatDrawer, MatDrawerContainer, MatDrawerContent } from '@angular/material/sidenav';
import { MatTooltip } from '@angular/material/tooltip';
import { ActivatedRoute } from '@angular/router';
import { Entry, TaskSummary } from '../../../frontend-client';
import { EntryStatus } from '../../../lib/entry';
import { getWeekNumber, toWeeklyDate } from '../../../lib/wall_date';
import { EntrySidebarComponent } from '../../components/entry-sidebar/entry-sidebar.component';
import { EntryComponent } from '../../components/entry/entry.component';
import { RenderLinksPipe } from '../../components/entry/render-links.pipe';
import { NavToolbarComponent } from '../../components/nav-toolbar/nav-toolbar.component';
import { NewEntryComponent } from '../../components/new-entry/new-entry.component';
import { StatusIconComponent } from '../../components/status-icon/status-icon.component';
import { TitleWrapperComponent } from '../../components/title-wrapper/title-wrapper.component';
import { EntryListService } from '../../service/entry-list.service';

const WEEKLY_COLLAPSED_KEY = 'dayView.weeklyCollapsed';

function readWeeklyCollapsed(): boolean {
  try {
    return localStorage.getItem(WEEKLY_COLLAPSED_KEY) === 'true';
  } catch {
    return false;
  }
}

@Component({
  selector: 'app-day-view',
  imports: [
    CdkDropList,
    EntryComponent,
    EntrySidebarComponent,
    MatDrawer,
    MatDrawerContainer,
    MatDrawerContent,
    NavToolbarComponent,
    NewEntryComponent,
    TitleWrapperComponent,
    MatIconButton,
    RenderLinksPipe,
    MatTooltip,
    StatusIconComponent,
    MatIcon,
    MatMenu,
    MatMenuItem,
    MatMenuTrigger,
    NgTemplateOutlet,
  ],
  templateUrl: './day-view.component.html',
  standalone: true,
  styleUrl: './day-view.component.scss'
})
export class DayViewComponent implements OnInit {

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)
  weeklyEntryList = signal<Entry[]>([])
  weeklyDoneEntryList = signal<Entry[]>([])

  weeklyDateString = computed(() => toWeeklyDate(this.entryListService.todayDate()))
  weeklyLabel = computed(() => `Week ${getWeekNumber(this.weeklyDateString())}`)
  weeklyCollapsed = signal<boolean>(readWeeklyCollapsed())

  private readonly route = inject(ActivatedRoute);

  constructor(public entryListService: EntryListService) {
    effect(() => {
      this.entryListService.refreshTrigger();
      // The weekly list is fetched again once the rail is expanded
      if (this.weeklyCollapsed()) return;
      this.entryListService.getWeeklyTasks().subscribe((response) => {
        this.weeklyEntryList.set(response.pending || [])
        this.weeklyDoneEntryList.set(response.done || [])
      })
    });
  }

  ngOnInit() {
    this.entryListService.todayDate.set(this.route.snapshot.params["date"])
    this.entryListService.loadDate().subscribe()
  }

  handleDrop(e: CdkDragDrop<string[]>) {
    let targetIndex = e.currentIndex
    let currentRank = e.item.data;
    this.entryListService.handleDrop(targetIndex, currentRank).subscribe()
  }

  openUpdates(entryId: number) {
    this.entryListService.getTaskSummary(entryId)
      .subscribe((ts) => {
        this.editedTaskSummary.set(ts)
        this.sideDrawer.open()
      })
  }

  reloadTaskSummary(taskId: number) {
    // Update task summary
    this.entryListService.getTaskSummary(taskId)
      .subscribe((ts) => { this.editedTaskSummary.set(ts) })

    // Update task list, e.g. to update the status icons
    this.entryListService.refreshTasks().subscribe()
  }

  deleteTaskFromSidebar(taskId: number) {
    return this.entryListService.deleteEntry(taskId)
      .subscribe(() => { this.sideDrawer.close() })
  }

  readonly EntryStatus = EntryStatus;

  migrateWeeklyToToday(weeklyEntry: Entry) {
    return this.entryListService.migrateWeeklyTaskToToday(weeklyEntry.id)
      .subscribe()
  }

  toggleWeekly() {
    this.weeklyCollapsed.update(collapsed => !collapsed)
    try {
      localStorage.setItem(WEEKLY_COLLAPSED_KEY, String(this.weeklyCollapsed()))
    } catch {
      // Storage unavailable (e.g. private mode) - the choice just isn't remembered
    }
  }
}
