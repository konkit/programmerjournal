import {Component, computed, inject, signal, ViewChild} from '@angular/core';
import {CdkDragDrop, CdkDropList} from '@angular/cdk/drag-drop';
import {EntrySidebarComponent} from '../../components/entry-sidebar/entry-sidebar.component';
import {MatDrawer, MatDrawerContainer, MatDrawerContent} from '@angular/material/sidenav';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {Entry, TaskService, TaskSummary} from '../../../frontend-client';
import {EntryComponent} from '../../components/entry/entry.component';
import {NewEntryComponent} from '../../components/new-entry/new-entry.component';
import {TitleWrapperComponent} from '../../components/title-wrapper/title-wrapper.component';
import {ActivatedRoute} from '@angular/router';
import {EntryService} from '../../../frontend-client/api/entry.service';
import {Today} from '../../../lib/wall_date';
import {switchMap, tap} from 'rxjs';
import {WeeklyEntryListService} from '../../service/weekly-entry-list.service';

@Component({
  selector: 'app-week-view',
  imports: [EntryComponent, EntrySidebarComponent, MatDrawer, MatDrawerContainer, MatDrawerContent, NavToolbarComponent, NewEntryComponent, TitleWrapperComponent, CdkDropList],
  templateUrl: './week-view.component.html',
  standalone: true,
  styleUrl: './week-view.component.scss'
})
export class WeekViewComponent {

  weeklyEntryList = computed(() => this.entryListService.weeklyEntryList())
  dailyEntryList = computed(() => this.entryListService.dailyEntryList())

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)

  private readonly route = inject(ActivatedRoute);

  constructor(public entryListService: WeeklyEntryListService,
              private entryService: EntryService,
              private taskService: TaskService) {
  }

  ngOnInit() {
    this.entryListService.todayDate.set(this.route.snapshot.params["date"])
    this.entryListService.refreshTasks().subscribe()
  }

  handleDrop(e: CdkDragDrop<Entry[]>) {
    if (e.previousContainer === e.container) {
      // Reordering within the same list
      let targetIndex = e.currentIndex
      let currentRank = e.item.data;
      // If it's the main weekly list
      if (e.container.id === 'weeklyList') {
         this.entryListService.handleDropWeekly(targetIndex, currentRank).subscribe()
      } else {
         this.entryListService.handleDropDaily(targetIndex, currentRank).subscribe()
      }
    } else {
      // Moving between lists
      if (e.previousContainer.id === 'weeklyList' && e.container.id === 'dailyList') {
          // We need to find the entry in the source list
          const entryToMove = this.weeklyEntryList()[e.previousIndex];
          // Call API to migrate to daily log
          this.taskService.migrateTaskToDailyLog(entryToMove.id, { date: Today() })
            .pipe(
              switchMap(() => this.entryListService.refreshTasks()), // Refresh weekly list
            ).subscribe()

      } else if (e.previousContainer.id === 'dailyList' && e.container.id === 'weeklyList') {
          // Moving from Today to Weekly
          // We need to find the entry in the source list
          const entryToMove = this.dailyEntryList()[e.previousIndex];

          this.taskService.migrateTaskToWeeklyLog(entryToMove.id, { date: this.entryListService.todayDate() })
            .pipe(
               switchMap(() => this.entryListService.refreshTasks()),
            ).subscribe()
      }
    }
  }

  openUpdates(entryId: number) {
    this.entryListService.getTaskSummary(entryId)
      .subscribe((ts) => {
        this.editedTaskSummary.set(ts)
        this.sideDrawer.open()
      })
  }

  reloadTaskSummary(taskId: number) {
    this.entryListService.getTaskSummary(taskId)
      .subscribe((ts) => {
        this.editedTaskSummary.set(ts)
      })
  }

  deleteTaskFromSidebar(taskId: number) {
    return this.entryListService.deleteEntry(taskId)
      .subscribe(() => {
        this.sideDrawer.close()
      })
  }
}
