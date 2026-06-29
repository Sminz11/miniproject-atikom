import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzCardModule } from 'ng-zorro-antd/card';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzMessageService } from 'ng-zorro-antd/message';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    NzButtonModule,
    NzCardModule,
    NzIconModule,
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})
export class LoginComponent {
  loading = false;

  constructor(
    private authService: AuthService,
    private router: Router,
    private message: NzMessageService
  ) {}

  onLogin() {
    this.loading = true;

    // Step 1: Get Authorization Code
    this.authService.getAuthorizationCode().subscribe({
      next: (res) => {
        if (res.code === '0000') {
          const code = res.data.authorizationCode;

          // Step 2: Exchange Token
          this.authService.exchangeToken(code).subscribe({
            next: (tokenRes) => {
              this.loading = false;
              if (tokenRes.code === '0000') {
                this.message.success('เข้าสู่ระบบสำเร็จ!');
                this.router.navigate(['/dashboard']);
              } else {
                this.message.error(tokenRes.message);
              }
            },
            error: () => {
              this.loading = false;
              this.message.error('เข้าสู่ระบบไม่สำเร็จ');
            }
          });
        }
      },
      error: () => {
        this.loading = false;
        this.message.error('เกิดข้อผิดพลาด');
      }
    });
  }
}