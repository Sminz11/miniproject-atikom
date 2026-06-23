import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { UploadService } from '../../services/upload';

@Component({
  selector: 'app-detail',
  imports: [CommonModule],
  templateUrl: './detail.html',
  styleUrl: './detail.scss'
})
export class Detail implements OnInit {
  uploadId: number = 0;
  items: any[] = [];
  total: number = 0;
  page: number = 1;
  pageSize: number = 10;
  loading: boolean = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private uploadService: UploadService
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

  prevPage() {
    if (this.page > 1) { this.page--; this.loadData(); }
  }

  nextPage() {
    if (this.page * this.pageSize < this.total) { this.page++; this.loadData(); }
  }

  goBack() {
    this.router.navigate(['/history']);
  }
}