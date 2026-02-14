import {Component, inject, input, output} from '@angular/core';
import {MatIcon, MatIconRegistry} from '@angular/material/icon';
import {MatIconButton} from '@angular/material/button';
import {MatMenu, MatMenuItem, MatMenuTrigger} from '@angular/material/menu';
import {Entry} from '../../../frontend-client';
import { EntryStatus } from '../../../lib/entry';
import {MatTooltip} from '@angular/material/tooltip';
import {DomSanitizer} from '@angular/platform-browser';
import {StatusIconComponent} from '../status-icon/status-icon.component';

@Component({
  selector: 'app-status-button',
  imports: [
    MatIcon,
    MatIconButton,
    MatMenu,
    MatMenuItem,
    MatMenuTrigger,
    MatTooltip,
    StatusIconComponent
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
  onTaskCanceled = output()
  onTaskToMonthly = output()
  onTaskToWeekly = output()
  onTaskToDaily = output()
  onEntryDelete = output()
  onAddUpdate = output()

  EntryStatus = EntryStatus;

  private matIconRegistry = inject(MatIconRegistry)
  private domSanitizer = inject(DomSanitizer)

  constructor() {
    this.matIconRegistry.addSvgIconSet(this.domSanitizer.bypassSecurityTrustResourceUrl('./assets/mdi.svg'));
  }

  markTaskAsCreated() {
    this.onTaskAsCreated.emit()
  }

  markTaskAsDone() {
    this.onTaskDone.emit()
  }

  snoozeTask() {
    this.onTaskSnoozed.emit()
  }

  cancelTask() {
    this.onTaskCanceled.emit()
  }

  migrateToMonthly() {
    this.onTaskToMonthly.emit()
  }

  migrateToWeekly() {
    this.onTaskToWeekly.emit()
  }


  migrateToDaily() {
    this.onTaskToDaily.emit()
  }

  isMonthEntry(entry: Entry): boolean {
    return entry.createdDate.length === 7; // 2024-12
  }

  isDayEntry(entry: Entry): boolean {
    return entry.createdDate.length === 10; // 2024-12-12
  }

  deleteEntry() {
    this.onEntryDelete.emit()
  }

  addUpdate() {
    this.onAddUpdate.emit()
  }
}
