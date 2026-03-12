import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TitleWrapperComponent } from './title-wrapper.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('TitleWrapperComponent', () => {
  let component: TitleWrapperComponent;
  let fixture: ComponentFixture<TitleWrapperComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TitleWrapperComponent, HttpClientTestingModule]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TitleWrapperComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
