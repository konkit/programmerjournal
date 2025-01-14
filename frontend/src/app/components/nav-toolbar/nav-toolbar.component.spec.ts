import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NavToolbarComponent } from './nav-toolbar.component';

describe('ToolbarComponent', () => {
  let component: NavToolbarComponent;
  let fixture: ComponentFixture<NavToolbarComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NavToolbarComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(NavToolbarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
