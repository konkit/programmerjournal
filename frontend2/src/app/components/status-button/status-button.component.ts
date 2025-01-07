import {Component, input, output} from '@angular/core';
import {MatIcon} from '@angular/material/icon';
import {MatIconButton} from '@angular/material/button';
import {MatMenu, MatMenuItem, MatMenuTrigger} from '@angular/material/menu';
import {Entry} from '../../../frontend-client';

enum EntryStatus {
  TaskCreated = "TaskCreated",
  TaskDone = "TaskDone",
  TaskSnoozed = "TaskSnoozed",
  TaskMigrated = "TaskMigrated",
  TaskCancelled = "TaskCancelled",
  Note = "Note",
}

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
  entry = input<Entry>()

  onTaskAsCreated = output()
  onTaskDone = output()
  onTaskSnoozed = output()
  onTaskToMonthly = output()

  EntryStatus = EntryStatus;

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
