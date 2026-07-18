import { post } from '@/utils/request'

export function getStickyNoteList<T>() {
  return post<T>({ url: '/panel/stickyNote/getList', data: {} })
}

export function createStickyNote<T>(data: any) {
  return post<T>({ url: '/panel/stickyNote/create', data })
}

export function updateStickyNote<T>(data: any) {
  return post<T>({ url: '/panel/stickyNote/update', data })
}

export function deleteStickyNote<T>(data: { id: number }) {
  return post<T>({ url: '/panel/stickyNote/delete', data })
}