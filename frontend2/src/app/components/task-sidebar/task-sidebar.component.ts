import {Component, EventEmitter, Input, Output, signal} from '@angular/core';
import {MatButton, MatIconButton} from '@angular/material/button';
import {MatFormField, MatLabel} from '@angular/material/form-field';
import {MatInput} from '@angular/material/input';
import {TaskService, TaskSummary} from '../../../service/task.service';
import {JsonPipe} from '@angular/common';
import {FormControl, ReactiveFormsModule} from '@angular/forms';
import {tap} from 'rxjs';
import {MatIcon} from '@angular/material/icon';
import {MatMenuItem} from '@angular/material/menu';
import { Task } from '../../../lib/task';
import {MarkdownPipe} from '../markdown.pipe';

@Component({
  selector: 'app-task-sidebar',
  imports: [
    MatButton,
    MatFormField,
    MatInput,
    MatLabel,
    JsonPipe,
    ReactiveFormsModule,
    MatIcon,
    MatMenuItem,
    MatIconButton,
    MarkdownPipe
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

  editing: boolean = false
  editingDescription = signal<boolean>(false)

  editedTaskUpdates = signal<string[]>([])

  updateFormControl = new FormControl("")
  updateDescriptionFormControl = new FormControl();

  constructor(private taskService: TaskService) {
  }

  isEditingTask(): boolean {
    return this.editing
  }

  isEditingDescription() {
    return this.editingDescription()
  }

  isIdle(): boolean {
    return !this.editing;
  }

  goToChangeUpdateState(e: Event) {
    e.preventDefault()
    this.updateFormControl.setValue(this.taskSummary?.task.update || "")
    this.editing = true;
  }

  goToChangeDescriptionState(e: Event) {
    e.preventDefault()
    this.updateDescriptionFormControl.setValue(this.taskSummary?.task.description || "")
    this.editingDescription.set(true);
  }

  submitUpdateChange(e: Event) {
    e.preventDefault()
    this.taskService.setTaskUpdate(this.taskSummary!.task.id, this.updateFormControl.value || "")
      .pipe(tap(() => this.editing = false))
      .subscribe(() => {
        this.onSubmit.emit(this.taskSummary?.task.id)
        console.log("submitUpdateChange - subscribe")
        // this.editing = false
      })
  }

  submitUpdateDescriptionChange(e: Event) {
    e.preventDefault()
    this.taskService.setTaskDescription(this.taskSummary!.task.id, this.updateDescriptionFormControl.value || "")
      .pipe(tap(() => this.editingDescription.set(false)))
      .subscribe(() => {
        this.onSubmit.emit(this.taskSummary?.task.id)
        console.log("submitUpdateChange - subscribe")
        // this.editing = false
      })
  }

  deleteTask(task: Task) {
    this.onTaskDelete.emit(task.id)
  }
}
