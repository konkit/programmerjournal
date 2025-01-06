import {Component, EventEmitter, Input, Output, signal} from '@angular/core';
import {MatButton, MatIconButton} from '@angular/material/button';
import {MatFormField, MatLabel} from '@angular/material/form-field';
import {MatInput} from '@angular/material/input';
import {FormControl, ReactiveFormsModule} from '@angular/forms';
import {tap} from 'rxjs';
import {MatIcon} from '@angular/material/icon';
import { Task } from '../../../lib/task';
import {MarkdownPipe} from '../markdown.pipe';
import {TaskService, TaskSummary} from '../../../frontend-client';
import {MatCard, MatCardContent, MatCardHeader} from '@angular/material/card';
import {MatToolbar} from '@angular/material/toolbar';

@Component({
  selector: 'app-task-sidebar',
  imports: [
    MatButton,
    MatFormField,
    MatInput,
    MatLabel,
    ReactiveFormsModule,
    MatIcon,
    MatIconButton,
    MarkdownPipe,
    MatToolbar
  ],
  templateUrl: './task-sidebar.component.html',
  styleUrl: './task-sidebar.component.scss',
  standalone: true,
})
export class TaskSidebarComponent {

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

  constructor(private taskService: TaskService) {
  }

  goToChangeUpdateState(e: Event) {
    e.preventDefault()
    this.updateFormControl.setValue(this.taskSummary?.task.update || "")
    this.editingUpdate.set(true)
  }

  goToChangeDescriptionState(e: Event) {
    e.preventDefault()
    this.updateDescriptionFormControl.setValue(this.taskSummary?.task.description || "")
    this.editingDescription.set(true);
  }

  submitUpdateChange(e: Event) {
    e.preventDefault()
    this.taskService.setTaskUpdate(this.taskSummary!.task.id, {update: this.updateFormControl.value || ""})
      .pipe(tap(() => this.editingUpdate.set(false)))
      .subscribe(() => {
        this.onSubmit.emit(this.taskSummary?.task.id)
        console.log("submitUpdateChange - subscribe")
      })
  }

  submitUpdateDescriptionChange(e: Event) {
    e.preventDefault()
    this.taskService.setTaskDescription(this.taskSummary!.task.id, {description: this.updateDescriptionFormControl.value || ""})
      .pipe(tap(() => this.editingDescription.set(false)))
      .subscribe(() => {
        this.onSubmit.emit(this.taskSummary?.task.id)
        console.log("submitUpdateChange - subscribe")
      })
  }

  deleteTask(task: Task) {
    this.onTaskDelete.emit(task.id)
  }

  closeSidebar() {
    this.onCancel.emit()
  }
}
