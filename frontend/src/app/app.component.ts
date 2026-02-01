import {Component} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {BackendStatusSpinnerComponent} from './components/backend-status-spinner/backend-status-spinner.component';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, BackendStatusSpinnerComponent],
  templateUrl: './app.component.html',
  standalone: true,
  styleUrl: './app.component.scss'
})
export class AppComponent {
}
