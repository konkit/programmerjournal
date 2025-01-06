import {Component, input, Input, InputSignal, output} from '@angular/core';
import {MatIcon} from '@angular/material/icon';
import {MatIconButton} from '@angular/material/button';
import {MatMenu, MatMenuItem, MatMenuTrigger} from '@angular/material/menu';
import {Task} from '../../../lib/task'

@Component({
  selector: 'app-status-button',
  imports: [
    MatIcon,
    MatIconButton,
    MatMenu,
    MatMenuItem,
    MatMenuTrigger
  ],
  templateUrl: './status-button.component.html',
  standalone: true,
  styleUrl: './status-button.component.scss'
})
export class StatusButtonComponent {
  @Input()
  task!: Task

  onTaskAsCreated = output()
  onTaskDone = output()
  onTaskSnoozed = output()
  onTaskToMonthly = output()

  markTaskAsCreated() {
    this.onTaskAsCreated.emit()
  }

  markTaskAsDone() {
    this.onTaskDone.emit()
  }

  snoozeTask() {
    this.onTaskSnoozed.emit()
  }

  migrateToMonthly() {
    this.onTaskToMonthly.emit()
  }
}
