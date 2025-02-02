import {Component, computed, inject} from '@angular/core';

import {FormBuilder, ReactiveFormsModule, Validators} from '@angular/forms';
import {MatInputModule} from '@angular/material/input';
import {MatButtonModule} from '@angular/material/button';
import {MatSelectModule} from '@angular/material/select';
import {MatRadioModule} from '@angular/material/radio';
import {MatCardModule} from '@angular/material/card';
import {MatDrawer, MatDrawerContainer, MatDrawerContent} from "@angular/material/sidenav";
import {MatIcon} from "@angular/material/icon";
import {NavToolbarComponent} from "../../components/nav-toolbar/nav-toolbar.component";
import {MatCheckbox} from '@angular/material/checkbox';
import {JsonPipe} from '@angular/common';
import {CreateRecurringTaskInputBody, RecurringTaskService} from '../../../frontend-client';
import {Router} from '@angular/router';


@Component({
  selector: 'app-recurring-task-add-view',
  templateUrl: './recurring-task-add-view.component.html',
  styleUrl: './recurring-task-add-view.component.scss',
  standalone: true,
  imports: [
    MatInputModule,
    MatButtonModule,
    MatSelectModule,
    MatRadioModule,
    MatCardModule,
    ReactiveFormsModule,
    NavToolbarComponent,
    MatCheckbox,
    JsonPipe
  ]
})
export class RecurringTaskAddViewComponent {

  private fb = inject(FormBuilder);

  recurringTaskForm = this.fb.group({
    title: ["", Validators.required],
    description: "",
    freqMonday: false,
    freqTuesday: false,
    freqWednesday: false,
    freqThursday: false,
    freqFriday: false,
    freqSaturday: false,
    freqSunday: false,
  });

  weekdayAndFormControl = computed(() => {
    return [
      {day: "Monday",    value: "MON", control: this.recurringTaskForm.controls.freqMonday},
      {day: "Tuesday",   value: "TUE", control: this.recurringTaskForm.controls.freqTuesday},
      {day: "Wednesday", value: "WED", control: this.recurringTaskForm.controls.freqWednesday},
      {day: "Thursday",  value: "THU", control: this.recurringTaskForm.controls.freqThursday},
      {day: "Friday",    value: "FRI", control: this.recurringTaskForm.controls.freqFriday},
      {day: "Saturday",  value: "SAT", control: this.recurringTaskForm.controls.freqSaturday},
      {day: "Sunday",    value: "SUN", control: this.recurringTaskForm.controls.freqSunday},
    ]
  })

  constructor(private recurringTaskService: RecurringTaskService, private router: Router) {
  }

  onSubmit(): void {
    let byWeekDay: string = this.weekdayAndFormControl()
      .filter(e => e.control.value == true)
      .map(e => e.value)
      .join(",")
    let createInput: CreateRecurringTaskInputBody = {
      taskTitle: this.recurringTaskForm.controls.title.value || "",
      taskDescription: this.recurringTaskForm.controls.description.value || "",
      freqByWeekDay: byWeekDay
    }
    this.recurringTaskService.create(createInput)
      .subscribe(() => this.router.navigate(["/recurring"]))
  }
}
