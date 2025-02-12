import {Component, inject, OnInit, signal} from '@angular/core';
import {Entry, Tag, TagService} from '../../../frontend-client';
import {JsonPipe} from '@angular/common';
import {ActivatedRoute} from '@angular/router';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';

@Component({
  selector: 'app-tag-view',
  imports: [
    JsonPipe,
    NavToolbarComponent
  ],
  templateUrl: './tag-view.component.html',
  standalone: true,
  styleUrl: './tag-view.component.scss'
})
export class TagViewComponent implements OnInit {

  tagName = signal<string>("")
  entries = signal<Entry[]>([])

  private readonly route = inject(ActivatedRoute);

  constructor(private tagService: TagService) {
  }

  ngOnInit() {
    const tagName = this.route.snapshot.paramMap.get('tag');
    this.tagName.set(tagName!)
    this.tagService.getTag(tagName!).subscribe(res => {
      this.entries.set(res)
    })
  }

}
