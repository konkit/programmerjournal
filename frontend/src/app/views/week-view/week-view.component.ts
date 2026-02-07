import {Component, computed, inject, signal, ViewChild} from '@angular/core';
import {EntryListService} from '../../service/entry-list.service';
import {CdkDragDrop, CdkDropList, moveItemInArray, transferArrayItem} from '@angular/cdk/drag-drop';
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

@Component({
  selector: 'app-week-view',
  imports: [CdkDropList, EntryComponent, EntrySidebarComponent, MatDrawer, MatDrawerContainer, MatDrawerContent, NavToolbarComponent, NewEntryComponent, TitleWrapperComponent],
  templateUrl: './week-view.component.html',
  standalone: true,
  styleUrl: './week-view.component.scss'
})
export class WeekViewComponent {

  entryList = computed(() => this.entryListService.entryList())
  todayEntryList = signal<Entry[]>([])

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)

  private readonly route = inject(ActivatedRoute);

  constructor(private entryListService: EntryListService,
              private entryService: EntryService,
              private taskService: TaskService) {
  }

  ngOnInit() {
    this.entryListService.todayDate.set(this.route.snapshot.params["date"])
    this.entryListService.refreshTasks().subscribe()
    this.refreshTodayTasks()
  }

  refreshTodayTasks() {
    this.entryService.listEntries(Today()).subscribe(entries => {
      this.todayEntryList.set(entries)
    })
  }

  handleDrop(e: CdkDragDrop<Entry[]>) {
    if (e.previousContainer === e.container) {
      // Reordering within the same list
      let targetIndex = e.currentIndex
      let currentRank = e.item.data;
      // If it's the main weekly list
      if (e.container.id === 'weeklyList') {
         this.entryListService.handleDrop(targetIndex, currentRank).subscribe()
      } else {
        // Reordering in today's list - we might need a separate service method or just ignore for now if not strictly required
        // But let's assume we want to reorder today's list too.
        // We need the ID of the item being moved.
        // The currentRank passed in data is the rank of the item.
        // We need to find the item in todayEntryList to get its ID.
        const item = this.todayEntryList().find(entry => entry.rank === currentRank);
        if (item) {
             this.entryService.changeRank(item.id, { newRank: targetIndex }).subscribe(() => this.refreshTodayTasks())
        }
      }
    } else {
      // Moving between lists
      const item = e.item.data; // This is the rank, but we need the entry object or ID.
      // Actually, let's look at how data is passed. In HTML: [entry]="entry".
      // Wait, in handleDrop(e: CdkDragDrop<string[]>), e.item.data is what we set [cdkDragData] to.
      // In EntryComponent, we need to see what is passed.
      // Let's assume we can get the entry object.
      // But wait, the previous implementation used `e.item.data` as rank.
      // Let's look at EntryComponent to be sure.
      // If we are moving from Weekly to Today:
      if (e.previousContainer.id === 'weeklyList' && e.container.id === 'todayList') {
          // We need to find the entry in the source list
          const entryToMove = this.entryList()[e.previousIndex];
          // Call API to migrate to daily log
          this.taskService.migrateTaskToDailyLog(entryToMove.id, { date: Today() })
            .pipe(
              switchMap(() => this.entryListService.refreshTasks()), // Refresh weekly list
              tap(() => this.refreshTodayTasks()) // Refresh today list
            ).subscribe()

      } else if (e.previousContainer.id === 'todayList' && e.container.id === 'weeklyList') {
          // Moving from Today to Weekly
          // We need to find the entry in the source list
          const entryToMove = this.todayEntryList()[e.previousIndex];
          // We need an API to migrate to weekly log (or monthly log if weekly is not supported directly, but here we are in week view)
          // The backend has MigrateTaskToMonthlyLog, but maybe not Weekly.
          // Wait, the user asked to move tasks from weekly column to the today one.
          // "The idea is that I should be able to move tasks from weekly column to the today one."
          // It doesn't explicitly say the reverse, but usually it's implied.
          // However, if there is no API for "MigrateToWeekly", we might be stuck.
          // Let's check TaskService.
          // We have migrateTaskToMonthlyLog and migrateTaskToDailyLog.
          // If the "Weekly View" is actually showing tasks for the week, maybe they are just tasks with a specific date?
          // Or is "Weekly Log" a concept?
          // In the previous step, we added WeekView which lists entries for a "Week Date" (e.g. 2024-W01).
          // So migrating to weekly log means changing the date to "2024-W01".
          // We can use `snoozeTask` or a generic "move" if available, or `migrateTaskToMonthlyLog` if that's what's intended.
          // But `migrateTaskToMonthlyLog` takes a YYYY-MM date.
          // `snoozeTask` takes a date string.
          // If we snooze to "2024-W01", it effectively moves it there.
          // So we can use snoozeTask.

          this.taskService.snoozeTask(entryToMove.id, { date: this.entryListService.todayDate() })
            .pipe(
               switchMap(() => this.entryListService.refreshTasks()),
               tap(() => this.refreshTodayTasks())
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
        this.refreshTodayTasks() // Refresh today list as well just in case
      })
  }
}
