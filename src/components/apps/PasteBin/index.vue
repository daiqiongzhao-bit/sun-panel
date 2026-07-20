<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import { NButton, NCheckbox, NDataTable, NDatePicker, NInput, NInputGroup, NModal, NSelect, NTag, NUpload, useMessage } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { t } from '@/locales'
import { createPasteBin, updatePasteBin, getMyPasteBinList, getPasteBinAccessUrl } from '@/api/panel/pasteBin'

const ms = useMessage()

interface PasteBinItem {
  id: number
  type: number
  title: string
  content: string
  fileName: string
  fileSize: number
  code: string
  expireAt: string
  accessCnt: number
  status: number
  createdAt: string
  burnAfterRead: number
}

const list = ref<PasteBinItem[]>([])
const loading = ref(false)
const keyWord = ref('')
const dateRange = ref<[number, number] | null>(null)
const pagination = ref({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => {
    pagination.value.page = page
    loadData()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.value.pageSize = pageSize
    pagination.value.page = 1
    loadData()
  },
})

// Create Modal
const showModal = ref(false)
const formTitle = ref('')
const formType = ref<number | null>(1)
const formContent = ref('')
const formExpire = ref<number | null>(1)
const formPassword = ref('')
const formBurnAfterRead = ref(false)
const uploadFile = ref<UploadFileInfo[]>([])
const uploadedUrl = ref('')

// Edit Modal
const showEditModal = ref(false)
const editId = ref(0)
const editPassword = ref('')
const editExpire = ref<number | null>(1)
const editBurnAfterRead = ref(false)
const editCode = ref('')

// Result after creation
const showResult = ref(false)
const resultCode = ref('')
const resultUrl = ref('')

const typeOptions = [
  { label: t('apps.pasteBin.text'), value: 1 },
  { label: t('apps.pasteBin.file'), value: 2 },
]

const expireOptions = [
  { label: t('apps.pasteBin.1h'), value: 1 },
  { label: t('apps.pasteBin.6h'), value: 6 },
  { label: t('apps.pasteBin.1d'), value: 24 },
  { label: t('apps.pasteBin.3d'), value: 72 },
  { label: t('apps.pasteBin.7d'), value: 168 },
]

function loadData() {
  loading.value = true
  const params: Record<string, any> = {
    page: pagination.value.page,
    pageSize: pagination.value.pageSize,
  }
  if (keyWord.value)
    params.keyword = keyWord.value
  if (dateRange.value && dateRange.value.length === 2) {
    params.startTime = Math.floor(dateRange.value[0] / 1000).toString()
    params.endTime = Math.floor(dateRange.value[1] / 1000).toString()
  }
  getMyPasteBinList<Common.ListResponse<PasteBinItem[]>>(params).then(({ code, data }) => {
    if (code === 0) {
      list.value = data.list || []
      pagination.value.itemCount = data.count || 0
    }
  }).finally(() => {
    loading.value = false
  })
}

function resetSearch() {
  keyWord.value = ''
  dateRange.value = null
  pagination.value.page = 1
  loadData()
}

function openCreateModal() {
  formTitle.value = ''
  formType.value = 1
  formContent.value = ''
  formExpire.value = 1
  formPassword.value = ''
  formBurnAfterRead.value = false
  uploadFile.value = []
  showResult.value = false
  resultCode.value = ''
  resultUrl.value = ''
  showModal.value = true
}

function handleUploadChange(options: { file: UploadFileInfo; fileList: Array<UploadFileInfo> }) {
  uploadFile.value = options.fileList
  if (options.file.file) {
    formTitle.value = options.file.name
  }
}

function handleUploadFinish(options: { file: UploadFileInfo; event?: any }) {
  // Parse upload response to get URL
  try {
    const response = JSON.parse((options.event?.target as XMLHttpRequest)?.responseText || '{}')
    if (response.code === 0 && response.data) {
      uploadedUrl.value = response.data.url || ''
      formContent.value = response.data.url || ''
    }
  } catch (e) {
    // fallback
    formContent.value = options.file.name
  }
}

function handleCreate() {
  if (!formTitle.value.trim()) {
    ms.warning('请输入标题')
    return
  }
  createPasteBin<Common.Response<PasteBinItem>>({
    title: formTitle.value,
    type: formType.value,
    content: formContent.value,
    fileName: uploadFile.value.length > 0 ? uploadFile.value[0].name : '',
    fileSize: uploadFile.value.length > 0 ? (uploadFile.value[0].file as File)?.size || 0 : 0,
    expireH: formExpire.value,
    password: formPassword.value,
    burnAfterRead: formBurnAfterRead.value ? 1 : 0,
  }).then(async ({ code, data }) => {
    if (code === 0) {
      resultCode.value = data.code
      try {
        const urlRes = await getPasteBinAccessUrl<Common.Response<any>>({ code: data.code })
        if (urlRes.code === 0 && urlRes.data) {
          resultUrl.value = urlRes.data.url || ''
        }
      }
      catch { resultUrl.value = '' }
      showResult.value = true
      ms.success(t('apps.pasteBin.created'))
      loadData()
    }
  })
}

function openEditModal(row: PasteBinItem) {
  editId.value = row.id
  editPassword.value = ''
  editExpire.value = 1
  editBurnAfterRead.value = row.burnAfterRead === 1
  editCode.value = row.code
  showEditModal.value = true
}

function handleUpdate() {
  updatePasteBin<Common.Response<any>>({
    id: editId.value,
    password: editPassword.value,
    expireH: editExpire.value || 0,
    burnAfterRead: editBurnAfterRead.value ? 1 : 0,
  }).then(async ({ code }) => {
    if (code === 0) {
      ms.success(t('common.saveSuccess'))
      showEditModal.value = false
      loadData()
    }
  })
}

async function copyAccessUrl(row: PasteBinItem) {
  try {
    const urlRes = await getPasteBinAccessUrl<Common.Response<any>>({ code: row.code })
    if (urlRes.code === 0 && urlRes.data) {
      navigator.clipboard.writeText(urlRes.data.url || '').then(() => {
        ms.success(t('apps.pasteBin.urlCopied'))
      })
    }
  }
  catch { ms.error(t('common.failed')) }
}

function copyToClipboard(text: string, messageKey: string) {
  navigator.clipboard.writeText(text).then(() => {
    ms.success(t(`apps.pasteBin.${messageKey}` as any))
  }).catch(() => {
    ms.error(t('common.networkError'))
  })
}

const columns = [
  {
    title: t('apps.pasteBin.title'),
    key: 'title',
    ellipsis: { tooltip: true },
  },
  {
    title: t('apps.pasteBin.type'),
    key: 'type',
    width: 90,
    render(row: PasteBinItem) {
      return h(NTag, { type: row.type === 1 ? 'info' : 'warning', size: 'small' },
        { default: () => (row.type === 1 ? t('apps.pasteBin.text') : t('apps.pasteBin.file')) },
      )
    },
  },
  {
    title: t('apps.pasteBin.accessCode'),
    key: 'code',
    width: 140,
    render(row: PasteBinItem) {
      return h(NTag, { copyable: true, size: 'small' }, { default: () => row.code })
    },
  },
  {
    title: t('apps.pasteBin.accessCnt'),
    key: 'accessCnt',
    width: 90,
  },
  {
    title: t('apps.pasteBin.expireAt'),
    key: 'expireAt',
    width: 170,
    ellipsis: { tooltip: true },
  },
  {
    title: t('common.action'),
    key: 'action',
    width: 210,
    render(row: PasteBinItem) {
      return h('div', { class: 'flex items-center gap-1' }, [
        h(NButton, { size: 'tiny', tertiary: true, onClick: () => copyAccessUrl(row) },
          { default: () => t('apps.pasteBin.copyUrl') }),
        h(NButton, { size: 'tiny', tertiary: true, onClick: () => openEditModal(row) },
          { default: () => t('common.edit') }),
      ])
    },
  },
]

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="paste-bin-page">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-bold">{{ t('apps.pasteBin.appName') }}</h3>
      <NButton type="primary" @click="openCreateModal">
        {{ t('apps.pasteBin.create') }}
      </NButton>
    </div>

    <div class="flex flex-wrap items-center gap-2 mb-3">
      <NInput
        v-model:value="keyWord"
        clearable
        placeholder="搜索标题 / 内容 / 提取码"
        style="width: 240px;"
        @keyup.enter="loadData"
      />
      <NDatePicker
        v-model:value="dateRange"
        type="daterange"
        clearable
        placeholder="按创建时间筛选"
      />
      <NButton type="primary" size="small" @click="loadData">
        {{ t('common.search') }}
      </NButton>
      <NButton size="small" @click="resetSearch">
        {{ t('common.reset') }}
      </NButton>
    </div>

    <NDataTable
      :columns="columns"
      :data="list"
      :loading="loading"
      :pagination="pagination"
      :bordered="false"
      striped
      :scroll-x="700"
    />

    <!-- Create Modal -->
    <NModal
      v-model:show="showModal"
      preset="card"
      style="max-width: 560px;"
      :title="t('apps.pasteBin.create')"
      :bordered="true"
      size="small"
    >
      <div v-if="showResult" class="create-result">
        <div class="mb-3">
          <div class="mb-2 font-semibold">{{ $t('apps.pasteBin.accessCode') }}</div>
          <NInputGroup>
            <NInput :value="resultCode" readonly />
            <NButton type="primary" @click="copyToClipboard(resultCode, 'codeCopied')">
              {{ $t('apps.pasteBin.copyCode') }}
            </NButton>
          </NInputGroup>
        </div>
        <div v-if="resultUrl" class="mb-3">
          <div class="mb-2 font-semibold">访问链接</div>
          <NInputGroup>
            <NInput :value="resultUrl" readonly />
            <NButton type="primary" @click="copyToClipboard(resultUrl, 'urlCopied')">
              {{ $t('apps.pasteBin.copyUrl') }}
            </NButton>
          </NInputGroup>
          <div class="text-xs text-gray-400 mt-1">将此链接分享给他人即可访问中转内容</div>
        </div>
      </div>

      <div v-else class="create-form">
        <div class="mb-3">
          <div class="mb-2">{{ $t('apps.pasteBin.title') }}</div>
          <NInput v-model:value="formTitle" :placeholder="$t('apps.pasteBin.title')" />
        </div>
        <div class="mb-3">
          <div class="mb-2">{{ $t('apps.pasteBin.type') }}</div>
          <NSelect v-model:value="formType" :options="typeOptions" />
        </div>
        <div v-if="formType === 2" class="mb-3">
          <div class="mb-2">上传文件</div>
          <NUpload
            :max="1"
            action="/api/panel/pasteBin/uploadFile"
            name="file"
            :show-file-list="false"
            @change="handleUploadChange"
            @finish="handleUploadFinish"
          >
            <NButton>选择并上传文件</NButton>
          </NUpload>
          <p v-if="uploadFile.length > 0" class="text-xs text-gray-400 mt-1">已选择：{{ uploadFile[0].name }}</p>
          <p v-if="uploadedUrl" class="text-xs text-green-500 mt-1">文件已上传</p>
        </div>
        <div v-if="formType === 1" class="mb-3">
          <div class="mb-2">{{ $t('apps.pasteBin.content') }}</div>
          <NInput v-model:value="formContent" type="textarea" :rows="6" :placeholder="$t('apps.pasteBin.content')" />
        </div>
        <div class="mb-3">
          <div class="mb-2">{{ $t('apps.pasteBin.expireH') }}</div>
          <NSelect v-model:value="formExpire" :options="expireOptions" />
        </div>
        <div class="mb-3">
          <div class="mb-2">{{ $t('apps.pasteBin.password') }}</div>
          <NInput v-model:value="formPassword" type="text" :placeholder="$t('apps.pasteBin.passwordPlaceholder')" />
        </div>
        <div class="mb-3">
          <NCheckbox v-model:checked="formBurnAfterRead">
            <span class="text-sm">{{ $t('apps.pasteBin.burnAfterRead') }}</span>
            <span class="text-xs text-gray-400 ml-1">({{ $t('apps.pasteBin.burnAfterReadDesc') }})</span>
          </NCheckbox>
        </div>
        <div class="flex justify-end">
          <NButton type="primary" @click="handleCreate">{{ $t('apps.pasteBin.create') }}</NButton>
        </div>
      </div>
    </NModal>

    <!-- Edit Modal -->
    <NModal v-model:show="showEditModal" preset="card" style="max-width: 420px;" :title="'修改中转 [' + editCode + ']'" :bordered="true" size="small">
      <div class="mb-3">
        <div class="mb-2">{{ $t('apps.pasteBin.password') }} <span class="text-xs text-gray-400">(留空=清除密码)</span></div>
        <NInput v-model:value="editPassword" type="text" :placeholder="$t('apps.pasteBin.passwordPlaceholder')" />
      </div>
      <div class="mb-3">
        <div class="mb-2">{{ $t('apps.pasteBin.expireH') }}</div>
        <NSelect v-model:value="editExpire" :options="expireOptions" />
      </div>
      <div class="mb-3">
        <NCheckbox v-model:checked="editBurnAfterRead">
          <span class="text-sm">{{ $t('apps.pasteBin.burnAfterRead') }}</span>
          <span class="text-xs text-gray-400 ml-1">({{ $t('apps.pasteBin.burnAfterReadDesc') }})</span>
        </NCheckbox>
      </div>
      <div class="flex justify-end">
        <NButton type="primary" @click="handleUpdate">{{ $t('common.save') }}</NButton>
      </div>
    </NModal>
  </div>
</template>

<style scoped>
.paste-bin-page { padding: 10px; }
.create-result { padding: 10px 0; }
.create-form { padding: 10px 0; }
</style>
