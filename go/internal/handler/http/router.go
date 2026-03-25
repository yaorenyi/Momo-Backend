package http

import (
	"momo-backend-go/internal/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 验证管理 Key
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Authorization Header
		authHeader := c.GetHeader("Authorization")

		// 2. 检查格式是否以 "Bearer " 开头
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization header is missing or format is invalid",
			})
			return
		}

		// 3. 提取 Token 部分（去除 "Bearer " 后的字符串）
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 4. 验证 Token 是否有效
		if token == "" || !utils.IsTokenValid(token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid or expired token",
			})
			return
		}

		// 验证通过，继续后续流程
		c.Next()
	}
}

func RegisterRoutes(r *gin.Engine, h *CommentHandler) {
	// 1. 用户端公开接口
	api := r.Group("/api")
	{
		api.POST("/comments", h.PostComment)
		api.GET("/comments", h.GetComments)
	}

	// 2. 管理员接口
	admin := r.Group("/admin")
	{
		// 登录接口无需校验
		admin.POST("/login", h.Login)

		// 评论管理接口：需 Authorization 校验
		auth := admin.Group("/")
		auth.Use(AuthMiddleware())
		{
			auth.GET("/comments/list", h.ListAllComments)
			auth.PUT("/comments/status", h.UpdateCommentStatus)
		}
	}
}
