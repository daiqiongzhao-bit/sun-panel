import { post } from '@/utils/request'

export function getGalleryList<T>(data: { pageNum: number; pageSize: number; categoryId?: number; type?: number; keyword?: string }) {
  return post<T>({
    url: '/gallery/list',
    data,
  })
}

export function getGalleryById<T>(data: { id: number }) {
  return post<T>({
    url: '/gallery/get',
    data,
  })
}

export function uploadGallery<T>(data: FormData) {
  return post<T>({
    url: '/gallery/upload',
    data,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function updateGallery<T>(data: { id: number; name?: string; categoryId?: number }) {
  return post<T>({
    url: '/gallery/update',
    data,
  })
}

export function deleteGallery<T>(data: { id: number }) {
  return post<T>({
    url: '/gallery/delete',
    data,
  })
}

export function batchDeleteGallery<T>(data: { ids: number[] }) {
  return post<T>({
    url: '/gallery/batchDelete',
    data,
  })
}

export function getGalleryCategories<T>(data: { type?: number }) {
  return post<T>({
    url: '/gallery/category/list',
    data,
  })
}

export function createGalleryCategory<T>(data: { name: string; type?: number; description?: string; sort?: number }) {
  return post<T>({
    url: '/gallery/category/create',
    data,
  })
}

export function updateGalleryCategory<T>(data: { id: number; name: string; description?: string; sort?: number }) {
  return post<T>({
    url: '/gallery/category/update',
    data,
  })
}

export function deleteGalleryCategory<T>(data: { id: number }) {
  return post<T>({
    url: '/gallery/category/delete',
    data,
  })
}