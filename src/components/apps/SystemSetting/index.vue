<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue'
import { NButton, NInput, NSlider, NSpace, NSwitch, NSelect, NCard, NGrid, NGridItem, useMessage } from 'naive-ui'
import { getLogoConfig, setLogoConfig, uploadLogo, getBackgroundConfig, setBackgroundConfig, uploadBackground, getPresetBackgrounds } from '@/api/system/systemSetting'
import { t } from '@/locales'

const message = useMessage()

const activeTab = ref<'logo' | 'background'>('logo')

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
  </div>
</template>