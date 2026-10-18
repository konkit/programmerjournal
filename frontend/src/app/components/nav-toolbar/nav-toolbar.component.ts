import { Component } from '@angular/core';
import {MatAnchor, MatIconButton} from '@angular/material/button';
import {MatIcon} from '@angular/material/icon';
import {MatMenu, MatMenuItem, MatMenuTrigger} from '@angular/material/menu';
import {MatToolbar} from '@angular/material/toolbar';
import {RouterLink, RouterLinkActive} from '@angular/router';

@Component({
  selector: 'app-nav-toolbar',
  imports: [
    MatAnchor,
    MatIconButton,
    MatIcon,
    MatMenu,
    MatMenuItem,
    MatMenuTrigger,
    MatToolbar,
    RouterLink,
    RouterLinkActive
  ],
  templateUrl: './nav-toolbar.component.html',
  standalone: true,
  styleUrl: './nav-toolbar.component.scss'
})
export class NavToolbarComponent {

  readonly links = [
    {path: '/day', label: 'Daily Log'},
    {path: '/week', label: 'Weekly Log'},
    {path: '/month', label: 'Monthly Log'},
    {path: '/notes', label: 'Notes'},
    {path: '/tags', label: 'Tags'},
    {path: '/weekSummary', label: 'Weekly Summary'},
    {path: '/recurring', label: 'Recurring tasks'},
  ]
}
