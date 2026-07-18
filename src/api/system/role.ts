import { post } from '@/utils/request'

export function getRoleList<T>(data: { pageNum: number; pageSize: number; name?: string }) {
  return post<T>({
    url: '/role/list',
    data,
  })
}

export function getRoleById<T>(data: { id: number }) {
  return post<T>({
    url: '/role/get',
    data,
  })
}

export function createRole<T>(data: { name: string; description?: string; status?: number }) {
  return post<T>({
    url: '/role/create',
    data,
  })
}

export function updateRole<T>(data: { id: number; name?: string; description?: string; status?: number }) {
  return post<T>({
    url: '/role/update',
    data,
  })
}

export function deleteRole<T>(data: { id: number }) {
  return post<T>({
    url: '/role/delete',
    data,
  })
}