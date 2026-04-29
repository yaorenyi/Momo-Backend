<template>
  <AdminLayout :baseUrl="apiUrl" @logout="logout" @refresh="loadSettings">
    <div class="max-w-4xl mx-auto space-y-8">
      <div class="flex items-center justify-between">
        <h1 class="text-2xl font-bold text-gray-800">系统设置</h1>
        <div class="flex items-center gap-3">
          <span v-if="saved" class="text-sm text-green-600 font-medium">保存成功</span>
          <button @click="saveSettings" :disabled="loading"
            class="px-5 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors text-sm font-medium">
            {{ loading ? '保存中...' : '保存设置' }}
          </button>
        </div>
      </div>

      <!-- 站点设置 -->
      <section class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
        <h2 class="text-lg font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <i class="fa-solid fa-globe text-blue-500"></i> 站点设置
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">站点名称</label>
            <input v-model="form.site_name" type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">管理员邮箱</label>
            <input v-model="form.admin_email" type="email"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">允许的跨域来源 (CORS)</label>
            <input v-model="form.allow_origin" type="text" placeholder="http://localhost:4321,https://example.com"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
            <p class="text-xs text-gray-400 mt-1">多个域名用逗号分隔</p>
          </div>
        </div>
      </section>

      <!-- SMTP 邮件设置 -->
      <section class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-gray-800 flex items-center gap-2">
            <i class="fa-solid fa-envelope text-blue-500"></i> 邮件通知设置
          </h2>
          <label class="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" v-model="form.email_enabled" class="sr-only peer" true-value="true" false-value="false">
            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
            <span class="ms-3 text-sm font-medium text-gray-700">
              {{ form.email_enabled === 'true' ? '已开启' : '已关闭' }}
            </span>
          </label>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">SMTP 服务器</label>
            <input v-model="form.smtp_host" type="text" placeholder="smtp.example.com"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">SMTP 端口</label>
            <input v-model="form.smtp_port" type="text" placeholder="465"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">邮箱用户名</label>
            <input v-model="form.email_user" type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">邮箱密码</label>
            <input v-model="form.email_password" type="password" placeholder="留空则不修改"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">安全连接 (SSL/TLS)</label>
            <select v-model="form.email_secure"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm bg-white">
              <option value="true">是 (端口 465)</option>
              <option value="false">否 (端口 587)</option>
            </select>
          </div>
        </div>

        <div class="mt-4 flex items-center gap-3 pt-4 border-t border-gray-100">
          <button @click="handleTestEmail" :disabled="testingEmail"
            class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 disabled:opacity-50 transition-colors text-sm font-medium flex items-center gap-2">
            <i v-if="testingEmail" class="fa-solid fa-circle-notch animate-spin"></i>
            <i v-else class="fa-solid fa-paper-plane"></i>
            {{ testingEmail ? '发送中...' : '发送测试邮件' }}
          </button>
          <p class="text-xs text-gray-400">保存设置后，点击测试 SMTP 配置是否正确</p>
        </div>
      </section>

      <!-- 邮件模板 -->
      <section class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
        <h2 class="text-lg font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <i class="fa-solid fa-pen-to-square text-blue-500"></i> 邮件模板
        </h2>
        <p class="text-sm text-gray-500 mb-4">
          留空则使用默认模板。可用占位符：
          <code class="bg-gray-100 px-1 rounded">&#123;&#123;toName&#125;&#125;</code>
          <code class="bg-gray-100 px-1 rounded">&#123;&#123;replyAuthor&#125;&#125;</code>
          <code class="bg-gray-100 px-1 rounded">&#123;&#123;postTitle&#125;&#125;</code>
          <code class="bg-gray-100 px-1 rounded">&#123;&#123;postUrl&#125;&#125;</code>
          <code class="bg-gray-100 px-1 rounded">&#123;&#123;commentContent&#125;&#125;</code>
          <code class="bg-gray-100 px-1 rounded">&#123;&#123;parentComment&#125;&#125;</code>
          <code class="bg-gray-100 px-1 rounded">&#123;&#123;replyContent&#125;&#125;</code>
        </p>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">回复通知模板</label>
            <textarea v-model="form.reply_template" rows="6"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm font-mono"
              placeholder="留空使用默认模板"></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">新评论通知模板</label>
            <textarea v-model="form.notification_template" rows="6"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm font-mono"
              placeholder="留空使用默认模板"></textarea>
          </div>
        </div>
      </section>

      <!-- 账户安全 -->
      <section class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
        <h2 class="text-lg font-semibold text-gray-800 mb-4 flex items-center gap-2">
          <i class="fa-solid fa-shield-halved text-blue-500"></i> 账户安全
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">当前用户名</label>
            <input :value="form.admin_name || adminName" type="text" disabled
              class="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-gray-50 text-gray-500 cursor-not-allowed" />
          </div>
          <div></div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">新用户名</label>
            <input v-model="passwordForm.new_name" type="text" placeholder="留空则不修改"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">新密码</label>
            <input v-model="passwordForm.new_password" type="password" placeholder="留空则不修改" minlength="4"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">确认新密码</label>
            <input v-model="passwordForm.confirm_password" type="password" placeholder="再次输入新密码"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">当前密码 <span class="text-red-500">*</span></label>
            <input v-model="passwordForm.old_password" type="password" placeholder="输入当前密码以确认修改"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm" />
          </div>
        </div>
        <div class="mt-4 flex items-center gap-3">
          <button @click="changePassword" :disabled="changingPassword"
            class="px-5 py-2 bg-amber-600 text-white rounded-lg hover:bg-amber-700 disabled:opacity-50 transition-colors text-sm font-medium">
            {{ changingPassword ? '更新中...' : '修改凭据' }}
          </button>
          <p class="text-xs text-gray-400">修改后需要重新登录</p>
        </div>
      </section>
    </div>
  </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import request from '../utils/request'
import toast from '../utils/toast'
import AdminLayout from '../components/AdminLayout.vue'

const router = useRouter()
const apiUrl = ref(localStorage.getItem('apiUrl') || window.location.origin)
const adminName = ref(localStorage.getItem('admin_name') || '')

const loading = ref(false)
const saved = ref(false)

const form = reactive({
  site_name: '',
  admin_email: '',
  allow_origin: '',
  smtp_host: '',
  smtp_port: '',
  email_user: '',
  email_password: '',
  email_secure: 'true',
  email_enabled: 'true',
  reply_template: '',
  notification_template: '',
})

const passwordForm = reactive({
  old_password: '',
  new_name: '',
  new_password: '',
  confirm_password: '',
})

const loadSettings = async () => {
  try {
    const res = await request.get('/admin/settings')
    if (res.code === 200 && res.data) {
      Object.assign(form, res.data)
    }
  } catch (e) {
    console.error('Failed to load settings:', e)
  }
}

onMounted(() => {
  loadSettings()
})

const saveSettings = async () => {
  loading.value = true
  saved.value = false
  try {
    const payload = { ...form }
    if (!payload.email_password) {
      delete payload.email_password
    }
    const res = await request.put('/admin/settings', payload)
    if (res.code === 200) {
      saved.value = true
      toast.success('设置已保存')
      setTimeout(() => { saved.value = false }, 3000)
    }
  } catch (e) {
    console.error('Failed to save settings:', e)
  } finally {
    loading.value = false
  }
}

const changingPassword = ref(false)
const testingEmail = ref(false)

const handleTestEmail = async () => {
  testingEmail.value = true
  try {
    const res = await request.post('/admin/settings/test-email')
    if (res.code === 200) {
      toast.success('测试邮件已发送，请查收')
    }
  } catch (e) {
    toast.error(e.message || '测试邮件发送失败')
  } finally {
    testingEmail.value = false
  }
}

const changePassword = async () => {
  if (!passwordForm.old_password) {
    toast.warning('请输入当前密码')
    return
  }
  if (!passwordForm.new_name && !passwordForm.new_password) {
    toast.warning('请填写新用户名或新密码')
    return
  }
  if (passwordForm.new_password && passwordForm.new_password.length < 4) {
    toast.warning('密码长度不能少于4位')
    return
  }
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    toast.warning('两次输入的密码不一致')
    return
  }

  changingPassword.value = true
  try {
    const currentName = form.admin_name || adminName.value
    const payload = {
      old_name: currentName,
      old_password: passwordForm.old_password,
      new_name: passwordForm.new_name || currentName,
      new_password: passwordForm.new_password || passwordForm.old_password,
    }
    const res = await request.put('/admin/password', payload)
    if (res.code === 200) {
      toast.success('凭据已更新，请重新登录')
      localStorage.removeItem('token')
      setTimeout(() => router.push('/login'), 1500)
    }
  } catch (e) {
    toast.error(e.message || '密码修改失败')
  } finally {
    changingPassword.value = false
  }
}

const logout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}
</script>
