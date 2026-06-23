import { Component, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { UploadService } from '../../services/upload';

@Component({
  selector: 'app-upload',
  imports: [CommonModule],
  templateUrl: './upload.html',
  styleUrl: './upload.scss'
})
export class Upload {
  selectedFile = signal<File | null>(null);
  result = signal<any>(null);
  error = signal<string>('');
  loading = signal<boolean>(false);

  constructor(private uploadService: UploadService, private router: Router) {}

  onFileSelect(event: any) {
    const file = event.target.files[0];
    if (file && file.name.endsWith('.txt')) {
      this.selectedFile.set(file);
      this.error.set('');
    } else {
      this.error.set('กรุณาเลือกไฟล์ .txt เท่านั้น');
      this.selectedFile.set(null);
    }
  }

  onUpload() {
    const file = this.selectedFile();
    if (!file) {
      this.error.set('กรุณาเลือกไฟล์ก่อน');
      return;
    }
    this.loading.set(true);
    this.result.set(null);
    this.error.set('');

    this.uploadService.uploadFile(file).subscribe({
      next: (res) => {
        this.loading.set(false);
        if (res.code === '0000') {
          this.result.set(res.data);
        } else {
          this.error.set(res.message);
        }
      },
      error: (err) => {
        this.loading.set(false);
        this.error.set('เกิดข้อผิดพลาด: ' + err.message);
      }
    });
  }

  goToHistory() {
    this.router.navigate(['/history']);
  }
}