import { ComponentFixture, TestBed } from '@angular/core/testing';

import { WeeklySummaryViewComponent } from './weekly-summary-view.component';

describe('WeeklySummaryViewComponent', () => {
  let component: WeeklySummaryViewComponent;
  let fixture: ComponentFixture<WeeklySummaryViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [WeeklySummaryViewComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(WeeklySummaryViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
