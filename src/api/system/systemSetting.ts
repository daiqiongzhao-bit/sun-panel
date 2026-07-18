import { post } from '@/utils/request'

export function getLogoConfig<T>() {
  return post<T>({
    url: '/systemSetting/logo/get',
  })
}

export function setLogoConfig<T>(data: { imageUrl: string; size: number; useCDN?: boolean; cdnUrl?: string }) {
  return post<T>({
    url: '/systemSetting/logo/set',
    data,
  })
}

export function uploadLogo<T>(data: FormData) {
  return post<T>({
    url: '/systemSetting/logo/upload',
    data,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function getBackgroundConfig<T>() {
  return post<T>({
    url: '/systemSetting/background/get',
  })
}

export function setBackgroundConfig<T>(data: { imageUrl: string; displayMode: string; useCustomUrl?: boolean; customUrl?: string }) {
  return post<T>({
    url: '/systemSetting/background/set',
    data,
  })
}

export function uploadBackground<T>(data: FormData) {
  return post<T>({
    url: '/systemSetting/background/upload',
    data,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function getPresetBackgrounds<T>() {
  return post<T>({
    url: '/systemSetting/background/preset',
  })
}