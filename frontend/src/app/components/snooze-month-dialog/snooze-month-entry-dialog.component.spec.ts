import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SnoozeMonthEntryDialogComponent } from './snooze-month-entry-dialog.component';

describe('MonthSnoozeDialogComponent', () => {
  let component: SnoozeMonthEntryDialogComponent;
  let fixture: ComponentFixture<SnoozeMonthEntryDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SnoozeMonthEntryDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SnoozeMonthEntryDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
