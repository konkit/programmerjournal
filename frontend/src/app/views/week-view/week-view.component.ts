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
import {EntryListService} from '../../service/entry-list.service';
import {MatBadge} from '@angular/material/badge';
import {MatButton} from '@angular/material/button';

@Component({
  selector: 'app-week-view',
  imports: [EntryComponent, EntrySidebarComponent, MatDrawer, MatDrawerContainer, MatDrawerContent, NavToolbarComponent, NewEntryComponent, TitleWrapperComponent, CdkDropList, MatBadge, MatButton],
  templateUrl: './week-view.component.html',
  standalone: true,
  styleUrl: './week-view.component.scss'
})
export class WeekViewComponent {

  entryList = computed(() => this.entryListService.entryList())

  pendingImportsCount = computed(() => this.entryListService.pendingImportsCount())

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)

  private readonly route = inject(ActivatedRoute);

  constructor(public entryListService: EntryListService,
              private entryService: EntryService,
              private taskService: TaskService) {
  }

  ngOnInit() {
    this.entryListService.todayDate.set(this.route.snapshot.params["date"])
    this.entryListService.refreshTasks().subscribe()
  }

  handleDrop(e: CdkDragDrop<Entry[]>) {
      let targetIndex = e.currentIndex
      let currentRank = e.item.data;
      this.entryListService.handleDrop(targetIndex, currentRank).subscribe();
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

  importPastTasks() {
    this.entryListService.importPastTasks().subscribe()
  }
}
