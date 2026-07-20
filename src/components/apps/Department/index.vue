<script lang="ts" setup>import { h, onMounted, reactive, ref } from 'vue';
import { NButton, NTree, NModal, NForm, NFormItem, NInput, NSelect, NSpace, NSpin, NEmpty, NTag, useMessage } from 'naive-ui';
import { getDepartmentTree, createDepartment, updateDepartment } from '@/api/system/department';
import { getAuthInfo } from '@/api/system/user';
import { SvgIcon } from '@/components/common';
import { t } from '@/locales';
const message = useMessage();
const treeData = ref<any[]>([]);
const userList = ref<any[]>([]);
const showModal = ref<boolean>(false);
const isEdit = ref<boolean>(false);
const loading = ref<boolean>(false);
const deptForm = reactive({
  id: 0,
  name: '',
  code: '',
  parentId: 0,
  description: '',
  leaderId: 0,
  sort: 0,
});
const expandedKeys = ref<number[]>([]);
const selectedKeys = ref<number[]>([]);
function transformTreeData(nodes: any[]): any[] {
  return nodes.map((node: any) => ({
    key: node.id,
    label: node.name,
    code: node.code,
    parentId: node.parentId,
    description: node.description,
    leaderId: node.leaderId,
    status: node.status,
    sort: node.sort,
    children: node.children ? transformTreeData(node.children) : [],
  }))
}

async function loadDepartmentTree() {
  loading.value = true;
  try {
    const { data } = await getDepartmentTree<any[]>();
    treeData.value = transformTreeData(data || []);
  }
  finally {
    loading.value = false;
  }
}
async function loadUserList() {
  const { data } = await getAuthInfo<any>();
  if (data && data.users) {
    userList.value = data.users;
  }
}
// 兼容 Naive UI 不同版本 render 签名：可能直接传 option，也可能传 { option }
function renderTreeNode(payload: any) {
  const node = payload?.option || payload;
  if (!node) return null;
  const disabled = node.status === 0;
  return h('div', { class: `dept-node${disabled ? ' dept-node--disabled' : ''}` }, [
    h('div', { class: 'dept-node__icon' }, [
      h(SvgIcon, { icon: 'ic-baseline-account-tree', size: 18 }),
    ]),
    h('div', { class: 'dept-node__body' }, [
      h('div', { class: 'dept-node__title' }, [
        node.label,
        h(NTag, { size: 'small', type: disabled ? 'default' : 'success', bordered: false, class: 'dept-node__tag' }, {
          default: () => disabled ? t('common.disabled') : t('common.enabled'),
        }),
      ]),
      h('div', { class: 'dept-node__meta' }, node.code || ''),
    ]),
    h('div', { class: 'dept-node__actions' }, [
      h(NButton, {
        size: 'tiny',
        tertiary: true,
        onClick: (e: MouseEvent) => {
          e.stopPropagation();
          handleEdit(node);
        },
      }, { default: () => t('common.edit') }),
    ]),
  ]);
}
function handleAdd(parentId: number = 0) {
  isEdit.value = false;
  deptForm.id = 0;
  deptForm.name = '';
  deptForm.code = '';
  deptForm.parentId = parentId;
  deptForm.description = '';
  deptForm.leaderId = 0;
  deptForm.sort = 0;
  showModal.value = true;
}
function handleEdit(node: any) {
  isEdit.value = true;
  deptForm.id = node.key;
  deptForm.name = node.label;
  deptForm.code = node.code || '';
  deptForm.parentId = node.parentId || 0;
  deptForm.description = node.description || '';
  deptForm.leaderId = node.leaderId || 0;
  deptForm.sort = node.sort || 0;
  showModal.value = true;
}

async function handleSave() {
  if (!deptForm.name.trim()) {
    message.error(t('common.nameRequired'));
    return;
  }
  if (!deptForm.code.trim()) {
    message.error(t('admin.setting.department.codeRequired'));
    return;
  }
  let result;
  if (isEdit.value) {
    result = await updateDepartment(deptForm);
  }
  else {
    result = await createDepartment(deptForm);
  }
  if (result.code === 0) {
    message.success(t('common.success'));
    showModal.value = false;
    loadDepartmentTree();
  }
}

onMounted(() => {
  loadDepartmentTree();
  loadUserList();
});
</script>

<template>
  <div class="dept-page overflow-auto h-full flex flex-col">
    <div class="dept-header">
      <div class="dept-header__title">
        <span class="dept-header__name">{{ $t('admin.setting.department.label') }}</span>
        <span class="dept-header__count">{{ treeData.length }}</span>
      </div>
      <NButton type="primary" size="small" @click="handleAdd(0)">
        <template #icon>
          <SvgIcon icon="ic-baseline-add" />
        </template>
        {{ $t('common.add') }}
      </NButton>
    </div>

    <div class="dept-body">
      <NSpin :show="loading">
        <NEmpty
          v-if="!loading && treeData.length === 0"
          :description="$t('common.noData')"
          class="dept-empty"
        >
          <template #extra>
            <NButton size="small" @click="handleAdd(0)">{{ $t('common.add') }}</NButton>
          </template>
        </NEmpty>

        <NTree
          v-else
          block-line
          :data="treeData"
          :expanded-keys="expandedKeys"
          :selected-keys="selectedKeys"
          :render="renderTreeNode"
          :default-expand-all="false"
          @update:expanded-keys="expandedKeys = $event"
          @update:selected-keys="selectedKeys = $event"
        />
      </NSpin>
    </div>

    <NModal v-model:show="showModal" preset="card" :title="isEdit ? $t('common.edit') : $t('common.add')" style="width: 500px">
      <NForm>
        <NFormItem :label="$t('common.name')" required>
          <NInput v-model:value="deptForm.name" :placeholder="$t('common.namePlaceholder')" />
        </NFormItem>
        <NFormItem :label="$t('admin.setting.department.code')" required>
          <NInput v-model:value="deptForm.code" :placeholder="$t('admin.setting.department.codePlaceholder')" />
        </NFormItem>
        <NFormItem :label="$t('admin.setting.department.parent')">
          <NSelect
            v-model:value="deptForm.parentId"
            :options="treeData.map(node => ({ label: node.label, value: node.key }))"
            :placeholder="$t('common.none')"
            clearable
          />
        </NFormItem>
        <NFormItem :label="$t('admin.setting.department.leader')">
          <NSelect
            v-model:value="deptForm.leaderId"
            :options="userList.map(user => ({ label: user.name || user.username, value: user.id }))"
            :placeholder="$t('common.select')"
            filterable
            clearable
          />
        </NFormItem>
        <NFormItem :label="$t('common.description')">
          <NInput v-model:value="deptForm.description" :placeholder="$t('common.descriptionPlaceholder')" />
        </NFormItem>
        <NFormItem :label="$t('common.sort')">
          <NInput v-model:value="deptForm.sort" type="number" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace>
          <NButton @click="showModal = false">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" @click="handleSave">{{ $t('common.confirm') }}</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style scoped>
.dept-page {
  padding: 4px 2px 16px;
}
.dept-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 10px 14px;
}
.dept-header__title {
  display: flex;
  align-items: center;
  gap: 10px;
}
.dept-header__name {
  font-size: 15px;
  font-weight: 600;
  color: var(--n-title-text-color, #1f2225);
}
.dept-header__count {
  min-width: 22px;
  height: 22px;
  padding: 0 7px;
  border-radius: 11px;
  background: var(--n-border-color, #eaecef);
  color: #8a8f99;
  font-size: 12px;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.dept-body {
  flex: 1;
  min-height: 120px;
}
.dept-empty {
  padding-top: 48px;
}
/* 树节点卡片样式 */
.dept-node {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 8px 10px;
  border-radius: 8px;
  transition: background-color 0.2s ease;
}
.dept-node:hover {
  background: var(--n-border-color, #f3f4f6);
}
.dept-node__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 8px;
  background: rgba(99, 102, 241, 0.12);
  color: #6366f1;
  flex-shrink: 0;
}
.dept-node__body {
  flex: 1;
  min-width: 0;
}
.dept-node__title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #1f2225;
  line-height: 1.3;
}
.dept-node__tag {
  font-weight: 500;
}
.dept-node__meta {
  font-size: 12px;
  color: #8a8f99;
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.dept-node__actions {
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.2s ease;
}
.dept-node:hover .dept-node__actions {
  opacity: 1;
}
/* 停用部门置灰 */
.dept-node--disabled .dept-node__icon {
  background: #eceef1;
  color: #b0b5bd;
}
.dept-node--disabled .dept-node__title {
  color: #9aa0a8;
  font-weight: 500;
}
.dept-node--disabled .dept-node__meta {
  color: #b8bdc5;
}
</style>
