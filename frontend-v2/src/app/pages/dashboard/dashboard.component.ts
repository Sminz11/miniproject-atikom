import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzStatisticModule } from 'ng-zorro-antd/statistic';
import { NzGridModule } from 'ng-zorro-antd/grid';
import { NzTableModule } from 'ng-zorro-antd/table';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { UploadService } from '../../services/upload.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    NzCardModule,
    NzStatisticModule,
    NzGridModule,
    NzTableModule,
    NzTagModule,
    NzIconModule,
  ],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss'
})
export class DashboardComponent implements OnInit {
  summary: any = null;
  loading = true;

  constructor(private uploadService: UploadService) {}

  ngOnInit() {
    this.loadSummary();
  }

  loadSummary() {
    this.loading = true;
    this.uploadService.getDashboardSummary().subscribe({
      next: (res) => {
        this.loading = false;
        if (res.code === '0000') {
          this.summary = res.data;
        }
      },
      error: () => { this.loading = false; }
    });
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