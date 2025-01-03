import {Component, OnInit} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {TaskListComponent} from './components/task-list/task-list.component';
import {MatButton} from "@angular/material/button";
import {NavigationComponent} from './navigation/navigation.component';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, TaskListComponent, MatButton, NavigationComponent],
  templateUrl: './app.component.html',
  standalone: true,
  styleUrl: './app.component.scss'
})
export class AppComponent {
}
