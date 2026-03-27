import {Component, input, OnInit, output, signal} from '@angular/core';
import {Entry, TaskService} from '../../../frontend-client';
import {CdkDrag, CdkDragHandle} from '@angular/cdk/drag-drop';
import {MatIcon} from '@angular/material/icon';
import {StatusButtonComponent} from '../status-button/status-button.component';
import {EntryListService} from '../../service/entry-list.service';
import {MatIconButton} from '@angular/material/button';
import {FormControl, ReactiveFormsModule} from '@angular/forms';
import {MatTooltip} from '@angular/material/tooltip';
import {EntryStatus, isTask} from '../../../lib/entry';
import {MarkdownPipe} from '../markdown.pipe';
import {RenderLinksPipe} from './render-links.pipe';

@Component({
  selector: 'app-entry',
  imports: [
    CdkDrag,
    MatIcon,
    StatusButtonComponent,
    CdkDragHandle,
    MatIconButton,
    ReactiveFormsModule,
    MatTooltip,
    MarkdownPipe,
    RenderLinksPipe
  ],
  templateUrl: './entry.component.html',
  standalone: true,
  styleUrl: './entry.component.scss'
})
export class EntryComponent implements OnInit {
  entry = input<Entry>()

  onOpenUpdatesSidebar = output<void>()

  fc = new FormControl("")
  updateFC = new FormControl("")


  isEdited = signal<boolean>(false)
  updatesVisible = signal<boolean>(false)
  isUpdateEdited = signal<boolean>(false)

  constructor(private entryListService: EntryListService, private taskService: TaskService) {
  }

  ngOnInit() {
    if (this.entry != null) {
      this.updatesVisible.set(this.entry()?.taskUpdate != null && this.entry()!.taskUpdate?.length > 0)
    }
  }

  markTaskAsDone() {
    return this.entryListService.markTaskAsDone(this.entry()!.id).subscribe()
  }

  markTaskAsCreated() {
    return this.entryListService.markTaskAsCreated(this.entry()!.id).subscribe()
  }

  snoozeTask() {
    this.entryListService.snoozeTask(this.entry()!).subscribe()
  }

  cancelTask() {
    this.entryListService.markTaskCancelled(this.entry()!.id).subscribe()
  }

  migrateToMonthly() {
    this.entryListService.migrateToMonthly(this.entry()!).subscribe()
  }

  migrateToWeekly() {
    this.entryListService.migrateToWeekly(this.entry()!).subscribe()
  }

  migrateToDaily() {
    this.entryListService.migrateToDaily(this.entry()!).subscribe()
  }

  deleteEntry() {
    this.entryListService.deleteEntry(this.entry()!.id).subscribe()
  }

  addUpdate() {
    this.updatesVisible.set(true)
    this.setUpdateEdited(true)
  }

  submitTitleEditWithValue(e: Event) {
    // console.log("newValue: ", this.fc.value)
    // let newValue = (e.target as HTMLDivElement).innerText
    this.entryListService.setTitle(this.fc.value || "", this.entry()!)
      .subscribe(() => {
        console.log("entryListService.setTitle subscribe")
        this.isEdited.set(false)
      })
  }

  submitUpdateEditWithValue(e: Event) {
    this.entryListService.setTaskUpdate(this.updateFC.value || "", this.entry()!)
      .subscribe(() => {
        console.log("taskService.setTaskUpdate subscribe")
        this.isUpdateEdited.set(false)
      })
  }

  setEdited(newValue: boolean) {
    this.fc.setValue(this.entry()?.title || "")
    this.isEdited.set(newValue);
  }

  setUpdateEdited(newValue: boolean) {
    this.updateFC.setValue(this.entry()?.taskUpdate || "")
    this.isUpdateEdited.set(newValue);
  }

  openUpdatesSidebar() {
    this.onOpenUpdatesSidebar.emit()
  }

  toggleUpdates() {
    this.updatesVisible.set( !this.updatesVisible() )
  }

  cancelEdit($event: MouseEvent) {
    this.fc.setValue("")
    this.isEdited.set(false);
  }

  cancelUpdateEdit($event: MouseEvent) {
    this.updateFC.setValue("")
    this.isUpdateEdited.set(false);
  }

  isTask() {
    return isTask(this.entry()?.status as EntryStatus)
  }
}
