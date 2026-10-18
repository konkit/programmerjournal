import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EntryComponent } from './entry.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { By } from '@angular/platform-browser';
import { CdkDrag } from '@angular/cdk/drag-drop';

describe('EntryComponent', () => {
  let component: EntryComponent;
  let fixture: ComponentFixture<EntryComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [EntryComponent, HttpClientTestingModule]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EntryComponent);
    component = fixture.componentInstance;
    fixture.componentRef.setInput('entry', {
      id: 1,
      title: 'Test Entry',
      createdDate: '2024-01-01',
      status: 'TODO',
      description: '',
      rank: 0,
      recurringTaskID: 0,
      taskID: '1',
      taskSnoozedUntil: '',
      taskUpdate: ''
    });
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

describe('EntryComponent dragging', () => {
  let fixture: ComponentFixture<EntryComponent>;

  function render(status: string) {
    fixture = TestBed.createComponent(EntryComponent);
    fixture.componentRef.setInput('entry', {
      id: 1,
      title: 'Test Entry',
      createdDate: '2024-01-01',
      status,
      description: '',
      rank: 0,
      recurringTaskID: 0,
      taskID: '1',
      taskSnoozedUntil: '',
      taskUpdate: ''
    });
    fixture.detectChanges();
    return {
      drag: fixture.debugElement.query(By.directive(CdkDrag)).injector.get(CdkDrag),
      handle: fixture.nativeElement.querySelector('.drag-handle'),
    };
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [EntryComponent, HttpClientTestingModule]
    })
    .compileComponents();
  });

  for (const status of ['TaskCreated', 'TaskDone', 'Note']) {
    it(`can reorder ${status} entries`, () => {
      const {drag, handle} = render(status);
      expect(drag.disabled).toBeFalse();
      expect(handle).not.toBeNull();
    });
  }

  for (const status of ['TaskSnoozed', 'TaskMigrated', 'TaskCancelled']) {
    it(`can't drag ${status} entries`, () => {
      const {drag, handle} = render(status);
      expect(drag.disabled).toBeTrue();
      expect(handle).toBeNull();
    });
  }
});
