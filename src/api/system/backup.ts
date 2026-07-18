import { post } from '@/utils/request'

export function createBackup<T>(data: { mode: number; name?: string }) {
  return post<T>({
    url: '/backup/create',
    data,
  })
}

export function getBackupList<T>(data: { pageNum: number; pageSize: number; mode?: number }) {
  return post<T>({
    url: '/backup/list',
    data,
  })
}

export function getBackupById<T>(data: { id: number }) {
  return post<T>({
    url: '/backup/get',
    data,
  })
}

export function restoreBackup<T>(data: { id: number }) {
  return post<T>({
    url: '/backup/restore',
    data,
  })
}

export function deleteBackup<T>(data: { id: number }) {
  return post<T>({
    url: '/backup/delete',
    data,
  })
}

export function exportBackup<T>(data: { id: number }) {
  return post<T>({
    url: '/backup/export',
    data,
    responseType: 'blob',
  })
}

export function importBackup<T>(data: FormData) {
  return post<T>({
    url: '/backup/import',
    data,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function getBackupTaskList<T>() {
  return post<T>({
    url: '/backup/task/list',
  })
}

export function createBackupTask<T>(data: { name: string; mode: number; cronExpr: string; retentionDays?: number; status?: number }) {
  return post<T>({
    url: '/backup/task/create',
    data,
  })
}

export function updateBackupTask<T>(data: { id: number; name: string; mode: number; cronExpr: string; retentionDays?: number; status?: number }) {
  return post<T>({
    url: '/backup/task/update',
    data,
  })
}

export function deleteBackupTask<T>(data: { id: number }) {
  return post<T>({
    url: '/backup/task/delete',
    data,
  })
}