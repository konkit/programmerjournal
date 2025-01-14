import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SnoozeDayEntryDialogComponent } from './snooze-day-entry-dialog.component';

describe('SnoozeDialogComponent', () => {
  let component: SnoozeDayEntryDialogComponent;
  let fixture: ComponentFixture<SnoozeDayEntryDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SnoozeDayEntryDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SnoozeDayEntryDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
