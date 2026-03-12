import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EntryComponent } from './entry.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';

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
