<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NDataTable, NInput, NModal, NPopconfirm, NTag, useDialog, useMessage } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import EditJob from './EditJob/index.vue'
import { deleteJob, getJobList, getJobLogList, pauseJob, runJobNow, startJob } from '@/api/system/job'
import { t } from '@/locales'

const message = useMessage()
const dialog = useDialog()
const tableIsLoading = ref<boolean>(false)
const editDialogShow = ref<boolean>(false)
const editData = ref<JobTask.Info | null>(null)
const keyWord = ref<string>('')

// 日志弹窗相关
const logModalShow = ref<boolean>(false)
const logTableIsLoading = ref<boolean>(false)
const currentLogJobId = ref<number>(0)
const logList = ref<JobLog.Info[]>()

const logPagination = reactive({
  page: 1,
  showSizePicker: true,
  pageSizes: [10, 30, 50, 100],
  pageSize: 10,
  itemCount: 0,
  onChange: (page: number) => {
    logPagination.page = page
    getLogList()
  },
  onUpdatePageSize: (pageSize: number) => {
    logPagination.pageSize = pageSize
    logPagination.page = 1
    getLogList()
  },
  prefix(item: PaginationProps) {
    return `Total ${item.itemCount}`
  },
})

const logColumns = createLogColumns()

function createLogColumns(): DataTableColumns<JobLog.Info> {
  return [
    {
      title: t('apps.jobManage.startTime'),
      key: 'startTime',
      width: 180,
    },
    {
      title: t('apps.jobManage.duration'),
      key: 'duration',
      width: 100,
      render(row) {
        return `${row.duration}ms`
      },
    },
    {
      title: t('apps.jobManage.status'),
      key: 'status',
      width: 80,
      render(row) {
        if (row.status === 1)
          return h(NTag, { type: 'success', size: 'small' }, { default: () => t('apps.jobManage.success') })
        return h(NTag, { type: 'error', size: 'small' }, { default: () => t('apps.jobManage.failed') })
      },
    },
    {
      title: t('apps.jobManage.errorMsg'),
      key: 'errorMsg',
      ellipsis: { tooltip: true },
      render(row) {
        if (row.status === 2)
          return row.errorMsg || '-'
        return '-'
      },
    },
  ]
}

async function getLogList() {
  logTableIsLoading.value = true
  const { data } = await getJobLogList<Common.ListResponse<JobLog.Info[]>>({
    jobId: currentLogJobId.value,
    page: logPagination.page,
    pageSize: logPagination.pageSize,
  })
  logPagination.itemCount = data.count
  if (data.list)
    logList.value = data.list
  logTableIsLoading.value = false
}

function handleViewLogs(row: JobTask.Info) {
  currentLogJobId.value = row.id as number
  logPagination.page = 1
  logModalShow.value = true
  getLogList()
}

const createColumns = ({
  pause,
  start,
  runNow,
  viewLogs,
  edit,
  del,
}: {
  pause: (row: JobTask.Info) => void
  start: (row: JobTask.Info) => void
  runNow: (row: JobTask.Info) => void
  viewLogs: (row: JobTask.Info) => void
  edit: (row: JobTask.Info) => void
  del: (row: JobTask.Info) => void
}): DataTableColumns<JobTask.Info> => {
  return [
    {
      title: t('apps.jobManage.name'),
      key: 'name',
      width: 160,
      ellipsis: { tooltip: true },
    },
    {
      title: t('apps.jobManage.cronExpr'),
      key: 'cronExpr',
      width: 150,
      ellipsis: { tooltip: true },
    },
    {
      title: t('apps.jobManage.nextRunAt'),
      key: 'nextRunAt',
      width: 180,
      render(row) {
        return row.nextRunAt || '-'
      },
    },
    {
      title: t('apps.jobManage.status'),
      key: 'status',
      width: 100,
      render(row) {
        if (row.status === 1)
          return h(NTag, { type: 'success', size: 'small' }, { default: () => t('apps.jobManage.running') })
        return h(NTag, { type: 'warning', size: 'small' }, { default: () => t('apps.jobManage.paused') })
      },
    },
    {
      title: t('common.action'),
      key: 'actions',
      width: 260,
      render(row) {
        const buttons: any[] = []

        if (row.status === 1) {
          // 运行中：暂停、立即执行、查看日志、编辑
          buttons.push(
            h(NButton, { size: 'small', onClick: () => pause(row) }, { default: () => t('apps.jobManage.paused') }),
          )
          buttons.push(
            h(NButton, { type: 'primary', size: 'small', onClick: () => runNow(row) }, { default: () => t('apps.jobManage.runNow') }),
          )
        }
        else {
          // 暂停中：启动
          buttons.push(
            h(NButton, { type: 'success', size: 'small', onClick: () => start(row) }, { default: () => t('apps.jobManage.running') }),
          )
        }

        // 查看日志
        buttons.push(
          h(NButton, { size: 'small', onClick: () => viewLogs(row) }, { default: () => t('apps.jobManage.viewLogs') }),
        )

        // 编辑
        buttons.push(
          h(NButton, { size: 'small', onClick: () => edit(row) }, { default: () => t('common.edit') }),
        )

        // 删除
        buttons.push(
          h(NPopconfirm, {
            onPositiveClick: () => del(row),
          }, {
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => t('common.delete') }),
            default: () => t('apps.jobManage.deleteConfirm'),
          }),
        )

        return h('div', { class: 'flex items-center gap-1' }, buttons)
      },
    },
  ]
}

const jobList = ref<JobTask.Info[]>()

const columns = createColumns({
  pause(row: JobTask.Info) {
    dialog.warning({
      title: t('common.warning'),
      content: t('apps.jobManage.pauseConfirm'),
      positiveText: t('common.confirm'),
      negativeText: t('common.cancel'),
      onPositiveClick: () => {
        pauseJob<null>({ id: row.id as number }).then(({ code, msg }) => {
          if (code === 0) {
            message.success(t('common.success'))
            getList()
          }
          else {
            message.error(msg || t('common.failed'))
          }
        })
      },
    })
  },
  start(row: JobTask.Info) {
    dialog.warning({
      title: t('common.warning'),
      content: t('apps.jobManage.startConfirm'),
      positiveText: t('common.confirm'),
      negativeText: t('common.cancel'),
      onPositiveClick: () => {
        startJob<null>({ id: row.id as number }).then(({ code, msg }) => {
          if (code === 0) {
            message.success(t('common.success'))
            getList()
          }
          else {
            message.error(msg || t('common.failed'))
          }
        })
      },
    })
  },
  runNow(row: JobTask.Info) {
    dialog.info({
      title: t('common.warning'),
      content: t('apps.jobManage.runNowConfirm'),
      positiveText: t('common.confirm'),
      negativeText: t('common.cancel'),
      onPositiveClick: () => {
        runJobNow<null>({ id: row.id as number }).then(({ code, msg }) => {
          if (code === 0) {
            message.success(t('common.success'))
          }
          else {
            message.error(msg || t('common.failed'))
          }
        })
      },
    })
  },
  viewLogs(row: JobTask.Info) {
    handleViewLogs(row)
  },
  edit(row: JobTask.Info) {
    editData.value = row
    editDialogShow.value = true
  },
  del(row: JobTask.Info) {
    deleteJob<null>({ id: row.id as number }).then(({ code, msg }) => {
      if (code === 0) {
        message.success(t('common.deleteSuccess'))
        getList()
      }
      else {
        message.error(`${t('common.deleteFail')}:${msg}`)
      }
    })
  },
})

const pagination = reactive({
  page: 1,
  showSizePicker: true,
  pageSizes: [10, 30, 50, 100],
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

function handleAdd() {
  editData.value = null
  editDialogShow.value = true
}

function handleDone() {
  editDialogShow.value = false
  getList()
}

function handleSearch() {
  pagination.page = 1
  getList()
}

async function getList() {
  tableIsLoading.value = true
  const { data } = await getJobList<Common.ListResponse<JobTask.Info[]>>({
    page: pagination.page,
    pageSize: pagination.pageSize,
    keyword: keyWord.value || undefined,
  })
  pagination.itemCount = data.count
  if (data.list)
    jobList.value = data.list
  tableIsLoading.value = false
}

onMounted(() => {
  getList()
})
</script>

<template>
  <div class="overflow-auto pt-2">
    <div class="my-[10px] flex items-center gap-2">
      <NInput
        v-model:value="keyWord"
        clearable
        :placeholder="t('common.inputPlaceholderByText', { text: t('apps.jobManage.name') })"
        style="width: 250px;"
        @keyup.enter="handleSearch"
      />
      <NButton type="primary" size="small" @click="handleSearch">
        {{ t('common.search') }}
      </NButton>
      <NButton type="primary" size="small" ghost @click="handleAdd">
        {{ t('apps.jobManage.addJob') }}
      </NButton>
    </div>
    <NDataTable
      :columns="columns"
      :data="jobList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
    />
    <EditJob v-model:visible="editDialogShow" :edit-data="editData" @done="handleDone" />

    <!-- 日志弹窗 -->
    <NModal
      v-model:show="logModalShow"
      preset="card"
      :title="t('apps.jobManage.viewLogs')"
      style="max-width: 800px; width: 90vw;"
      :bordered="true"
    >
      <NDataTable
        :columns="logColumns"
        :data="logList"
        :pagination="logPagination"
        :bordered="false"
        :loading="logTableIsLoading"
        :remote="true"
      />
      <template v-if="!logList || logList.length === 0" #empty>
        {{ t('apps.jobManage.noLogs') }}
      </template>
    </NModal>
  </div>
</template>
