import { post } from '@/utils/request'

export function localSearch<T>(data: { keyword: string; limit?: number }) {
  return post<T>({ url: '/panel/search/local', data })
}