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
      <button mat-button [mat-dialog-close]="true" class="delete-button">Delete</button>
    </mat-dialog-actions>
  `,
  styles: `
    @use '@angular/material' as mat;

    // color="warn" has no effect with M3 themes
    .delete-button {
      @include mat.button-overrides((
        text-label-text-color: var(--mat-sys-error),
        text-state-layer-color: var(--mat-sys-error),
      ));
    }
  `,
})
export class DeleteConfirmationDialogComponent {
  readonly data = inject<DeleteConfirmationData | null>(MAT_DIALOG_DATA, {optional: true});
}
