import { Routes } from '@angular/router';
import {RecurringTasksViewComponent} from './views/recurring-tasks-view/recurring-tasks-view.component';
import {DayViewComponent} from './views/day-view/day-view.component';
import {MonthViewComponent} from './views/month-view/month-view.component';
import {WeeklySummaryViewComponent} from './views/weekly-summary-view/weekly-summary-view.component';

export const routes: Routes = [
  { path: '', redirectTo: '/day', pathMatch: 'full' },
  { path: 'day', component: DayViewComponent },
  { path: 'weekSummary', component: WeeklySummaryViewComponent },
  { path: 'month', component: MonthViewComponent },
  { path: 'recurring', component: RecurringTasksViewComponent }
];
