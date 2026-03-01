import {Component, OnInit, signal} from '@angular/core';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {Entry, NoteService} from '../../../frontend-client';
import {RenderLinksPipe} from '../../components/entry/render-links.pipe';
import {MatTooltip} from '@angular/material/tooltip';

@Component({
  selector: 'app-note-list-view',
  standalone: true,
  imports: [
    NavToolbarComponent,
    RenderLinksPipe,
    MatTooltip
  ],
  templateUrl: './note-list-view.component.html',
  styleUrl: './note-list-view.component.scss'
})
export class NoteListViewComponent implements OnInit {

  noteList = signal<Entry[]>([])

  constructor(private noteService: NoteService) {
  }

  ngOnInit(): void {
    this.noteService.listNotes().subscribe(resp => {
      this.noteList.set(resp.notes || [])
    })
  }
}
