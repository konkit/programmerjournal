import { Routes } from '@angular/router';
import {DayViewComponent} from './components/day-view/day-view.component';
import {MonthViewComponent} from './components/month-view/month-view.component';

export const routes: Routes = [
  { path: '', redirectTo: '/day', pathMatch: 'full' },
  { path: 'day', component: DayViewComponent },
  { path: 'month', component: MonthViewComponent },
];
