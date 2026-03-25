## Momo Backend Go

Momo Backend Go 是一个基于 Go 语言开发的博客评论系统，使用 Gin +  实现。


## 技术栈

* HPPT服务器：Gin
* 数据库：SQLite


```bash
├── main.go                  # 程序入口，初始化并启动服务
├── internal/
│   ├── model/
│   │   └── comment.go       # 数据库模型与 API 结构体定义
│   ├── repository/
│   │   ├── sqlite/
│   │   │   └── comment.go   # SQLite 数据库的具体实现 (CRUD)
│   │   └── repository.go    # 存储层接口定义
│   ├── service/
│   │   └── comment.go       # 核心业务逻辑（如评论审核、树状结构递归）
│   ├── handler/
│   │   ├── http/
│   │   │   ├── admin.go     # 管理员接口
│   │   │   ├── comment.go   # 评论接口
│   │   │   └── handler.go   # 路由处理
│   │   └── middleware/
│   └── pkg/
│        └── utils/          # 工具库（如 Logger, Gravatar 生成、HTML 转义）
├── config/
│   └── config.go            # 配置管理
└── go.mod
```

# 安装 Web 框架
go get -u github.com/gin-gonic/gin

# 安装 数据库 工具
go get github.com/jmoiron/sqlx
go get modernc.org/sqlite