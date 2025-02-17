import {Component, input, output, signal} from '@angular/core';
import {Entry} from '../../../frontend-client';
import {CdkDrag, CdkDragHandle} from '@angular/cdk/drag-drop';
import {MatIcon} from '@angular/material/icon';
import {StatusButtonComponent} from '../status-button/status-button.component';
import {EntryListService} from '../../service/entry-list.service';
import {MatIconButton} from '@angular/material/button';
import {RenderLinksPipe} from './render-links.pipe';
import {FormControl, ReactiveFormsModule} from '@angular/forms';
import {MatFormField, MatLabel} from '@angular/material/form-field';
import {MatInput} from '@angular/material/input';
import {MatTooltip} from '@angular/material/tooltip';

@Component({
  selector: 'app-entry',
  imports: [
    CdkDrag,
    MatIcon,
    StatusButtonComponent,
    CdkDragHandle,
    MatIconButton,
    RenderLinksPipe,
    ReactiveFormsModule,
    MatTooltip
  ],
  templateUrl: './entry.component.html',
  standalone: true,
  styleUrl: './entry.component.scss'
})
export class EntryComponent {
  entry = input<Entry>()

  onOpenUpdates = output<void>()

  fc = new FormControl("")

  isEdited = signal<boolean>(false)

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
    // console.log("newValue: ", this.fc.value)
    // let newValue = (e.target as HTMLDivElement).innerText
    this.entryListService.setTitle(this.fc.value || "", this.entry()!)
      .subscribe(() => {
        console.log("entryListService.setTitle subscribe")
        this.isEdited.set(false)
      })
  }

  setEdited(newValue: boolean) {
    this.fc.setValue(this.entry()?.title || "")
    this.isEdited.set(newValue);
  }

  openUpdates() {
    this.onOpenUpdates.emit()
  }

  cancelEdit($event: MouseEvent) {
    this.fc.setValue("")
    this.isEdited.set(false);
  }
}
