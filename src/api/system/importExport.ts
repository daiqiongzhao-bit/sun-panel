import { post } from '@/utils/request'

export function exportUsers<T>() {
  return post<T>({ url: '/system/exportUsers' })
}

export function importUsers<T>(data: { users: any[] }) {
  return post<T>({ url: '/system/importUsers', data })
}

export function exportRoles<T>() {
  return post<T>({ url: '/system/exportRoles' })
}

export function importRoles<T>(data: { roles: any[] }) {
  return post<T>({ url: '/system/importRoles', data })
}