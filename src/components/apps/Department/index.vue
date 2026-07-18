<script lang="ts" setup>import { h, onMounted, reactive, ref } from 'vue';
import { NButton, NTree, NModal, NForm, NFormItem, NInput, NSelect, NSpace, useMessage } from 'naive-ui';
import { getDepartmentTree, createDepartment, updateDepartment } from '@/api/system/department';
import { getAuthInfo } from '@/api/system/user';
import { SvgIcon } from '@/components/common';
import { t } from '@/locales';
const message = useMessage();
const treeData = ref<any[]>([]);
const userList = ref<any[]>([]);
const showModal = ref<boolean>(false);
const isEdit = ref<boolean>(false);
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
 const { data } = await getDepartmentTree<any[]>();
 treeData.value = transformTreeData(data || []);
}
async function loadUserList() {
 const { data } = await getAuthInfo<any>();
 if (data && data.users) {
 userList.value = data.users;
 }
}
function renderTreeNode(node: any) {
 return h('span', null, [
 h(SvgIcon, { icon: 'ic-baseline-add-business', size: 14, class: 'mr-1' }),
 node.label,
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
  <div class="overflow-auto pt-2 h-full flex flex-col">
    <div class="my-[10px]">
      <NButton type="primary" size="small" ghost @click="handleAdd(0)">
        {{ $t('common.add') }}
      </NButton>
    </div>

    <div class="flex-1">
      <NTree
        :data="treeData"
        :expanded-keys="expandedKeys"
        :selected-keys="selectedKeys"
        :render="renderTreeNode"
        @update:expanded-keys="expandedKeys = $event"
        @update:selected-keys="selectedKeys = $event"
        @click="(e: any, info: any) => handleEdit(info.node)"
      />
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