import { post } from '@/utils/request'

export function getNoticeList<T>(data: { page: number; pageSize: number; keyword?: string; noticeType?: number }) {
  return post<T>({ url: '/noticeManage/getList', data })
}

export function getNoticeById<T>(data: { id: number }) {
  return post<T>({ url: '/noticeManage/get', data })
}

export function createNotice<T>(data: any) {
  return post<T>({ url: '/noticeManage/create', data })
}

export function updateNotice<T>(data: any) {
  return post<T>({ url: '/noticeManage/update', data })
}

export function deleteNotice<T>(data: { id: number }) {
  return post<T>({ url: '/noticeManage/delete', data })
}

export function getVisibleNotices<T>(data: { displayType: number[] }) {
  return post<T>({ url: '/noticeManage/getVisibleNotices', data })
}