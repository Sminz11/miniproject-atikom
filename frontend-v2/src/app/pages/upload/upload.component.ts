import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { NzUploadModule } from 'ng-zorro-antd/upload';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzMessageService } from 'ng-zorro-antd/message';
import { NzAlertModule } from 'ng-zorro-antd/alert';
import { NzDescriptionsModule } from 'ng-zorro-antd/descriptions';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { UploadService } from '../../services/upload.service';

@Component({
  selector: 'app-upload',
  standalone: true,
  imports: [
    CommonModule,
    NzUploadModule,
    NzButtonModule,
    NzCardModule,
    NzIconModule,
    NzAlertModule,
    NzDescriptionsModule,
    NzTagModule,
  ],
  templateUrl: './upload.component.html',
  styleUrl: './upload.component.scss'
})
export class UploadComponent {
  selectedFile: File | null = null;
  result: any = null;
  error: string = '';
  loading = false;

  constructor(
    private uploadService: UploadService,
    private router: Router,
    private message: NzMessageService
  ) {}

  onDragEnter(event: MouseEvent): void {
    const target = event.target as HTMLElement;
    target.style.borderColor = '#1890ff';
  }

  onDragLeave(event: MouseEvent): void {
    const target = event.target as HTMLElement;
    target.style.borderColor = '#d9d9d9';
  }

  onFileSelect(event: any) {
    const file = event.target.files[0];
    if (file && file.name.endsWith('.txt')) {
      this.selectedFile = file;
      this.error = '';
    } else {
      this.error = 'กรุณาเลือกไฟล์ .txt เท่านั้น';
      this.selectedFile = null;
    }
  }

  onUpload() {
    if (!this.selectedFile) {
      this.error = 'กรุณาเลือกไฟล์ก่อน';
      return;
    }
    this.loading = true;
    this.result = null;
    this.error = '';

    this.uploadService.uploadFile(this.selectedFile).subscribe({
      next: (res) => {
        this.loading = false;
        if (res.code === '0000') {
          this.result = res.data;
          this.message.success('อัปโหลดสำเร็จ!');
        } else {
          this.error = res.message;
        }
      },
      error: (err) => {
        this.loading = false;
        this.error = 'เกิดข้อผิดพลาด: ' + err.message;
      }
    });
  }

  goToHistory() {
    this.router.navigate(['/history']);
  }
}