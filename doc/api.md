# API 接口定义

## API

| 方法 | 路由 | 功能 |
| --- | --- | --- |
| POST | `/api/comments` | 提交评论 |
| GET | `/api/comments` | 获取评论 |
| GET | `/admin/comments/list` | 获取所有评论 |
| PUT | `/admin/comments/status` | 修改评论状态 |
| POST | `/admin/login` | 登录 |
| GET | `/admin/stats/overview` | 统计概览 |
| GET | `/admin/stats/users` | 用户列表 |
| GET | `/admin/stats/users/comments` | 用户的评论 |

**接口说明**

* 每次请求会返回一个状态码 `code`，请求成果为 200，失败为 400
* 每次登录的时候会返回一个 token，用于后续的 API 请求
* 管理员接口请求头格式：`Authorization: Bearer <token>`
* 错误处理：如果key无效，则返回状态码 `401 Unauthorized`
    ```json
    {
        "code": 401,
        "message": "Invalid token"
    }
    ```

## 用户接口

### 提交评论（POST `/api/comments`）

**请求体**：
```json
{
  "post_slug": "/posts/my-article",
  "author": "张三",
  "email": "zhangsan@example.com",
  "url": "https://example.com",
  "content": "写得真好！",
  "parent_id": null,
  "post_url": "https://blog.example.com/posts/my-article",
  "post_title": "我的文章"
}
```

**响应（成功）**：
```json
{
  "code": 200,
  "message": "Comment submitted successfully"
}
```

**响应（失败）**：
```json
{
  "code": 400,
  "message": "Invalid request body"
}
```
```json
{
  "code": 400,
  "message": "Time limit exceeded"
}
```

### 获取评论（GET `/api/comments`）

> 🔒 仅返回 `status = 'approved'` 的评论

**查询参数**：
- `post_slug`：博客文章唯一标识（必需）
- `page`：查询页数（默认 1）
- `limit`：每页的评论数量（默认 20，最大 50）
- `nested`：评论是否使用嵌套结果返回（默认 true）

**响应（成功）**：
`GET /api/comments?post_slug=...&nested=false`

```json
{
  "code": 200,
  "message": "Comments fetched successfully",
  "data": {
    "comments": [
      {
        "id": 123,
        "author": "张三",
        "url": "https://example.com",
        "avatar": "https://example.com/avatar.png",
        "contentText": "写得真好！",
        "contentHtml": "<p>写得真好！</p>",
        "pubDate": "2025-10-23T10:00:00Z",
        "parentId": null
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "totalPage": 1
    }
  }
}
```

`GET /api/comments?post_slug=...&nested=true`

```json
{
  "code": 200,
  "message": "Comments fetched successfully",
  "data": {
    "comments": [
      {
        "id": 123,
        "author": "张三",
        "url": "https://example.com",
        "avatar": "https://example.com/avatar.png",
        "contentText": "写得真好！",
        "contentHtml": "<p>写得真好！</p>",
        "pubDate": "2025-10-23T10:00:00Z",
        "replies": [
          {
            "id": 124,
            "author": "李四",
            "url": "https://example.com",
            "avatar": "https://example.com/avatar.png",
            "contentText": "同意",
            "contentHtml": "<p>同意</p>",
            "pubDate": "2025-10-23T11:00:00Z",
            "replies": []
          }
        ]
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 1
    }
  }
}
```

**响应（失败）**：

```json
{
  "code": 400,
  "message": "Invalid query parameters"
}
```
---

## 管理员接口

> 🚧 需要 `Authorization: Bearer <token>`

### 登录 (POST `/admin/login`)

**请求体**：
```json
{
  "name": "admin",
  "password": "password"
}
```

**响应（成功）**：

如果登录成功，返回一个key
```json
{
  "code": 200,
  "message": "Login successful",
  "token": "<token>"
}
```

**响应（失败）**：

```json
{
  "code": 400,
  "message": "Invalid username or password"
}
```

```json
{
  "code": 400,
  "message": "IP is blocked due to multiple failed login attempts"
}
```

### 修改评论状态 (PUT `/admin/comments/status`)

**请求参数**：
- `id`：评论ID（必需）
- `status`：评论状态，包括`approved`、`pending`、`deleted`（必需）

**响应（成功）**：
`PUT/admin/comments/status?id=...&status=...`

```json
{
  "code": 200,
  "message": "Comment status updated"
}
```

**响应（失败）**：

```json
{
  "code": 400,
  "message": "Invalid request parameters"
}
```

### 获取所有评论 (GET `/admin/comments/list`)

**查询参数**：
- `page`：查询页数（默认 1）
- `status`：按状态筛选（可选，取值：`approved`、`pending`、`deleted`，为空返回全部）

**响应（成功）**：
`GET /admin/comments/list&page=1`
或过滤：`GET /admin/comments/list?page=1&status=pending`

```json
{
  "code": 200,
  "message": "Comments fetched successfully",
  "data": {
    "comments": [
      {
        "id": 123,
        "pubDate": "2025-10-23T10:00:00Z",
        "postSlug": "/posts/my-article",
        "author": "张三",
        "email": "zhangsan@example.com",
        "url": "https://example.com",
        "ipAddress": "192.168.1.1",
        "os": "Windows 10",
        "browser": "Chrome 96.0.4664.110",
        "contentText": "写得真好！",
        "contentHtml": "<p>写得真好！</p>",
        "status": "approved",
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "totalPage": 1
    }
  }
}
```

**响应（失败）**：

```json
{
  "code": 400,
  "message": "Invalid query parameters"
}
```

### 统计概览 (GET `/admin/stats/overview`)

> 获取整体数据统计，包括评论数、用户数、状态分布、趋势等

**查询参数**：无

**响应（成功）**：
```json
{
  "code": 200,
  "message": "Stats fetched successfully",
  "data": {
    "totalComments": 100,
    "totalUsers": 25,
    "totalPosts": 10,
    "statusDistribution": {
      "approved": 80,
      "pending": 15,
      "deleted": 5
    },
    "recentComments": [
      { "date": "2026-04-21", "count": 5 },
      { "date": "2026-04-22", "count": 3 },
      { "date": "2026-04-23", "count": 8 },
      { "date": "2026-04-24", "count": 2 },
      { "date": "2026-04-25", "count": 7 },
      { "date": "2026-04-26", "count": 4 },
      { "date": "2026-04-27", "count": 6 }
    ],
    "topCommenters": [
      { "author": "张三", "email": "zhangsan@example.com", "count": 15, "lastCommentDate": "2026-04-27T10:00:00.000Z" },
      { "author": "李四", "email": "lisi@example.com", "count": 10, "lastCommentDate": "2026-04-26T08:00:00.000Z" }
    ]
  }
}
```

**响应（失败）**：
```json
{
  "code": 401,
  "message": "Invalid token"
}
```

### 用户列表 (GET `/admin/stats/users`)

> 按用户名+邮箱唯一标识用户，显示每个用户的评论统计

**查询参数**：
- `page`：查询页数（默认 1）
- `limit`：每页用户数（默认 20）

**响应（成功）**：
```json
{
  "code": 200,
  "message": "Users fetched successfully",
  "data": {
    "users": [
      {
        "author": "张三",
        "email": "zhangsan@example.com",
        "commentCount": 15,
        "approvedCount": 12,
        "pendingCount": 2,
        "deletedCount": 1,
        "firstCommentDate": "2024-01-01T00:00:00.000Z",
        "lastCommentDate": "2026-04-27T10:00:00.000Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "totalPage": 2
    }
  }
}
```

### 用户的评论列表 (GET `/admin/stats/users/comments`)

> 获取指定用户的所有评论详情

**查询参数**：
- `author`：作者昵称（必需）
- `email`：邮箱（必需）
- `page`：查询页数（默认 1）

**响应（成功）**：
```json
{
  "code": 200,
  "message": "User comments fetched successfully",
  "data": {
    "comments": [
      {
        "id": 123,
        "pubDate": "2025-10-23T10:00:00Z",
        "postSlug": "/posts/my-article",
        "author": "张三",
        "email": "zhangsan@example.com",
        "url": null,
        "ipAddress": "192.168.1.1",
        "os": "Windows 10",
        "browser": "Chrome 96",
        "contentText": "写得真好！",
        "contentHtml": "<p>写得真好！</p>",
        "status": "approved"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "totalPage": 2
    }
  }
}
```

**响应（失败）**：
```json
{
  "code": 400,
  "message": "author and email are required"
}
```

