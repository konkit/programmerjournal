import {Component, computed, inject, OnInit, signal, ViewChild} from '@angular/core';
import {EntryListService} from '../../service/entry-list.service';
import {MatDrawer, MatDrawerContainer, MatDrawerContent} from '@angular/material/sidenav';
import {CdkDragDrop, CdkDropList} from '@angular/cdk/drag-drop';
import {Entry, TaskSummary} from '../../../frontend-client';
import {EntryComponent} from '../../components/entry/entry.component';
import {EntrySidebarComponent} from '../../components/entry-sidebar/entry-sidebar.component';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {NewEntryComponent} from '../../components/new-entry/new-entry.component';
import {TitleWrapperComponent} from '../../components/title-wrapper/title-wrapper.component';
import {ActivatedRoute} from '@angular/router';
import {MatButton} from '@angular/material/button';

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
    MatButton
  ],
  templateUrl: './day-view.component.html',
  standalone: true,
  styleUrl: './day-view.component.scss'
})
export class DayViewComponent implements OnInit {

  entryPriority1 = computed(() => this.entryListService.entryList().filter(e => e.rank < 0))
  nonPriority = computed(() => this.entryListService.entryList().filter(x => x.rank >= 0))

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  @ViewChild('weeklyDrawer') weeklyDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)
  weeklyEntryList = signal<Entry[]>([])

  private readonly route = inject(ActivatedRoute);

  constructor(private entryListService: EntryListService) {
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
    // Update task summary
    this.entryListService.getTaskSummary(taskId)
      .subscribe((ts) => {this.editedTaskSummary.set(ts)})

    // Update task list, e.g. to update the status icons
    this.entryListService.refreshTasks().subscribe()
  }

  deleteTaskFromSidebar(taskId: number) {
    return this.entryListService.deleteEntry(taskId)
      .subscribe(() => {this.sideDrawer.close()})
  }

  openWeeklySidebar() {
    this.entryListService.getWeeklyTasks().subscribe((entries) => {
      this.weeklyEntryList.set(entries)
      this.weeklyDrawer.open()
    })
  }
}
