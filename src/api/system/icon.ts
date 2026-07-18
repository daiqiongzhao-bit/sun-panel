import { post } from '@/utils/request'

export function getIconById<T>(data: { id: number }) {
  return post<T>({
    url: '/icon/get',
    data,
  })
}

export function getIconByName<T>(data: { name: string }) {
  return post<T>({
    url: '/icon/getByName',
    data,
  })
}

export function getIconList<T>(data: { pageNum: number; pageSize: number; categoryId?: number; keyword?: string }) {
  return post<T>({
    url: '/icon/list',
    data,
  })
}

export function getIconBatch<T>(data: { ids?: number[]; categoryId?: number }) {
  return post<T>({
    url: '/icon/batch',
    data,
  })
}

export function createIcon<T>(data: { name: string; path: string; categoryId?: number; description?: string }) {
  return post<T>({
    url: '/icon/create',
    data,
  })
}

export function updateIcon<T>(data: { id: number; name: string; path: string; categoryId?: number; description?: string }) {
  return post<T>({
    url: '/icon/update',
    data,
  })
}

export function deleteIcon<T>(data: { id: number }) {
  return post<T>({
    url: '/icon/delete',
    data,
  })
}

export function toggleIconFavorite<T>(data: { iconId: number }) {
  return post<T>({
    url: '/icon/favorite',
    data,
  })
}

export function getFavoriteIcons<T>() {
  return post<T>({
    url: '/icon/favorites',
  })
}

export function getIconCategories<T>() {
  return post<T>({
    url: '/icon/category/list',
  })
}

export function createIconCategory<T>(data: { name: string; description?: string; sort?: number }) {
  return post<T>({
    url: '/icon/category/create',
    data,
  })
}

export function updateIconCategory<T>(data: { id: number; name: string; description?: string; sort?: number }) {
  return post<T>({
    url: '/icon/category/update',
    data,
  })
}

export function deleteIconCategory<T>(data: { id: number }) {
  return post<T>({
    url: '/icon/category/delete',
    data,
  })
}