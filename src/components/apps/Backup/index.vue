<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NDataTable, NDropdown, NInput, NModal, NForm, NFormItem, NSpace, NTag, NSelect, NSwitch, useDialog, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { createBackup, getBackupList, restoreBackup, deleteBackup, exportBackup, importBackup, getBackupTaskList, createBackupTask, updateBackupTask, deleteBackupTask } from '@/api/system/backup'
import { SvgIcon } from '@/components/common'
import { t } from '@/locales'

const message = useMessage()
const dialog = useDialog()

const activeTab = ref<'backup' | 'task'>('backup')

const tableIsLoading = ref<boolean>(false)
const backupList = ref<any[]>([])
const taskList = ref<any[]>([])


const showCreateModal = ref<boolean>(false)
const showTaskModal = ref<boolean>(false)
const isEditTask = ref<boolean>(false)

const backupForm = reactive({
  mode: 1,
  name: '',
})

const taskForm = reactive({
  id: 0,
  name: '',
  mode: 1,
  cronExpr: '',
  retentionDays: 7,
  status: 1,
})

const backupModeOptions = [
  { label: '数据库备份', value: 1 },
  { label: '全部数据备份', value: 2 },
]

const backupColumns = createColumns({
  restore: (row: any) => {
    dialog.warning({
      title: t('common.warning'),
      content: t('admin.setting.backup.restoreConfirm'),
      positiveText: t('common.confirm'),
      negativeText: t('common.cancel'),
      onPositiveClick: async () => {
        const result = await restoreBackup({ id: row.id })
        if (result.code === 0) {
          message.success(t('common.success'))
          loadBackupList()
        }
      },
    })
  },
  export: async (row: any) => {
    const result = await exportBackup<{ code: number; data: Blob }>({ id: row.id })
    if (result.code === 0) {
      const blob = new Blob([result.data as unknown as BlobPart], { type: 'application/octet-stream' })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = row.name || 'backup'
      a.click()
      window.URL.revokeObjectURL(url)
    }
  },
  delete: (row: any) => {
    dialog.warning({
      title: t('common.warning'),
      content: t('common.deletePrompt', { name: row.name }),
      positiveText: t('common.confirm'),
      negativeText: t('common.cancel'),
      onPositiveClick: async () => {
        const result = await deleteBackup({ id: row.id })
        if (result.code === 0) {
          message.success(t('common.deleteSuccess'))
          loadBackupList()
        }
      },
    })
  },
})

const taskColumns = createTaskColumns({
  edit: (row: any) => {
    isEditTask.value = true
    taskForm.id = row.id
    taskForm.name = row.name
    taskForm.mode = row.mode
    taskForm.cronExpr = row.cronExpr
    taskForm.retentionDays = row.retentionDays
    taskForm.status = row.status
    showTaskModal.value = true
  },
  delete: (row: any) => {
    dialog.warning({
      title: t('common.warning'),
      content: t('common.deletePrompt', { name: row.name }),
      positiveText: t('common.confirm'),
      negativeText: t('common.cancel'),
      onPositiveClick: async () => {
        const result = await deleteBackupTask({ id: row.id })
        if (result.code === 0) {
          message.success(t('common.deleteSuccess'))
          loadTaskList()
        }
      },
    })
  },
})

function createColumns({ restore, export: exp, delete: del }: { restore: (row: any) => void; export: (row: any) => void; delete: (row: any) => void }) {
  return [
    {
      title: t('common.name'),
      key: 'name',
      render(row: any) {
        return row.name || row.filePath.split('/').pop()
      },
    },
    {
      title: t('admin.setting.backup.mode'),
      key: 'mode',
      render(row: any) {
        return h(NTag, row.mode === 1 ? { type: 'info' } : { type: 'success' }, row.mode === 1 ? '数据库备份' : '全部数据备份')
      },
    },
    {
      title: t('common.size'),
      key: 'size',
      render(row: any) {
        if (row.size < 1024) return row.size + ' B'
        if (row.size < 1024 * 1024) return (row.size / 1024).toFixed(2) + ' KB'
        return (row.size / (1024 * 1024)).toFixed(2) + ' MB'
      },
    },
    {
      title: t('common.status'),
      key: 'status',
      render(row: any) {
        return h(NTag, row.status === 1 ? { type: 'success' } : row.status === 2 ? { type: 'info' } : { type: 'error' },
          row.status === 1 ? '正常' : row.status === 2 ? '已恢复' : '失败')
      },
    },
    {
      title: t('common.createTime'),
      key: 'createTime',
      render(row: any) {
        if (!row.createTime) return '-'
        const d = new Date(row.createTime)
        if (isNaN(d.getTime())) return row.createTime
        return d.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' })
      },
    },
    {
      title: t('common.action'),
      key: '',
      render(row: any) {
        const btn = h(NButton, { strong: true, tertiary: true, size: 'small' }, { default: () => h(SvgIcon, { icon: 'mingcute:more-1-fill' }) })
        return h(
          NDropdown,
          {
            trigger: 'click',
            onSelect(key: string | number) {
              switch (key) {
                case 'restore':
                  restore(row)
                  break
                case 'export':
                  exp(row)
                  break
                case 'delete':
                  del(row)
                  break
              }
            },
            options: [
              { label: t('admin.setting.backup.restore'), key: 'restore' },
              { label: t('admin.setting.backup.export'), key: 'export' },
              { label: t('common.delete'), key: 'delete' },
            ],
          },
          { default: () => btn },
        )
      },
    },
  ] as DataTableColumns<any>
}

function createTaskColumns({ edit, delete: del }: { edit: (row: any) => void; delete: (row: any) => void }) {
  return [
    {
      title: t('common.name'),
      key: 'name',
    },
    {
      title: t('admin.setting.backup.mode'),
      key: 'mode',
      render(row: any) {
        return h(NTag, row.mode === 1 ? { type: 'info' } : { type: 'success' }, row.mode === 1 ? '数据库备份' : '全部数据备份')
      },
    },
    {
      title: t('admin.setting.backup.cronExpr'),
      key: 'cronExpr',
    },
    {
      title: t('admin.setting.backup.retentionDays'),
      key: 'retentionDays',
      render(row: any) {
        return row.retentionDays + ' ' + t('common.day')
      },
    },
    {
      title: t('common.status'),
      key: 'status',
      render(row: any) {
        return h(NTag, row.status === 1 ? { type: 'success' } : { type: 'error' }, row.status === 1 ? t('common.active') : t('common.inactive'))
      },
    },
    {
      title: t('common.action'),
      key: '',
      render(row: any) {
        const btn = h(NButton, { strong: true, tertiary: true, size: 'small' }, { default: () => h(SvgIcon, { icon: 'mingcute:more-1-fill' }) })
        return h(
          NDropdown,
          {
            trigger: 'click',
            onSelect(key: string | number) {
              switch (key) {
                case 'edit':
                  edit(row)
                  break
                case 'delete':
                  del(row)
                  break
              }
            },
            options: [
              { label: t('common.edit'), key: 'edit' },
              { label: t('common.delete'), key: 'delete' },
            ],
          },
          { default: () => btn },
        )
      },
    },
  ] as DataTableColumns<any>
}

const backupPagination = reactive({
  page: 1,
  showSizePicker: true,
  pageSizes: [10, 30, 50],
  pageSize: 10,
  itemCount: 0,
  onChange: (page: number) => {
    backupPagination.page = page
    loadBackupList()
  },
  onUpdatePageSize: (pageSize: number) => {
    backupPagination.pageSize = pageSize
    backupPagination.page = 1
    loadBackupList()
  },
})

async function loadBackupList() {
  tableIsLoading.value = true
  const { data } = await getBackupList<Common.ListResponse<any>>({
    pageNum: backupPagination.page,
    pageSize: backupPagination.pageSize,
  })
  if (data) {
    backupList.value = data.list || []
    backupPagination.itemCount = (data as any).total || 0
  }
  tableIsLoading.value = false
}

async function loadTaskList() {
  tableIsLoading.value = true
  const { data } = await getBackupTaskList<any[]>()
  if (data) {
    taskList.value = data
  }
  tableIsLoading.value = false
}

async function handleCreateBackup() {
  const result = await createBackup({ mode: backupForm.mode, name: backupForm.name })
  if (result.code === 0) {
    message.success(t('common.success'))
    showCreateModal.value = false
    loadBackupList()
  }
}

async function handleSaveTask() {
  if (!taskForm.name.trim()) {
    message.error(t('common.nameRequired'))
    return
  }
  if (!taskForm.cronExpr.trim()) {
    message.error(t('admin.setting.backup.cronRequired'))
    return
  }
  let result
  if (isEditTask.value) {
    result = await updateBackupTask(taskForm)
  } else {
    result = await createBackupTask(taskForm)
  }
  if (result.code === 0) {
    message.success(t('common.success'))
    showTaskModal.value = false
    loadTaskList()
  }
}

async function handleImportBackup(e: any) {
  const file = e.target.files[0]
  if (!file) return
  const formData = new FormData()
  formData.append('file', file)
  const result = await importBackup(formData)
  if (result.code === 0) {
    message.success(t('common.importSuccess'))
    loadBackupList()
  }
}

function triggerImport() {
  const el = document.getElementById('importBackup') as HTMLInputElement
  el?.click()
}

function handleAddTask() {
  isEditTask.value = false
  taskForm.id = 0
  taskForm.name = ''
  taskForm.mode = 1
  taskForm.cronExpr = ''
  taskForm.retentionDays = 7
  taskForm.status = 1
  showTaskModal.value = true
}

onMounted(() => {
  loadBackupList()
  loadTaskList()
})
</script>

<template>
  <div class="overflow-auto pt-2">
    <div class="flex gap-2 mb-4">
      <NButton :type="activeTab === 'backup' ? 'primary' : 'default'" @click="activeTab = 'backup'">
        {{ $t('admin.setting.backup.backupList') }}
      </NButton>
      <NButton :type="activeTab === 'task' ? 'primary' : 'default'" @click="activeTab = 'task'">
        {{ $t('admin.setting.backup.taskList') }}
      </NButton>
    </div>

    <div v-if="activeTab === 'backup'">
      <div class="my-[10px] flex gap-2">
        <NButton type="primary" size="small" ghost @click="showCreateModal = true">
          {{ $t('common.create') }}
        </NButton>
        <input type="file" accept=".sql,.zip" @change="handleImportBackup" class="hidden" id="importBackup" />
        <NButton size="small" ghost @click="triggerImport">
        {{ $t('admin.setting.backup.import') }}
      </NButton>
      </div>

      <NDataTable
        :columns="backupColumns"
        :data="backupList"
        :pagination="backupPagination"
        :bordered="false"
        :loading="tableIsLoading"
        :remote="true"
      />

      <NModal v-model:show="showCreateModal" preset="card" :title="$t('admin.setting.backup.create')" style="width: 400px">
        <NForm>
          <NFormItem :label="$t('admin.setting.backup.mode')" required>
            <NSelect v-model:value="backupForm.mode" :options="backupModeOptions" />
          </NFormItem>
          <NFormItem :label="$t('common.name')">
            <NInput v-model:value="backupForm.name" :placeholder="$t('admin.setting.backup.namePlaceholder')" />
          </NFormItem>
        </NForm>
        <template #footer>
          <NSpace>
            <NButton @click="showCreateModal = false">{{ $t('common.cancel') }}</NButton>
            <NButton type="primary" @click="handleCreateBackup">{{ $t('common.confirm') }}</NButton>
          </NSpace>
        </template>
      </NModal>
    </div>

    <div v-if="activeTab === 'task'">
      <div class="my-[10px]">
        <NButton type="primary" size="small" ghost @click="handleAddTask">
          {{ $t('common.add') }}
        </NButton>
      </div>

      <NDataTable
        :columns="taskColumns"
        :data="taskList"
        :bordered="false"
        :loading="tableIsLoading"
      />

      <NModal v-model:show="showTaskModal" preset="card" :title="isEditTask ? $t('common.edit') : $t('common.add')" style="width: 450px">
        <NForm>
          <NFormItem :label="$t('common.name')" required>
            <NInput v-model:value="taskForm.name" :placeholder="$t('common.namePlaceholder')" />
          </NFormItem>
          <NFormItem :label="$t('admin.setting.backup.mode')" required>
            <NSelect v-model:value="taskForm.mode" :options="backupModeOptions" />
          </NFormItem>
          <NFormItem :label="$t('admin.setting.backup.cronExpr')" required>
            <NInput v-model:value="taskForm.cronExpr" :placeholder="$t('admin.setting.backup.cronPlaceholder')" />
          </NFormItem>
          <NFormItem :label="$t('admin.setting.backup.retentionDays')">
            <NInput v-model:value="taskForm.retentionDays" type="number" />
          </NFormItem>
          <NFormItem :label="$t('common.status')">
            <NSwitch v-model:value="taskForm.status" />
          </NFormItem>
        </NForm>
        <template #footer>
          <NSpace>
            <NButton @click="showTaskModal = false">{{ $t('common.cancel') }}</NButton>
            <NButton type="primary" @click="handleSaveTask">{{ $t('common.confirm') }}</NButton>
          </NSpace>
        </template>
      </NModal>
    </div>
  </div>
</template>