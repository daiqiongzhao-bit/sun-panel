<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue'
import { NButton, NDataTable, NDropdown, NInput, NTag, useDialog, useMessage } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'
import EditNotice from './EditNotice/index.vue'
import { deleteNotice, getNoticeList } from '@/api/system/noticeManage'
import { t } from '@/locales'

const message = useMessage()
const dialog = useDialog()
const tableIsLoading = ref<boolean>(false)
const editDialogShow = ref<boolean>(false)
const editData = ref<Notice.NoticeInfo | null>(null)
const keyWord = ref<string>('')

const createColumns = ({
  edit,
  del,
}: {
  edit: (row: Notice.NoticeInfo) => void
  del: (row: Notice.NoticeInfo) => void
}): DataTableColumns<Notice.NoticeInfo> => {
  return [
    {
      title: t('apps.noticeManage.title'),
      key: 'title',
    },
    {
      title: t('apps.noticeManage.noticeType'),
      key: 'noticeType',
      width: 120,
      render(row) {
        if (row.noticeType === 1)
          return h(NTag, { type: 'info' }, { default: () => t('apps.noticeManage.announcement') })
        return h(NTag, { type: 'warning' }, { default: () => t('apps.noticeManage.message') })
      },
    },
    {
      title: t('apps.noticeManage.displayType'),
      key: 'displayType',
      width: 120,
      render(row) {
        if (row.displayType === 1)
          return t('apps.noticeManage.displayLogin')
        return t('apps.noticeManage.displayHome')
      },
    },
    {
      title: t('apps.noticeManage.status'),
      key: 'status',
      width: 100,
      render(row) {
        if (row.status === 1)
          return h(NTag, { type: 'success' }, { default: () => t('apps.noticeManage.enabled') })
        return h(NTag, { type: 'error' }, { default: () => t('apps.noticeManage.disabled') })
      },
    },
    {
      title: t('common.action'),
      key: 'actions',
      width: 100,
      render(row) {
        return h(NDropdown, {
          trigger: 'click',
          onSelect(key: string | number) {
            switch (key) {
              case 'edit':
                edit(row)
                break
              case 'delete':
                del(row)
                break
              default:
                break
            }
          },
          options: [
            { label: t('common.edit'), key: 'edit' },
            { label: t('common.delete'), key: 'delete' },
          ],
        }, {
          default: () => h(NButton, { strong: true, tertiary: true, size: 'small' }, { default: () => '...' }),
        })
      },
    },
  ]
}

const noticeList = ref<Notice.NoticeInfo[]>()

const columns = createColumns({
  edit(row: Notice.NoticeInfo) {
    editData.value = row
    editDialogShow.value = true
  },
  del(row: Notice.NoticeInfo) {
    dialog.warning({
      title: t('common.warning'),
      content: t('apps.noticeManage.deleteConfirm'),
      positiveText: t('common.confirm'),
      negativeText: t('common.cancel'),
      onPositiveClick: () => {
        deleteNotice<null>({ id: row.id as number }).then(({ code, msg }) => {
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
    return `${t('common.title')} ${item.itemCount}`
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
  const { data } = await getNoticeList<Common.ListResponse<Notice.NoticeInfo[]>>({
    page: pagination.page,
    pageSize: pagination.pageSize,
    keyword: keyWord.value || undefined,
  })
  pagination.itemCount = data.count
  if (data.list)
    noticeList.value = data.list
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
        :placeholder="t('common.inputPlaceholderByText', { text: t('apps.noticeManage.title') })"
        style="width: 250px;"
        @keyup.enter="handleSearch"
      />
      <NButton type="primary" size="small" @click="handleSearch">
        {{ t('common.search') }}
      </NButton>
      <NButton type="primary" size="small" ghost @click="handleAdd">
        {{ t('common.add') }}
      </NButton>
    </div>
    <NDataTable
      :columns="columns"
      :data="noticeList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
    />
    <EditNotice v-model:visible="editDialogShow" :edit-data="editData" @done="handleDone" />
  </div>
</template>