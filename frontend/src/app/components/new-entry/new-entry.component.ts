import { Component, EventEmitter, Input, Output } from '@angular/core';
import {EntryListService} from '../../service/entry-list.service';
import {FormControl, ReactiveFormsModule} from '@angular/forms';
import {MatFormField, MatLabel} from '@angular/material/form-field';
import {MatInput} from '@angular/material/input';
import {MatButton} from '@angular/material/button';

export enum EditorStateEnum {
  IDLE,
  EDITING_NEW_TASK,
  EDITING_NEW_NOTE,
}

@Component({
  selector: 'app-new-entry',
  imports: [
    MatFormField,
    MatInput,
    MatLabel,
    MatButton,
    ReactiveFormsModule
  ],
  templateUrl: './new-entry.component.html',
  standalone: true,
  styleUrl: './new-entry.component.scss'
})
export class NewEntryComponent {

  @Input() date?: string;
  @Output() onEntryCreated = new EventEmitter<void>();

  EditorStateEnum = EditorStateEnum;
  editorState = EditorStateEnum.IDLE

  newEntryFormControl = new FormControl("");

  constructor(private entryListService: EntryListService) {
  }

  submitNewTask() {
    let taskValue = this.newEntryFormControl.value || "";
    const date = this.date || this.entryListService.todayDate();
    return this.entryListService.createTask(taskValue, date)
      .subscribe(() => {
        this.editorState = EditorStateEnum.IDLE;
        this.onEntryCreated.emit();
      })
  }

  submitNewNote() {
    let taskValue = this.newEntryFormControl.value || "";
    const date = this.date || this.entryListService.todayDate();

    return this.entryListService.createNote(taskValue, date)
      .subscribe(() => {
        this.editorState = EditorStateEnum.IDLE;
        this.onEntryCreated.emit();
      })
  }

  cancelEdit() {
    this.newEntryFormControl.setValue("")
    this.editorState = EditorStateEnum.IDLE;
  }

  setCreatingNewTaskState() {
    this.newEntryFormControl.setValue("")
    this.editorState = EditorStateEnum.EDITING_NEW_TASK;
  }

  setCreatingNewNoteState() {
    this.newEntryFormControl.setValue("")
    this.editorState = EditorStateEnum.EDITING_NEW_NOTE;
  }

}
