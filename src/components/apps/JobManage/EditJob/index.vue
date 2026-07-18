<script lang="ts" setup>
import { computed, reactive, ref, watch } from 'vue'
import { NButton, NForm, NFormItem, NInput, NSelect, useMessage } from 'naive-ui'
import RoundCardModal from '@/components/common/RoundCardModal/index.vue'
import { createJob, updateJob, previewCron } from '@/api/system/job'
import { t } from '@/locales'

const props = defineProps<{
  visible: boolean
  editData: JobTask.Info | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'done'): void
}>()

const message = useMessage()
const loading = ref(false)
const previewText = ref<string>('')
const previewLoading = ref(false)
const schedulePreset = ref<string>('')

const showModal = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  jobType: 1,
  cronExpr: '',
  content: '',
  status: 1,
})

const statusOptions = [
  { label: t('apps.jobManage.running'), value: 1 },
  { label: t('apps.jobManage.paused'), value: 2 },
]

const schedulePresets = [
  { label: '每分钟', value: '* * * * * ?' },
  { label: '每5分钟', value: '0 */5 * * * ?' },
  { label: '每10分钟', value: '0 */10 * * * ?' },
  { label: '每30分钟', value: '0 */30 * * * ?' },
  { label: '每小时', value: '0 0 * * * ?' },
  { label: '每天 00:00', value: '0 0 0 * * ?' },
  { label: '每天 08:00', value: '0 0 8 * * ?' },
  { label: '每天 12:00', value: '0 0 12 * * ?' },
  { label: '每周一 08:00', value: '0 0 8 * * 1' },
  { label: '自定义(Cron)', value: '' },
]

function applySchedulePreset(value: string) {
  if (value) {
    formData.cronExpr = value
  }
}

const isEdit = computed(() => !!formData.id)

watch(() => props.visible, (val) => {
  if (val) {
    previewText.value = ''
    schedulePreset.value = ''
    if (props.editData) {
      formData.id = props.editData.id
      formData.name = props.editData.name || ''
      formData.jobType = props.editData.jobType || 1
      formData.cronExpr = props.editData.cronExpr || ''
      formData.content = props.editData.content || ''
      formData.status = props.editData.status || 1
    }
    else {
      formData.id = undefined
      formData.name = ''
      formData.jobType = 1
      formData.cronExpr = ''
      formData.content = ''
      formData.status = 1
    }
  }
})

function handlePreview() {
  if (!formData.cronExpr.trim()) {
    message.warning(t('form.required'))
    return
  }
  previewLoading.value = true
  previewCron<string[]>({ cronExpr: formData.cronExpr }).then(({ code, data }) => {
    previewLoading.value = false
    if (code === 0 && data) {
      if (Array.isArray(data)) {
        previewText.value = data.slice(0, 5).join('\n')
      }
      else {
        previewText.value = String(data)
      }
    }
    else {
      previewText.value = '-'
    }
  }).catch(() => {
    previewLoading.value = false
    previewText.value = '-'
  })
}

function handleSubmit() {
  if (!formData.name.trim()) {
    message.warning(t('form.required'))
    return
  }
  if (!formData.cronExpr.trim()) {
    message.warning(t('form.required'))
    return
  }
  loading.value = true
  const reqData = {
    id: formData.id,
    name: formData.name,
    jobType: formData.jobType,
    cronExpr: formData.cronExpr,
    content: formData.content,
    status: formData.status,
  }
  const api = formData.id ? updateJob : createJob
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
    :title="editData ? t('apps.jobManage.editJob') : t('apps.jobManage.addJob')"
    size="small"
    style="max-width: 600px;"
  >
    <NForm label-placement="left" label-width="120">
      <NFormItem :label="t('apps.jobManage.name')">
        <NInput v-model:value="formData.name" :placeholder="t('common.inputPlaceholderByText', { text: t('apps.jobManage.name') })" />
      </NFormItem>
      <NFormItem label="调度方式">
        <NSelect
          v-model:value="schedulePreset"
          :options="schedulePresets"
          placeholder="选择调度方式"
          @update:value="applySchedulePreset"
        />
      </NFormItem>
      <NFormItem :label="t('apps.jobManage.cronExpr')">
        <div class="flex items-center gap-2 w-full">
          <NInput
            v-model:value="formData.cronExpr"
            class="flex-1"
            :placeholder="t('common.inputPlaceholderByText', { text: t('apps.jobManage.cronExpr') })"
          />
          <NButton size="small" :loading="previewLoading" @click="handlePreview">
            {{ t('apps.jobManage.preview') }}
          </NButton>
        </div>
      </NFormItem>
      <NFormItem v-if="previewText" :label="t('apps.jobManage.nextRunPreview')">
        <div class="text-sm text-gray-600 whitespace-pre-line">
          {{ previewText }}
        </div>
      </NFormItem>
      <NFormItem :label="t('apps.jobManage.content')">
        <NInput
          v-model:value="formData.content"
          type="textarea"
          :rows="4"
          :placeholder="t('common.inputPlaceholderByText', { text: t('apps.jobManage.content') })"
        />
      </NFormItem>
      <NFormItem v-if="isEdit" :label="t('apps.jobManage.status')">
        <NSelect v-model:value="formData.status" :options="statusOptions" />
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
