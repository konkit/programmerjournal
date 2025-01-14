import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RecurringTasksViewComponent } from './recurring-tasks-view.component';

describe('RecurringTasksViewComponent', () => {
  let component: RecurringTasksViewComponent;
  let fixture: ComponentFixture<RecurringTasksViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RecurringTasksViewComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RecurringTasksViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
