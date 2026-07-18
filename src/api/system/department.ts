import { post } from '@/utils/request'

export function getDepartmentList<T>() {
  return post<T>({
    url: '/department/list',
  })
}

export function getDepartmentTree<T>() {
  return post<T>({
    url: '/department/tree',
  })
}

export function getDepartmentById<T>(data: { id: number }) {
  return post<T>({
    url: '/department/get',
    data,
  })
}

export function createDepartment<T>(data: { name: string; code: string; parentId?: number; description?: string; leaderId?: number; sort?: number }) {
  return post<T>({
    url: '/department/create',
    data,
  })
}

export function updateDepartment<T>(data: { id: number; name: string; code: string; parentId?: number; description?: string; leaderId?: number; sort?: number }) {
  return post<T>({
    url: '/department/update',
    data,
  })
}

export function deleteDepartment<T>(data: { id: number }) {
  return post<T>({
    url: '/department/delete',
    data,
  })
}