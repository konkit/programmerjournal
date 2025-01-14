import { Component } from '@angular/core';
import {MatButton} from '@angular/material/button';
import {MatToolbar} from '@angular/material/toolbar';
import {RouterLink} from '@angular/router';

@Component({
  selector: 'app-nav-toolbar',
  imports: [
    MatButton,
    MatToolbar,
    RouterLink
  ],
  templateUrl: './nav-toolbar.component.html',
  standalone: true,
  styleUrl: './nav-toolbar.component.scss'
})
export class NavToolbarComponent {

}
