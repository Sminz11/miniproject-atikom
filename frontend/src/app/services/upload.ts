import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class UploadService {
  private apiUrl = 'http://localhost:8080/api/v1';

  constructor(private http: HttpClient) {}

  uploadFile(file: File): Observable<any> {
    const formData = new FormData();
    formData.append('file', file);
    return this.http.post(`${this.apiUrl}/uploads`, formData);
  }

  getHistory(params: any): Observable<any> {
    let httpParams = new HttpParams();
    if (params.fileName) httpParams = httpParams.set('fileName', params.fileName);
    if (params.status) httpParams = httpParams.set('status', params.status);
    httpParams = httpParams.set('page', params.page || 1);
    httpParams = httpParams.set('pageSize', params.pageSize || 10);
    return this.http.get(`${this.apiUrl}/uploads`, { params: httpParams });
  }

  getDetail(uploadId: number, page: number = 1, pageSize: number = 10): Observable<any> {
    const params = new HttpParams()
      .set('page', page)
      .set('pageSize', pageSize);
    return this.http.get(`${this.apiUrl}/uploads/${uploadId}/details`, { params });
  }
}