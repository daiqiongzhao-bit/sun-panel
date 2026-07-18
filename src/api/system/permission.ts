import { post } from '@/utils/request'

export function getPermissionList<T>(data: { pageNum?: number; pageSize?: number }) {
  return post<T>({
    url: '/permission/list',
    data,
  })
}

export function getPermissionByModule<T>(data: { module: string }) {
  return post<T>({
    url: '/permission/byModule',
    data,
  })
}

export function getPermissionMatrix<T>() {
  return post<T>({
    url: '/permission/matrix',
  })
}

export function saveRolePermissions<T>(data: { roleId: number; permissionIds: number[] }) {
  return post<T>({
    url: '/permission/saveRolePermissions',
    data,
  })
}

export function getRolePermissions<T>(data: { roleId: number }) {
  return post<T>({
    url: '/permission/getRolePermissions',
    data,
  })
}

export function createPermission<T>(data: { module: string; name: string; code: string; description?: string }) {
  return post<T>({
    url: '/permission/create',
    data,
  })
}

export function updatePermission<T>(data: { id: number; name: string; description?: string }) {
  return post<T>({
    url: '/permission/update',
    data,
  })
}