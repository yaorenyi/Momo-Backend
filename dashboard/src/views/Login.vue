<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 via-stone-50 to-slate-100 relative overflow-hidden font-sans">
    <!-- Abstract geometric background -->
    <div class="absolute inset-0 pointer-events-none">
      <div class="absolute -top-32 -right-32 w-[500px] h-[500px] rounded-full bg-gradient-to-br from-blue-50/60 to-slate-100/30 blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 w-[450px] h-[450px] rounded-full bg-gradient-to-tr from-stone-100/50 to-slate-50/30 blur-3xl"></div>
      <div class="absolute top-1/3 left-1/2 w-72 h-72 rounded-full bg-gradient-to-b from-slate-100/30 to-transparent blur-2xl"></div>
      <svg class="absolute top-[15%] right-[20%] w-24 h-24 opacity-[0.06]" viewBox="0 0 100 100" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round">
        <circle cx="50" cy="50" r="40"/>
        <circle cx="50" cy="50" r="25"/>
        <circle cx="50" cy="50" r="10"/>
        <line x1="50" y1="10" x2="50" y2="90"/>
        <line x1="10" y1="50" x2="90" y2="50"/>
      </svg>
      <svg class="absolute bottom-[20%] left-[15%] w-20 h-20 opacity-[0.05]" viewBox="0 0 100 100" fill="none" stroke="currentColor" stroke-width="1">
        <rect x="10" y="10" width="80" height="80" rx="8"/>
        <rect x="25" y="25" width="50" height="50" rx="4"/>
      </svg>
    </div>

    <!-- Login card -->
    <div class="relative w-full max-w-sm mx-4">
      <div class="bg-white rounded-2xl shadow-[0_20px_60px_-12px_rgba(0,0,0,0.08)] px-10 py-12">
        <div class="text-center mb-10">
          <h1 class="text-[22px] font-semibold text-gray-900 tracking-tight">登录</h1>
          <p class="mt-2 text-sm text-gray-400 font-normal">评论管理系统</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <input v-model="form.name" type="text" required placeholder="用户名"
              class="w-full px-4 py-[11px] border border-gray-200 rounded-xl placeholder:text-gray-400
                     focus:outline-none focus:ring-2 focus:ring-slate-400/20 focus:border-slate-300
                     transition-all text-sm bg-gray-50/60 hover:bg-white hover:border-gray-300" />
          </div>
          <div>
            <input v-model="form.password" type="password" required placeholder="密码"
              class="w-full px-4 py-[11px] border border-gray-200 rounded-xl placeholder:text-gray-400
                     focus:outline-none focus:ring-2 focus:ring-slate-400/20 focus:border-slate-300
                     transition-all text-sm bg-gray-50/60 hover:bg-white hover:border-gray-300" />
          </div>
          <button type="submit" :disabled="loading"
            class="w-full py-[11px] bg-slate-900 text-white text-sm font-medium rounded-xl
                   hover:bg-slate-800 focus:outline-none focus:ring-2 focus:ring-slate-400/20
                   transition-all disabled:opacity-50 disabled:cursor-not-allowed">
            {{ loading ? '验证中...' : '登录' }}
          </button>
        </form>
      </div>
    </div>

    <!-- 修改默认密码弹窗 -->
    <div v-if="showPasswordModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl max-w-sm w-full mx-4 p-7">
        <h3 class="text-lg font-semibold text-gray-900 mb-1">首次登录安全提醒</h3>
        <p class="text-sm text-gray-400 mb-6">
          您正在使用默认用户名和密码，请立即修改。
        </p>

        <form @submit.prevent="handleChangePassword" class="space-y-4">
          <div>
            <input v-model="passwordForm.new_name" type="text" required placeholder="新用户名"
              class="w-full px-4 py-[11px] border border-gray-200 rounded-xl placeholder:text-gray-400
                     focus:outline-none focus:ring-2 focus:ring-slate-400/20 focus:border-slate-300
                     transition-all text-sm bg-gray-50/60" />
          </div>
          <div>
            <input v-model="passwordForm.new_password" type="password" required minlength="4" placeholder="新密码"
              class="w-full px-4 py-[11px] border border-gray-200 rounded-xl placeholder:text-gray-400
                     focus:outline-none focus:ring-2 focus:ring-slate-400/20 focus:border-slate-300
                     transition-all text-sm bg-gray-50/60" />
          </div>
          <div>
            <input v-model="passwordForm.confirm_password" type="password" required minlength="4" placeholder="确认新密码"
              class="w-full px-4 py-[11px] border border-gray-200 rounded-xl placeholder:text-gray-400
                     focus:outline-none focus:ring-2 focus:ring-slate-400/20 focus:border-slate-300
                     transition-all text-sm bg-gray-50/60" />
          </div>
          <button type="submit" :disabled="changingPassword"
            class="w-full py-[11px] bg-slate-900 text-white text-sm font-medium rounded-xl
                   hover:bg-slate-800 focus:outline-none focus:ring-2 focus:ring-slate-400/20
                   transition-all disabled:opacity-50">
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
    toast.error('密码更新失败')
  } finally {
    changingPassword.value = false
  }
}

const handleLogin = async () => {
  if (!form.name || !form.password) {
    toast.warning('请填写用户名和密码')
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
    toast.error('登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}
</script>
