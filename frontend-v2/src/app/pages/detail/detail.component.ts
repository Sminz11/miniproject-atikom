import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { NzTableModule } from 'ng-zorro-antd/table';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzMessageService } from 'ng-zorro-antd/message';
import { NzPopconfirmModule } from 'ng-zorro-antd/popconfirm';
import { UploadService } from '../../services/upload.service';

@Component({
  selector: 'app-detail',
  standalone: true,
  imports: [
    CommonModule,
    NzTableModule,
    NzCardModule,
    NzButtonModule,
    NzTagModule,
    NzIconModule,
    NzPopconfirmModule,
  ],
  templateUrl: './detail.component.html',
  styleUrl: './detail.component.scss'
})
export class DetailComponent implements OnInit {
  uploadId = 0;
  items: any[] = [];
  total = 0;
  page = 1;
  pageSize = 10;
  loading = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private uploadService: UploadService,
    private message: NzMessageService
  ) {}

  ngOnInit() {
    this.uploadId = Number(this.route.snapshot.paramMap.get('id'));
    this.loadData();
  }

  loadData() {
    this.loading = true;
    this.uploadService.getDetail(this.uploadId, this.page, this.pageSize).subscribe({
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

  onPageChange(page: number) {
    this.page = page;
    this.loadData();
  }

  onRetry(detailId: number) {
    this.uploadService.retryDetail(this.uploadId, detailId).subscribe({
      next: (res) => {
        if (res.code === '0000') {
          this.message.success('Retry สำเร็จ รอ Batch ประมวลผล');
          this.loadData();
        } else {
          this.message.error(res.message);
        }
      },
      error: () => this.message.error('Retry ไม่สำเร็จ')
    });
  }

  onExportCSV() {
    this.uploadService.exportCSV(this.uploadId).subscribe({
      next: (blob) => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `upload_${this.uploadId}.csv`;
        a.click();
        window.URL.revokeObjectURL(url);
        this.message.success('Export CSV สำเร็จ');
      },
      error: () => this.message.error('Export CSV ไม่สำเร็จ')
    });
  }

  getStatusColor(status: string): string {
    switch (status) {
      case 'SUCCESS': return 'success';
      case 'FAILED': return 'error';
      case 'PENDING': return 'default';
      default: return 'default';
    }
  }

  goBack() {
    this.router.navigate(['/history']);
  }
}