import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RecurringTaskAddViewComponent } from './recurring-task-add-view.component';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';

describe('RecurringTaskAddViewComponent', () => {
  let component: RecurringTaskAddViewComponent;
  let fixture: ComponentFixture<RecurringTaskAddViewComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [NoopAnimationsModule, HttpClientTestingModule, RecurringTaskAddViewComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            params: of({ id: 'test-id' })
          }
        }
      ]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RecurringTaskAddViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should compile', () => {
    expect(component).toBeTruthy();
  });
});
