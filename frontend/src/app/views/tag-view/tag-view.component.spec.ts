import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TagViewComponent } from './tag-view.component';

describe('TagsViewComponent', () => {
  let component: TagViewComponent;
  let fixture: ComponentFixture<TagViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TagViewComponent]
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
