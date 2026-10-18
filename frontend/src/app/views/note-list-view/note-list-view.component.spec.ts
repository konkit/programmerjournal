import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';
import { of } from 'rxjs';

import { NoteListViewComponent } from './note-list-view.component';
import { NoteService } from '../../../frontend-client';

describe('NoteListViewComponent', () => {
  let fixture: ComponentFixture<NoteListViewComponent>;
  let el: HTMLElement;

  const note = (id: number, createdDate: string, title: string) =>
    ({id, createdDate, title, status: 'Note', rank: 0, taskID: `n${id}`} as any);

  async function create(notes: any[]) {
    await TestBed.configureTestingModule({
      imports: [NoteListViewComponent, HttpClientTestingModule],
      providers: [provideRouter([])]
    }).compileComponents();

    spyOn(TestBed.inject(NoteService), 'listNotes').and.returnValue(of({notes} as any));
    fixture = TestBed.createComponent(NoteListViewComponent);
    fixture.detectChanges();
    el = fixture.nativeElement;
  }

  function filterBy(text: string) {
    const input = el.querySelector('.filter input') as HTMLInputElement;
    input.value = text;
    input.dispatchEvent(new Event('input'));
    fixture.detectChanges();
  }

  it('groups notes by date with a link to that log', async () => {
    await create([
      note(1, '2026-10-18', 'Standup moved'),
      note(2, '2026-10-18', 'Deploy at 3pm'),
      note(3, '2026-W42', 'Weekly retro notes'),
    ]);

    const headings = Array.from(el.querySelectorAll<HTMLAnchorElement>('.date-heading'));
    expect(headings.map(h => h.textContent!.trim())).toEqual(['Sunday, 18 Oct 2026', 'Week 42, 2026']);
    expect(headings.map(h => h.getAttribute('href'))).toEqual(['/day/2026-10-18', '/week/2026-W42']);
    expect(el.querySelectorAll('.note').length).toBe(3);
    expect(el.querySelector('.page-subtitle')!.textContent!.trim()).toBe('3 notes');
  });

  it('sorts weekly and monthly notes by date among daily ones', async () => {
    // The order the backend returns them in: created_date desc, compared as strings
    await create([
      note(1, '2026-W42', 'Week 42 note'),
      note(2, '2026-10-18', 'Sunday note'),
      note(3, '2026-10-16', 'Friday note'),
      note(4, '2026-10', 'October note'),
    ]);

    const headings = Array.from(el.querySelectorAll('.date-heading')).map(h => h.textContent!.trim());
    // Week 42 starts on Monday 12 Oct, October on the 1st
    expect(headings).toEqual(['Sunday, 18 Oct 2026', 'Friday, 16 Oct 2026', 'Week 42, 2026', 'October 2026']);
  });

  it('renders markdown and links', async () => {
    await create([note(1, '2026-10-18', 'See **this** https://example.com')]);

    const rendered = el.querySelector('.note')!;
    expect(rendered.querySelector('strong')!.textContent).toBe('this');
    expect(rendered.querySelector('a')!.getAttribute('href')).toBe('https://example.com');
  });

  it('filters notes', async () => {
    await create([note(1, '2026-10-18', 'Standup moved'), note(2, '2026-10-17', 'Deploy at 3pm')]);

    filterBy('deploy');
    expect(Array.from(el.querySelectorAll('.note')).map(n => n.textContent!.trim())).toEqual(['Deploy at 3pm']);

    filterBy('nothing like this');
    expect(el.querySelector('.empty-state')!.textContent).toContain('No notes match');
  });

  it('shows an empty state without notes', async () => {
    await create([]);

    expect(el.querySelector('.filter')).toBeNull();
    expect(el.querySelector('.empty-state')!.textContent).toContain('No notes yet');
  });
});
