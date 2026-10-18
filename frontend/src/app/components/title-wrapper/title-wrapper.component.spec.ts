import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TitleWrapperComponent } from './title-wrapper.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { EntryListService } from '../../service/entry-list.service';
import { EMPTY } from 'rxjs';
import { Today } from '../../../lib/wall_date';

describe('TitleWrapperComponent', () => {
  let component: TitleWrapperComponent;
  let fixture: ComponentFixture<TitleWrapperComponent>;
  let service: EntryListService;

  function render(date: string): HTMLElement {
    service.todayDate.set(date);
    fixture.detectChanges();
    return fixture.nativeElement;
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TitleWrapperComponent, HttpClientTestingModule]
    })
    .compileComponents();

    service = TestBed.inject(EntryListService);
    fixture = TestBed.createComponent(TitleWrapperComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('formats a day', () => {
    const el = render('2026-10-18');
    expect(el.querySelector('h1')!.textContent!.trim()).toBe('Sunday, 18 Oct 2026');
    expect(el.querySelector('.subtitle')).toBeNull();
  });

  it('formats a week with its date range', () => {
    const el = render('2026-W43');
    expect(el.querySelector('h1')!.textContent!.trim()).toBe('Week 43');
    expect(el.querySelector('.subtitle')!.textContent!.trim()).toBe('19 – 25 Oct 2026');
  });

  it('formats week ranges spanning months and years', () => {
    expect(render('2026-W44').querySelector('.subtitle')!.textContent!.trim()).toBe('26 Oct – 1 Nov 2026');
    expect(render('2026-W01').querySelector('.subtitle')!.textContent!.trim()).toBe('29 Dec 2025 – 4 Jan 2026');
  });

  it('formats a month', () => {
    expect(render('2026-10').querySelector('h1')!.textContent!.trim()).toBe('October 2026');
  });

  it('shows a chip for the current date', () => {
    const el = render(Today());
    expect(el.querySelector('.current-chip')!.textContent!.trim()).toBe('Today');
    expect(el.querySelector('.current-button')).toBeNull();
  });

  it('offers to jump back to the current date', () => {
    const goToCurrent = spyOn(service, 'goToCurrent').and.returnValue(EMPTY);
    const el = render('2020-W10');

    const button = el.querySelector<HTMLButtonElement>('.current-button')!;
    expect(button.textContent!.trim()).toBe('Go to this week');
    button.click();

    expect(goToCurrent).toHaveBeenCalled();
  });

  it('labels the arrow buttons', () => {
    const labels = Array.from(render('2026-10').querySelectorAll('.title-row button'))
      .map(b => b.getAttribute('aria-label'));
    expect(labels).toEqual(['Previous month', 'Next month']);
  });

  it('shows the import banner only when there are past tasks', () => {
    expect(render('2026-10-18').querySelector('.import-banner')).toBeNull();

    const importPastTasks = spyOn(service, 'importPastTasks').and.returnValue(EMPTY);
    service.pendingImportsCount.set(3);
    fixture.detectChanges();

    const banner = fixture.nativeElement.querySelector('.import-banner') as HTMLElement;
    expect(banner.textContent).toContain('3 unfinished tasks from earlier');
    banner.querySelector('button')!.click();
    expect(importPastTasks).toHaveBeenCalled();
  });
});
