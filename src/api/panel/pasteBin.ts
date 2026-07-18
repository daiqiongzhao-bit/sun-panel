import { post } from '@/utils/request'

export function createPasteBin<T>(data: any) {
  return post<T>({ url: '/panel/pasteBin/create', data })
}

export function updatePasteBin<T>(data: any) {
  return post<T>({ url: '/panel/pasteBin/update', data })
}

export function getPasteBinByCode<T>(data: { code: string; password?: string }) {
  return post<T>({ url: '/panel/pasteBin/getByCode', data })
}

export function getMyPasteBinList<T>(data: { page: number; pageSize: number }) {
  return post<T>({ url: '/panel/pasteBin/myList', data })
}

export function deletePasteBin<T>(data: { id: number }) {
  return post<T>({ url: '/panel/pasteBin/delete', data })
}

export function getPasteBinAccessUrl<T>(data: { code: string }) {
  return post<T>({ url: '/panel/pasteBin/accessUrl', data })
}

export function downloadPasteBinFile(code: string, password?: string) {
  return post<Blob>({
    url: '/panel/pasteBin/downloadFile',
    data: { code, password },
    responseType: 'blob',
  } as any)
}