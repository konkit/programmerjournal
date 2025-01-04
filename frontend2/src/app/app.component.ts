import {Component} from '@angular/core';
import {TaskListComponent} from './components/task-list/task-list.component';
import {RouterOutlet} from '@angular/router';

@Component({
  selector: 'app-root',
  imports: [TaskListComponent, RouterOutlet],
  templateUrl: './app.component.html',
  standalone: true,
  styleUrl: './app.component.scss'
})
export class AppComponent {
}
