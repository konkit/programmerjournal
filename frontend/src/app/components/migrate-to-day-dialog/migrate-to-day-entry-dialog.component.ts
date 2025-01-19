import {Component, inject, Injectable, input, model} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogModule,} from '@angular/material/dialog';
import {MatCardModule} from '@angular/material/card'
import {MatButtonModule} from '@angular/material/button'
import {MatDatepickerModule} from '@angular/material/datepicker';
import {MatInputModule} from '@angular/material/input'
import {MatFormFieldModule} from '@angular/material/form-field'
import {MAT_DATE_LOCALE, MatNativeDateModule, provideNativeDateAdapter} from '@angular/material/core';

@Component({
  selector: 'app-migrate-to-day-entry-dialog',
  imports: [MatDialogModule, MatButtonModule, MatCardModule, MatDatepickerModule, MatInputModule, MatFormFieldModule, MatNativeDateModule],
  templateUrl: './migrate-to-day-entry-dialog.component.html',
  styleUrl: './migrate-to-day-entry-dialog.component.scss',
  providers: [
    provideNativeDateAdapter(),
    {provide: MAT_DATE_LOCALE, useValue: "en-GB"}
  ],
  standalone: true,
})
export class MigrateToDayEntryDialogComponent {
  selected = model<Date | null>(new Date());
}
