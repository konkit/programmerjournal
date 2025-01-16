import {Component, EventEmitter, Input, Output, signal} from '@angular/core';
import {MatButton, MatIconButton} from '@angular/material/button';
import {MatFormField, MatLabel} from '@angular/material/form-field';
import {MatInput} from '@angular/material/input';
import {FormControl, ReactiveFormsModule} from '@angular/forms';
import {tap} from 'rxjs';
import {MatIcon} from '@angular/material/icon';
import {MarkdownPipe} from '../markdown.pipe';
import {Entry, EntryService, TaskService, TaskSummary} from '../../../frontend-client';
import {MatToolbar} from '@angular/material/toolbar';
import {EntryStatus, isTask} from '../../../lib/entry';
import {MatMenu, MatMenuItem, MatMenuTrigger} from '@angular/material/menu';

@Component({
  selector: 'app-entry-sidebar',
  imports: [
    MatButton,
    MatFormField,
    MatInput,
    MatLabel,
    ReactiveFormsModule,
    MatIcon,
    MatIconButton,
    MarkdownPipe,
    MatToolbar,
    MatMenu,
    MatMenuItem,
    MatMenuTrigger
  ],
  templateUrl: './entry-sidebar.component.html',
  styleUrl: './entry-sidebar.component.scss',
  standalone: true,
})
export class EntrySidebarComponent {

  @Input()
  taskSummary: TaskSummary | null = null

  @Output()
  onSubmit = new EventEmitter<number>()

  @Output()
  onTaskDelete = new EventEmitter<number>()

  @Output()
  onCancel = new EventEmitter<void>()

  editingUpdate = signal<boolean>(false)
  editingDescription = signal<boolean>(false)

  updateFormControl = new FormControl("")
  updateDescriptionFormControl = new FormControl();

  constructor(private taskService: TaskService,
              private entryService: EntryService) {
  }

  goToChangeUpdateState(e: Event) {
    e.preventDefault()
    this.updateFormControl.setValue(this.taskSummary?.taskEntry.taskUpdate || "")
    this.editingUpdate.set(true)
  }

  goToChangeDescriptionState(e: Event) {
    e.preventDefault()
    this.updateDescriptionFormControl.setValue(this.taskSummary?.taskEntry.description || "")
    this.editingDescription.set(true);
  }

  submitUpdateChange(e: Event) {
    e.preventDefault()
    this.taskService.setTaskUpdate(this.taskSummary!.taskEntry.id, {update: this.updateFormControl.value || ""})
      .pipe(tap(() => this.editingUpdate.set(false)))
      .subscribe(() => {
        this.onSubmit.emit(this.taskSummary?.taskEntry.id)
        console.log("submitUpdateChange - subscribe")
      })
  }

  submitUpdateDescriptionChange(e: Event) {
    e.preventDefault()
    this.entryService.setDescription(this.taskSummary!.taskEntry.id, {description: this.updateDescriptionFormControl.value || ""})
      .pipe(tap(() => this.editingDescription.set(false)))
      .subscribe(() => {
        this.onSubmit.emit(this.taskSummary?.taskEntry.id)
        console.log("submitUpdateChange - subscribe")
      })
  }

  deleteTask(entry: Entry) {
    this.onTaskDelete.emit(entry.id)
  }

  closeSidebar() {
    this.onCancel.emit()
  }

  isTaskEntry(taskSummary: TaskSummary | null) {
    let status: EntryStatus = taskSummary?.taskEntry?.status as EntryStatus
    return isTask(status)
  }
}
