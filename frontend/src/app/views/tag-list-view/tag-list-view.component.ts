import {Component, OnInit, signal} from '@angular/core';
import {Tag, TagService} from '../../../frontend-client';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {Router} from '@angular/router';
import {MatList, MatListItem} from '@angular/material/list';

@Component({
  selector: 'app-tag-list-view',
  imports: [
    NavToolbarComponent,
    MatList,
    MatListItem
  ],
  templateUrl: './tag-list-view.component.html',
  standalone: true,
  styleUrl: './tag-list-view.component.scss'
})
export class TagListViewComponent implements OnInit {

  tags = signal<Tag[]>([])

  constructor(private tagService: TagService, private router: Router) {
  }

  ngOnInit() {
    this.tagService.listTags().subscribe(res => {
      this.tags.set(res)
    })
  }

  navigateTo(tag: Tag) {
    this.router.navigate(["tags", tag.Name])
  }
}
