import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Task} from "../lib/task";


export interface TaskUpdate {
  date:   string,
  update: string
}

export interface TaskSummary {
  task: Task,
  updates: TaskUpdate[]
}

@Injectable({
  providedIn: 'root'
})
export class TaskService {

  constructor(private http: HttpClient) { }

  loadTaskList(todayDate: string) {
    return this.http.get<Task[]>('/api/tasks/list/' + todayDate)
  }

  loadTaskSummary(id: number) {
    return this.http.get<TaskSummary>(`/api/tasks/${id}/summary`)
  }

  setTaskTitle(id: number, newValue: string) {
    const payload = {
      title: newValue,
    }
    return this.http.patch(`/api/tasks/${id}/setTitle`, payload)
  }

  setTaskUpdate(id: number, newValue: string) {
    const payload = {
      update: newValue,
    }
    return this.http.patch(`/api/tasks/${id}/setUpdate`, payload)
  }

  setTaskDescription(id: number, newValue: string) {
    const payload = {
      description: newValue,
    }
    return this.http.patch(`/api/tasks/${id}/setDescription`, payload)
  }

  setTaskDone(id: number, task: Task) {
    const payload = {
      done: true,
    }

    return this.http.patch(`/api/tasks/${id}/setDone`, payload)
  }

  setTaskCreated(id: number, task: Task) {
    const payload = {
      done: false,
    }

    return this.http.patch(`/api/tasks/${id}/setDone`, payload)
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

  handleDrop(id: number, currentIndex: number) {
    const payload = {
      newRank: currentIndex,
    }

    return this.http.patch(`/api/tasks/${id}/changeRank`, JSON.stringify(payload))
  }

  importPastTasks(todayDate: string) {
    return this.http.post(`/api/tasks/importPastTasks/${todayDate}`, JSON.stringify({}))
  }
}
