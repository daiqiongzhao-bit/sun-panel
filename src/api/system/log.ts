import { post } from '@/utils/request'

export function getOperationLogList<T>(data: { page: number; pageSize: number; keyword?: string; module?: string; startTime?: number; endTime?: number }) {
  return post<T>({ url: '/log/operationList', data })
}

export function getLoginLogList<T>(data: { page: number; pageSize: number; keyword?: string; startTime?: number; endTime?: number }) {
  return post<T>({ url: '/log/loginList', data })
}

export function clearOperationLog<T>(data?: Record<string, never>) {
  return post<T>({ url: '/log/clearOperation', data: data || {} })
}

export function clearLoginLog<T>(data?: Record<string, never>) {
  return post<T>({ url: '/log/clearLogin', data: data || {} })
}