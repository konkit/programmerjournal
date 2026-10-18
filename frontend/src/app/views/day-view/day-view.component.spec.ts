import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, provideRouter } from '@angular/router';
import { provideNoopAnimations } from '@angular/platform-browser/animations';
import { of } from 'rxjs';

import { DayViewComponent } from './day-view.component';
import { EntryListService } from '../../service/entry-list.service';

describe('DayViewComponent', () => {
  let component: DayViewComponent;
  let fixture: ComponentFixture<DayViewComponent>;
  let service: EntryListService;
  let loadDate: jasmine.Spy;
  let getWeeklyTasks: jasmine.Spy;

  beforeEach(async () => {
    try { localStorage.removeItem('dayView.weeklyCollapsed'); } catch {}
    await TestBed.configureTestingModule({
      imports: [DayViewComponent, HttpClientTestingModule],
      providers: [
        provideRouter([]),
        provideNoopAnimations(),
        {provide: ActivatedRoute, useValue: {snapshot: {params: {date: '2026-10-18'}}}},
      ]
    })
    .compileComponents();

    service = TestBed.inject(EntryListService);
    loadDate = spyOn(service, 'loadDate').and.returnValue(of({pending: [], done: []} as any));
    getWeeklyTasks = spyOn(service, 'getWeeklyTasks').and.returnValue(of({
      pending: [{id: 3, title: 'Weekly thing', status: 'TaskCreated', createdDate: '2026-W42', rank: 0, taskUpdate: ''} as any],
      done: [],
    } as any));

    fixture = TestBed.createComponent(DayViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('loads the date from the route', () => {
    expect(service.todayDate()).toBe('2026-10-18');
    expect(loadDate).toHaveBeenCalled();
  });

  it('opens the updates sidebar for a task', () => {
    const summary = {taskEntry: {id: 7, title: 'Task'}} as any;
    spyOn(service, 'getTaskSummary').and.returnValue(of(summary));
    const open = spyOn(component.sideDrawer, 'open').and.resolveTo('open');

    component.openUpdates(7);

    expect(open).toHaveBeenCalled();
    expect(component.editedTaskSummary()).toBe(summary);
  });

  it('shows an empty state and a quick-add input', () => {
    const el: HTMLElement = fixture.nativeElement;
    expect(el.querySelector('.hide-on-drag')!.textContent).toContain('Nothing planned yet');
    expect(el.querySelector('app-new-entry input')).not.toBeNull();
  });

  it('lists done entries under "Others"', () => {
    service.doneEntryList.set([{id: 1, title: 'Done thing', status: 'TaskDone', createdDate: '2026-10-18', rank: 0} as any]);
    fixture.detectChanges();

    const summary = fixture.nativeElement.querySelector('details summary') as HTMLElement;
    expect(summary.textContent!.trim()).toBe('Others (1)');
  });

  it('shows the week in the side rail', () => {
    const rail = fixture.nativeElement.querySelector('.weekly-rail') as HTMLElement;
    expect(rail.querySelector('.rail-title')!.textContent!.trim()).toBe('Week 42');
    expect(rail.textContent).toContain('Weekly thing');
  });

  it('collapses the weekly rail, remembers it and stops fetching it', () => {
    const toggle = fixture.nativeElement.querySelector('.collapse-toggle') as HTMLButtonElement;
    expect(toggle.getAttribute('aria-expanded')).toBe('true');

    toggle.click();
    fixture.detectChanges();
    getWeeklyTasks.calls.reset();
    service.triggerRefresh();
    fixture.detectChanges();

    expect(fixture.nativeElement.querySelector('.weekly-list')).toBeNull();
    expect(toggle.getAttribute('aria-expanded')).toBe('false');
    expect(localStorage.getItem('dayView.weeklyCollapsed')).toBe('true');
    expect(getWeeklyTasks).not.toHaveBeenCalled();

    toggle.click();
    fixture.detectChanges();
    expect(getWeeklyTasks).toHaveBeenCalled();
  });

  afterEach(() => {
    try { localStorage.removeItem('dayView.weeklyCollapsed'); } catch {}
  });
});
