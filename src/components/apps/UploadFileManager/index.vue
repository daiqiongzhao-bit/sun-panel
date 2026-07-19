<script setup lang="ts">
import { NAlert, NButton, NButtonGroup, NCard, NEllipsis, NGrid, NGridItem, NImage, NImageGroup, NSpin, NUpload, useDialog, useMessage, NTag } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { deletes, getList } from '@/api/system/file'
import { set as savePanelConfig } from '@/api/panel/userConfig'
import { setBackgroundConfig as saveLoginBg, setLogoConfig } from '@/api/system/systemSetting'
import { RoundCardModal, SvgIcon } from '@/components/common'
import { copyToClipboard, timeFormat } from '@/utils/cmn'
import { t } from '@/locales'
import { usePanelState } from '@/store'

const FILE_CATEGORIES = [
  { label: '全部文件', value: 'all' },
  { label: '🖼️ 壁纸', value: 'wallpaper' },
  { label: '🎨 图标', value: 'icon' },
]

interface InfoModalState {
  title: string
  show: boolean
  fileInfo: File.Info | null
}
const imageList = ref<File.Info[]>([])
const ms = useMessage()
const dialog = useDialog()
const panelStore = usePanelState()
const loading = ref(false)
const uploadLoading = ref(false)
const selectedCategory = ref('all')
const infoModalState = ref<InfoModalState>({
  show: false,
  title: '',
  fileInfo: null,
})

// 根据文件路径和名称判断文件类型（壁纸 vs 图标两大类）
function getFileCategory(file: File.Info): string {
  const src = (file.src || '').toLowerCase()
  const fileName = (file.fileName || '').toLowerCase()

  // 图标类：明确的图标路径/名称特征
  if (src.includes('/icon/') || src.includes('/logo/') || fileName.includes('icon') || fileName.includes('favicon')
    || fileName.includes('.ico') || fileName.includes('logo'))
    return 'icon'

  // 图标类：小尺寸图片（<200KB）且文件名为 UUID/哈希格式（典型的 favicon 自动获取结果）
  const fileSize = Number(file.size) || 0
  const isSmallImage = fileSize > 0 && fileSize < 200 * 1024 // < 200KB
  const isHashLikeFileName = /^[a-f0-9\-]{8,}\.(png|jpg|jpeg|ico|svg|webp)$/.test(fileName)
  if (isSmallImage && isHashLikeFileName)
    return 'icon'

  // 壁纸类（background目录、登录背景、主页背景、wallpaper等，其余默认也归壁纸）
  if (src.includes('/background/') || fileName.includes('background') || fileName.includes('bg')
    || fileName.includes('loginbg') || fileName.includes('壁纸') || fileName.includes('wallpaper')
    || fileName.includes('风景') || fileName.includes('桌面') || fileName.includes('landscape'))
    return 'wallpaper'
  // 默认为壁纸类（大多数上传的是壁纸/背景图）
  return 'wallpaper'
}

// 获取分类的显示标签
function getCategoryLabel(category: string): string {
  return FILE_CATEGORIES.find(c => c.value === category)?.label || '未知类型'
}

// 获取分类对应的标签颜色
function getCategoryColor(category: string): string {
  switch (category) {
    case 'icon': return 'info'
    case 'wallpaper': return 'default'
    default: return 'default'
  }
}

const filteredImageList = computed(() => {
  if (selectedCategory.value === 'all')
    return imageList.value
  return imageList.value.filter(item => getFileCategory(item) === selectedCategory.value)
})

async function handleBatchUpload(options: { file: UploadFileInfo; fileList: Array<UploadFileInfo> }) {
  if (options.file.status === 'finished') {
    getFileList()
  }
}

function handleUploadError() {
  ms.error(t('common.uploadFail'))
}

async function getFileList() {
  loading.value = true
  const { data } = await getList<Common.ListResponse<File.Info[]>>()
  imageList.value = data.list
  loading.value = false
}

async function copyImageUrl(text: string) {
  const res = await copyToClipboard(text)
  if (res)
    ms.success(t('apps.uploadsFileManager.copySuccess'))

  else
    ms.error(t('apps.uploadsFileManager.copyFailed'))
}

function handleDelete(id: number) {
  dialog.warning({
    title: t('common.warning'),
    content: t('apps.uploadsFileManager.deleteWarningText'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      deletesImges(id)
    },
  })
}

async function deletesImges(id: number) {
  try {
    const { code, msg } = await deletes([id])
    if (code === 0) {
      getFileList()
      ms.success(t('common.success'))
    }
    else {
      ms.error(`${t('common.failed')}:${msg}`)
    }
  }
  catch (error) {
    ms.error(t('common.failed'))
  }
}

function handleInfoClick(fileInfo: File.Info) {
  infoModalState.value.fileInfo = fileInfo
  infoModalState.value.show = true
}

function handleSetWallpaper(imgSrc: string) {
  panelStore.panelConfig.backgroundImageSrc = imgSrc
  savePanelConfig({ panel: panelStore.panelConfig })
  ms.success('已设置为主页背景')
}

function handleSetLoginBg(imgSrc: string) {
  saveLoginBg({ imageUrl: imgSrc, displayMode: 'cover', useCustomUrl: false, customUrl: '' }).then((res) => {
    if (res.code === 0)
      ms.success('已设置为登录页背景')
    else
      ms.error('设置登录页背景失败')
  })
}

function handleSetLogo(imgSrc: string) {
  setLogoConfig({ imageUrl: imgSrc, size: 80, useCDN: false, cdnUrl: '' }).then((res: any) => {
    if (res.code === 0)
      ms.success('已设置为系统Logo')
    else
      ms.error('设置Logo失败')
  })
}

function handleQuickSet(item: File.Info) {
  const category = getFileCategory(item)
  switch (category) {
    case 'icon':
      handleSetLogo(item.src)
      break
    default:
      // 壁纸类：根据原始特征判断设为主页背景还是登录背景
      const src = (item.src || '').toLowerCase()
      const fileName = (item.fileName || '').toLowerCase()
      if (src.includes('/background/') || fileName.includes('loginbg') || fileName.includes('login'))
        handleSetLoginBg(item.src)
      else
        handleSetWallpaper(item.src)
      break
  }
}

onMounted(() => {
  getFileList()
})
</script>

<template>
  <div class="bg-slate-200 dark:bg-zinc-900 p-2 h-full">
    <NSpin v-show="loading" size="small" />
    <NAlert type="info" :bordered="false">
      {{ $t('apps.uploadsFileManager.alertText') }}
    </NAlert>
    <div class="my-2">
      <NUpload
        multiple
        accept="image/*"
        action="/api/file/upload"
        :show-file-list="false"
        @finish="getFileList"
      >
        <NButton type="primary" ghost size="small">
          <template #icon>
            <SvgIcon icon="tabler:file-upload" />
          </template>
          批量上传
        </NButton>
      </NUpload>
    </div>
    <!-- 分类筛选（标签按钮） -->
    <div class="flex justify-center flex-wrap gap-2 mt-2 mb-3">
      <NButton
        v-for="cat in FILE_CATEGORIES"
        :key="cat.value"
        :type="selectedCategory === cat.value ? 'primary' : 'default'"
        :ghost="selectedCategory !== cat.value"
        size="small"
        round
        @click="selectedCategory = cat.value"
      >
        {{ cat.label }}
        <span v-if="cat.value !== 'all'" class="ml-1 text-xs opacity-70">
          ({{ imageList.filter(item => getFileCategory(item) === cat.value).length }})
        </span>
      </NButton>
    </div>
    <div class="flex justify-center mt-2">
      <div v-if="filteredImageList.length === 0 && !loading" class="flex">
        {{ $t('apps.uploadsFileManager.nothingText') }}
      </div>
      <NImageGroup v-else>
        <NGrid cols="2 300:2 600:4 900:6 1100:9" :x-gap="5" :y-gap="5">
          <NGridItem v-for=" item, index in filteredImageList" :key="index">
            <NCard size="small" style="border-radius: 5px;" :bordered="true">
              <template #cover>
                <div class="card transparent-grid">
                  <NImage :lazy="true" style="object-fit: contain;height: 100%;" :src="item.src" />
                </div>
              </template>
              <template #footer>
                <span class="text-xs">
                  <NEllipsis>
                    {{ item.fileName }}
                  </NEllipsis>
                </span>
                <div class="flex justify-center mt-[6px]">
                  <NTag :type="getCategoryColor(getFileCategory(item))" size="small" round>
                    {{ getCategoryLabel(getFileCategory(item)) }}
                  </NTag>
                </div>
                <div class="flex justify-center mt-[10px]">
                  <NButtonGroup>
                    <NButton size="tiny" tertiary style="cursor: pointer;" :title="$t('apps.uploadsFileManager.copyLink')" @click="copyImageUrl(item.src)">
                      <template #icon>
                        <SvgIcon icon="ion-copy" />
                      </template>
                    </NButton>
                    <NButton size="tiny" tertiary style="cursor: pointer;" :title="timeFormat(item.createTime)" @click="handleInfoClick(item)">
                      <template #icon>
                        <SvgIcon icon="mdi-information-box-outline" />
                      </template>
                    </NButton>
                    <NButton size="tiny" tertiary style="cursor: pointer;" :title="$t('apps.uploadsFileManager.setWallpaper')" @click="handleQuickSet(item)">
                      <template #icon>
                        <SvgIcon icon="lucide:wallpaper" />
                      </template>
                    </NButton>
                    <NButton size="tiny" tertiary type="error" style="cursor: pointer;" :title="$t('common.delete')" @click="handleDelete(item.id as number)">
                      <template #icon>
                        <SvgIcon icon="material-symbols-delete" />
                      </template>
                    </NButton>
                  </NButtonGroup>
                </div>
              </template>
            </NCard>
          </NGridItem>
        </NGrid>
      </NImageGroup>
    </div>

    <RoundCardModal v-model:show="infoModalState.show" style="max-width: 300px;" size="small" :title="$t('apps.uploadsFileManager.infoTitle')">
      <div>
        <div>
          <div class="mb-2">
            <span class="text-slate-500">
              {{ $t('apps.uploadsFileManager.fileName') }}
            </span>
            <div class="text-xs">
              {{ infoModalState.fileInfo?.fileName }}
            </div>
          </div>
          <div class="mb-2">
            <span class="text-slate-500">
              {{ $t('apps.uploadsFileManager.path') }}
            </span>
            <div class="text-xs">
              {{ infoModalState.fileInfo?.src }}
            </div>
          </div>
          <div class="mb-2">
            <span class="text-slate-500">
              {{ $t('apps.uploadsFileManager.uploadTime') }}
            </span>
            <div class="text-xs">
              {{ timeFormat(infoModalState.fileInfo?.createTime) }}
            </div>
          </div>
        </div>
      </div>
    </RoundCardModal>
  </div>
</template>

<style scoped>
.card {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 80px;
}

.transparent-grid {
  background-image: linear-gradient(45deg, #f0f0f0 25%, transparent 25%, transparent 75%, #f0f0f0 75%),
    linear-gradient(45deg, #f0f0f0 25%, transparent 25%, transparent 75%, #f0f0f0 75%);
  background-size: 16px 16px;
  background-position: 0 0, 8px 8px;
}
</style>
