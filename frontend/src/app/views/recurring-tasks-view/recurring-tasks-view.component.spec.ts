import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RecurringTasksViewComponent } from './recurring-tasks-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';

describe('RecurringTasksViewComponent', () => {
  let component: RecurringTasksViewComponent;
  let fixture: ComponentFixture<RecurringTasksViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RecurringTasksViewComponent, HttpClientTestingModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            params: of({ id: 'test-id' })
          }
        }
      ]
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
