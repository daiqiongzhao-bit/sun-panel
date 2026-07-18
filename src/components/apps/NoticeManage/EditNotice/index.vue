<script lang="ts" setup>
import { computed, reactive, ref, watch } from 'vue'
import { NButton, NForm, NFormItem, NInput, NSelect, NSwitch, useMessage } from 'naive-ui'
import RoundCardModal from '@/components/common/RoundCardModal/index.vue'
import { createNotice, updateNotice } from '@/api/system/noticeManage'
import { t } from '@/locales'

const props = defineProps<{
  visible: boolean
  editData: Notice.NoticeInfo | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}>()

const message = useMessage()
const loading = ref(false)

const showModal = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

const formData = reactive({
  id: undefined as number | undefined,
  title: '',
  content: '',
  noticeType: 1,
  displayType: 1,
  oneRead: false,
  url: '',
  status: 1,
  targetUserIds: '',
})

const noticeTypeOptions = [
  { label: t('apps.noticeManage.announcement'), value: 1 },
  { label: t('apps.noticeManage.message'), value: 2 },
]

const displayTypeOptions = [
  { label: t('apps.noticeManage.displayLogin'), value: 1 },
  { label: t('apps.noticeManage.displayHome'), value: 2 },
]

const statusOptions = [
  { label: t('apps.noticeManage.enabled'), value: 1 },
  { label: t('apps.noticeManage.disabled'), value: 2 },
]

const isMessage = computed(() => formData.noticeType === 2)

watch(() => props.visible, (val) => {
  if (val) {
    if (props.editData) {
      formData.id = props.editData.id
      formData.title = props.editData.title || ''
      formData.content = props.editData.content || ''
      formData.noticeType = props.editData.noticeType || 1
      formData.displayType = props.editData.displayType || 1
      formData.oneRead = props.editData.oneRead === 1
      formData.url = props.editData.url || ''
      formData.status = props.editData.status || 1
      formData.targetUserIds = props.editData.targetUserIds || ''
    }
    else {
      formData.id = undefined
      formData.title = ''
      formData.content = ''
      formData.noticeType = 1
      formData.displayType = 1
      formData.oneRead = false
      formData.url = ''
      formData.status = 1
      formData.targetUserIds = ''
    }
  }
})

function handleSubmit() {
  if (!formData.title.trim()) {
    message.warning(t('form.required'))
    return
  }
  loading.value = true
  const reqData = {
    id: formData.id,
    title: formData.title,
    content: formData.content,
    noticeType: formData.noticeType,
    displayType: formData.displayType,
    oneRead: formData.oneRead ? 1 : 0,
    url: formData.url,
    status: formData.status,
    targetUserIds: formData.targetUserIds,
  }
  const api = formData.id ? updateNotice : createNotice
  api<null>(reqData).then(({ code, msg }) => {
    loading.value = false
    if (code === 0) {
      message.success(t('common.success'))
      showModal.value = false
      emit('done')
    }
    else {
      message.error(msg || t('common.failed'))
    }
  }).catch(() => {
    loading.value = false
  })
}
</script>

<template>
  <RoundCardModal
    v-model:show="showModal"
    :title="editData ? t('common.edit') : t('common.add')"
    size="small"
    style="max-width: 600px;"
  >
    <NForm label-placement="left" label-width="120">
      <NFormItem :label="t('apps.noticeManage.title')">
        <NInput v-model:value="formData.title" :placeholder="t('common.inputPlaceholderByText', { text: t('apps.noticeManage.title') })" />
      </NFormItem>
      <NFormItem :label="t('apps.noticeManage.content')">
        <NInput v-model:value="formData.content" type="textarea" :rows="6" :placeholder="t('common.inputPlaceholderByText', { text: t('apps.noticeManage.content') })" />
      </NFormItem>
      <NFormItem :label="t('apps.noticeManage.noticeType')">
        <NSelect v-model:value="formData.noticeType" :options="noticeTypeOptions" />
      </NFormItem>
      <NFormItem :label="t('apps.noticeManage.displayType')">
        <NSelect v-model:value="formData.displayType" :options="displayTypeOptions" />
      </NFormItem>
      <NFormItem :label="t('apps.noticeManage.oneRead')">
        <NSwitch v-model:value="formData.oneRead" />
      </NFormItem>
      <NFormItem :label="t('apps.noticeManage.url')">
        <NInput v-model:value="formData.url" :placeholder="t('common.inputPlaceholderByText', { text: t('apps.noticeManage.url') })" />
      </NFormItem>
      <NFormItem :label="t('apps.noticeManage.status')">
        <NSelect v-model:value="formData.status" :options="statusOptions" />
      </NFormItem>
      <NFormItem v-if="isMessage" :label="t('apps.noticeManage.targetUserIds')">
        <NInput v-model:value="formData.targetUserIds" :placeholder="t('apps.noticeManage.targetUserIds')" />
      </NFormItem>
    </NForm>
    <template #action>
      <div class="flex justify-end gap-2">
        <NButton @click="showModal = false">{{ t('common.cancel') }}</NButton>
        <NButton type="primary" :loading="loading" @click="handleSubmit">{{ t('common.confirm') }}</NButton>
      </div>
    </template>
  </RoundCardModal>
</template>