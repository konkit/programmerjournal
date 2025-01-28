import {Component, input, output} from '@angular/core';
import {Entry} from '../../../../frontend-client';
import {CdkDrag, CdkDragHandle} from '@angular/cdk/drag-drop';
import {MatIcon} from '@angular/material/icon';
import {StatusButtonComponent} from '../../status-button/status-button.component';
import {EntryListService} from '../../../service/entry-list.service';
import {MatIconButton} from '@angular/material/button';

@Component({
  selector: 'app-entry',
  imports: [
    CdkDrag,
    MatIcon,
    StatusButtonComponent,
    CdkDragHandle,
    MatIconButton
  ],
  templateUrl: './entry.component.html',
  standalone: true,
  styleUrl: './entry.component.scss'
})
export class EntryComponent {
  entry = input<Entry>()

  onOpenUpdates = output<void>()

  constructor(private entryListService: EntryListService) {
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

  migrateToMonthly() {
    this.entryListService.migrateToMonthly(this.entry()!).subscribe()
  }

  migrateToDaily() {
    this.entryListService.migrateToDaily(this.entry()!).subscribe()
  }

  submitTitleEditWithValue(e: Event) {
    let newValue = (e.target as HTMLDivElement).innerText
    this.entryListService.setTitle(newValue, this.entry()!);
  }

  openUpdates() {
    this.onOpenUpdates.emit()
  }
}
