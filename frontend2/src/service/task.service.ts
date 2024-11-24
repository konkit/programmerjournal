import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Task} from "../lib/task";

@Injectable({
  providedIn: 'root'
})
export class TaskService {

  constructor(private http: HttpClient) { }

  loadTaskList(todayDate: string) {
    return this.http.get<Task[]>('/api/tasks/list/' + todayDate)
  }

  setTaskTitle(id: number, newValue: string) {
    const payload = {
      title: newValue,
    }
    return this.http.patch(`/api/tasks/${id}/setTitle`, payload)
  }

  setTaskDone(id: number, task: Task) {
    let currentValue = task.status == "Done"
    let newValue = !currentValue

    const payload = {
      done: newValue,
    }

    return this.http.patch(`/api/tasks/${id}/setDone`, payload)
  }

  setTaskDescription(id: number, date: string, newValue: string) {
    const payload = {
      update: newValue,
    }
    return this.http.patch(`/api/tasks/${id}/setUpdate`, payload)
  }

  createTask(title: string, date: string) {
    const payload = {
      "title": title,
      "createdDate": date,
    }

    return this.http.post('/api/tasks/create', payload)
  }

  deleteTask(taskId: number) {
    return this.http.delete(`/api/tasks/${taskId}/delete`)
  }

  snoozeTask(taskId: number, date: string) {
    const payload = {
      date: date,
    }

    return this.http.patch(`/api/tasks/${taskId}/snooze`, JSON.stringify(payload))
  }
}
