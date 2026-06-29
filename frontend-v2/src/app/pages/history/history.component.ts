import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { NzTableModule } from 'ng-zorro-antd/table';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzInputModule } from 'ng-zorro-antd/input';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzSpaceModule } from 'ng-zorro-antd/space';
import { UploadService } from '../../services/upload.service';

@Component({
  selector: 'app-history',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    NzTableModule,
    NzCardModule,
    NzButtonModule,
    NzInputModule,
    NzSelectModule,
    NzTagModule,
    NzIconModule,
    NzSpaceModule,
  ],
  templateUrl: './history.component.html',
  styleUrl: './history.component.scss'
})
export class HistoryComponent implements OnInit {
  items: any[] = [];
  total = 0;
  page = 1;
  pageSize = 10;
  filterFileName = '';
  filterStatus = '';
  loading = false;

  constructor(private uploadService: UploadService, private router: Router) {}

  ngOnInit() {
    this.loadData();
  }

  loadData() {
    this.loading = true;
    this.uploadService.getHistory({
      fileName: this.filterFileName,
      status: this.filterStatus,
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
    this.filterFileName = '';
    this.filterStatus = '';
    this.page = 1;
    this.loadData();
  }

  onPageChange(page: number) {
    this.page = page;
    this.loadData();
  }

  goToDetail(uploadId: number) {
    this.router.navigate(['/history', uploadId]);
  }

  getStatusColor(status: string): string {
    switch (status) {
      case 'COMPLETED': return 'success';
      case 'UPLOADED': return 'processing';
      case 'PROCESSING': return 'warning';
      default: return 'default';
    }
  }
}