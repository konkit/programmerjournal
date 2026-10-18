import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';
import { provideNoopAnimations } from '@angular/platform-browser/animations';
import { MatDialog } from '@angular/material/dialog';
import { of } from 'rxjs';

import { RecurringTasksViewComponent } from './recurring-tasks-view.component';
import { RecurringTaskService } from '../../../frontend-client';

describe('RecurringTasksViewComponent', () => {
  let fixture: ComponentFixture<RecurringTasksViewComponent>;
  let el: HTMLElement;
  let service: RecurringTaskService;
  let deleteTask: jasmine.Spy;

  const task = {id: 4, taskTitle: 'Standup', taskDescription: 'Daily sync', freqByWeekDay: 'MON,TUE,WED,THU,FRI', generatedUntil: ''};

  async function create(tasks: any[]) {
    await TestBed.configureTestingModule({
      imports: [RecurringTasksViewComponent, HttpClientTestingModule],
      providers: [provideRouter([]), provideNoopAnimations()]
    }).compileComponents();

    service = TestBed.inject(RecurringTaskService);
    spyOn(service, 'list').and.returnValue(of(tasks) as any);
    deleteTask = spyOn(service, '_delete').and.returnValue(of({}) as any);
    fixture = TestBed.createComponent(RecurringTasksViewComponent);
    fixture.detectChanges();
    el = fixture.nativeElement;
  }

  function confirmDialog(answer: boolean | undefined) {
    spyOn(TestBed.inject(MatDialog), 'open').and.returnValue({afterClosed: () => of(answer)} as any);
    (el.querySelector('.delete-button') as HTMLButtonElement).click();
  }

  it('shows each task with its schedule', async () => {
    await create([task]);

    const card = el.querySelector('.recurring-task')!;
    expect(card.querySelector('.title')!.textContent!.trim()).toBe('Standup');
    expect(card.querySelector('.description')!.textContent!.trim()).toBe('Daily sync');
    expect(card.querySelector('.schedule-text')!.textContent!.trim()).toBe('Weekdays');
    const pills = Array.from(card.querySelectorAll('.weekday-pill')).map(p => p.classList.contains('on'));
    expect(pills).toEqual([true, true, true, true, true, false, false]);
  });

  it('deletes only after confirmation', async () => {
    await create([task]);

    confirmDialog(true);

    expect(deleteTask).toHaveBeenCalledWith(4);
  });

  it('keeps the task when the dialog is cancelled', async () => {
    await create([task]);

    confirmDialog(undefined);

    expect(deleteTask).not.toHaveBeenCalled();
  });

  it('shows an empty state with a way to add one', async () => {
    await create([]);

    const empty = el.querySelector('.empty-state')!;
    expect(empty.textContent).toContain('No recurring tasks yet');
    expect(empty.querySelector('a')!.getAttribute('href')).toBe('/recurring/add');
  });
});
