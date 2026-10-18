import { Component, EventEmitter, Input, Output, signal } from '@angular/core';
import {EntryListService} from '../../service/entry-list.service';
import {FormControl, ReactiveFormsModule} from '@angular/forms';
import {MatButton} from '@angular/material/button';
import {MatIcon} from '@angular/material/icon';

// Entries starting with this prefix are created as notes instead of tasks
export const NOTE_PREFIX = '-';

@Component({
  selector: 'app-new-entry',
  imports: [
    MatButton,
    MatIcon,
    ReactiveFormsModule
  ],
  templateUrl: './new-entry.component.html',
  standalone: true,
  styleUrl: './new-entry.component.scss'
})
export class NewEntryComponent {

  @Input() date?: string;
  @Output() onEntryCreated = new EventEmitter<void>();

  newEntryFormControl = new FormControl("", {nonNullable: true});
  submitting = signal(false);

  constructor(private entryListService: EntryListService) {
  }

  submit() {
    const value = this.newEntryFormControl.value.trim();
    if (!value || this.submitting()) {
      return;
    }

    const date = this.date || this.entryListService.todayDate();
    const isNote = value.startsWith(NOTE_PREFIX);
    const noteText = value.substring(NOTE_PREFIX.length).trim();
    if (isNote && !noteText) {
      return;
    }
    const request = isNote
      ? this.entryListService.createNote(noteText, date)
      : this.entryListService.createTask(value, date);

    this.submitting.set(true);
    request.subscribe({
      next: () => {
        this.newEntryFormControl.setValue("");
        this.onEntryCreated.emit();
      },
      complete: () => this.submitting.set(false),
      error: () => this.submitting.set(false),
    })
  }

  cancelEdit() {
    this.newEntryFormControl.setValue("")
  }
}
