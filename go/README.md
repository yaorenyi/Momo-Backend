## Momo Backend Go

Momo Backend Go 是一个基于 Go 语言开发的博客评论系统，使用 Gin + SQLite 实现。

## 部署条件

* 拥有一个支持 Linux/Windows 的服务器
* 不需要任何运行时环境，直接使用编译好的二进制文件

## 快速部署

#### 1. 下载二进制文件

从 [Release](https://github.com/Motues/Momo-Backend/releases/latest) 下载最新的 Go 版本压缩包，根据你的系统选择对应的文件：

* **Linux**: `momo-backend-linux-amd64.tar.gz` 
* **Windows**: `momo-backend-windows-amd64.zip`

以 Linux 为例：

```bash
wget https://github.com/Motues/Momo-Backend/releases/latest/download/momo-backend-linux-amd64.tar.gz
tar -xzf momo-backend-linux-amd64.tar.gz
./momo-backend-go
```

#### 2. 设置环境变量

允许之后会生成一个 `./config/config.yaml` 文件，请根据需要修改。

```bash
vim ./config/config.yaml
# 根据实际情况修改环境变量
```

#### 3. 启动服务

直接运行二进制文件：

```bash
./momo-backend
```

启动成功后，访问 `http://localhost:17171`

**建议使用 systemd 或 pm2 管理进程**

* **使用 systemd (推荐)**

创建服务文件 `/etc/systemd/system/momo-backend.service`：

```ini
[Unit]
Description=Momo Backend Go Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/momo-backend
EnvironmentFile=/path/to/momo-backend/.env
ExecStart=/path/to/momo-backend/momo-backend
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

然后启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable momo-backend
sudo systemctl start momo-backend
sudo systemctl status momo-backend
```

* **使用 pm2**

```bash
npm install -g pm2
pm2 start ./momo-backend --name momo-backend-go
pm2 save
pm2 startup
```

## 环境变量

| 变量名 | 说明 | 
| ------ | ---- | 
| `PORT` | 端口号，默认为 17172 |
| `DATABASE_URL` | 数据库连接地址，默认为 `./data/momo.db` |
| `ALLOW_ORIGIN` | 允许跨域访问的域名，多个域名用逗号分隔 |
| `RESEND_API_KEY` | Resend API Key，用于启用邮箱通知功能；**如不开启，请设置为空** |
| `RESEND_FROM_EMAIL` | Resend 邮件发送通知的邮箱，需要在 Resend 中认证；**如不开启，请设置为空** |
| `EMAIL_ADDRESS` | 管理员邮件接收通知的邮箱；**如不开启，请设置为空** |
| [ADMIN_NAME](file://c:\Data\Desktop\Github\Motues\Momo-Backend\worker\src\api\admin\login.ts#L21-L21) | 管理员登录账号 |
| [ADMIN_PASSWORD](file://c:\Data\Desktop\Github\Motues\Momo-Backend\worker\src\api\admin\login.ts#L22-L22) | 管理员登录密码 |

**注:** [Resend 官网](https://resend.com/)

## 其他

* 日志通过标准输出打印，可通过 `journalctl -u momo-backend -f` (systemd) 或 `pm2 logs` (pm2) 查看
* 超过 5 次错误登录，则锁定 IP 30 分钟
* Token 有效期为 20 分钟
* 数据库文件默认存放在 `./data/momo.db`
* 静态文件（管理页面）已内嵌到二进制文件中，无需单独部署

## Ngnix 配置

如果使用 Ngnix 和 Cloudflare 做反向代理，需要按照下面的配置 ngnix，确保可以获取到正确的 IP 地址

```nginx
server {
    server_name api.example.com; # 这里修改为你的域名

    location / {
        proxy_pass http://localhost:17172; # 这里修改为实际端口号
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header CF-Connecting-IP $http_cf_connecting_ip;
        proxy_set_header True-Client-IP $http_true_client_ip;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## 自行编译

如果你想要自行编译源代码，需要安装 Go 1.21+ 环境：

```bash
cd go
go mod download
go build -o momo-backend main.go
```

编译完成后，将生成的二进制文件移动到合适的位置，并按照上述步骤进行配置。

## 从 Node.js 版本迁移

如果你正在使用 Node.js 版本，可以直接替换为 Go 版本：

1. 备份数据库文件（通常位于 `nodejs/prisma/dev.db` 或你自定义的路径）
2. 下载 Go 版本的二进制文件
3. 配置相同的环境变量
4. 将备份的数据库文件放到 `./data/momo.db` 或通过 `DATABASE_URL` 指定路径
5. 启动 Go 版本服务

Go 版本具有更好的性能、更低的内存占用和更简单的部署流程。