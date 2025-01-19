import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MigrateToDayEntryDialogComponent } from './migrate-to-day-entry-dialog.component';

describe('SnoozeDialogComponent', () => {
  let component: MigrateToDayEntryDialogComponent;
  let fixture: ComponentFixture<MigrateToDayEntryDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MigrateToDayEntryDialogComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(MigrateToDayEntryDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
