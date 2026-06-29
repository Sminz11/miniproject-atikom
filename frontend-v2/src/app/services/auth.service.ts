import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private apiUrl = 'http://localhost:8080/api/v1';

  constructor(private http: HttpClient) {}

  getAuthorizationCode(): Observable<any> {
    return this.http.get(`${this.apiUrl}/oauth/authorize`);
  }

  exchangeToken(code: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/oauth/token`, { code }).pipe(
      tap((res: any) => {
        if (res.code === '0000') {
          localStorage.setItem('access_token', res.data.accessToken);
          localStorage.setItem('username', 'intern_user');
        }
      })
    );
  }

  getMe(): Observable<any> {
    return this.http.get(`${this.apiUrl}/me`);
  }

  logout(): Observable<any> {
    return this.http.post(`${this.apiUrl}/logout`, {}).pipe(
      tap(() => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('username');
      })
    );
  }

  isLoggedIn(): boolean {
    return !!localStorage.getItem('access_token');
  }

  getToken(): string {
    return localStorage.getItem('access_token') || '';
  }

  getUsername(): string {
    return localStorage.getItem('username') || '';
  }
}