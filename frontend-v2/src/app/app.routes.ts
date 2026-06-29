import { Routes } from '@angular/router';
import { authGuard } from './guards/auth.guard';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'login',
    pathMatch: 'full'
  },
  {
    path: 'login',
    loadComponent: () =>
      import('./pages/login/login.component').then(m => m.LoginComponent)
  },
  {
    path: 'dashboard',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./pages/dashboard/dashboard.component').then(m => m.DashboardComponent)
  },
  {
    path: 'upload',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./pages/upload/upload.component').then(m => m.UploadComponent)
  },
  {
    path: 'history',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./pages/history/history.component').then(m => m.HistoryComponent)
  },
  {
    path: 'history/:id',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./pages/detail/detail.component').then(m => m.DetailComponent)
  },
  {
    path: 'audit-log',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./pages/audit-log/audit-log.component').then(m => m.AuditLogComponent)
  },
  {
    path: '**',
    redirectTo: 'login'
  }
];