import { TestBed } from '@angular/core/testing';
import {EntryListService} from './entry-list.service';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { provideRouter, Router } from '@angular/router';
import { of } from 'rxjs';
import { EntryService, TaskService } from '../../frontend-client';
import { ThisMonth, ThisWeek, Today } from '../../lib/wall_date';

describe('EntryListService', () => {
  let service: EntryListService;
  let listEntries: jasmine.Spy;
  let countPastTasks: jasmine.Spy;
  let navigate: jasmine.Spy;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [provideRouter([])]
    });
    service = TestBed.inject(EntryListService);

    listEntries = spyOn(TestBed.inject(EntryService), 'listEntries')
      .and.returnValue(of({pending: [], done: []}) as any);
    countPastTasks = spyOn(TestBed.inject(TaskService), 'countPastTasks')
      .and.returnValue(of({count: 3}) as any);
    navigate = spyOn(TestBed.inject(Router), 'navigate').and.resolveTo(true);

    service.todayDate.set('2026-10-18');
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('refreshTasks reloads only the entry list', () => {
    service.refreshTasks().subscribe();

    expect(listEntries).toHaveBeenCalledOnceWith('2026-10-18');
    expect(countPastTasks).not.toHaveBeenCalled();
  });

  it('loadDate also reloads the past tasks count', () => {
    service.loadDate().subscribe();

    expect(listEntries).toHaveBeenCalledOnceWith('2026-10-18');
    expect(countPastTasks).toHaveBeenCalledOnceWith('2026-10-18');
    expect(service.pendingImportsCount()).toBe(3);
  });

  it('does not refetch the past tasks count after a mutation', () => {
    spyOn(TestBed.inject(TaskService), 'setTaskDone').and.returnValue(of({}) as any);

    service.markTaskAsDone(1).subscribe();

    expect(listEntries).toHaveBeenCalledTimes(1);
    expect(countPastTasks).not.toHaveBeenCalled();
  });

  it('moves forward and back by a day', () => {
    service.dateForward().subscribe();
    expect(service.todayDate()).toBe('2026-10-19');
    expect(navigate).toHaveBeenCalledWith(['/day', '2026-10-19']);

    service.dateBackward().subscribe();
    service.dateBackward().subscribe();
    expect(service.todayDate()).toBe('2026-10-17');
  });

  it('moves by a week across the year boundary', () => {
    service.todayDate.set('2027-W01');

    service.dateBackward().subscribe();

    expect(service.todayDate()).toBe('2026-W53');
    expect(navigate).toHaveBeenCalledWith(['/week', '2026-W53']);
  });

  it('moves by a month across the year boundary', () => {
    service.todayDate.set('2026-12');

    service.dateForward().subscribe();

    expect(service.todayDate()).toBe('2027-01');
    expect(navigate).toHaveBeenCalledWith(['/month', '2027-01']);
  });

  it('reloads the past tasks count when the date changes', () => {
    service.dateForward().subscribe();

    expect(countPastTasks).toHaveBeenCalledOnceWith('2026-10-19');
  });

  it('goToCurrent jumps to the current day, week or month', () => {
    service.todayDate.set('2020-01-01');
    service.goToCurrent().subscribe();
    expect(service.todayDate()).toBe(Today());

    service.todayDate.set('2020-W10');
    service.goToCurrent().subscribe();
    expect(service.todayDate()).toBe(ThisWeek());
    expect(navigate).toHaveBeenCalledWith(['/week', ThisWeek()]);

    service.todayDate.set('2020-03');
    service.goToCurrent().subscribe();
    expect(service.todayDate()).toBe(ThisMonth());
  });
});
