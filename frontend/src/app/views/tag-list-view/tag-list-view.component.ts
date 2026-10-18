import {Component, computed, OnInit, signal} from '@angular/core';
import {Tag, TagService} from '../../../frontend-client';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {RouterLink} from '@angular/router';
import {MatIcon} from '@angular/material/icon';

@Component({
  selector: 'app-tag-list-view',
  imports: [
    NavToolbarComponent,
    RouterLink,
    MatIcon
  ],
  templateUrl: './tag-list-view.component.html',
  standalone: true,
  styleUrl: './tag-list-view.component.scss'
})
export class TagListViewComponent implements OnInit {

  tags = signal<Tag[]>([])
  loaded = signal(false)
  filter = signal('')

  filteredTags = computed(() => {
    const query = this.filter().trim().replace(/^#/, '').toLowerCase()
    return this.tags()
      .filter(tag => tag.Name.toLowerCase().includes(query))
      .sort((a, b) => a.Name.localeCompare(b.Name))
  })

  constructor(private tagService: TagService) {
  }

  ngOnInit() {
    this.tagService.listTags().subscribe(res => {
      this.tags.set(res || [])
      this.loaded.set(true)
    })
  }

  setFilter(event: Event) {
    this.filter.set((event.target as HTMLInputElement).value)
  }
}
