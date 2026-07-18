import { post } from '@/utils/request'

export function getJobList<T>(data: { page: number; pageSize: number; keyword?: string }) {
  return post<T>({ url: '/job/getList', data })
}

export function createJob<T>(data: any) {
  return post<T>({ url: '/job/create', data })
}

export function updateJob<T>(data: any) {
  return post<T>({ url: '/job/update', data })
}

export function pauseJob<T>(data: { id: number }) {
  return post<T>({ url: '/job/pause', data })
}

export function startJob<T>(data: { id: number }) {
  return post<T>({ url: '/job/start', data })
}

export function runJobNow<T>(data: { id: number }) {
  return post<T>({ url: '/job/runNow', data })
}

export function deleteJob<T>(data: { id: number }) {
  return post<T>({ url: '/job/delete', data })
}

export function getJobLogList<T>(data: { jobId: number; page: number; pageSize: number }) {
  return post<T>({ url: '/job/getLogList', data })
}

export function previewCron<T>(data: { cronExpr: string }) {
  return post<T>({ url: '/job/previewCron', data })
}
