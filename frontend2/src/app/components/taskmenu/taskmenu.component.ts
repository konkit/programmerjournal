import { Component } from '@angular/core';
import {MatCardModule} from '@angular/material/card';
import {MatMenuModule} from '@angular/material/menu'
import {MatButtonModule } from '@angular/material/button'
import {MatIconModule} from '@angular/material/icon'

@Component({
  selector: 'app-taskmenu',
  imports: [MatIconModule, MatMenuModule, MatButtonModule],
  templateUrl: './taskmenu.component.html',
  styleUrl: './taskmenu.component.scss',
  standalone: true,

})
export class TaskmenuComponent {

}
