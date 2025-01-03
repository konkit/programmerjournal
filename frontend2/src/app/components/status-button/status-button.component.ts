import {Component, Input} from '@angular/core';
import {MatIcon} from '@angular/material/icon';
import {MatIconButton} from '@angular/material/button';

@Component({
  selector: 'app-status-button',
  imports: [
    MatIcon,
    MatIconButton
  ],
  templateUrl: './status-button.component.html',
  standalone: true,
  styleUrl: './status-button.component.scss'
})
export class StatusButtonComponent {
  @Input()
  status!: string
}
