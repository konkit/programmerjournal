import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { provideRouter, Router } from '@angular/router';
import { provideNoopAnimations } from '@angular/platform-browser/animations';
import { of } from 'rxjs';

import { RecurringTaskAddViewComponent } from './recurring-task-add-view.component';
import { RecurringTaskService } from '../../../frontend-client';

describe('RecurringTaskAddViewComponent', () => {
  let component: RecurringTaskAddViewComponent;
  let fixture: ComponentFixture<RecurringTaskAddViewComponent>;
  let el: HTMLElement;
  let create: jasmine.Spy;
  let navigate: jasmine.Spy;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RecurringTaskAddViewComponent, HttpClientTestingModule],
      providers: [provideRouter([]), provideNoopAnimations()]
    }).compileComponents();

    create = spyOn(TestBed.inject(RecurringTaskService), 'create').and.returnValue(of({}) as any);
    navigate = spyOn(TestBed.inject(Router), 'navigate').and.resolveTo(true);
    fixture = TestBed.createComponent(RecurringTaskAddViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    el = fixture.nativeElement;
  });

  function submit() {
    el.querySelector('form')!.dispatchEvent(new Event('submit'));
    fixture.detectChanges();
  }

  it('does not show the form value as JSON', () => {
    expect(el.textContent).not.toContain('freqMonday');
    expect(el.textContent).not.toContain('{');
  });

  it('requires a title and at least one day', () => {
    submit();

    expect(create).not.toHaveBeenCalled();
    expect(el.textContent).toContain('Give the task a title');
    expect(el.querySelector('.days-error')!.textContent).toContain('Pick at least one day');
  });

  it('fills in presets', () => {
    component.selectWorkdays();
    fixture.detectChanges();
    expect(el.querySelector('.schedule-summary')!.textContent!.trim()).toBe('Weekdays');

    component.selectEveryDay();
    fixture.detectChanges();
    expect(el.querySelector('.schedule-summary')!.textContent!.trim()).toBe('Every day');
  });

  it('creates the task with days in week order and goes back to the list', () => {
    component.recurringTaskForm.setValue({title: '  Weekly report ', description: '', days: ['FRI', 'MON']});

    submit();

    expect(create).toHaveBeenCalledOnceWith({taskTitle: 'Weekly report', taskDescription: '', freqByWeekDay: 'MON,FRI'});
    expect(navigate).toHaveBeenCalledWith(['/recurring']);
  });

  it('toggles days by clicking them', () => {
    const toggles = Array.from(el.querySelectorAll<HTMLButtonElement>('mat-button-toggle button'));
    expect(toggles.length).toBe(7);

    toggles[2].click();
    toggles[0].click();
    fixture.detectChanges();

    expect(component.recurringTaskForm.controls.days.value).toEqual(['WED', 'MON']);
    expect(el.querySelector('.schedule-summary')!.textContent!.trim()).toBe('Mon, Wed');
  });
});
