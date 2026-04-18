import { Component, computed, effect, inject, OnInit, signal, ViewChild } from '@angular/core';
import { EntryListService } from '../../service/entry-list.service';
import { MatDrawer, MatDrawerContainer, MatDrawerContent } from '@angular/material/sidenav';
import { CdkDragDrop, CdkDropList } from '@angular/cdk/drag-drop';
import { Entry, TaskSummary } from '../../../frontend-client';
import { EntryComponent } from '../../components/entry/entry.component';
import { EntrySidebarComponent } from '../../components/entry-sidebar/entry-sidebar.component';
import { NavToolbarComponent } from '../../components/nav-toolbar/nav-toolbar.component';
import { NewEntryComponent } from '../../components/new-entry/new-entry.component';
import { TitleWrapperComponent } from '../../components/title-wrapper/title-wrapper.component';
import { ActivatedRoute } from '@angular/router';
import { MatButton } from '@angular/material/button';
import { RenderLinksPipe } from '../../components/entry/render-links.pipe';
import { MatTooltip } from '@angular/material/tooltip';
import { StatusIconComponent } from '../../components/status-icon/status-icon.component';
import { MatIcon } from '@angular/material/icon';
import { MatMenu, MatMenuItem, MatMenuTrigger } from '@angular/material/menu';
import { EntryStatus } from '../../../lib/entry';
import { switchMap, tap } from 'rxjs';
import { MatBadge } from '@angular/material/badge';
import { WeekViewComponent } from '../week-view/week-view.component';
import { getWeekString } from '../../../lib/wall_date';

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
    MatButton,
    RenderLinksPipe,
    MatTooltip,
    StatusIconComponent,
    MatIcon,
    MatMenu,
    MatMenuItem,
    MatMenuTrigger,
    MatBadge,
    WeekViewComponent
  ],
  templateUrl: './day-view.component.html',
  standalone: true,
  styleUrl: './day-view.component.scss'
})
export class DayViewComponent implements OnInit {

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  @ViewChild('weeklyDrawer') weeklyDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)
  weeklyEntryList = signal<Entry[]>([])

  weeklyDateString = computed(() => getWeekString(new Date(this.entryListService.todayDate())))

  pendingImportsCount = computed(() => this.entryListService.pendingImportsCount())

  private readonly route = inject(ActivatedRoute);

  constructor(public entryListService: EntryListService) {
    effect(() => {
      this.entryListService.refreshTrigger();
      this.entryListService.getWeeklyTasks().subscribe((entries) => {
        this.weeklyEntryList.set(entries)
      })
    });
  }

  ngOnInit() {
    console.log("ngOnInit");
    this.entryListService.todayDate.set(this.route.snapshot.params["date"])
    this.entryListService.refreshTasks().subscribe()
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

  openWeeklySidebar() {
    this.entryListService.getWeeklyTasks().subscribe((entries) => {
      this.weeklyEntryList.set(entries)
      this.weeklyDrawer.open()
    })
  }

  readonly EntryStatus = EntryStatus;

  migrateWeeklyToToday(weeklyEntry: Entry) {
    return this.entryListService.migrateWeeklyTaskToToday(weeklyEntry.id)
      .subscribe()
  }

  importPastTasks() {
    this.entryListService.importPastTasks().subscribe()
  }
}
