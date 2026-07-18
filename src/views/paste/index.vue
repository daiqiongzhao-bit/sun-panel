<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { NButton, NCard, NInput, NSpin, NTag, useMessage } from 'naive-ui'
import { downloadPasteBinFile, getPasteBinByCode } from '@/api/panel/pasteBin'

const route = useRoute()
const ms = useMessage()
const code = ref(route.params.code as string)
const loading = ref(false)
const password = ref('')
const needPassword = ref(false)
const showContent = ref(false)
const item = ref<any>(null)
const errorMsg = ref('')
const downloadLoading = ref(false)

const fileSizeStr = computed(() => {
  const size = item.value?.fileSize || 0
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  return (size / (1024 * 1024)).toFixed(1) + ' MB'
})

const isDownloadable = computed(() => {
  const content = item.value?.content || ''
  return content.startsWith('http://') || content.startsWith('https://') || content.startsWith('/')
})

async function fetchContent() {
  loading.value = true
  errorMsg.value = ''
  try {
    const { code: resCode, data, msg } = await getPasteBinByCode<any>({
      code: code.value,
      password: password.value || undefined,
    })
    if (resCode === 0) {
      item.value = data
      showContent.value = true
    }
    else if (resCode === 1004) {
      needPassword.value = true
    }
    else {
      errorMsg.value = msg || '内容不存在或已过期'
    }
  }
  catch {
    errorMsg.value = '请求失败，请稍后重试'
  }
  finally {
    loading.value = false
  }
}

function copyContent(text: string) {
  navigator.clipboard.writeText(text).then(() => {
    ms.success('已复制到剪贴板')
  })
}

function downloadFile() {
  if (isDownloadable.value && item.value) {
    // Direct URL download
    const url = item.value.content
    const a = document.createElement('a')
    a.href = url
    a.download = item.value.fileName || ''
    a.target = '_blank'
    a.rel = 'noopener noreferrer'
    a.click()
  }
  else if (item.value?.code) {
    // Use backend download endpoint
    const downloadUrl = `/api/panel/pasteBin/downloadFile`
    // Create form and submit
    const form = document.createElement('form')
    form.method = 'POST'
    form.action = downloadUrl
    form.target = '_blank'
    const codeInput = document.createElement('input')
    codeInput.type = 'hidden'
    codeInput.name = 'code'
    codeInput.value = item.value.code
    form.appendChild(codeInput)
    if (password.value) {
      const pwInput = document.createElement('input')
      pwInput.type = 'hidden'
      pwInput.name = 'password'
      pwInput.value = password.value
      form.appendChild(pwInput)
    }
    document.body.appendChild(form)
    form.submit()
    document.body.removeChild(form)
  }
}

onMounted(() => {
  fetchContent()
})
</script>

<template>
  <div class="paste-view-page min-h-screen flex items-center justify-center bg-gray-100 dark:bg-zinc-900 p-4">
    <NCard class="max-w-xl w-full" :bordered="true">
      <div class="text-center">
        <h2 class="text-xl font-bold mb-4">Sun-Panel 中转站</h2>
      </div>

      <NSpin :show="loading">
        <!-- Error -->
        <div v-if="errorMsg" class="text-center py-8">
          <div class="text-gray-400 text-5xl mb-4">📭</div>
          <p class="text-gray-500">{{ errorMsg }}</p>
        </div>

        <!-- Password prompt -->
        <div v-else-if="needPassword && !showContent" class="py-4">
          <p class="mb-3 text-center text-gray-600">此内容需要访问密码</p>
          <NInput
            v-model:value="password"
            type="password"
            placeholder="请输入访问密码"
            @keyup.enter="fetchContent"
          />
          <div class="text-center mt-4">
            <NButton type="primary" @click="fetchContent">验证</NButton>
          </div>
        </div>

        <!-- Content -->
        <div v-else-if="showContent && item" class="py-2">
          <div class="mb-3">
            <NTag :type="item.type === 1 ? 'info' : 'warning'" size="small">
              {{ item.type === 1 ? '文本' : '文件' }}
            </NTag>
            <span v-if="item.burnAfterRead === 1" class="ml-2">
              <NTag type="error" size="small">阅后即焚</NTag>
            </span>
          </div>

          <h3 class="text-lg font-semibold mb-3">{{ item.title }}</h3>

          <!-- Text content -->
          <div v-if="item.type === 1" class="mb-4">
            <div
              class="bg-gray-50 dark:bg-zinc-800 rounded-lg p-4 whitespace-pre-wrap break-all max-h-96 overflow-auto"
            >
              {{ item.content }}
            </div>
            <div class="mt-2 text-right">
              <NButton size="small" @click="copyContent(item.content)">复制内容</NButton>
            </div>
          </div>

          <!-- File content -->
          <div v-else-if="item.type === 2" class="mb-4 text-center">
            <div class="text-gray-400 text-4xl mb-2">📎</div>
            <p class="mb-1 font-semibold">{{ item.fileName || '文件' }}</p>
            <p class="text-xs text-gray-400 mb-2">大小: {{ fileSizeStr }}</p>
            <div class="mt-3 flex justify-center gap-2">
              <NButton v-if="item.content" type="primary" size="small" :loading="downloadLoading" @click="downloadFile">
                下载文件
              </NButton>
            </div>
          </div>

          <div class="text-xs text-gray-400 text-center mt-4">
            访问码: {{ item.code }} · 访问次数: {{ item.accessCnt }} · 过期时间: {{ item.expireAt }}
          </div>
        </div>
      </NSpin>
    </NCard>
  </div>
</template>

<style scoped>
.paste-view-page {
  min-height: 100vh;
}
</style>
