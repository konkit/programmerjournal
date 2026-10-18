import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, provideRouter } from '@angular/router';
import { provideNoopAnimations } from '@angular/platform-browser/animations';
import { of } from 'rxjs';

import { WeekViewComponent } from './week-view.component';
import { EntryListService } from '../../service/entry-list.service';

describe('WeekViewComponent', () => {
  let component: WeekViewComponent;
  let fixture: ComponentFixture<WeekViewComponent>;
  let service: EntryListService;
  let loadDate: jasmine.Spy;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [WeekViewComponent, HttpClientTestingModule],
      providers: [
        provideRouter([]),
        provideNoopAnimations(),
        {provide: ActivatedRoute, useValue: {snapshot: {params: {date: '2026-W43'}}}},
      ]
    })
    .compileComponents();

    service = TestBed.inject(EntryListService);
    loadDate = spyOn(service, 'loadDate').and.returnValue(of({pending: [], done: []} as any));

    fixture = TestBed.createComponent(WeekViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('loads the date from the route', () => {
    expect(service.todayDate()).toBe('2026-W43');
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
    expect(el.querySelector('.hide-on-drag')!.textContent).toContain('Nothing planned for this week yet');
    expect(el.querySelector('app-new-entry input')).not.toBeNull();
  });

  it('lists done entries under "Others"', () => {
    service.doneEntryList.set([{id: 1, title: 'Done thing', status: 'TaskDone', createdDate: '2026-W43', rank: 0} as any]);
    fixture.detectChanges();

    const summary = fixture.nativeElement.querySelector('details summary') as HTMLElement;
    expect(summary.textContent!.trim()).toBe('Others (1)');
  });
});
