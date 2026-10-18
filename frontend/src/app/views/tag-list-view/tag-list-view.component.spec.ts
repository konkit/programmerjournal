import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';
import { of } from 'rxjs';

import { TagListViewComponent } from './tag-list-view.component';
import { TagService } from '../../../frontend-client';

describe('TagListViewComponent', () => {
  let fixture: ComponentFixture<TagListViewComponent>;
  let el: HTMLElement;

  async function create(names: string[]) {
    await TestBed.configureTestingModule({
      imports: [TagListViewComponent, HttpClientTestingModule],
      providers: [provideRouter([])]
    }).compileComponents();

    spyOn(TestBed.inject(TagService), 'listTags')
      .and.returnValue(of(names.map((Name, i) => ({ID: i + 1, Name}))) as any);
    fixture = TestBed.createComponent(TagListViewComponent);
    fixture.detectChanges();
    el = fixture.nativeElement;
  }

  const chips = () => Array.from(el.querySelectorAll<HTMLAnchorElement>('.tag-chip'));

  it('lists tags alphabetically as links', async () => {
    await create(['work', 'angular', 'go']);

    expect(chips().map(c => c.textContent!.trim())).toEqual(['#angular', '#go', '#work']);
    expect(chips()[0].getAttribute('href')).toBe('/tags/angular');
    expect(el.querySelector('.page-subtitle')!.textContent!.trim()).toBe('3 tags');
  });

  it('filters tags, ignoring a leading #', async () => {
    await create(['work', 'angular', 'go']);

    const input = el.querySelector('.filter input') as HTMLInputElement;
    input.value = '#ang';
    input.dispatchEvent(new Event('input'));
    fixture.detectChanges();

    expect(chips().map(c => c.textContent!.trim())).toEqual(['#angular']);
  });

  it('shows an empty state without tags', async () => {
    await create([]);

    expect(el.querySelector('.empty-state')!.textContent).toContain('No tags yet');
  });
});
