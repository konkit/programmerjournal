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
import {getWeekOfYear} from '../../../lib/wall_date';

@Component({
  selector: 'app-snooze-week-dialog',
  imports: [MatDialogModule, MatButtonModule, MatCardModule, MatDatepickerModule, MatInputModule, MatFormFieldModule, MatNativeDateModule, MatOption, MatSelect, FormsModule],
  templateUrl: './snooze-week-dialog.component.html',
  styleUrl: './snooze-week-dialog.component.scss',
  standalone: true,
})
export class SnoozeWeekDialogComponent {
  yearValues = [2025, 2026, 2027]
  weekValues = Array.from({length: 53}, (_, i) => i + 1)

  yearValue = model<string>(new Date().getFullYear().toString());
  weekValue = model<string>(getWeekOfYear(new Date()).toString());

  selected = computed(() => `${this.yearValue()}-W${this.weekValue().padStart(2, '0')}`)
}
