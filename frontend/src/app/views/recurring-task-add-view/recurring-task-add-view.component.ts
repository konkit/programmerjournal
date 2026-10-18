import {Component, inject, signal} from '@angular/core';

import {AbstractControl, FormBuilder, ReactiveFormsModule, ValidationErrors, Validators} from '@angular/forms';
import {MatInputModule} from '@angular/material/input';
import {MatButtonModule} from '@angular/material/button';
import {MatButtonToggleModule} from '@angular/material/button-toggle';
import {NavToolbarComponent} from "../../components/nav-toolbar/nav-toolbar.component";
import {CreateRecurringTaskInputBody, RecurringTaskService} from '../../../frontend-client';
import {Router, RouterLink} from '@angular/router';
import {describeWeekdays, WEEKDAYS, WORKDAYS} from '../../../lib/weekdays';

function atLeastOneDay(control: AbstractControl): ValidationErrors | null {
  return (control.value?.length ?? 0) > 0 ? null : {noDays: true};
}

@Component({
  selector: 'app-recurring-task-add-view',
  templateUrl: './recurring-task-add-view.component.html',
  styleUrl: './recurring-task-add-view.component.scss',
  standalone: true,
  imports: [
    MatInputModule,
    MatButtonModule,
    MatButtonToggleModule,
    ReactiveFormsModule,
    NavToolbarComponent,
    RouterLink
  ]
})
export class RecurringTaskAddViewComponent {

  private fb = inject(FormBuilder);

  recurringTaskForm = this.fb.nonNullable.group({
    title: ["", Validators.required],
    description: "",
    days: this.fb.nonNullable.control<string[]>([], atLeastOneDay),
  });

  submitting = signal(false);

  protected readonly weekdays = WEEKDAYS

  constructor(private recurringTaskService: RecurringTaskService, private router: Router) {
  }

  scheduleSummary(): string {
    const days = this.recurringTaskForm.controls.days.value
    return days.length > 0 ? describeWeekdays(days.join(",")) : ''
  }

  setDays(days: string[]) {
    this.recurringTaskForm.controls.days.setValue(days)
    this.recurringTaskForm.controls.days.markAsTouched()
  }

  selectWorkdays() {
    this.setDays([...WORKDAYS])
  }

  selectEveryDay() {
    this.setDays(WEEKDAYS.map(d => d.value))
  }

  onSubmit(): void {
    if (this.recurringTaskForm.invalid || this.submitting()) {
      this.recurringTaskForm.markAllAsTouched()
      return
    }

    // Keep the days in week order, whatever order they were clicked in
    const days = this.recurringTaskForm.controls.days.value
    const createInput: CreateRecurringTaskInputBody = {
      taskTitle: this.recurringTaskForm.controls.title.value.trim(),
      taskDescription: this.recurringTaskForm.controls.description.value.trim(),
      freqByWeekDay: WEEKDAYS.filter(d => days.includes(d.value)).map(d => d.value).join(","),
    }

    this.submitting.set(true)
    this.recurringTaskService.create(createInput)
      .subscribe({
        next: () => this.router.navigate(["/recurring"]),
        error: () => this.submitting.set(false),
      })
  }
}
