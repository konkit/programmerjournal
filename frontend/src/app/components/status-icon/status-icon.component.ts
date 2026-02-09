import {Component, input} from '@angular/core';
import {MatIcon} from '@angular/material/icon';
import {MatIconButton} from '@angular/material/button';
import {EntryStatus} from '../../../lib/entry';
import {Entry} from '../../../frontend-client';
import {MatTooltip} from '@angular/material/tooltip';

@Component({
  selector: 'app-status-icon',
  imports: [
    MatIcon,
    MatIconButton,
    MatTooltip
  ],
  templateUrl: './status-icon.component.html',
  styleUrl: './status-icon.component.scss'
})
export class StatusIconComponent {

  entry = input<Entry>()

  protected readonly EntryStatus = EntryStatus;
}
