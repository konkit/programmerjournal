import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap, provideRouter } from '@angular/router';
import { of } from 'rxjs';

import { TagViewComponent } from './tag-view.component';
import { TagService } from '../../../frontend-client';

describe('TagViewComponent', () => {
  let fixture: ComponentFixture<TagViewComponent>;
  let el: HTMLElement;
  let getTag: jasmine.Spy;

  const entry = (id: number, createdDate: string, title: string, status: string, taskID: string) =>
    ({id, createdDate, title, status, taskID, rank: 0} as any);

  async function create(entries: any[]) {
    await TestBed.configureTestingModule({
      imports: [TagViewComponent, HttpClientTestingModule],
      providers: [
        provideRouter([]),
        {provide: ActivatedRoute, useValue: {snapshot: {paramMap: convertToParamMap({tag: 'work'})}}},
      ]
    }).compileComponents();

    getTag = spyOn(TestBed.inject(TagService), 'getTag').and.returnValue(of(entries) as any);
    fixture = TestBed.createComponent(TagViewComponent);
    fixture.detectChanges();
    el = fixture.nativeElement;
  }

  it('loads the tag from the route', async () => {
    await create([]);

    expect(getTag).toHaveBeenCalledWith('work');
    expect(el.querySelector('h1')!.textContent!.trim()).toBe('#work');
    expect(el.querySelector('.empty-state')!.textContent).toContain('Nothing is tagged #work');
  });

  it('shows each task once, at its latest date, newest first', async () => {
    await create([
      entry(1, '2026-10-16', 'Review PR #work', 'TaskMigrated', 'a'),
      entry(2, '2026-10-18', 'Review PR #work', 'TaskDone', 'a'),
      entry(3, '2026-10-17', 'Write ADR #work', 'TaskCreated', 'b'),
      entry(4, '2026-10-17', 'Idea #work', 'Note', ''),
    ]);

    const headings = Array.from(el.querySelectorAll('.date-heading')).map(h => h.textContent!.trim());
    expect(headings).toEqual(['Sunday, 18 Oct 2026', 'Saturday, 17 Oct 2026']);

    const rows = Array.from(el.querySelectorAll('.tagged-entry'));
    expect(rows.length).toBe(3);
    expect(rows[0].classList).toContain('done');
    expect(rows[0].querySelector('.status')!.getAttribute('aria-label')).toBe('Done');
    expect(el.querySelector('.page-subtitle')!.textContent!.trim()).toBe('3 entries');
  });
});
