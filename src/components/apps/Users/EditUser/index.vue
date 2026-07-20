<script setup lang="ts">
import { computed, defineEmits, defineProps, onMounted, ref, watch } from 'vue'
import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NSelect, useMessage } from 'naive-ui'
import { edit as userManageEdit } from '@/api/panel/users'
import { twofaStatus, twofaEnable, twofaConfirm, twofaDisable } from '@/api'
import { getDepartmentList } from '@/api/system/department'
import { getRoleList } from '@/api/system/role'
import { RoundCardModal } from '@/components/common'
import { t } from '@/locales'

interface Props {
  visible: boolean
  userId?: number
  userInfo?: User.Info
}

const props = defineProps<Props>()
const emit = defineEmits<Emit>()
const message = useMessage()

interface Emit {
  (e: 'update:visible', visible: boolean): void
  (e: 'done', id: number): void// 创建完成
}

const formInitValue = {
  name: '',
  username: '',
  role: 2,
  status: 3,
  departmentId: 0,
}

const model = ref<User.Info>(formInitValue)
const formRef = ref<FormInst | null>(null)

const roleOtions = ref<SelectOption[]>([])

const departmentOptions = ref<SelectOption[]>([])

async function loadRoles() {
  try {
    const { data } = await getRoleList<Common.ListResponse<any>>({ pageNum: 1, pageSize: 1000 })
    if (data && data.list && data.list.length > 0) {
      roleOtions.value = data.list.map((role: any) => ({
        label: role.name,
        value: role.id,
      }))
    }
    else {
      // 无角色数据时提供基础选项
      roleOtions.value = [
        { label: t('common.role.admin'), value: 1 },
        { label: t('common.role.regularUser'), value: 2 },
        { label: t('common.role.deptAdmin'), value: 3 },
      ]
    }
  }
  catch (e) {
    // 角色列表加载失败时提供基础选项
    roleOtions.value = [
      { label: t('common.role.admin'), value: 1 },
      { label: t('common.role.regularUser'), value: 2 },
      { label: t('common.role.deptAdmin'), value: 3 },
    ]
  }
}

async function loadDepartments() {
  try {
    const { data } = await getDepartmentList<Common.Response<Department.Info[]>>()
    if (data) {
      departmentOptions.value = data.map((dept: Department.Info) => ({
        label: dept.name,
        value: dept.id,
      }))
      departmentOptions.value.unshift({ label: t('common.noDepartment'), value: 0 })
    }
  } catch (e) {
    // 部门列表加载失败不阻塞表单
  }
}

const rules: FormRules = {
  username: [
    {
      required: true,
      trigger: 'blur',
      message: t('adminSettingUsers.formRules.usernameRequired'),
      min: 5,
    },
  ],
  role: {
    required: true,
    trigger: 'blur',
    type: 'number',
    message: t('adminSettingUsers.formRules.roleRequired'),
  },
  // status: {
  //   required: true,
  //   trigger: 'blur',
  //   type: 'number',
  //   message: '请选择账号状态',
  // },
  password: {
    trigger: 'blur',
    min: 6,
    max: 20,
    message: t('adminSettingUsers.formRules.passwordLimit'),
  },
}

// 更新值父组件传来的值
const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})

watch(show, (newValue, oldValue) => {
  if (props.userInfo?.id)
    model.value = { ...formInitValue, ...props.userInfo } as User.Info
  else
    model.value = { ...formInitValue }
  loadTwoFAStatus()
})

const add = async () => {
  const res = await userManageEdit<User.Info>(model.value)
  if (res.code === 0)
    emit('done', res.data.id as number)

  else if (res.code !== -1)
    message.warning(t('common.failed'))
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors)
      add()
    else
      console.log(errors)
  })
}

// ---- 两步验证(2FA)管理（作用于当前登录账号）----
const twofaEnabled = ref(false)
const twofaOtpauth = ref('')
const twofaSecret = ref('')
const twofaCode = ref('')
const twofaBusy = ref(false)

async function loadTwoFAStatus() {
  try {
    const res = await twofaStatus<Common.Response<any>>()
    if (res.code === 0)
      twofaEnabled.value = !!res.data.enabled
  }
  catch { /* ignore */ }
}

async function enableTwoFA() {
  twofaBusy.value = true
  try {
    const res = await twofaEnable<Common.Response<any>>()
    if (res.code === 0) {
      twofaOtpauth.value = res.data.otpauth
      twofaSecret.value = res.data.secret
      message.info('请使用身份验证器扫描二维码，或手动输入下方密钥')
    }
  }
  finally {
    twofaBusy.value = false
  }
}

async function confirmTwoFA() {
  if (!twofaCode.value || twofaCode.value.length !== 6) {
    message.warning('请输入6位验证码')
    return
  }
  twofaBusy.value = true
  try {
    const res = await twofaConfirm<Common.Response<any>>({ code: twofaCode.value })
    if (res.code === 0) {
      twofaEnabled.value = true
      twofaOtpauth.value = ''
      twofaSecret.value = ''
      twofaCode.value = ''
      message.success('两步验证已启用')
    }
    else if (res.code !== -1) {
      message.warning(res.msg || '验证码错误')
    }
  }
  finally {
    twofaBusy.value = false
  }
}

async function disableTwoFA() {
  if (!twofaCode.value || twofaCode.value.length !== 6) {
    message.warning('请输入6位验证码')
    return
  }
  twofaBusy.value = true
  try {
    const res = await twofaDisable<Common.Response<any>>({ code: twofaCode.value })
    if (res.code === 0) {
      twofaEnabled.value = false
      twofaCode.value = ''
      message.success('两步验证已关闭')
    }
    else if (res.code !== -1) {
      message.warning(res.msg || '验证码错误')
    }
  }
  finally {
    twofaBusy.value = false
  }
}

onMounted(() => {
  loadDepartments()
  loadRoles()
  loadTwoFAStatus()
})
</script>

<template>
  <RoundCardModal v-model:show="show" size="small" preset="card" style="width: 400px" :title="`${userInfo?.id ? $t('common.edit') : $t('common.add')}`">
    <NForm ref="formRef" :model="model" :rules="rules">
      <NFormItem path="username" :label="$t('common.username')">
        <NInput v-model:value="model.username" type="text" :placeholder="$t('common.inputPlaceholder')" />
      </NFormItem>

      <NFormItem path="name" :label="$t('common.nikeName')">
        <NInput v-model:value="model.name" type="text" :placeholder="$t('common.inputPlaceholder')" />
      </NFormItem>

      <NFormItem path="role" :label="$t('adminSettingUsers.role')">
        <NSelect
          v-model:value="model.role"
          :options="roleOtions"
        />
      </NFormItem>

      <NFormItem path="departmentId" :label="$t('common.department')">
        <NSelect
          v-model:value="model.departmentId"
          :options="departmentOptions"
        />
      </NFormItem>

      <NFormItem path="password" :label="$t('common.password')">
        <NInput v-model:value="model.password" :maxlength="20" type="password" :placeholder="`${userInfo?.id ? $t('adminSettingUsers.EditpasswordPlaceholder') : $t('adminSettingUsers.passwordPlaceholder')}`" />
      </NFormItem>

      <NFormItem label="两步验证(2FA)">
        <div class="w-full">
          <div class="text-sm">
            <span v-if="!twofaEnabled" class="text-slate-500">当前登录账号未启用。启用后登录需输入动态验证码。</span>
            <span v-else class="text-green-600">当前登录账号已启用两步验证。</span>
            <NButton v-if="!twofaEnabled && !twofaSecret" size="small" type="primary" class="ml-2" :loading="twofaBusy" @click="enableTwoFA">获取密钥</NButton>
          </div>
          <div v-if="twofaSecret" class="mt-2 p-2 bg-slate-100 rounded text-xs break-all">
            <div>密钥(手动输入): <b>{{ twofaSecret }}</b></div>
            <div class="mt-1">或复制链接到身份验证器: <span class="break-all">{{ twofaOtpauth }}</span></div>
          </div>
          <div v-if="twofaSecret || twofaEnabled" class="mt-2 flex items-center gap-2">
            <NInput v-model:value="twofaCode" :maxlength="6" placeholder="6位验证码" />
            <NButton v-if="twofaSecret" size="small" type="success" :loading="twofaBusy" @click="confirmTwoFA">确认启用</NButton>
            <NButton v-if="twofaEnabled" size="small" type="warning" :loading="twofaBusy" @click="disableTwoFA">关闭</NButton>
          </div>
        </div>
      </NFormItem>
    </NForm>

    <template #footer>
      <div class="float-right">
        <NButton type="success" size="small" @click="handleValidateButtonClick">
          {{ $t('common.save') }}
        </NButton>
      </div>
    </template>
  </RoundCardModal>
</template>