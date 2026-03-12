import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SnoozeWeekDialogComponent } from './snooze-week-dialog.component';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';

describe('SnoozeWeekDialogComponent', () => {
  let component: SnoozeWeekDialogComponent;
  let fixture: ComponentFixture<SnoozeWeekDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SnoozeWeekDialogComponent, NoopAnimationsModule]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SnoozeWeekDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
