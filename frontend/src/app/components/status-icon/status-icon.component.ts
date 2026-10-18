import {Component, computed, input} from '@angular/core';
import {MatIcon} from '@angular/material/icon';
import {MatIconButton} from '@angular/material/button';
import {EntryStatus} from '../../../lib/entry';
import {Entry} from '../../../frontend-client';
import {MatTooltip} from '@angular/material/tooltip';
import {registerJournalSvgIcons} from '../../icons/journal-icons';

@Component({
  selector: 'app-status-icon',
  imports: [
    MatIcon,
    MatIconButton,
    MatTooltip
  ],
  templateUrl: './status-icon.component.html',
  styleUrl: './status-icon.component.scss'
})
export class StatusIconComponent {

  entry = input<Entry>()

  constructor() {
    registerJournalSvgIcons()
  }

  statusLabel = computed(() => {
    const status = this.entry()?.status
    const label = {
      [EntryStatus.TaskCreated]: 'To do',
      [EntryStatus.TaskDone]: 'Done',
      [EntryStatus.TaskSnoozed]: 'Snoozed',
      [EntryStatus.TaskMigrated]: 'Migrated',
      [EntryStatus.TaskCancelled]: 'Cancelled',
      [EntryStatus.Note]: 'Note',
    }[status as EntryStatus] ?? 'Entry'
    return `${label}. Change status`
  })

  protected readonly EntryStatus = EntryStatus;
}
