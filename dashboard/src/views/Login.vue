<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 px-4 sm:px-6 lg:px-8 font-sans">
    <div class="max-w-md w-full space-y-8 bg-white p-10 rounded-2xl shadow-sm border border-gray-100">
      <div>
        <h2 class="mt-2 text-center text-3xl font-extrabold text-gray-900 tracking-tight">
          评论管理系统
        </h2>
        <p class="mt-2 text-center text-sm text-gray-500">
          请输入您的凭据以访问后台
        </p>
      </div>

      <form class="mt-8 space-y-5" @submit.prevent="handleLogin">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">后端 API 地址</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400 text-xs">
              API
            </div>
            <input 
              v-model="form.apiUrl" 
              type="text" 
              required 
              placeholder="https://api.example.com"
              class="appearance-none block w-full pl-12 pr-3 py-2.5 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all sm:text-sm"
            />
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
          <input 
            v-model="form.name" 
            type="text" 
            required 
            placeholder=""
            class="appearance-none block w-full px-3 py-2.5 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all sm:text-sm"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">密码</label>
          <input 
            v-model="form.password" 
            type="password" 
            required 
            placeholder=""
            class="appearance-none block w-full px-3 py-2.5 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all sm:text-sm"
          />
        </div>

        <div>
          <button 
            type="submit" 
            :disabled="loading"
            class="group relative w-full flex justify-center py-2.5 px-4 border border-transparent text-sm font-semibold rounded-lg text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors disabled:opacity-70 disabled:cursor-not-allowed"
          >
            <span v-if="loading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg class="animate-spin h-5 w-5 text-indigo-300" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            {{ loading ? '正在验证...' : '立即登录' }}
          </button>
        </div>
      </form>
    </div>

    <!-- 修改默认密码弹窗 -->
    <div v-if="showPasswordModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl max-w-md w-full mx-4 p-8">
        <div class="text-center mb-6">
          <div class="w-14 h-14 mx-auto bg-amber-100 rounded-full flex items-center justify-center mb-3">
            <i class="fa-solid fa-shield-halved text-amber-600 text-2xl"></i>
          </div>
          <h3 class="text-xl font-bold text-gray-900">首次登录安全提醒</h3>
          <p class="text-sm text-gray-500 mt-2">
            您正在使用默认用户名和密码，请立即修改以保护您的站点安全。
          </p>
        </div>

        <form @submit.prevent="handleChangePassword" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">新用户名</label>
            <input v-model="passwordForm.new_name" type="text" required
              class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">新密码</label>
            <input v-model="passwordForm.new_password" type="password" required minlength="4"
              class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">确认新密码</label>
            <input v-model="passwordForm.confirm_password" type="password" required minlength="4"
              class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <button type="submit" :disabled="changingPassword"
            class="w-full py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors text-sm font-semibold">
            {{ changingPassword ? '更新中...' : '立即更新' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import request from '../utils/request'
import toast from '../utils/toast'

const router = useRouter()
const loading = ref(false)
const showPasswordModal = ref(false)
const changingPassword = ref(false)

const form = reactive({
  apiUrl: localStorage.getItem('apiUrl') || window.location.origin,
  name: '',
  password: ''
})

const passwordForm = reactive({
  new_name: '',
  new_password: '',
  confirm_password: ''
})

const handleChangePassword = async () => {
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    toast.warning('两次输入的密码不一致')
    return
  }
  if (passwordForm.new_password.length < 4) {
    toast.warning('密码长度不能少于4位')
    return
  }

  changingPassword.value = true
  try {
    const res = await request.put('/admin/password', {
      old_name: form.name,
      old_password: form.password,
      new_name: passwordForm.new_name,
      new_password: passwordForm.new_password
    })
    if (res.code === 200) {
      toast.success('管理员凭据已更新，请重新登录')
      localStorage.removeItem('token')
      showPasswordModal.value = false
      form.name = ''
      form.password = ''
    }
  } catch (e) {
    console.error('Password change failed:', e)
  } finally {
    changingPassword.value = false
  }
}

const handleLogin = async () => {
  // 基础表单验证
  if (!form.apiUrl || !form.name || !form.password) {
    toast.warning('请完整填写登录信息')
    return
  }

  loading.value = true
  const formattedApiUrl = form.apiUrl.replace(/\/$/, '')
  localStorage.setItem('apiUrl', formattedApiUrl)

  try {
    const payload = {
      name: form.name,
      password: form.password
    }
    const res = await request.post('/admin/login', payload)
    
    if (res.code === 200) {
      toast.success('登录成功')
      localStorage.setItem('token', res.token)
      localStorage.setItem('admin_name', form.name)
      if (res.needChangePassword) {
        showPasswordModal.value = true
      } else {
        router.push('/')
      }
    }
  } catch (error) {
    console.error('Login failed:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* 大部分样式已通过 Tailwind 实现，此处只需处理特定交互 */
input::placeholder {
  color: #9ca3af;
}
</style>