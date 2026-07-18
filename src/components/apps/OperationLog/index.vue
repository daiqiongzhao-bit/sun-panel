<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NDataTable, NDatePicker, NInput, NPopconfirm, NSelect, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import { getOperationLogList, clearOperationLog } from '@/api/system/log'
import { t } from '@/locales'

const message = useMessage()
const tableIsLoading = ref<boolean>(false)
const keyWord = ref<string>('')
const moduleValue = ref<string | null>(null)
const timeRange = ref<[number, number] | null>(null)

const createColumns = (): DataTableColumns<OperationLog.Info> => {
  return [
    {
      title: t('common.username'),
      key: 'username',
      width: 120,
    },
    {
      title: t('apps.operationLog.module'),
      key: 'module',
      width: 120,
    },
    {
      title: t('apps.operationLog.action'),
      key: 'action',
      width: 120,
    },
    {
      title: 'Method',
      key: 'method',
      width: 80,
      render(row) {
        const method = (row.method || '').toUpperCase()
        const colorMap: Record<string, string> = {
          GET: 'success',
          POST: 'info',
          PUT: 'warning',
          DELETE: 'error',
        }
        return h(NTag, { type: (colorMap[method] as any) || 'default', size: 'small' }, { default: () => method })
      },
    },
    {
      title: t('apps.operationLog.path'),
      key: 'path',
      ellipsis: { tooltip: true },
    },
    {
      title: t('apps.operationLog.ip'),
      key: 'ip',
      width: 140,
    },
    {
      title: t('apps.operationLog.statusCode'),
      key: 'responseCode',
      width: 100,
      render(row) {
        const code = row.responseCode
        if (code >= 200 && code < 300)
          return h(NTag, { type: 'success', size: 'small' }, { default: () => String(code) })
        if (code >= 400)
          return h(NTag, { type: 'error', size: 'small' }, { default: () => String(code) })
        return h(NTag, { type: 'warning', size: 'small' }, { default: () => String(code) })
      },
    },
    {
      title: t('apps.operationLog.duration'),
      key: 'duration',
      width: 100,
      render(row) {
        return `${row.duration}ms`
      },
    },
    {
      title: t('common.createTime'),
      key: 'createTime',
      width: 180,
    },
  ]
}

const logList = ref<OperationLog.Info[]>()

const columns = createColumns()

const pagination = reactive({
  page: 1,
  showSizePicker: true,
  pageSizes: [10, 30, 50, 100, 200],
  pageSize: 10,
  itemCount: 0,
  onChange: (page: number) => {
    pagination.page = page
    getList()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    getList()
  },
  prefix(item: PaginationProps) {
    return `Total ${item.itemCount}`
  },
})

function handleSearch() {
  pagination.page = 1
  getList()
}

function handleClear() {
  clearOperationLog<null>({}).then(({ code }) => {
    if (code === 0) {
      message.success(t('common.deleteSuccess'))
      getList()
    }
  })
}

async function getList() {
  tableIsLoading.value = true
  try {
    const params: any = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      keyword: keyWord.value || undefined,
      module: moduleValue.value || undefined,
    }
    if (timeRange.value && timeRange.value[0] && timeRange.value[1]) {
      params.startTime = Math.floor(timeRange.value[0] / 1000)
      params.endTime = Math.floor(timeRange.value[1] / 1000)
    }
    const { data } = await getOperationLogList<Common.ListResponse<OperationLog.Info[]>>(params)
    pagination.itemCount = data?.count || 0
    logList.value = data?.list || []
  } catch {
    logList.value = []
    pagination.itemCount = 0
  }
  tableIsLoading.value = false
}

onMounted(() => {
  getList()
})
</script>

<template>
  <div class="overflow-auto pt-2">
    <div class="my-[10px] flex items-center gap-2 flex-wrap">
      <NInput
        v-model:value="keyWord"
        clearable
        :placeholder="t('common.inputPlaceholder')"
        style="width: 180px;"
        @keyup.enter="handleSearch"
      />
      <NDatePicker
        v-model:value="timeRange"
        type="datetimerange"
        clearable
        style="width: 340px;"
        :placeholder="t('common.timeRange')"
      />
      <NButton type="primary" size="small" @click="handleSearch">
        {{ t('common.search') }}
      </NButton>
      <NPopconfirm @positive-click="handleClear">
        <template #trigger>
          <NButton type="error" size="small" ghost>
            {{ t('common.clear') }}
          </NButton>
        </template>
        {{ t('common.deleteConfirm') }}
      </NPopconfirm>
    </div>
    <NDataTable
      :columns="columns"
      :data="logList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
      :scroll-x="1000"
      :max-height="550"
      :virtual-scroll="true"
    />
  </div>
</template>