import {Component, computed, model} from '@angular/core';
import {MatDialogModule} from '@angular/material/dialog';
import {MatButtonModule} from '@angular/material/button';
import {MatCardModule} from '@angular/material/card';
import {MatDatepickerModule} from '@angular/material/datepicker';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatNativeDateModule, MatOption} from '@angular/material/core';
import {MatSelect} from '@angular/material/select';
import {FormsModule} from '@angular/forms';

@Component({
  selector: 'app-snooze-month-dialog',
  imports: [MatDialogModule, MatButtonModule, MatCardModule, MatDatepickerModule, MatInputModule, MatFormFieldModule, MatNativeDateModule, MatOption, MatSelect, FormsModule],
  templateUrl: './snooze-month-entry-dialog.component.html',
  styleUrl: './snooze-month-entry-dialog.component.scss',
  standalone: true,
})
export class SnoozeMonthEntryDialogComponent {
  yearValues = [2025, 2026, 2027]
  monthValues = [
    { value: "01", label: 'January' },
    { value: "02", label: 'February' },
    { value: "03", label: 'March' },
    { value: "04", label: 'April' },
    { value: "05", label: 'May' },
    { value: "06", label: 'June' },
    { value: "07", label: 'July' },
    { value: "08", label: 'August' },
    { value: "09", label: 'September' },
    { value: "10", label: 'October' },
    { value: "11", label: 'November' },
    { value: "12", label: 'December' },
  ];

  yearValue = model<string>(new Date().getFullYear().toString());
  monthValue = model<string>((new Date().getMonth() + 1).toString().padStart(2, "0"));

  // selected = model<Date | null>(new Date());
  selected = computed(() => `${this.yearValue()}-${this.monthValue()}`)
}
