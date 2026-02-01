import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { interval, switchMap, catchError, of, map, retry } from 'rxjs';
import { BASE_PATH } from '../../frontend-client';
import { Inject } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class BackendStatusService {
  isBackendDown = signal<boolean>(false);

  constructor(private http: HttpClient, @Inject(BASE_PATH) private basePath: string) {
    this.startHealthCheck();
  }

  private startHealthCheck() {
    interval(5000) // Check every 5 seconds
      .pipe(
        switchMap(() => this.checkBackendHealth())
      )
      .subscribe();
  }

  private checkBackendHealth() {
    // We use a lightweight call to check connectivity.
    // Using a known endpoint that returns data (even if empty) is better than relying on 404s.
    // /api/entries/list/2024-01-01 is likely to return an empty list quickly.
    // Note: In development, BASE_PATH might be http://localhost:4200 (frontend) but proxy forwards to 8080.
    // If BASE_PATH is set to backend URL directly, it works too.
    // The current app configuration sets BASE_PATH to current location (frontend), so it relies on proxy.
    return this.http.get(this.basePath + '/api/entries/list/2024-01-01', { observe: 'response' })
      .pipe(
        map(() => {
          this.isBackendDown.set(false);
          return true;
        }),
        catchError(err => {
            // If status is 0, it usually means network error or server down (CORS, connection refused, etc.)
            // 502, 503, 504 are also indicators of backend issues.
            if (err.status === 0 || err.status == 500 || err.status === 502 || err.status === 503 || err.status === 504) {
                this.isBackendDown.set(true);
            } else {
                // 4xx or other 500s might still mean backend is reachable but erroring on logic
                this.isBackendDown.set(false);
            }
            return of(false);
        })
      );
  }
}
