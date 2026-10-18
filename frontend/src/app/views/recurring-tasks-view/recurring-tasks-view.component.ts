import {Component, OnInit, signal} from '@angular/core';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {MatIcon} from '@angular/material/icon';
import {MatAnchor, MatIconButton} from '@angular/material/button';
import {MatTooltip} from '@angular/material/tooltip';
import {MatDialog} from '@angular/material/dialog';
import {RouterLink} from '@angular/router';
import {RecurringTask, RecurringTaskService} from '../../../frontend-client';
import {filter, switchMap, tap} from 'rxjs';
import {
  DeleteConfirmationData,
  DeleteConfirmationDialogComponent
} from '../../components/delete-confirmation-dialog/delete-confirmation-dialog.component';
import {describeWeekdays, WEEKDAYS} from '../../../lib/weekdays';

@Component({
  selector: 'app-recurring-tasks-view',
  imports: [
    NavToolbarComponent,
    MatIcon,
    MatIconButton,
    MatAnchor,
    MatTooltip,
    RouterLink,
  ],
  templateUrl: './recurring-tasks-view.component.html',
  standalone: true,
  styleUrl: './recurring-tasks-view.component.scss'
})
export class RecurringTasksViewComponent implements OnInit {

  recTaskList = signal<RecurringTask[]>([])
  loaded = signal(false)

  protected readonly weekdays = WEEKDAYS
  protected readonly describeWeekdays = describeWeekdays

  constructor(private service: RecurringTaskService, private dialog: MatDialog) {
  }

  ngOnInit() {
    this.refreshRecurringTasks().subscribe()
  }

  isOn(recTask: RecurringTask, day: string) {
    return recTask.freqByWeekDay.includes(day)
  }

  deleteRecurringTask(recTask: RecurringTask) {
    const data: DeleteConfirmationData = {
      title: 'Delete recurring task',
      message: `Stop creating “${recTask.taskTitle}”? Tasks that were already created stay in your logs.`,
    }
    this.dialog.open(DeleteConfirmationDialogComponent, {data, width: '400px'})
      .afterClosed()
      .pipe(
        filter(confirmed => confirmed === true),
        switchMap(() => this.service._delete(recTask.id)),
        switchMap(() => this.refreshRecurringTasks())
      )
      .subscribe()
  }

  refreshRecurringTasks() {
    return this.service.list()
      .pipe(
        tap(entries => {
          this.recTaskList.set(entries || [])
          this.loaded.set(true)
        })
      )
  }
}
