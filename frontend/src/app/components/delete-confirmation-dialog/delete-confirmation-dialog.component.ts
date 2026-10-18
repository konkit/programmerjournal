import {Component, inject} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogModule} from '@angular/material/dialog';
import {MatButtonModule} from '@angular/material/button';

export interface DeleteConfirmationData {
  title?: string;
  message?: string;
}

@Component({
  selector: 'app-delete-confirmation-dialog',
  standalone: true,
  imports: [MatDialogModule, MatButtonModule],
  template: `
    <h2 mat-dialog-title>{{ data?.title ?? 'Delete Task' }}</h2>
    <mat-dialog-content>
      {{ data?.message ?? 'Are you sure you want to delete this task?' }}
    </mat-dialog-content>
    <mat-dialog-actions align="end">
      <button mat-button mat-dialog-close>Cancel</button>
      <button mat-button [mat-dialog-close]="true" color="warn">Delete</button>
    </mat-dialog-actions>
  `,
})
export class DeleteConfirmationDialogComponent {
  readonly data = inject<DeleteConfirmationData | null>(MAT_DIALOG_DATA, {optional: true});
}
