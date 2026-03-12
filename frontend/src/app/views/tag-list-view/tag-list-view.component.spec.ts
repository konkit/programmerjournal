import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TagListViewComponent } from './tag-list-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute } from '@angular/router';
import { of } from 'rxjs';

describe('TagListViewComponent', () => {
  let component: TagListViewComponent;
  let fixture: ComponentFixture<TagListViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TagListViewComponent, HttpClientTestingModule],
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

    fixture = TestBed.createComponent(TagListViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
