import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TaskmenuComponent } from './taskmenu.component';

describe('TaskmenuComponent', () => {
  let component: TaskmenuComponent;
  let fixture: ComponentFixture<TaskmenuComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TaskmenuComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TaskmenuComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
