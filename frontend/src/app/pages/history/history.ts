import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { UploadService } from '../../services/upload';

@Component({
  selector: 'app-history',
  imports: [CommonModule, FormsModule],
  templateUrl: './history.html',
  styleUrl: './history.scss'
})
export class History implements OnInit {
  items = signal<any[]>([]);
  total = signal<number>(0);
  page = signal<number>(1);
  pageSize = 10;
  filterFileName: string = '';
  filterStatus: string = '';
  loading = signal<boolean>(false);

  constructor(private uploadService: UploadService, private router: Router) {}

  ngOnInit() {
    this.loadData();
  }

  loadData() {
    this.loading.set(true);
    this.uploadService.getHistory({
      fileName: this.filterFileName,
      status: this.filterStatus,
      page: this.page(),
      pageSize: this.pageSize
    }).subscribe({
      next: (res) => {
        this.loading.set(false);
        if (res.code === '0000') {
          this.items.set(res.data.items || []);
          this.total.set(res.data.total);
        }
      },
      error: () => { this.loading.set(false); }
    });
  }

  onSearch() {
    this.page.set(1);
    this.loadData();
  }

  onReset() {
    this.filterFileName = '';
    this.filterStatus = '';
    this.page.set(1);
    this.loadData();
  }

  goToDetail(uploadId: number) {
    this.router.navigate(['/history', uploadId]);
  }

  prevPage() {
    if (this.page() > 1) { this.page.set(this.page() - 1); this.loadData(); }
  }

  nextPage() {
    if (this.page() * this.pageSize < this.total()) { this.page.set(this.page() + 1); this.loadData(); }
  }
}