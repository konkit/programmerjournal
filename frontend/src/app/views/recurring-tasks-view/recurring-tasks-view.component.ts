import {Component, OnInit, signal} from '@angular/core';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {MatIcon} from '@angular/material/icon';
import {MatFabButton, MatIconButton} from '@angular/material/button';
import {EntrySidebarComponent} from '../../components/entry-sidebar/entry-sidebar.component';
import {MatDrawer, MatDrawerContainer, MatDrawerContent} from '@angular/material/sidenav';
import {RouterLink} from '@angular/router';
import {RecurringTask, RecurringTaskService} from '../../../frontend-client';
import {JsonPipe} from '@angular/common';
import {MatActionList, MatList, MatListItem} from '@angular/material/list';
import {switchMap, tap} from 'rxjs';

@Component({
  selector: 'app-recurring-tasks-view',
  imports: [
    NavToolbarComponent,
    MatIcon,
    MatIconButton,
    MatFabButton,
    RouterLink,
  ],
  templateUrl: './recurring-tasks-view.component.html',
  standalone: true,
  styleUrl: './recurring-tasks-view.component.scss'
})
export class RecurringTasksViewComponent implements OnInit {

  recTaskList = signal<RecurringTask[]>([])

  constructor(private service: RecurringTaskService) {
  }

  ngOnInit() {
    this.refreshRecurringTasks().subscribe()
  }

  deleteRecurringTask(id: number) {
    this.service._delete(id)
      .pipe(switchMap(() => this.refreshRecurringTasks()))
      .subscribe()
  }

  refreshRecurringTasks() {
    console.log(`Refreshing recurring tasks`)
    return this.service.list()
      .pipe(
        tap(entries => {
          this.recTaskList.set(entries)
        })
      )
  }
}
