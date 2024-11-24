import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SnoozeDialogComponent } from './snooze-dialog.component';

describe('SnoozeDialogComponent', () => {
  let component: SnoozeDialogComponent;
  let fixture: ComponentFixture<SnoozeDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SnoozeDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SnoozeDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
