<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue'
import { NButton, NInput, NSlider, NSpace, NSwitch, NSelect, NCard, NGrid, NGridItem, useMessage } from 'naive-ui'
import { getLogoConfig, setLogoConfig, uploadLogo, getBackgroundConfig, setBackgroundConfig, uploadBackground, getPresetBackgrounds } from '@/api/system/systemSetting'
import { getLoginConfig, setLoginConfig } from '@/api/openness'
import { t } from '@/locales'

const message = useMessage()

const activeTab = ref<'logo' | 'background' | 'stickyNote' | 'loginSecurity'>('logo')

// 登录安全配置
const loginSecurity = reactive({
  loginCaptcha: false,
  loginAllowIps: '',
})

const logoConfig = reactive({
  imageUrl: '/assets/logo.png',
  size: 80,
  useCDN: false,
  cdnUrl: '',
})

const backgroundConfig = reactive({
  imageUrl: '/assets/defaultBackground.webp',
  displayMode: 'cover',
  useCustomUrl: false,
  customUrl: '',
})

// 便签配置（localStorage 存储）
const STICKY_NOTE_STORAGE_KEY = 'sun-panel-sticky-note-config'
interface StickyNoteConfig {
  enabled: boolean
  transparent: boolean
}
const stickyNoteConfig = reactive<StickyNoteConfig>({
  enabled: true,
  transparent: false,
})

const presetBackgrounds = ref<any[]>([])

const displayModeOptions = [
  { label: '拉伸', value: 'cover' },
  { label: '平铺', value: 'repeat' },
  { label: '居中', value: 'center' },
  { label: '自适应', value: 'stretch' },
]



async function loadLogoConfig() {
  const { data } = await getLogoConfig<any>()
  if (data) {
    logoConfig.imageUrl = data.imageUrl || '/assets/logo.png'
    logoConfig.size = data.size || 80
    logoConfig.useCDN = data.useCDN || false
    logoConfig.cdnUrl = data.cdnUrl || ''
  }
}

async function loadBackgroundConfig() {
  const { data } = await getBackgroundConfig<any>()
  if (data) {
    backgroundConfig.imageUrl = data.imageUrl || '/assets/defaultBackground.webp'
    backgroundConfig.displayMode = data.displayMode || 'cover'
    backgroundConfig.useCustomUrl = data.useCustomUrl || false
    backgroundConfig.customUrl = data.customUrl || ''
  }
}

async function loadPresetBackgrounds() {
  const { data } = await getPresetBackgrounds<any[]>()
  if (data) {
    presetBackgrounds.value = data
  }
}

function loadStickyNoteConfig() {
  try {
    const saved = localStorage.getItem(STICKY_NOTE_STORAGE_KEY)
    if (saved) {
      const cfg = JSON.parse(saved)
      stickyNoteConfig.enabled = cfg.enabled !== false // 默认启用
      stickyNoteConfig.transparent = !!cfg.transparent
    }
  } catch { /* ignore */ }
}

async function loadLoginSecurity() {
  const { data } = await getLoginConfig<any>()
  if (data) {
    loginSecurity.loginCaptcha = !!data.loginCaptcha
    loginSecurity.loginAllowIps = data.loginAllowIps || ''
  }
}

async function handleSaveLoginSecurity() {
  const result = await setLoginConfig({
    loginCaptcha: loginSecurity.loginCaptcha,
    loginAllowIps: loginSecurity.loginAllowIps,
  })
  if (result.code === 0) {
    message.success(t('common.success'))
  }
}

function saveStickyNoteConfig() {
  localStorage.setItem(STICKY_NOTE_STORAGE_KEY, JSON.stringify({
    enabled: stickyNoteConfig.enabled,
    transparent: stickyNoteConfig.transparent,
  }))
  // 触发 storage 事件让 StickyNotes 组件响应（同页面通过 custom event）
  window.dispatchEvent(new CustomEvent('sticky-note-config-changed'))
  message.success(t('common.success'))
}

async function handleSaveLogo() {
  const result = await setLogoConfig({
    imageUrl: logoConfig.imageUrl,
    size: logoConfig.size,
    useCDN: logoConfig.useCDN,
    cdnUrl: logoConfig.cdnUrl,
  })
  if (result.code === 0) {
    message.success(t('common.success'))
  }
}

async function handleSaveBackground() {
  const result = await setBackgroundConfig({
    imageUrl: backgroundConfig.imageUrl,
    displayMode: backgroundConfig.displayMode,
    useCustomUrl: backgroundConfig.useCustomUrl,
    customUrl: backgroundConfig.customUrl,
  })
  if (result.code === 0) {
    message.success(t('common.success'))
  }
}

async function handleLogoUpload(e: any) {
  const file = e.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('file', file)

  const result: { code: number; data: { imageUrl: string } } = await uploadLogo(formData)
  if (result.code === 0) {
    logoConfig.imageUrl = result.data.imageUrl
    logoConfig.useCDN = false
    message.success(t('common.uploadSuccess'))
    // 上传后自动保存配置，避免用户忘记点保存按钮
    await handleSaveLogo()
  }
}

async function handleBackgroundUpload(e: any) {
  const file = e.target.files[0]
  if (!file) return

  const formData = new FormData()
  formData.append('file', file)

  const result: { code: number; data: { imageUrl: string } } = await uploadBackground(formData)
  if (result.code === 0) {
    backgroundConfig.imageUrl = result.data.imageUrl
    backgroundConfig.useCustomUrl = false
    message.success(t('common.uploadSuccess'))
  }
}

function selectPresetBackground(preset: any) {
  backgroundConfig.imageUrl = preset.imageUrl
  backgroundConfig.useCustomUrl = false
}

onMounted(() => {
  loadLogoConfig()
  loadBackgroundConfig()
  loadPresetBackgrounds()
  loadStickyNoteConfig()
  loadLoginSecurity()
})
</script>

<template>
  <div class="overflow-auto pt-2">
    <div class="flex gap-2 mb-4">
      <NButton :type="activeTab === 'logo' ? 'primary' : 'default'" @click="activeTab = 'logo'">
        {{ $t('admin.setting.logo') }}
      </NButton>
      <NButton :type="activeTab === 'background' ? 'primary' : 'default'" @click="activeTab = 'background'">
        {{ $t('admin.setting.background') }}（登录页背景）
      </NButton>
      <NButton :type="activeTab === 'stickyNote' ? 'primary' : 'default'" @click="activeTab = 'stickyNote'">
        便签
      </NButton>
      <NButton :type="activeTab === 'loginSecurity' ? 'primary' : 'default'" @click="activeTab = 'loginSecurity'">
        登录安全
      </NButton>
    </div>

    <div v-if="activeTab === 'logo'">
      <NGrid cols="2" :x-gap="24">
        <NGridItem>
          <NCard :title="$t('admin.setting.logoPreview')" class="h-full">
            <div class="flex flex-col items-center">
              <div 
                class="border-2 border-gray-200 rounded-lg p-8 mb-4"
                :style="{ width: `${logoConfig.size * 2}px`, height: `${logoConfig.size * 2}px`, display: 'flex', alignItems: 'center', justifyContent: 'center' }"
              >
                <img 
                  :src="logoConfig.useCDN && logoConfig.cdnUrl ? logoConfig.cdnUrl : logoConfig.imageUrl" 
                  :style="{ width: `${logoConfig.size}px`, height: `${logoConfig.size}px`, objectFit: 'contain' }"
                  :alt="$t('common.logo')"
                />
              </div>
              <span class="text-sm text-gray-500">{{ logoConfig.size }}px</span>
            </div>
          </NCard>
        </NGridItem>
        <NGridItem>
          <NCard :title="$t('admin.setting.logoSettings')" class="h-full">
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium mb-2">{{ $t('admin.setting.uploadLogo') }}</label>
                <input
                  type="file"
                  accept="image/png,image/jpg,image/jpeg"
                  @change="handleLogoUpload"
                  class="w-full"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium mb-2">{{ $t('admin.setting.logoSize') }}: {{ logoConfig.size }}px</label>
                <NSlider 
                  v-model:value="logoConfig.size" 
                  :min="50" 
                  :max="200" 
                  :step="5"
                  :mark-step="50"
                />
              </div>

              <div>
                <NSpace>
                  <NSwitch v-model:value="logoConfig.useCDN" />
                  <span>{{ $t('admin.setting.useCDN') }}</span>
                </NSpace>
              </div>

              <div v-if="logoConfig.useCDN">
                <label class="block text-sm font-medium mb-2">{{ $t('admin.setting.cdnUrl') }}</label>
                <NInput v-model:value="logoConfig.cdnUrl" :placeholder="$t('admin.setting.cdnUrlPlaceholder')" />
              </div>

              <div class="flex justify-end">
                <NButton type="primary" @click="handleSaveLogo">{{ $t('common.save') }}</NButton>
              </div>
            </div>
          </NCard>
        </NGridItem>
      </NGrid>
    </div>

    <div v-if="activeTab === 'background'">
      <NGrid cols="2" :x-gap="24">
        <NGridItem>
          <NCard :title="$t('admin.setting.backgroundPreview')" class="h-full">
            <div 
              class="w-full h-64 border-2 border-gray-200 rounded-lg"
              :style="{
                backgroundImage: backgroundConfig.useCustomUrl && backgroundConfig.customUrl 
                  ? `url(${backgroundConfig.customUrl})` 
                  : `url(${backgroundConfig.imageUrl})`,
                backgroundSize: backgroundConfig.displayMode === 'cover' ? 'cover' : 
                               backgroundConfig.displayMode === 'repeat' ? 'auto' :
                               backgroundConfig.displayMode === 'center' ? 'contain' : '100% 100%',
                backgroundRepeat: backgroundConfig.displayMode === 'repeat' ? 'repeat' : 'no-repeat',
                backgroundPosition: 'center',
              }"
            />
          </NCard>
        </NGridItem>
        <NGridItem>
          <NCard :title="$t('admin.setting.backgroundSettings')" class="h-full">
            <div class="space-y-4">
              <div class="bg-blue-50 dark:bg-blue-900 p-3 rounded-lg mb-2">
                <span class="text-sm text-blue-600 dark:text-blue-300">此背景用于登录页面显示。主页背景请在"样式设置"中修改。</span>
              </div>
              <div>
                <label class="block text-sm font-medium mb-2">{{ $t('admin.setting.uploadBackground') }}</label>
                <input
                  type="file"
                  accept="image/png,image/jpg,image/jpeg"
                  @change="handleBackgroundUpload"
                  class="w-full"
                />
              </div>

              <div>
                <NSpace>
                  <NSwitch v-model:value="backgroundConfig.useCustomUrl" />
                  <span>{{ $t('admin.setting.useCustomUrl') }}</span>
                </NSpace>
              </div>

              <div v-if="backgroundConfig.useCustomUrl">
                <label class="block text-sm font-medium mb-2">{{ $t('admin.setting.customUrl') }}</label>
                <NInput v-model:value="backgroundConfig.customUrl" :placeholder="$t('admin.setting.customUrlPlaceholder')" />
              </div>

              <div>
                <label class="block text-sm font-medium mb-2">{{ $t('admin.setting.displayMode') }}</label>
                <NSelect
                  v-model:value="backgroundConfig.displayMode"
                  :options="displayModeOptions"
                  :placeholder="$t('common.select')"
                />
              </div>

              <div class="flex justify-end">
                <NButton type="primary" @click="handleSaveBackground">{{ $t('common.save') }}</NButton>
              </div>
            </div>
          </NCard>
        </NGridItem>
      </NGrid>

      <NCard :title="$t('admin.setting.presetBackgrounds')" class="mt-4">
        <div class="flex flex-wrap gap-4">
          <div
            v-for="preset in presetBackgrounds"
            :key="preset.id"
            class="w-24 h-24 cursor-pointer border-2 rounded-lg overflow-hidden hover:border-primary"
            :class="{ 'border-primary': backgroundConfig.imageUrl === preset.imageUrl }"
            :style="{
              background: preset.isGradient ? preset.imageUrl : `url(${preset.imageUrl})`,
              backgroundSize: 'cover',
              backgroundPosition: 'center',
            }"
            @click="selectPresetBackground(preset)"
          />
        </div>
      </NCard>
    </div>

    <!-- 便签设置 -->
    <div v-if="activeTab === 'stickyNote'">
      <NCard title="便签设置" class="h-full">
        <div class="space-y-5">
          <div class="bg-amber-50 dark:bg-amber-900/20 p-3 rounded-lg">
            <span class="text-sm text-amber-700 dark:text-amber-300">
              便签功能可在主页上显示可拖拽的便利贴。关闭后主页将不再显示便签组件和悬浮按钮。
            </span>
          </div>

          <div class="flex items-center justify-between py-2">
            <div>
              <div class="font-medium">启用便签</div>
              <div class="text-xs text-gray-500 mt-0.5">关闭后便签将完全隐藏（已有便签数据保留）</div>
            </div>
            <NSwitch v-model:value="stickyNoteConfig.enabled" @update:value="saveStickyNoteConfig" />
          </div>

          <div class="flex items-center justify-between py-2">
            <div>
              <div class="font-medium">透明模式</div>
              <div class="text-xs text-gray-500 mt-0.5">开启后便签卡片呈现毛玻璃半透明效果</div>
            </div>
            <NSwitch
              :disabled="!stickyNoteConfig.enabled"
              v-model:value="stickyNoteConfig.transparent"
              @update:value="saveStickyNoteConfig"
            />
          </div>

          <div v-if="!stickyNoteConfig.enabled" class="mt-3 p-3 bg-red-50 dark:bg-red-900/20 rounded-lg text-sm text-red-600 dark:text-red-400">
            便签已关闭。如需重新启用，请打开上方开关。
          </div>
        </div>
      </NCard>
    </div>

    <!-- 登录安全 -->
    <div v-if="activeTab === 'loginSecurity'">
      <NGrid cols="2" :x-gap="24">
        <NGridItem>
          <NCard title="登录限制" class="h-full">
            <div class="space-y-5">
              <div class="bg-blue-50 dark:bg-blue-900/20 p-3 rounded-lg">
                <span class="text-sm text-blue-600 dark:text-blue-300">
                  开启登录验证码可防止爆破；配置"允许登录 IP"后，仅列表内的 IP 可以登录，
                  留空表示不限制。配置错误不会锁死已登录会话，仅拦截新的登录请求。
                </span>
              </div>

              <div class="flex items-center justify-between py-2">
                <div>
                  <div class="font-medium">登录验证码</div>
                  <div class="text-xs text-gray-500 mt-0.5">开启后登录页需输入图形验证码</div>
                </div>
                <NSwitch v-model:value="loginSecurity.loginCaptcha" />
              </div>

              <div>
                <label class="block text-sm font-medium mb-2">允许登录的 IP</label>
                <NInput
                  v-model:value="loginSecurity.loginAllowIps"
                  type="textarea"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  placeholder="留空=不限制。多个用逗号分隔，支持 CIDR，例如：&#10;192.168.1.100, 10.0.0.0/24"
                />
                <p class="text-xs text-gray-500 mt-1">仅限制登录入口，不影响已登录用户的正常访问。</p>
              </div>

              <div class="flex justify-end">
                <NButton type="primary" @click="handleSaveLoginSecurity">保存</NButton>
              </div>
            </div>
          </NCard>
        </NGridItem>
        <NGridItem>
          <NCard title="说明" class="h-full">
            <ul class="list-disc list-inside space-y-2 text-sm text-gray-600 dark:text-gray-300">
              <li>IP 白名单与登录验证码均在"系统设置 → 登录安全"中配置。</li>
              <li>支持单个 IP（如 <code>203.0.113.5</code>）或网段（如 <code>192.168.0.0/16</code>）。</li>
              <li>被拦截的 IP 在登录时会提示"当前 IP 不在允许登录的列表中"。</li>
              <li>如误配导致自己无法登录，可通过服务器直接清空该配置项恢复。</li>
            </ul>
          </NCard>
        </NGridItem>
      </NGrid>
    </div>
  </div>
</template>