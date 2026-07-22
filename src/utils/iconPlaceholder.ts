/**
 * 图标占位工具
 *
 * 当图标获取失败 / src 为空时，提供「站名首字母 + 由 URL 哈希决定的确定性配色」占位图标，
 * 避免出现空白。相同输入始终生成相同颜色，便于用户凭颜色辨识站点。
 */

/**
 * 根据字符串生成确定性的 HSL 颜色。
 * 使用简单的 32 位哈希（与 Java hashCode 同思路），相同 seed 始终得到同一色相。
 */
export function getDeterministicColor(seed: string): string {
  const s = seed || '?'
  let hash = 0
  for (let i = 0; i < s.length; i++) {
    hash = (hash << 5) - hash + s.charCodeAt(i)
    hash |= 0 // 转为 32 位整数
  }
  const hue = Math.abs(hash) % 360
  return `hsl(${hue}, 62%, 52%)`
}

/**
 * 生成占位图标的首字母：
 * 优先取名称/标题首字符；若为空则取 URL 主机名首字符；最终回退为「?」。
 */
export function getIconInitial(name?: string, url?: string): string {
  const text = (name || '').trim() || (url || '').trim() || '?'
  const first = text.charAt(0)
  return first ? first.toUpperCase() : '?'
}

/**
 * 判断条目是否应使用占位图标：
 * - 无图标对象
 * - 图片类图标（itemType=2）但 src 为空
 */
export function shouldUseIconPlaceholder(icon?: Panel.ItemIcon | null): boolean {
  if (!icon) return true
  if (icon.itemType === 2) {
    return !icon.src || icon.src.trim() === ''
  }
  return false
}

/**
 * 从 URL 中提取主机名（用于占位图标的配色种子与首字母）。
 */
export function getHostnameFromUrl(url?: string): string {
  if (!url) return ''
  try {
    const u = new URL(url)
    return u.hostname || url
  } catch {
    return url
  }
}
