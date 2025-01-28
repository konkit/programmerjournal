import {Component, computed, OnInit, output, signal, ViewChild} from '@angular/core';
import {Today} from '../../../lib/wall_date';
import {EntryListService} from '../../service/entry-list.service';
import {MatDrawer, MatDrawerContainer, MatDrawerContent} from '@angular/material/sidenav';
import {CdkDragDrop, CdkDropList} from '@angular/cdk/drag-drop';
import {TaskSummary} from '../../../frontend-client';
import {EntryComponent} from '../../components/entry-list/entry/entry.component';
import {EntrySidebarComponent} from '../../components/entry-sidebar/entry-sidebar.component';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {NewEntryComponent} from '../../components/entry-list/new-entry/new-entry.component';
import {TitleWrapperComponent} from '../../components/entry-list/title-wrapper/title-wrapper.component';

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
    TitleWrapperComponent
  ],
  templateUrl: './day-view.component.html',
  standalone: true,
  styleUrl: './day-view.component.scss'
})
export class DayViewComponent implements OnInit {

  entryPriority1 = computed(() => this.entryListService.entryList().filter(e => e.rank < 0))
  nonPriority = computed(() => this.entryListService.entryList().filter(x => x.rank >= 0))

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)

  constructor(private entryListService: EntryListService) {
  }

  ngOnInit() {
    this.entryListService.todayDate.set(Today())
    this.entryListService.refreshTasks().subscribe()
  }

  handleDrop(e: CdkDragDrop<string[]>) {
    let targetIndex = e.currentIndex
    let currentRank = e.item.data;
    this.entryListService.handleDrop(targetIndex, currentRank).subscribe()
  }

  handleDropToPriority(e: CdkDragDrop<string[]>) {
    let targetIndex = e.currentIndex
    let currentRank = e.item.data
    this.entryListService.handleDropToPriority(targetIndex, currentRank).subscribe()
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
      .subscribe((ts) => {this.editedTaskSummary.set(ts)})
  }

  deleteTaskFromSidebar(taskId: number) {
    return this.entryListService.deleteEntry(taskId)
      .subscribe(() => {this.sideDrawer.close()})
  }
}
