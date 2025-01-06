import {Component} from '@angular/core';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';

@Component({
  selector: 'app-recurring-tasks-view',
  imports: [
    NavToolbarComponent
  ],
  templateUrl: './recurring-tasks-view.component.html',
  standalone: true,
  styleUrl: './recurring-tasks-view.component.scss'
})
export class RecurringTasksViewComponent {

}
