<script setup lang="ts">
import { NAvatar, NImage } from 'naive-ui'
import { computed, ref } from 'vue'
import { SvgIconOnline } from '@/components/common'
import { getDeterministicColor, getIconInitial } from '@/utils/iconPlaceholder'

interface Prop {
  itemIcon?: Panel.ItemIcon | null
  size?: number // 默认70
  forceBackground?: string // 强制背景色
  // 用于占位图标的辅助信息：站点标题/名称与 URL（决定首字母与确定性配色）
  name?: string
  url?: string
}

const props = withDefaults(defineProps<Prop>(), { size: 70 })
const defaultBackground = '#2a2a2a6b'
const defaultStyle = ref({
  width: `${props.size}px`,
  height: `${props.size}px`,
})

// 取出扩展名时先去掉查询串/锚点，确保 .ico / .svg 等带参 URL 也能被正确识别为图标
const iconExt = computed(() => {
  const src = props.itemIcon?.src
  if (!src) return ''
  const clean = src.split(/[?#]/)[0]
  const ext = clean.split('.').pop()
  return ext ? ext.toLowerCase() : ''
})

// 占位图标：首字母 + 由 URL/名称哈希决定的确定性配色
const placeholderSeed = computed(() => props.url || props.name || props.itemIcon?.text || '?')
const placeholderColor = computed(() => getDeterministicColor(placeholderSeed.value))
const placeholderInitial = computed(() => getIconInitial(props.name || props.itemIcon?.text, props.url))

// 是否存在有效的图片类图标（itemType=2 且 src 非空）
const hasValidImageIcon = computed(() => props.itemIcon?.itemType === 2 && !!props.itemIcon?.src)
</script>

<template>
  <div class="item-icon" :style="defaultStyle">
    <slot>
      <template v-if="itemIcon">
        <template v-if="itemIcon?.itemType === 1">
          <NAvatar :size="props.size" :style="{ backgroundColor: (forceBackground ?? itemIcon?.backgroundColor) || defaultBackground }">
            {{ itemIcon.text }}
          </NAvatar>
        </template>

        <template v-else-if="hasValidImageIcon">
          <div v-if="iconExt === 'svg'" :style="{ backgroundColor: (forceBackground ?? itemIcon?.backgroundColor) || defaultBackground, ...defaultStyle }" class="flex justify-center items-center">
            <img :src="itemIcon?.src" class="w-[35px] h-[35px]">
          </div>
          <NImage v-else :style="{ backgroundColor: (forceBackground ?? itemIcon?.backgroundColor) || defaultBackground, ...defaultStyle }" :src="itemIcon?.src" preview-disabled />
        </template>

        <template v-else-if="itemIcon?.itemType === 3">
          <NAvatar :size="props.size" :style="{ backgroundColor: (forceBackground ?? itemIcon?.backgroundColor) || defaultBackground }">
            <SvgIconOnline style="font-size: 35px;" :icon="itemIcon.text" />
          </NAvatar>
        </template>

        <!-- 图片类图标缺 src 或未知类型：展示占位图标 -->
        <template v-else>
          <NAvatar :size="props.size" style="color:#fff" :style="{ backgroundColor: placeholderColor, ...defaultStyle }">
            {{ placeholderInitial }}
          </NAvatar>
        </template>
      </template>

      <!-- 无图标对象：展示占位图标，避免空白 -->
      <template v-else>
        <NAvatar :size="props.size" style="color:#fff" :style="{ backgroundColor: placeholderColor, ...defaultStyle }">
          {{ placeholderInitial }}
        </NAvatar>
      </template>
    </slot>
  </div>
</template>
