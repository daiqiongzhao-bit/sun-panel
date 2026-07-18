<script setup lang="ts">
import { computed, defineEmits, defineProps, onMounted, ref, watch } from 'vue'
import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import { NButton, NForm, NFormItem, NInput, NSelect, useMessage } from 'naive-ui'
import { edit as userManageEdit } from '@/api/panel/users'
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
    const { data } = await getRoleList<Common.ListResponse<any>>({ pageNum: 1, pageSize: 100 })
    if (data && data.list) {
      roleOtions.value = data.list.map((role: any) => ({
        label: role.name,
        value: role.id,
      }))
    }
  } catch (e) {
    // Fallback to hardcoded roles
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

onMounted(() => {
  loadDepartments()
  loadRoles()
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