import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '../router'

const service = axios.create({
  timeout: 10000 // 建议超时时间稍微设长一点，以应对不同网络环境
})

// 请求拦截器：动态配置 baseURL 和添加 token
service.interceptors.request.use(
  config => {
    // 动态获取 LocalStorage 中的 API 地址
    const apiUrl = localStorage.getItem('apiUrl')
    if (apiUrl) {
      config.baseURL = apiUrl
    }

    // 处理 Token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

// 响应拦截器：处理错误码和 Token 过期
service.interceptors.response.use(
  response => {
    const res = response.data
    // API 规定成功状态码为 200
    if (res.code && res.code !== 200 && res.code !== 401) {
      ElMessage.error(res.message || 'Error')
      return Promise.reject(new Error(res.message || 'Error'))
    }
    return res
  },
  error => {
    // 先检查返回的 JSON 中是否有 code 字段
    if (error.response && error.response.data) {
      const res = error.response.data
      // 如果 code 存在且不是 400，输出对应的 message
      if (res.code && res.code !== 400) {
        ElMessage.error(res.message || 'Error')
        return Promise.reject(error)
      }
      // 如果 code 是 400，按照错误处理（登录过期）
      if (res.code === 400 || error.response.status === 401) {
        ElMessage.error('登录已过期或凭证无效，请重新登录')
        localStorage.removeItem('token')
        router.push('/login')
        return Promise.reject(error)
      }
    }
    
    // 其他错误情况
    ElMessage.error(error.message || '网络或接口错误，请检查 API 地址是否正确')
    return Promise.reject(error)
  }
)

export default service