<script lang="ts" setup>
import { computed, h, onMounted, reactive, ref } from 'vue'
import { NButton, NDataTable, NDropdown, NInput, NModal, NForm, NFormItem, NCheckbox, NSpace, NTag, useDialog, useMessage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { getRoleList, createRole, updateRole, deleteRole } from '@/api/system/role'
import { getPermissionMatrix, saveRolePermissions, getRolePermissions } from '@/api/system/permission'
import { exportRoles, importRoles } from '@/api/system/importExport'
import { SvgIcon } from '@/components/common'
import { t } from '@/locales'

const message = useMessage()
const dialog = useDialog()

const tableIsLoading = ref<boolean>(false)
const roleList = ref<any[]>([])
const allPermissions = ref<any[]>([])
const selectedRoleId = ref<number>(0)
const selectedRolePermissions = ref<number[]>([])
const keyWord = ref<string>('')

// 新增角色时的权限选择
const newRolePermissions = ref<number[]>([])

// 将扁平的权限数据按模块分组
interface PermissionModule {
  moduleCode: string
  moduleName: string
  permissions: Array<{ id: number; name: string }>
}
const groupedPermissions = computed<PermissionModule[]>(() => {
  const map = new Map<string, PermissionModule>()
  for (const p of allPermissions.value) {
    const code = p.moduleCode || 'other'
    if (!map.has(code)) {
      map.set(code, { moduleCode: code, moduleName: p.moduleName || code, permissions: [] })
    }
    map.get(code)!.permissions.push({ id: p.id, name: p.permissionName || p.name || '' })
  }
  return Array.from(map.values())
})

const showRoleModal = ref<boolean>(false)
const isEditRole = ref<boolean>(false)
const roleForm = reactive({
  id: 0,
  name: '',
  description: '',
  status: 1,
})

const showPermissionModal = ref<boolean>(false)

const columns = createColumns({
  edit: (row: any) => {
    isEditRole.value = true
    roleForm.id = row.id
    roleForm.name = row.name
    roleForm.description = row.description || ''
    roleForm.status = row.status
    showRoleModal.value = true
  },
  permissions: (row: any) => {
    selectedRoleId.value = row.id
    loadRolePermissions(row.id)
    showPermissionModal.value = true
  },
  delete: (row: any) => {
    dialog.warning({
      title: t('common.warning'),
      content: t('common.deletePrompt', { name: row.name }),
      positiveText: t('common.confirm'),
      negativeText: t('common.cancel'),
      onPositiveClick: () => {
        handleDelete(row.id)
      },
    })
  },
  toggleStatus: (row: any) => {
    const newStatus = row.status === 1 ? 0 : 1
    const action = newStatus === 1 ? t('common.active') : t('common.inactive')
    dialog.warning({
      title: t('common.warning'),
      content: `${t('common.deletePrompt', { name: row.name })} (${action})`,
      positiveText: t('common.confirm'),
      negativeText: t('common.cancel'),
      onPositiveClick: () => {
        handleToggleStatus(row, newStatus)
      },
    })
  },
})

function createColumns({ edit, permissions, delete: del, toggleStatus }: { edit: (row: any) => void; permissions: (row: any) => void; delete: (row: any) => void; toggleStatus: (row: any) => void }) {
  return [
    {
      title: t('common.name'),
      key: 'name',
    },
    {
      title: t('common.description'),
      key: 'description',
      render(row: any) {
        return row.description || '-'
      },
    },
    {
      title: t('common.status'),
      key: 'status',
      render(row: any) {
        return h(NTag, row.status === 1 ? { type: 'success' } : {}, row.status === 1 ? t('common.active') : t('common.inactive'))
      },
    },
    {
      title: t('common.action'),
      key: '',
      render(row: any) {
        const btn = h(
          NButton,
          { strong: true, tertiary: true, size: 'small' },
          { default: () => h(SvgIcon, { icon: 'mingcute:more-1-fill' }) },
        )
        return h(
          NDropdown,
          {
            trigger: 'click',
            onSelect(key: string | number) {
              switch (key) {
                case 'edit':
                  edit(row)
                  break
                case 'permissions':
                  permissions(row)
                  break
                case 'delete':
                  del(row)
                  break
              }
            },
            options: [
              { label: t('common.edit'), key: 'edit' },
              { label: t('admin.setting.permissions'), key: 'permissions' },
              { label: row.status === 1 ? t('common.inactive') : t('common.active'), key: 'toggleStatus' },
              { label: t('common.delete'), key: 'delete' },
            ],
          },
          { default: () => btn },
        )
      },
    },
  ] as DataTableColumns<any>
}

const pagination = reactive({
  page: 1,
  showSizePicker: true,
  pageSizes: [10, 30, 50],
  pageSize: 10,
  itemCount: 0,
  onChange: (page: number) => {
    pagination.page = page
    loadRoleList()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    loadRoleList()
  },
})

async function loadRoleList() {
  tableIsLoading.value = true
  const { data } = await getRoleList<Common.ListResponse<any>>({
    pageNum: pagination.page,
    pageSize: pagination.pageSize,
    name: keyWord.value,
  })
  if (data) {
    roleList.value = data.list || []
    pagination.itemCount = (data as any).total || 0
  }
  tableIsLoading.value = false
}

async function loadPermissionMatrix(roleId: number = 0) {
  const { data } = await getPermissionMatrix<any>({ roleId: roleId || undefined })
  if (data && Array.isArray(data)) {
    allPermissions.value = data
    // 如果传了 roleId，从返回的数据中提取已选中的权限ID
    if (roleId > 0) {
      const checkedIds = data.filter((p: any) => p.checked).map((p: any) => p.id)
      selectedRolePermissions.value = checkedIds
    }
  }
}

async function loadRolePermissions(roleId: number) {
  // 加载带权限状态的权限矩阵
  await loadPermissionMatrix(roleId)
}

function handleAddRole() {
  isEditRole.value = false
  roleForm.id = 0
  roleForm.name = ''
  roleForm.description = ''
  roleForm.status = 1
  newRolePermissions.value = []
  showRoleModal.value = true
}

// 导出角色
async function handleExportRoles() {
  const { code, data } = await exportRoles<any>()
  if (code === 0) {
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `sun-panel-roles-${new Date().toISOString().slice(0, 10)}.json`
    a.click()
    URL.revokeObjectURL(url)
    message.success(t('admin.setting.exportRolesSuccess'))
  }
}

// 导入角色
const importRoleFileInput = ref<HTMLInputElement | null>(null)
function handleImportRolesClick() {
  importRoleFileInput.value?.click()
}
async function handleImportRolesFile(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  try {
    const text = await file.text()
    const roles = JSON.parse(text)
    const { code, data } = await importRoles<any>({ roles: Array.isArray(roles) ? roles : [roles] })
    if (code === 0) {
      message.success(t('admin.setting.importRolesSuccess', { imported: data.imported || 0 }))
      loadRoleList()
    }
  }
  catch {
    message.error(t('common.failed'))
  }
  finally {
    input.value = ''
  }
}

async function handleSaveRole() {
  if (!roleForm.name.trim()) {
    message.error(t('common.nameRequired'))
    return
  }
  let result
  if (isEditRole.value) {
    result = await updateRole({ id: roleForm.id, name: roleForm.name, description: roleForm.description, status: roleForm.status })
  } else {
    result = await createRole({ name: roleForm.name, description: roleForm.description, status: roleForm.status })
    // 新增角色后，如果有选择权限，自动保存
    if (result.code === 0 && result.data && newRolePermissions.value.length > 0) {
      await saveRolePermissions({ roleId: (result.data as any).id, permissionIds: newRolePermissions.value })
    }
  }
  if (result.code === 0) {
    message.success(t('common.success'))
    showRoleModal.value = false
    newRolePermissions.value = []
    loadRoleList()
  }
}

async function handleDelete(id: number) {
  const result = await deleteRole({ id })
  if (result.code === 0) {
    message.success(t('common.deleteSuccess'))
    loadRoleList()
  }
}

async function handleToggleStatus(row: any, newStatus: number) {
  const result = await updateRole({ id: row.id, status: newStatus })
  if (result.code === 0) {
    message.success(t('common.success'))
    loadRoleList()
  }
}

function toggleModulePermissions(module: PermissionModule) {
  const modulePermissionIds = module.permissions.map(p => p.id)
  const allSelected = modulePermissionIds.every(id => selectedRolePermissions.value.includes(id))

  if (allSelected) {
    selectedRolePermissions.value = selectedRolePermissions.value.filter(id => !modulePermissionIds.includes(id))
  } else {
    modulePermissionIds.forEach((id) => {
      if (!selectedRolePermissions.value.includes(id)) {
        selectedRolePermissions.value.push(id)
      }
    })
  }
}

// 新增角色时的模块权限切换
function toggleNewRoleModulePermissions(module: PermissionModule) {
  const modulePermissionIds = module.permissions.map(p => p.id)
  const allSelected = modulePermissionIds.every(id => newRolePermissions.value.includes(id))

  if (allSelected) {
    newRolePermissions.value = newRolePermissions.value.filter(id => !modulePermissionIds.includes(id))
  } else {
    modulePermissionIds.forEach((id) => {
      if (!newRolePermissions.value.includes(id)) {
        newRolePermissions.value.push(id)
      }
    })
  }
}

function togglePermission(permissionId: number) {
  const index = selectedRolePermissions.value.indexOf(permissionId)
  if (index > -1) {
    selectedRolePermissions.value.splice(index, 1)
  } else {
    selectedRolePermissions.value.push(permissionId)
  }
}

function toggleNewRolePermission(permissionId: number) {
  const index = newRolePermissions.value.indexOf(permissionId)
  if (index > -1) {
    newRolePermissions.value.splice(index, 1)
  } else {
    newRolePermissions.value.push(permissionId)
  }
}

async function savePermissions() {
  const result = await saveRolePermissions({ roleId: selectedRoleId.value, permissionIds: selectedRolePermissions.value })
  if (result.code === 0) {
    message.success(t('common.success'))
    showPermissionModal.value = false
  }
}

onMounted(() => {
  loadRoleList()
  loadPermissionMatrix()
})
</script>

<template>
  <div class="overflow-auto pt-2">
    <div class="my-[10px] flex justify-between items-center">
      <div>
        <NButton type="primary" size="small" ghost @click="handleAddRole">
          {{ $t('common.add') }}
        </NButton>
        <NButton size="small" class="ml-2" @click="handleExportRoles">
          {{ $t('admin.setting.exportRoles') }}
        </NButton>
        <NButton size="small" class="ml-2" @click="handleImportRolesClick">
          {{ $t('admin.setting.importRoles') }}
        </NButton>
        <input ref="importRoleFileInput" type="file" accept=".json" style="display:none" @change="handleImportRolesFile">
      </div>
      <NInput v-model:value="keyWord" placeholder="搜索角色" size="small" @keyup.enter="loadRoleList" />
    </div>

    <NDataTable
      :columns="columns"
      :data="roleList"
      :pagination="pagination"
      :bordered="false"
      :loading="tableIsLoading"
      :remote="true"
    />

    <!-- 新增/编辑角色弹窗 -->
    <NModal v-model:show="showRoleModal" preset="card" :title="isEditRole ? $t('common.edit') : $t('common.add')" style="width: 700px">
      <NForm>
        <NFormItem :label="$t('common.name')" required>
          <NInput v-model:value="roleForm.name" :placeholder="$t('common.namePlaceholder')" />
        </NFormItem>
        <NFormItem :label="$t('common.description')">
          <NInput v-model:value="roleForm.description" :placeholder="$t('common.descriptionPlaceholder')" />
        </NFormItem>
        <NFormItem :label="$t('common.status')">
          <NSwitch v-model:value="roleForm.status" :checked-value="1" :unchecked-value="0">
            <template #checked>{{ $t('common.enabled') }}</template>
            <template #unchecked>{{ $t('common.disabled') }}</template>
          </NSwitch>
        </NFormItem>
        <!-- 新增角色时可选择权限 -->
        <NFormItem v-if="!isEditRole && groupedPermissions.length > 0" label="权限配置">
          <div class="max-h-[300px] overflow-y-auto">
            <div v-for="module in groupedPermissions" :key="module.moduleCode" class="mb-3">
              <div class="flex items-center mb-1">
                <NCheckbox :checked="module.permissions.every((p: any) => newRolePermissions.includes(p.id))" @update:checked="toggleNewRoleModulePermissions(module)" />
                <span class="ml-2 font-bold">{{ module.moduleName }}</span>
              </div>
              <div class="ml-6 flex flex-wrap gap-2">
                <NCheckbox
                  v-for="perm in module.permissions"
                  :key="perm.id"
                  :checked="newRolePermissions.includes(perm.id)"
                  @update:checked="toggleNewRolePermission(perm.id)"
                  size="small"
                >
                  {{ perm.name }}
                </NCheckbox>
              </div>
            </div>
          </div>
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace>
          <NButton @click="showRoleModal = false">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" @click="handleSaveRole">{{ $t('common.confirm') }}</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- 权限配置弹窗 -->
    <NModal v-model:show="showPermissionModal" preset="card" title="权限配置" style="width: 700px">
      <div v-if="groupedPermissions.length > 0" class="max-h-[400px] overflow-y-auto">
        <div v-for="module in groupedPermissions" :key="module.moduleCode" class="mb-3">
          <div class="flex items-center mb-1">
            <NCheckbox :checked="module.permissions.every((p: any) => selectedRolePermissions.includes(p.id))" @update:checked="toggleModulePermissions(module)" />
            <span class="ml-2 font-bold">{{ module.moduleName }}</span>
          </div>
          <div class="ml-6 flex flex-wrap gap-2">
            <NCheckbox
              v-for="perm in module.permissions"
              :key="perm.id"
              :checked="selectedRolePermissions.includes(perm.id)"
              @update:checked="togglePermission(perm.id)"
              size="small"
            >
              {{ perm.name }}
            </NCheckbox>
          </div>
        </div>
      </div>
      <div v-else class="text-center py-4 text-gray-400">
        暂无权限数据
      </div>
      <template #footer>
        <NSpace>
          <NButton @click="showPermissionModal = false">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" @click="savePermissions">{{ $t('common.confirm') }}</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>