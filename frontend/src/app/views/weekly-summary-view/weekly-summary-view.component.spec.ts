import { ComponentFixture, TestBed } from '@angular/core/testing';

import { WeeklySummaryViewComponent } from './weekly-summary-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';

describe('WeeklySummaryViewComponent', () => {
  let component: WeeklySummaryViewComponent;
  let fixture: ComponentFixture<WeeklySummaryViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [WeeklySummaryViewComponent, HttpClientTestingModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            params: of({ date: '2024-01-01' }),
            snapshot: {
              params: {
                date: '2024-01-01'
              }
            }
          }
        }
      ]
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
