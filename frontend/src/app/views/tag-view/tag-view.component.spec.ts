import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TagViewComponent } from './tag-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';
import { RouterTestingModule } from '@angular/router/testing';

describe('TagsViewComponent', () => {
  let component: TagViewComponent;
  let fixture: ComponentFixture<TagViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TagViewComponent, HttpClientTestingModule, RouterTestingModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            snapshot: {
              paramMap: convertToParamMap({ tag: 'test-tag' })
            }
          }
        }
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TagViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
