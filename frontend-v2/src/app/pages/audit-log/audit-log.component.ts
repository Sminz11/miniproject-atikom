import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NzTableModule } from 'ng-zorro-antd/table';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzSpaceModule } from 'ng-zorro-antd/space';
import { UploadService } from '../../services/upload.service';

@Component({
  selector: 'app-audit-log',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    NzTableModule,
    NzCardModule,
    NzButtonModule,
    NzSelectModule,
    NzIconModule,
    NzTagModule,
    NzSpaceModule,
  ],
  templateUrl: './audit-log.component.html',
  styleUrl: './audit-log.component.scss'
})
export class AuditLogComponent implements OnInit {
  items: any[] = [];
  total = 0;
  page = 1;
  pageSize = 10;
  filterAction = '';
  loading = false;

  actionOptions = [
    { label: 'All Actions', value: '' },
    { label: 'LOGIN_SUCCESS', value: 'LOGIN_SUCCESS' },
    { label: 'LOGOUT', value: 'LOGOUT' },
    { label: 'UPLOAD_FILE', value: 'UPLOAD_FILE' },
    { label: 'VIEW_UPLOAD_DETAIL', value: 'VIEW_UPLOAD_DETAIL' },
    { label: 'RETRY_TRANSACTION', value: 'RETRY_TRANSACTION' },
    { label: 'EXPORT_CSV', value: 'EXPORT_CSV' },
    { label: 'BATCH_PROCESS', value: 'BATCH_PROCESS' },
  ];

  constructor(private uploadService: UploadService) {}

  ngOnInit() {
    this.loadData();
  }

  loadData() {
    this.loading = true;
    this.uploadService.getAuditLogs({
      action: this.filterAction,
      page: this.page,
      pageSize: this.pageSize
    }).subscribe({
      next: (res) => {
        this.loading = false;
        if (res.code === '0000') {
          this.items = res.data.items || [];
          this.total = res.data.total;
        }
      },
      error: () => { this.loading = false; }
    });
  }

  onSearch() {
    this.page = 1;
    this.loadData();
  }

  onReset() {
    this.filterAction = '';
    this.page = 1;
    this.loadData();
  }

  onPageChange(page: number) {
    this.page = page;
    this.loadData();
  }

  getActionColor(action: string): string {
    switch (action) {
      case 'LOGIN_SUCCESS': return 'success';
      case 'LOGOUT': return 'default';
      case 'UPLOAD_FILE': return 'processing';
      case 'RETRY_TRANSACTION': return 'warning';
      case 'EXPORT_CSV': return 'cyan';
      case 'BATCH_PROCESS': return 'purple';
      default: return 'default';
    }
  }
}