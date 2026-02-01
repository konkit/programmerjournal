import {Routes} from '@angular/router';
import {RecurringTasksViewComponent} from './views/recurring-tasks-view/recurring-tasks-view.component';
import {DayViewComponent} from './views/day-view/day-view.component';
import {MonthViewComponent} from './views/month-view/month-view.component';
import {WeeklySummaryViewComponent} from './views/weekly-summary-view/weekly-summary-view.component';
import {RecurringTaskAddViewComponent} from './views/recurring-task-add-view/recurring-task-add-view.component';
import {StartOfThisWeek, ThisMonth, Today} from '../lib/wall_date';
import {TagListViewComponent} from './views/tag-list-view/tag-list-view.component';
import {TagViewComponent} from './views/tag-view/tag-view.component';

export const routes: Routes = [
  { path: '', redirectTo: '/day', pathMatch: 'full' },
  {
    path: 'day',
    redirectTo: ({ queryParams }) => {
      let today = Today()
      return '/day/' + today
    },
    pathMatch: 'full'
  },
  { path: 'day/:date', component: DayViewComponent },
  {
    path: 'month',
    redirectTo: ({ queryParams }) => {
      let thisMonth = ThisMonth()
      return '/month/' + thisMonth
    },
    pathMatch: 'full'
  },
  { path: 'month/:date', component: MonthViewComponent },
  {
    path: 'weekSummary',
    redirectTo: ({ queryParams }) => {
      let startOfThisWeek = StartOfThisWeek()
      return '/weekSummary/' + startOfThisWeek
    },
    pathMatch: 'full'
  },
  { path: 'weekSummary/:date', component: WeeklySummaryViewComponent },
  {
    path: 'weekUpdates',
    redirectTo: ({ queryParams }) => {
      let startOfThisWeek = StartOfThisWeek()
      return '/weekUpdates/' + startOfThisWeek
    },
    pathMatch: 'full'
  },
  { path: 'recurring', component: RecurringTasksViewComponent },
  { path: 'recurring/add', component: RecurringTaskAddViewComponent },
  { path: "tags", component: TagListViewComponent },
  { path: "tags/:tag", component: TagViewComponent },
];
