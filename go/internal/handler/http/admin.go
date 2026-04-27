package http

import (
	"log"
	"momo-backend-go/internal/config"
	"momo-backend-go/internal/model"
	"momo-backend-go/internal/pkg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Login 优化后的登录逻辑
func (h *CommentHandler) Login(c *gin.Context) {
	ip := utils.GetClientIP(c)

	// 1. 检查 IP 封禁状态
	if utils.Limiter.IsIPBlocked(ip) {
		log.Printf("[WARN] Blocked IP attempted to login: %s", ip)
		c.JSON(http.StatusForbidden, gin.H{
			"code":    400,
			"message": "IP is blocked due to multiple failed"})
		return
	}

	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body"})
		return
	}

	// 2. 验证凭据
	if req.Name != config.GlobalConfig.AdminName || req.Password != config.GlobalConfig.AdminPassword {
		isBlocked := utils.Limiter.RecordAttempt(ip)
		log.Printf("[WARN] Login failed for IP: %s", ip)

		if isBlocked {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    400,
				"message": "IP is blocked due to multiple failed"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    400,
				"message": "Invalid username or password"})
		}
		return
	}

	// 3. 登录成功处理
	utils.Limiter.ResetAttempt(ip)
	log.Printf("[INFO] Login successful for IP: %s", ip)

	// 生成密钥并按文档格式返回
	tempKey := utils.GenerateTempKey(req.Name)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Login successful",
		"token":   tempKey,
	})
}

func (h *CommentHandler) ListAllComments(c *gin.Context) {
	// 1. 解析分页参数和状态筛选
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid query parameters",
		})
		return
	}

	status := c.DefaultQuery("status", "")
	limit := 10
	offset := (page - 1) * limit

	// 2. 从 Repo 调用获取数据
	comments, total, err := h.Repo.List(c.Request.Context(), offset, limit, status)
	if err != nil {
		log.Printf("[ERROR] Failed to list comments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to fetch comments",
		})
		return
	}

	// 3. 计算总页数
	totalPage := int((total + int64(limit) - 1) / int64(limit))
	if total == 0 {
		totalPage = 0
	}

	// 4. 构造响应数据 (处理时间格式及字段映射)
	// 假设 AdminCommentResponse 是为了匹配文档定义的 JSON 标签

	respComments := make([]model.AdminCommentResponse, 0, len(comments))
	for _, comm := range comments {
		respComments = append(respComments, model.AdminCommentResponse{
			ID:          comm.ID,
			PubDate:     time.UnixMilli(comm.PubDate).UTC().Format("2006-01-02T15:04:05.000Z"),
			PostSlug:    comm.PostSlug,
			Author:      comm.Author,
			Email:       comm.Email,
			URL:         comm.URL,
			IPAddress:   comm.IPAddress,
			OS:          comm.OS,
			Browser:     comm.Browser,
			ContentText: comm.ContentText,
			ContentHtml: comm.ContentHTML,
			Status:      comm.Status,
		})
	}

	// 5. 返回符合文档要求的 JSON
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Comments fetched successfully",
		"data": gin.H{
			"comments": respComments,
			"pagination": gin.H{
				"page":      page,
				"limit":     limit,
				"totalPage": totalPage,
			},
		},
	})
}

// UpdateCommentStatus 修改评论状态
func (h *CommentHandler) UpdateCommentStatus(c *gin.Context) {
	idStr := c.Query("id")
	status := c.Query("status")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.Repo.UpdateStatus(c.Request.Context(), id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    400,
			"message": "Invalid request parameters",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Comment status updated",
	})
}

// GetStatsOverview 统计概览
func (h *CommentHandler) GetStatsOverview(c *gin.Context) {
	stats, err := h.Repo.GetStatsOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to fetch stats",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Stats fetched successfully",
		"data":    stats,
	})
}

// GetUserList 用户列表
func (h *CommentHandler) GetUserList(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid query parameters"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit
	users, total, err := h.Repo.GetUserList(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch users"})
		return
	}

	totalPage := int64((total + int64(limit) - 1) / int64(limit))
	if total == 0 {
		totalPage = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Users fetched successfully",
		"data": gin.H{
			"users": users,
			"pagination": gin.H{
				"page":      page,
				"limit":     limit,
				"totalPage": totalPage,
			},
		},
	})
}

// GetUserComments 获取指定用户的评论
func (h *CommentHandler) GetUserComments(c *gin.Context) {
	author := c.Query("author")
	email := c.Query("email")
	if author == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "author and email are required"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid query parameters"})
		return
	}

	limit := 10
	offset := (page - 1) * limit

	comments, total, err := h.Repo.GetUserComments(c.Request.Context(), author, email, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to fetch comments"})
		return
	}

	totalPage := int64((total + int64(limit) - 1) / int64(limit))
	if total == 0 {
		totalPage = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "User comments fetched successfully",
		"data": gin.H{
			"comments": comments,
			"pagination": gin.H{
				"page":      page,
				"limit":     limit,
				"totalPage": totalPage,
			},
		},
	})
}
