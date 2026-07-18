<script setup lang="ts">
import { NButton, NCard, NForm, NFormItem, NGradientText, NInput, NSelect, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { login } from '@/api'
import { useAppStore, useAuthStore } from '@/store'
import { SvgIcon } from '@/components/common'
import { router } from '@/router'
import { t } from '@/locales'
import { languageOptions } from '@/utils/defaultData'
import { getBackgroundConfig, getLogoConfig } from '@/api/system/systemSetting'
import type { Language } from '@/store/modules/app/helper'

// const userStore = useUserStore()
const authStore = useAuthStore()
const appStore = useAppStore()
const ms = useMessage()
const loading = ref(false)
const languageValue = ref<Language>(appStore.language)

// 背景图配置
const bgImageUrl = ref('/assets/defaultBackground.webp')
const bgDisplayMode = ref('cover')

// Logo配置
const logoUrl = ref('/assets/logo.png')
const logoSize = ref(80)
const logoError = ref(false)

async function loadLoginConfig() {
  try {
    const bgRes = await getBackgroundConfig<any>()
    if (bgRes.code === 0 && bgRes.data) {
      const imgUrl = bgRes.data.imageUrl
      if (imgUrl && imgUrl.trim() !== '') {
        bgImageUrl.value = imgUrl
        bgDisplayMode.value = bgRes.data.displayMode || 'cover'
        // 预加载背景图
        const img = new Image()
        img.onload = () => { bgLoaded.value = true }
        img.onerror = () => { bgLoaded.value = true }
        img.src = imgUrl
      } else {
        bgLoaded.value = true
      }
    } else {
      bgLoaded.value = true
    }
  } catch {
    bgLoaded.value = true
  }
  try {
    const logoRes = await getLogoConfig<any>()
    if (logoRes.code === 0 && logoRes.data) {
      const imgUrl = logoRes.data.imageUrl
      if (imgUrl && imgUrl.trim() !== '') {
        logoUrl.value = imgUrl
        logoSize.value = logoRes.data.size || 80
      }
    }
  } catch { /* ignore */ }
}

function onLogoError() {
  logoError.value = true
  logoUrl.value = '/assets/logo.png'
}

// 背景图预加载状态
const bgLoaded = ref(false)

onMounted(() => {
  loadLoginConfig()
})

// const isShowCaptcha = ref<boolean>(false)
// const isShowRegister = ref<boolean>(false)

const form = ref<Login.LoginReqest>({
  username: '',
  password: '',
})

const loginPost = async () => {
  loading.value = true
  try {
    const res = await login<Login.LoginResponse>(form.value)
    if (res.code === 0) {
      authStore.setToken(res.data.token)
      authStore.setUserInfo(res.data)

      setTimeout(() => {
        ms.success(`Hi ${res.data.name},${t('login.welcomeMessage')}`)
        loading.value = false
        router.push({ path: '/' })
      }, 500)
    }
    else {
      loading.value = false
      // captchaRef.value.refresh()
    }
  }
  catch (error) {
    loading.value = false
    // 请检查网络或者服务器错误
    console.error('Login page init error:', error)
  }
}

function handleSubmit() {
  // 点击登录按钮触发
  loginPost()
}

function handleChangeLanuage(value: Language) {
  languageValue.value = value
  appStore.setLanguage(value)
}
</script>

<template>
  <div
    class="login-container"
    :class="{ 'bg-loaded': bgLoaded }"
    :style="{
      backgroundImage: bgImageUrl.startsWith('linear-gradient') ? bgImageUrl : `url(${bgImageUrl})`,
      backgroundSize: bgDisplayMode === 'cover' ? 'cover' : bgDisplayMode === 'repeat' ? 'auto' : bgDisplayMode === 'stretch' ? '100% 100%' : 'contain',
      backgroundRepeat: bgDisplayMode === 'repeat' ? 'repeat' : 'no-repeat',
      backgroundPosition: 'center center',
      backgroundColor: bgImageUrl.startsWith('linear-gradient') ? 'transparent' : '#f2f6ff',
    }"
  >
    <NCard class="login-card" style="border-radius: 20px;">
      <div class="mb-5 flex items-center justify-end">
        <div class="mr-2">
          <SvgIcon icon="ion-language" style="width: 20;height: 20;" />
        </div>
        <div class="min-w-[100px]">
          <NSelect v-model:value="languageValue" size="small" :options="languageOptions" @update-value="handleChangeLanuage" />
        </div>
      </div>

      <div class="login-title">
        <img
          v-if="logoUrl && !logoError"
          :src="logoUrl"
          :style="{ width: `${logoSize}px`, height: `${logoSize}px`, margin: '0 auto' }"
          class="object-contain"
          alt="Logo"
          @error="onLogoError"
        />
        <NGradientText :size="30" type="success" class="!font-bold mt-2">
          {{ $t('common.appName') }}
        </NGradientText>
      </div>
      <NForm :model="form" label-width="100px" @keydown.enter="handleSubmit">
        <NFormItem>
          <NInput v-model:value="form.username" :placeholder="$t('login.usernamePlaceholder')">
            <template #prefix>
              <SvgIcon icon="ph:user-bold" />
            </template>
          </NInput>
        </NFormItem>

        <NFormItem>
          <NInput v-model:value="form.password" type="password" :placeholder="$t('login.passwordPlaceholder')">
            <template #prefix>
              <SvgIcon icon="mdi:password-outline" />
            </template>
          </NInput>
        </NFormItem>

        <!-- <NFormItem v-if="isShowCaptcha">
          <div class="w-[120px] h-[34px] mr-[20px] rounded border flex cursor-pointer">
            <Captcha ref="captchaRef" src="/api/captcha/getImage" />
          </div>
          <NInput v-model:value="form.vcode" type="text" placeholder="请输入图像验证码" />
        </NFormItem> -->
        <NFormItem style="margin-top: 10px">
          <NButton type="primary" block :loading="loading" @click="handleSubmit">
            {{ $t('login.loginButton') }}
          </NButton>
        </NFormItem>

        <!-- <div class="flex justify-end">
          <NButton v-if="isShowRegister" quaternary type="info" class="flex" @click="$router.push({ path: '/register' })">
            注册
          </NButton>
          <NButton quaternary type="info" class="flex" @click="$router.push({ path: '/resetPassword' })">
            忘记密码?
          </NButton>
        </div> -->

        <div class="flex justify-center text-slate-300">
          Powered By Sun-Panel
        </div>
      </NForm>
    </NCard>
  </div>
</template>

  <style scoped>
    .login-container {
        padding: 20px;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        background-color: #f2f6ff;
        background-size: cover;
        background-position: center;
        opacity: 0;
        transition: opacity 0.5s ease;
    }
    .login-container.bg-loaded {
        opacity: 1;
    }
    .login-container::before {
      content: '';
      position: absolute;
      top: 0; left: 0;
      width: 100%; height: 100%;
      background: rgba(0,0,0,0.3);
      z-index: 0;
    }

    /* 夜间模式 */
    .dark .login-container{
      background-color: rgb(43, 43, 43);
    }

    @media (min-width: 600px) {
        .login-card {
            width: auto;
            margin: 0px 10px;
        }
        .login-button {
            width: 100%;
        }
    }

    .login-card {
        margin: 20px;
        min-width:400px;
        position: relative;
        z-index: 1;
    }

  .login-title{
    text-align: center;
    margin: 20px;
  }
  </style>
