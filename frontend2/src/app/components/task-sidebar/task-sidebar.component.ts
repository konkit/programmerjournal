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
    MatIconButton
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

  editedTaskUpdates = signal<string[]>([])

  updateFormControl = new FormControl("update")

  constructor(private taskService: TaskService) {
  }

  isEditingTask(): boolean {
    return this.editing
  }

  isIdle(): boolean {
    return !this.editing;
  }

  addUpdate(e: Event) {
    e.preventDefault()
    this.updateFormControl.setValue(this.taskSummary?.task.update || "")
    this.editing = true;
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

  deleteTask(task: Task) {
    this.onTaskDelete.emit(task.id)
  }
}
