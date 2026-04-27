import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { Bindings } from './bindings'
import { customCors } from './utils/cors'
import { adminAuth } from './utils/auth'

import { getComments } from './api/public/getComments'
import { postComment } from './api/public/postComment'
import { adminLogin } from './api/admin/login'
import { listComments } from './api/admin/listComments'
import { updateStatus } from './api/admin/updateStatus'
import { statsOverview } from './api/admin/statsOverview'
import { userList } from './api/admin/userList'
import { userComments } from './api/admin/userComments'

const app = new Hono<{ Bindings: Bindings }>()

// 跨域
app.use('/api/*', async (c, next) => {
  const corsMiddleware = customCors(c.env.ALLOW_ORIGIN)
  return corsMiddleware(c, next)
})

// API
app.get('/api/comments', getComments)
app.post('/api/comments', postComment)

app.post('/admin/login', adminLogin)
app.use('/admin/*', adminAuth)
app.get('/admin/comments/list', listComments);
app.put('/admin/comments/status', updateStatus);
app.get('/admin/stats/overview', statsOverview);
app.get('/admin/stats/users', userList);
app.get('/admin/stats/users/comments', userComments);

export default app