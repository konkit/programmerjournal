import { Component } from '@angular/core';
import {NavToolbarComponent} from '../../components/nav-toolbar/nav-toolbar.component';

@Component({
  selector: 'app-weekly-summary-view',
  imports: [
    NavToolbarComponent
  ],
  templateUrl: './weekly-summary-view.component.html',
  standalone: true,
  styleUrl: './weekly-summary-view.component.scss'
})
export class WeeklySummaryViewComponent {

}
