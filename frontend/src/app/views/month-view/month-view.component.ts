import {Component, computed, signal, ViewChild} from '@angular/core';
import {ThisMonth} from '../../../lib/wall_date';
import {EntryListService} from '../../service/entry-list.service';
import {CdkDragDrop, CdkDropList} from '@angular/cdk/drag-drop';
import {EntryComponent} from '../../components/entry-list/entry/entry.component';
import {EntrySidebarComponent} from '../../components/entry-sidebar/entry-sidebar.component';
import {MatDrawer, MatDrawerContainer, MatDrawerContent} from '@angular/material/sidenav';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {NewEntryComponent} from '../../components/entry-list/new-entry/new-entry.component';
import {TitleWrapperComponent} from '../../components/entry-list/title-wrapper/title-wrapper.component';
import {TaskSummary} from '../../../frontend-client';

@Component({
  selector: 'app-month-view',
  imports: [CdkDropList, EntryComponent, EntrySidebarComponent, MatDrawer, MatDrawerContainer, MatDrawerContent, NavToolbarComponent, NewEntryComponent, TitleWrapperComponent],
  templateUrl: './month-view.component.html',
  standalone: true,
  styleUrl: './month-view.component.scss'
})
export class MonthViewComponent {

  entryList = computed(() => this.entryListService.entryList())

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)

  constructor(private entryListService: EntryListService) {
  }

  ngOnInit() {
    this.entryListService.todayDate.set(ThisMonth())
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
