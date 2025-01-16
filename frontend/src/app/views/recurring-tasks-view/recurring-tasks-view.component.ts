import {Component} from '@angular/core';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';
import {MatIcon} from '@angular/material/icon';
import {MatIconButton} from '@angular/material/button';
import {EntrySidebarComponent} from '../../components/entry-sidebar/entry-sidebar.component';
import {MatDrawer, MatDrawerContainer, MatDrawerContent} from '@angular/material/sidenav';

@Component({
  selector: 'app-recurring-tasks-view',
  imports: [
    NavToolbarComponent,
    MatIcon,
    MatIconButton,
    EntrySidebarComponent,
    MatDrawer,
    MatDrawerContainer,
    MatDrawerContent
  ],
  templateUrl: './recurring-tasks-view.component.html',
  standalone: true,
  styleUrl: './recurring-tasks-view.component.scss'
})
export class RecurringTasksViewComponent {

}
