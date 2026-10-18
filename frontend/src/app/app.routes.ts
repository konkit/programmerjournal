import {Routes} from '@angular/router';
import {RecurringTasksViewComponent} from './views/recurring-tasks-view/recurring-tasks-view.component';
import {DayViewComponent} from './views/day-view/day-view.component';
import {MonthViewComponent} from './views/month-view/month-view.component';
import {RecurringTaskAddViewComponent} from './views/recurring-task-add-view/recurring-task-add-view.component';
import {ThisMonth, ThisWeek, Today} from '../lib/wall_date';
import {TagListViewComponent} from './views/tag-list-view/tag-list-view.component';
import {TagViewComponent} from './views/tag-view/tag-view.component';
import {WeekViewComponent} from './views/week-view/week-view.component';
import {NoteListViewComponent} from './views/note-list-view/note-list-view.component';

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
    path: 'week',
    redirectTo: ({ queryParams }) => {
      let thisWeek = ThisWeek()
      return '/week/' + thisWeek
    },
    pathMatch: 'full'
  },
  { path: 'week/:date', component: WeekViewComponent },
  { path: 'recurring', component: RecurringTasksViewComponent },
  { path: 'recurring/add', component: RecurringTaskAddViewComponent },
  { path: "tags", component: TagListViewComponent },
  { path: "tags/:tag", component: TagViewComponent },
  { path: "notes", component: NoteListViewComponent },
];
