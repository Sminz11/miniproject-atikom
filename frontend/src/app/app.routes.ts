import { Routes } from '@angular/router';
import { Upload } from './pages/upload/upload';
import { History } from './pages/history/history';
import { Detail } from './pages/detail/detail';

export const routes: Routes = [
  { path: '', redirectTo: 'upload', pathMatch: 'full' },
  { path: 'upload', component: Upload },
  { path: 'history', component: History },
  { path: 'history/:id', component: Detail }
];
