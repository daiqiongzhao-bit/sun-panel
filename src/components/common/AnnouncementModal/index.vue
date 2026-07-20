<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NButton, NModal } from 'naive-ui'
import { getVisibleNotices, markRead } from '@/api/system/noticeManage'
import { useNoticeStore } from '@/store'
import { t } from '@/locales'

const showModal = ref(false)
const currentNotice = ref<Notice.NoticeInfo | null>(null)
const noticeStore = useNoticeStore()

function markReadRemote(id: number) {
  // 服务端持久化已读（跨设备/清缓存同步）；本地 store 作为兜底
  markRead<Common.CommonResponse<unknown>>({ id }).catch(() => {})
  noticeStore.setReadByGlobal(id)
}

function handleMarkRead() {
  if (currentNotice.value?.id) {
    markReadRemote(currentNotice.value.id)
  }
  showModal.value = false
  currentNotice.value = null
  // 继续检查下一条
  checkNextUnread()
}

function handleClose() {
  // 关闭时也标记为已读，避免每次登录重复弹出
  if (currentNotice.value?.id) {
    markReadRemote(currentNotice.value.id)
  }
  showModal.value = false
  currentNotice.value = null
}

function checkNextUnread() {
  getVisibleNotices<Common.ListResponse<Notice.NoticeInfo[]>>({ displayType: [2] }).then(({ code, data }) => {
    if (code === 0 && data.list) {
      const unread = data.list.find(
        (notice: Notice.NoticeInfo) => notice.oneRead === 1 && !noticeStore.getReadByNoticeId(notice.id as number),
      )
      if (unread) {
        currentNotice.value = unread
        showModal.value = true
      }
    }
  })
}

onMounted(() => {
  checkNextUnread()
})
</script>

<template>
  <NModal
    v-model:show="showModal"
    preset="card"
    :title="currentNotice?.title || ''"
    style="max-width: 500px; border-radius: 1rem;"
    size="small"
    :mask-closable="true"
    :closable="true"
    @close="handleClose"
  >
    <div class="whitespace-pre-wrap text-base leading-relaxed" style="word-break: break-word;">
      {{ currentNotice?.content || '' }}
    </div>
    <template #footer>
      <div class="flex justify-end">
        <NButton type="primary" @click="handleMarkRead">
          {{ t('apps.noticeManage.markRead') }}
        </NButton>
      </div>
    </template>
  </NModal>
</template>