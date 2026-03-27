## Momo Backend Go

Momo Backend Go 是一个基于 Go 语言开发的博客评论系统，使用 Gin + SQLite 实现。

## 部署条件

* 拥有一个支持 Linux/Windows 的服务器
* 不需要任何运行时环境，直接使用编译好的二进制文件

## 快速部署

#### 1. 下载二进制文件

从 [Release](https://github.com/Motues/Momo-Backend/releases/latest) 下载最新的 Go 版本压缩包，根据你的系统选择对应的文件：

* **Linux**: `go-linux-amd64.tar.gz` 
* **Windows**: `go-windows-amd64.zip`

以 Linux 为例，可以使用自带的脚本进行部署：

```bash
wget https://github.com/Motues/Momo-Backend/releases/latest/download/go-linux-amd64.tar.gz
tar -xzf mgo-linux-amd64.tar.gz
chmod +x momo.bash
./momo.bash
# 然后根据提示进行操作
```

#### 2. 设置环境变量

运行之后会生成一个 `./config/config.yaml` 文件，可以参考[这里](#环境变量)据需要修改，修改后需要重启服务。

```bash
vim ./config/config.yaml
# 根据实际情况修改环境变量
./momo.bash
```

启动成功后，访问 `http://localhost:17171`

## 环境变量

```yaml
# ./config/config.yaml
# Development Settings
PORT: 17171                                      # server port
# allow origin for CORS, separated by comma
ALLOW_ORIGIN: https://blog.example.com,http://localhost:4321

# Site Settings
SITE_NAME: Momo Blog                             # site name
ADMIN_NAME: admin                                # admin name
ADMIN_PASSWORD: password                         # admin password
ADMIN_EMAIL: admin@example.com                   # admin email, used for receiving notifications

# Email Settings
SMTP_HOST: smtp.qq.com                           # email host, set to empty string to disable email notifications
SMTP_PORT: 465                                   # email port
EMAIL_USER: example@qq.com                       # email user
EMAIL_PASSWORD: xxxxxxxx                         # email password
EMAIL_SECURE: true                               # whether to use secure connection (SSL/TLS)
```

## 其他

* 日志存放在 `momo.log` 中，可以使用 `momo.bash` 脚本查看
* 超过 5 次错误登录，则锁定 IP 30 分钟
* Token 有效期为 20 分钟
* 数据库文件默认存放在 `./data/momo.db`
* 静态文件存放在 `./public` 中，如果只需要更新后台管理页面，则可以直接替换 `./public` 中的文件

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

如果你想要自行编译源代码，需要安装 Go 1.25 环境：

```bash
cd go
go mod download
go build -o ./build/momo-backend main.go
```

编译完成后，将生成的二进制文件移动到合适的位置，并按照上述步骤进行配置。

## 从 Node.js 版本迁移

如果你正在使用 Node.js 版本，可以直接替换为 Go 版本：

1. 备份数据库文件（通常位于 `nodejs/prisma/dev.db` 或你自定义的路径）
2. 下载 Go 版本的二进制文件
3. 配置相同的环境变量
4. 将备份的数据库文件放到 `./data/data.db` 
5. 启动 Go 版本服务
