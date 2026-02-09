import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StatusIconComponent } from './status-icon.component';

describe('StatusiconComponent', () => {
  let component: StatusIconComponent;
  let fixture: ComponentFixture<StatusIconComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StatusIconComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(StatusIconComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
