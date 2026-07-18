<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NDataTable, NDatePicker, NInput, NPopconfirm, NTag, useMessage } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import { getLoginLogList, clearLoginLog } from '@/api/system/log'
import { t } from '@/locales'

const message = useMessage()
const tableIsLoading = ref<boolean>(false)
const keyWord = ref<string>('')
const timeRange = ref<[number, number] | null>(null)

const createColumns = (): DataTableColumns<LoginLog.Info> => {
  return [
    {
      title: t('common.username'),
      key: 'username',
      width: 120,
    },
    {
      title: t('apps.operationLog.ip'),
      key: 'ip',
      width: 140,
    },
    {
      title: t('apps.operationLog.userAgent'),
      key: 'userAgent',
      width: 200,
      ellipsis: { tooltip: true },
      render(row) {
        return (row.userAgent || '').substring(0, 50)
      },
    },
    {
      title: t('apps.loginLog.status'),
      key: 'status',
      width: 100,
      render(row) {
        if (row.status === 1)
          return h(NTag, { type: 'success' }, { default: () => t('apps.loginLog.success') })
        return h(NTag, { type: 'error' }, { default: () => t('apps.loginLog.failed') })
      },
    },
    {
      title: t('apps.loginLog.remark'),
      key: 'remark',
      width: 150,
      ellipsis: { tooltip: true },
    },
    {
      title: t('common.createTime'),
      key: 'createTime',
      width: 180,
    },
  ]
}

const logList = ref<LoginLog.Info[]>()

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
  clearLoginLog<null>({}).then(({ code }) => {
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
    }
    if (timeRange.value && timeRange.value[0] && timeRange.value[1]) {
      params.startTime = Math.floor(timeRange.value[0] / 1000)
      params.endTime = Math.floor(timeRange.value[1] / 1000)
    }
    const { data } = await getLoginLogList<Common.ListResponse<LoginLog.Info[]>>(params)
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
      :scroll-x="800"
      :max-height="550"
      :virtual-scroll="true"
    />
  </div>
</template>