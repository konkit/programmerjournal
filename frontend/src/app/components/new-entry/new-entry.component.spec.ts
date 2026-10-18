import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NewEntryComponent } from './new-entry.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { EntryListService } from '../../service/entry-list.service';
import { of } from 'rxjs';

describe('NewEntryComponent', () => {
  let component: NewEntryComponent;
  let fixture: ComponentFixture<NewEntryComponent>;
  let createTask: jasmine.Spy;
  let createNote: jasmine.Spy;

  function submit(value: string) {
    component.newEntryFormControl.setValue(value);
    fixture.nativeElement.querySelector('form').dispatchEvent(new Event('submit'));
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NewEntryComponent, HttpClientTestingModule]
    })
    .compileComponents();

    const service = TestBed.inject(EntryListService);
    service.todayDate.set('2026-10-18');
    createTask = spyOn(service, 'createTask').and.returnValue(of({} as any));
    createNote = spyOn(service, 'createNote').and.returnValue(of({} as any));

    fixture = TestBed.createComponent(NewEntryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('creates a task and clears the input', () => {
    const created = jasmine.createSpy('created');
    component.onEntryCreated.subscribe(created);

    submit('  Review PR #41 ');

    expect(createTask).toHaveBeenCalledOnceWith('Review PR #41', '2026-10-18');
    expect(createNote).not.toHaveBeenCalled();
    expect(component.newEntryFormControl.value).toBe('');
    expect(created).toHaveBeenCalled();
  });

  it('creates a note when the text starts with a dash', () => {
    submit('- Standup moved to 10:30');

    expect(createNote).toHaveBeenCalledOnceWith('Standup moved to 10:30', '2026-10-18');
    expect(createTask).not.toHaveBeenCalled();
  });

  it('ignores empty input', () => {
    submit('   ');
    submit('-');

    expect(createTask).not.toHaveBeenCalled();
    expect(createNote).not.toHaveBeenCalled();
  });

  it('uses the date input when given', () => {
    fixture.componentRef.setInput('date', '2026-W43');

    submit('Plan the week');

    expect(createTask).toHaveBeenCalledOnceWith('Plan the week', '2026-W43');
  });

  it('clears the input on Escape', () => {
    component.newEntryFormControl.setValue('draft');
    const input = fixture.nativeElement.querySelector('input') as HTMLInputElement;

    input.dispatchEvent(new KeyboardEvent('keydown', {key: 'Escape'}));

    expect(component.newEntryFormControl.value).toBe('');
  });
});
