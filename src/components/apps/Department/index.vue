<script lang="ts" setup>import { computed, h, onMounted, reactive, ref } from 'vue';
import { NButton, NTree, NModal, NForm, NFormItem, NInput, NSelect, NSpace, NSpin, NEmpty, NTag, useMessage } from 'naive-ui';
import { getDepartmentTree, createDepartment, updateDepartment } from '@/api/system/department';
import { getAuthInfo } from '@/api/system/user';
import { SvgIcon } from '@/components/common';
import { getDeterministicColor } from '@/utils/iconPlaceholder';
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

// 用户 id -> 名称 映射，用于在节点上展示负责人
const userMap = computed<Record<number, string>>(() => {
  const map: Record<number, string> = {};
  for (const u of userList.value) {
    map[u.id] = u.name || u.username || `#${u.id}`;
  }
  return map;
});

// 部门头像的确定性配色：同一部门始终为同一颜色
function deptColor(node: any): string {
  return getDeterministicColor(`${node.key ?? ''}-${node.label ?? ''}`);
}
function leaderName(leaderId?: number): string {
  if (!leaderId) return '';
  return userMap.value[leaderId] || '';
}

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
  const color = deptColor(node);
  const leader = leaderName(node.leaderId);
  return     h('div', { class: `dept-node${disabled ? ' dept-node--disabled' : ''}` }, [
    h('div', { class: 'dept-node__avatar' }, [
      h(SvgIcon, { icon: 'ic-baseline-account-tree', size: 18 }),
    ]),
    h('div', { class: 'dept-node__body' }, [
      h('div', { class: 'dept-node__title' }, [
        h('span', { class: 'dept-node__name' }, node.label),
        h(NTag, { size: 'small', type: disabled ? 'default' : 'success', bordered: false, class: 'dept-node__tag' }, {
          default: () => disabled ? t('common.disabled') : t('common.enabled'),
        }),
      ]),
      h('div', { class: 'dept-node__meta' }, [
        node.code ? h('span', { class: 'dept-node__code' }, node.code) : null,
        leader ? h('span', { class: 'dept-node__leader' }, [
          h(SvgIcon, { icon: 'ic-baseline-person', size: 13 }),
          ` ${leader}`,
        ]) : null,
      ]),
    ]),
    h('div', { class: 'dept-node__actions' }, [
      h(NButton, {
        size: 'tiny',
        quaternary: true,
        title: t('common.add'),
        onClick: (e: MouseEvent) => {
          e.stopPropagation();
          handleAdd(node.key);
        },
      }, { default: () => h(SvgIcon, { icon: 'ic-baseline-add', size: 16 }) }),
      h(NButton, {
        size: 'tiny',
        tertiary: true,
        title: t('common.edit'),
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

// 可选的上级部门（排除自身，避免循环引用）
const parentOptions = computed(() => {
  const selfId = isEdit.value ? deptForm.id : 0;
  const walk = (nodes: any[]): any[] =>
    nodes.flatMap((n) => {
      if (n.key === selfId) return [];
      const children = n.children && n.children.length ? walk(n.children) : [];
      return [{ label: n.label, value: n.key }, ...children];
    });
  return walk(treeData.value);
});

onMounted(() => {
  loadDepartmentTree();
  loadUserList();
});
</script>

<template>
  <div class="dept-page overflow-auto h-full flex flex-col">
    <div class="dept-header">
      <div class="dept-header__title">
        <div class="dept-header__icon">
          <SvgIcon icon="ic-baseline-account-tree" :size="18" />
        </div>
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

        <div v-else class="dept-tree-card">
          <NTree
            block-line
            :data="treeData"
            :expanded-keys="expandedKeys"
            :selected-keys="selectedKeys"
            :render="renderTreeNode"
            :default-expand-all="false"
            @update:expanded-keys="expandedKeys = $event"
            @update:selected-keys="selectedKeys = $event"
          />
        </div>
      </NSpin>
    </div>

    <NModal v-model:show="showModal" preset="card" :title="isEdit ? $t('common.edit') : $t('common.add')" style="width: 520px">
      <div class="dept-form">
        <NForm>
          <div class="dept-form__grid">
            <NFormItem :label="$t('common.name')" required>
              <NInput v-model:value="deptForm.name" :placeholder="$t('common.namePlaceholder')" />
            </NFormItem>
            <NFormItem :label="$t('admin.setting.department.code')" required>
              <NInput v-model:value="deptForm.code" :placeholder="$t('admin.setting.department.codePlaceholder')" />
            </NFormItem>
            <NFormItem :label="$t('admin.setting.department.parent')">
              <NSelect
                v-model:value="deptForm.parentId"
                :options="parentOptions"
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
          </div>
          <NFormItem :label="$t('common.description')">
            <NInput v-model:value="deptForm.description" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" :placeholder="$t('common.descriptionPlaceholder')" />
          </NFormItem>
          <NFormItem :label="$t('common.sort')">
            <NInput v-model:value="deptForm.sort" type="number" />
          </NFormItem>
        </NForm>
      </div>
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
.dept-header__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 9px;
  color: #fff;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
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
/* 树形容器：干净的白底卡片，杜绝大面积彩色块 */
.dept-tree-card {
  background: #ffffff;
  border: 1px solid var(--n-border-color, #e8eaed);
  border-radius: 16px;
  padding: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}
/* 树节点：独立白色小卡片，hover 才浮现阴影 */
.dept-node {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  margin: 6px 0;
  padding: 10px 12px;
  border-radius: 12px;
  background: #ffffff;
  border: 1px solid transparent;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
}
.dept-node:hover {
  background: #fafbff;
  border-color: var(--n-border-color, #e0e4ea);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.08);
}
.dept-node__avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  color: #fff;
  flex-shrink: 0;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  box-shadow: 0 2px 6px rgba(99, 102, 241, 0.25);
}
.dept-node__body {
  flex: 1;
  min-width: 0;
}
.dept-node__title {
  display: flex;
  align-items: center;
  gap: 8px;
  line-height: 1.3;
}
.dept-node__name {
  font-size: 14px;
  font-weight: 600;
  color: #1f2225;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.dept-node__tag {
  font-weight: 500;
  flex-shrink: 0;
}
.dept-node__meta {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  color: #8a8f99;
  margin-top: 3px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.dept-node__code {
  display: inline-flex;
  align-items: center;
  padding: 1px 7px;
  border-radius: 6px;
  background: rgba(99, 102, 241, 0.10);
  color: #6366f1;
  font-weight: 600;
}
.dept-node__leader {
  display: inline-flex;
  align-items: center;
  gap: 3px;
}
.dept-node__actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.2s ease;
}
.dept-node:hover .dept-node__actions {
  opacity: 1;
}
/* 停用部门置灰 */
.dept-node--disabled .dept-node__avatar {
  background: #eceef1 !important;
  color: #b0b5bd;
  box-shadow: none;
}
.dept-node--disabled .dept-node__name {
  color: #9aa0a8;
  font-weight: 500;
}
.dept-node--disabled .dept-node__meta {
  color: #b8bdc5;
}
.dept-node--disabled .dept-node__code {
  background: #eceef1;
  color: #9aa0a8;
}
/* 编辑弹窗表单 */
.dept-form__grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 16px;
}
@media (max-width: 520px) {
  .dept-form__grid {
    grid-template-columns: 1fr;
  }
}
</style>
