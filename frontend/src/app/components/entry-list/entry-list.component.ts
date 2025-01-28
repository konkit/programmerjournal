import {Component, computed, output, signal, ViewChild} from '@angular/core';
import {Entry, TaskSummary} from '../../../frontend-client';
import {FormControl, FormsModule, ReactiveFormsModule} from '@angular/forms';
import {CommonModule} from '@angular/common';
import {MatCardModule} from '@angular/material/card';
import {MatMenuModule} from '@angular/material/menu'
import {MatChipsModule} from '@angular/material/chips'
import {MatButtonModule} from '@angular/material/button'
import {MatIconModule} from '@angular/material/icon'
import {MatInputModule} from '@angular/material/input'
import {MatFormFieldModule} from '@angular/material/form-field'
import {getDayOfWeekFromDate, getMonthFromDate, getYearFromDate} from "../../../lib/wall_date";
import {CdkDrag, CdkDragDrop, CdkDragHandle, CdkDropList} from '@angular/cdk/drag-drop';
import {MatSelectModule} from '@angular/material/select';
import {MatDrawer, MatSidenavModule} from '@angular/material/sidenav';
import {EntrySidebarComponent} from '../entry-sidebar/entry-sidebar.component';
import {StatusButtonComponent} from '../status-button/status-button.component';
import {NavToolbarComponent} from '../nav-toolbar/nav-toolbar.component';
import {MatList, MatListItem} from '@angular/material/list';
import {EntryListService} from '../../service/entry-list.service';

export enum EditorStateEnum {
  IDLE,
  EDITING_NEW_TASK,
  EDITING_NEW_NOTE,
}

@Component({
  imports: [
    MatSidenavModule,
    MatFormFieldModule,
    MatSelectModule,
    MatButtonModule,
    CdkDropList,
    CdkDrag,
    CdkDragHandle,
    MatChipsModule,
    MatCardModule,
    MatButtonModule,
    FormsModule,
    CommonModule,
    MatButtonModule,
    MatMenuModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    EntrySidebarComponent,
    ReactiveFormsModule,
    StatusButtonComponent,
    NavToolbarComponent,
    MatList,
    MatListItem
  ],
  selector: 'app-entry-list',
  standalone: true,
  styleUrl: './entry-list.component.scss',
  templateUrl: './entry-list.component.html',
})
export class EntryListComponent {

  dateForward = output<void>()
  dateBackward = output<void>()

  todayDate = computed(() => this.entryListService.todayDate())
  entryPriority1 = computed(() => this.entryListService.entryList().filter(e => e.rank < 0))
  nonPriority = computed(() => this.entryListService.entryList().filter(x => x.rank >= 0))

  currentDateString = computed<string>(() => {
    let isMonthlyDate = this.todayDate().length == 7;
    if (isMonthlyDate) {
      return `Month: ${getMonthFromDate(this.todayDate())} ${getYearFromDate(this.todayDate())}`
    } else {
      return `Day: ${getDayOfWeekFromDate(this.todayDate())}, ${this.todayDate()}`
    }
  })

  @ViewChild('drawer') sideDrawer!: MatDrawer;
  editedTaskSummary = signal<TaskSummary | null>(null)

  EditorStateEnum = EditorStateEnum;
  editorState = EditorStateEnum.IDLE

  newEntryFormControl = new FormControl("");

  constructor(private entryListService: EntryListService) {
  }

  changeDateForward() {
    this.dateForward.emit()
  }

  changeDateBackward() {
    this.dateBackward.emit()
  }

  submitTitleEditWithValue(entry: Entry, e: Event) {
    let newValue = (e.target as HTMLDivElement).innerText
    this.entryListService.setTitle(newValue, entry);
  }

  setCreatingNewTaskState() {
    this.newEntryFormControl.setValue("")
    this.editorState = EditorStateEnum.EDITING_NEW_TASK;
  }

  setCreatingNewNoteState() {
    this.newEntryFormControl.setValue("")
    this.editorState = EditorStateEnum.EDITING_NEW_NOTE;
  }

  submitNewTask() {
    let taskValue = this.newEntryFormControl.value || "";
    return this.entryListService.createTask(taskValue)
      .subscribe(() => this.editorState = EditorStateEnum.IDLE)
  }

  submitNewNote() {
    let taskValue = this.newEntryFormControl.value || "";

    return this.entryListService.createNote(taskValue)
      .subscribe(() => this.editorState = EditorStateEnum.IDLE)
  }

  cancelEdit() {
    this.newEntryFormControl.setValue("")
    this.editorState = EditorStateEnum.IDLE;
  }

  markTaskAsDone(entry: Entry) {
    return this.entryListService.markTaskAsDone(entry.id).subscribe()
  }

  markTaskAsCreated(entry: Entry) {
    return this.entryListService.markTaskAsCreated(entry.id).subscribe()
  }

  snoozeTask(entry: Entry) {
    this.entryListService.snoozeTask(entry).subscribe()
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

  importPastTasks() {
    this.entryListService.importPastTasks().subscribe()
  }

  openUpdates(entry: Entry) {
    this.entryListService.getTaskSummary(entry.id)
      .subscribe((ts) => {
        console.log("openUpdates")
        this.editedTaskSummary.set(ts)
        this.sideDrawer.open()
      })
  }

  reloadTaskSummary(taskId: number) {
    this.entryListService.getTaskSummary(taskId)
      .subscribe((ts) => {
        console.log("reloadTaskSummary")
        this.editedTaskSummary.set(ts)
      })
  }

  deleteTaskFromSidebar(taskId: number) {
    return this.entryListService.deleteEntry(taskId)
      .subscribe(() => {
        this.sideDrawer.close()
      })
  }

  migrateToMonthly(entry: Entry) {
    this.entryListService.migrateToMonthly(entry).subscribe()
  }

  migrateToDaily(entry: Entry) {
    this.entryListService.migrateToDaily(entry).subscribe()
  }
}

