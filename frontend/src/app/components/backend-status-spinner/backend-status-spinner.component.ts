import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { BackendStatusService } from '../../service/backend-status.service';

@Component({
  selector: 'app-backend-status-spinner',
  standalone: true,
  imports: [CommonModule, MatProgressSpinnerModule],
  templateUrl: './backend-status-spinner.component.html',
  styleUrl: './backend-status-spinner.component.scss'
})
export class BackendStatusSpinnerComponent {
  backendStatusService = inject(BackendStatusService);
}
