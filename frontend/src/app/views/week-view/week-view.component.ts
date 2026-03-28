import {Component, computed, effect, inject, input, OnInit, signal, ViewChild} from '@angular/core';
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
import {EntryListService} from '../../service/entry-list.service';
import {MatBadge} from '@angular/material/badge';
import {MatButton} from '@angular/material/button';
import {toWeeklyDate} from '../../../lib/wall_date';

@Component({
  selector: 'app-week-view',
  imports: [EntryComponent, EntrySidebarComponent, MatDrawer, MatDrawerContainer, MatDrawerContent, NavToolbarComponent, NewEntryComponent, TitleWrapperComponent, CdkDropList, MatBadge, MatButton],
  templateUrl: './week-view.component.html',
  standalone: true,
  styleUrl: './week-view.component.scss'
})
export class WeekViewComponent implements OnInit {

  selectedDate = input<string>();
  hideToolbar = input<boolean>(false);

  entryList = signal<Entry[]>([])

  pendingImportsCount = computed(() => this.entryListService.pendingImportsCount())

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)

  private readonly route = inject(ActivatedRoute);

  constructor(public entryListService: EntryListService,
              private entryService: EntryService,
              private taskService: TaskService) {
    effect(() => {
      this.refreshWeekTasks();
    });
  }

  ngOnInit() {
    if (!this.selectedDate()) {
      this.entryListService.todayDate.set(this.route.snapshot.params["date"])
    }
  }

  refreshWeekTasks() {
    const date = this.selectedDate() || this.entryListService.todayDate();
    if (!date) return;

    const weekString = toWeeklyDate(date);
    this.entryService.listEntries(weekString).subscribe(entries => {
      this.entryList.set(entries);
    });
  }

  handleDrop(e: CdkDragDrop<Entry[]>) {
      let targetIndex = e.currentIndex
      let currentRank = e.item.data;
      this.entryService.changeRank(this.entryList().find(e => e.rank === currentRank)!.id, { newRank: targetIndex })
        .subscribe(() => this.refreshWeekTasks());
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
    this.refreshWeekTasks();
  }

  deleteTaskFromSidebar(taskId: number) {
    return this.entryService.deleteEntry(taskId)
      .subscribe(() => {
        this.sideDrawer.close()
        this.refreshWeekTasks();
      })
  }

  importPastTasks() {
    const date = this.selectedDate() || this.entryListService.todayDate();
    const weekString = toWeeklyDate(date);
    this.taskService.importPastTasks(weekString).subscribe(() => this.refreshWeekTasks());
  }

  protected readonly toWeeklyDate = toWeeklyDate;
}
