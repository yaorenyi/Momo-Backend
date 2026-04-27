package main

import (
	"fmt"
	"log"
	"net/http"

	// _ "net/http/pprof" // 隐式初始化 pprof 路由
	"os"
	"path/filepath"
	"strings"
	"time"

	"momo-backend-go/internal/config"
	h "momo-backend-go/internal/handler/http"
	"momo-backend-go/internal/pkg/utils"
	"momo-backend-go/internal/repository/sqlite"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const Version = "1.2.0"

func main() {

	// go func() {
	// 	// 启动一个独立的端口供 pprof 访问
	// 	http.ListenAndServe("0.0.0.0:6060", nil)
	// }()

	// 处理命令行参数
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("momo-backend version %s\n", Version)
			os.Exit(0)
		default:
			fmt.Printf("未知的参数: %s\n", os.Args[1])
			fmt.Printf("使用 --version 或 -v 查看版本信息\n")
			os.Exit(1)
		}
	}

	gin.SetMode(gin.ReleaseMode)

	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	// 2. 准备数据库环境
	dbPath := "./data/data.db"
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("无法创建数据库目录: %v", err)
	}

	// 3. 连接数据库并初始化表结构
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// SQLite 连接池优化：SQLite 仅支持单写入者，限制连接数避免浪费
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	// SQLite 运行时优化 PRAGMA
	db.MustExec("PRAGMA journal_mode = WAL")
	db.MustExec("PRAGMA synchronous = NORMAL")
	db.MustExec("PRAGMA cache_size = -8000")    // 8MB 缓存
	db.MustExec("PRAGMA temp_store = MEMORY")   // 临时表存内存
	db.MustExec("PRAGMA mmap_size = 268435456") // 256MB 内存映射
	db.MustExec("PRAGMA busy_timeout = 5000")   // 锁等待 5s

	if err := sqlite.InitSchema(db); err != nil {
		log.Fatalf("初始化表结构失败: %v", err)
	}

	// 4. 启动后台清理任务：每 10 分钟清理过期 token 和登录限流记录
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			utils.CleanupExpired()
		}
	}()

	// 5. 初始化
	repo := sqlite.NewCommentRepository(db)
	handler := &h.CommentHandler{Repo: repo}

	// 6. 设置 Gin 引擎（使用 gin.New() 避免 Logger 中间件额外分配）
	r := gin.New()
	r.Use(gin.Recovery())

	// 全局中间件：跨域处理（预解析 allowed origins 避免每次请求重复分配）
	allowedOrigins := strings.Split(cfg.AllowOrigin, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		isAllowed := false
		for _, o := range allowedOrigins {
			if o == origin {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// 7. 注册路由
	h.RegisterRoutes(r, handler)

	// 只要访问 /admin 及其子路径，都尝试返回前端页面
	// 假设你的管理后台打包后放在 ./dist 目录下
	r.Static("/assets", "./public/assets") // 如果打包后的资源路径有前缀

	// 如果你的前端管理页面是完全独立的单页应用
	r.NoRoute(func(c *gin.Context) {
		// 只有非 API 请求才返回 index.html
		path := c.Request.URL.Path
		if !strings.HasPrefix(path, "/api") && !strings.HasPrefix(path, "/admin") {
			c.File("./public/index.html")
		}
	})

	// 8. 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("--- 评论系统后端已启动 ---\n")
	fmt.Printf("监听地址: %s\n", addr)
	fmt.Printf("数据库路径: %s\n", dbPath)
	fmt.Printf("管理员名称: %s\n", cfg.AdminName)
	fmt.Printf("版本: %s\n", Version)

	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
