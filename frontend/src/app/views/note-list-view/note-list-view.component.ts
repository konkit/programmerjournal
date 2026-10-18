import {Component, computed, OnInit, signal} from '@angular/core';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {Entry, NoteService} from '../../../frontend-client';
import {RenderLinksPipe} from '../../components/entry/render-links.pipe';
import {MarkdownPipe} from '../../components/markdown.pipe';
import {MatIcon} from '@angular/material/icon';
import {RouterLink} from '@angular/router';
import {groupByDate} from '../../../lib/entry';
import {formatWallDate, logRouteFor, wallDateSortKey} from '../../../lib/wall_date';

@Component({
  selector: 'app-note-list-view',
  standalone: true,
  imports: [
    NavToolbarComponent,
    RenderLinksPipe,
    MarkdownPipe,
    MatIcon,
    RouterLink
  ],
  templateUrl: './note-list-view.component.html',
  styleUrl: './note-list-view.component.scss'
})
export class NoteListViewComponent implements OnInit {

  noteList = signal<Entry[]>([])
  loaded = signal(false)
  filter = signal('')

  filteredNotes = computed(() => {
    const query = this.filter().trim().toLowerCase()
    return query
      ? this.noteList().filter(note => (note.title || '').toLowerCase().includes(query))
      : this.noteList()
  })

  // Newest first. The backend sorts dates as plain strings, which puts weekly and monthly notes
  // after every day of the year
  groups = computed(() => groupByDate(
    [...this.filteredNotes()].sort((a, b) => wallDateSortKey(b.createdDate).localeCompare(wallDateSortKey(a.createdDate)))
  ))

  protected readonly formatWallDate = formatWallDate
  protected readonly logRouteFor = logRouteFor

  constructor(private noteService: NoteService) {
  }

  ngOnInit(): void {
    this.noteService.listNotes().subscribe(resp => {
      this.noteList.set(resp.notes || [])
      this.loaded.set(true)
    })
  }

  setFilter(event: Event) {
    this.filter.set((event.target as HTMLInputElement).value)
  }
}
