import {Component, inject, output} from '@angular/core';
import {MatMenuModule} from '@angular/material/menu'
import {MatButtonModule } from '@angular/material/button'
import {MatIconModule} from '@angular/material/icon'
import {MatDialog} from '@angular/material/dialog';
import {DeleteConfirmationDialogComponent} from '../delete-confirmation-dialog/delete-confirmation-dialog.component';

@Component({
  selector: 'app-taskmenu',
  imports: [MatIconModule, MatMenuModule, MatButtonModule],
  templateUrl: './taskmenu.component.html',
  styleUrl: './taskmenu.component.scss',
  standalone: true,

})
export class TaskmenuComponent {
  taskDeleted = output<void>()
  taskSnoozed = output<void>()
  taskDone = output<void>()

  readonly dialog = inject(MatDialog);

  openDeleteConfirmationDialog() {
    const dialogRef = this.dialog.open(DeleteConfirmationDialogComponent, {
      width: '300px',
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.taskDeleted.emit();
      }
    });
  }
}
