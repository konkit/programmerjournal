import {Component, computed, inject, OnInit, signal} from '@angular/core';
import {Entry, TagService} from '../../../frontend-client';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {MatIcon} from '@angular/material/icon';
import {MarkdownPipe} from '../../components/markdown.pipe';
import {RenderLinksPipe} from '../../components/entry/render-links.pipe';
import {EntryStatus, groupByDate} from '../../../lib/entry';
import {formatWallDate, logRouteFor, wallDateSortKey} from '../../../lib/wall_date';

const STATUS_ICONS: Record<string, {icon: string, label: string}> = {
  [EntryStatus.TaskCreated]: {icon: 'check_box_outline_blank', label: 'To do'},
  [EntryStatus.TaskDone]: {icon: 'check_box', label: 'Done'},
  [EntryStatus.TaskSnoozed]: {icon: 'arrow_forward_ios', label: 'Snoozed'},
  [EntryStatus.TaskMigrated]: {icon: 'arrow_back_ios_new', label: 'Migrated'},
  [EntryStatus.TaskCancelled]: {icon: 'indeterminate_check_box', label: 'Cancelled'},
  [EntryStatus.Note]: {icon: 'remove', label: 'Note'},
}

@Component({
  selector: 'app-tag-view',
  imports: [
    NavToolbarComponent,
    RouterLink,
    MatIcon,
    MarkdownPipe,
    RenderLinksPipe
  ],
  templateUrl: './tag-view.component.html',
  standalone: true,
  styleUrl: './tag-view.component.scss'
})
export class TagViewComponent implements OnInit {

  tagName = signal<string>("")
  entries = signal<Entry[]>([])
  loaded = signal(false)

  // A task that was migrated or snoozed has an entry for every day it was on; show only its latest one
  groups = computed(() => {
    const sorted = [...this.entries()]
      .sort((a, b) => wallDateSortKey(b.createdDate).localeCompare(wallDateSortKey(a.createdDate)))
    const seenTasks = new Set<string>()
    const latest = sorted.filter(e => {
      if (!e.taskID) return true
      if (seenTasks.has(e.taskID)) return false
      seenTasks.add(e.taskID)
      return true
    })
    return groupByDate(latest)
  })

  count = computed(() => this.groups().reduce((n, g) => n + g.entries.length, 0))

  protected readonly formatWallDate = formatWallDate
  protected readonly logRouteFor = logRouteFor

  private readonly route = inject(ActivatedRoute);

  constructor(private tagService: TagService) {
  }

  ngOnInit() {
    const tagName = this.route.snapshot.paramMap.get('tag');
    this.tagName.set(tagName!)
    this.tagService.getTag(tagName!).subscribe(res => {
      this.entries.set(res || [])
      this.loaded.set(true)
    })
  }

  statusIcon(entry: Entry) {
    return STATUS_ICONS[entry.status] ?? STATUS_ICONS[EntryStatus.Note]
  }
}
