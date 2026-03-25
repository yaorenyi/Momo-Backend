package model

// CommentRequest 提交评论请求体
type CommentRequest struct {
	PostSlug  string `json:"post_slug" binding:"required"`
	Author    string `json:"author" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	URL       string `json:"url"`
	Content   string `json:"content" binding:"required"`
	ParentID  *int64 `json:"parent_id"`
	PostURL   string `json:"post_url"`
	PostTitle string `json:"post_title"`
}

// LoginRequest 登录请求体
type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CommentResponse 评论响应体
type CommentResponse struct {
	ID          int64              `json:"id"`
	Author      string             `json:"author"`
	URL         *string            `json:"url"`
	Avatar      string             `json:"avatar"`
	ContentText string             `json:"contentText"`
	ContentHTML string             `json:"contentHtml"`
	PubDate     string             `json:"pubDate"`
	ParentID    *int64             `json:"parentId,omitempty"` // 非嵌套模式显示
	Replies     []*CommentResponse `json:"replies,omitempty"`  // 嵌套模式显示
}

type AdminCommentResponse struct {
	ID          int64   `json:"id"`
	PubDate     string  `json:"pubDate"`
	PostSlug    string  `json:"postSlug"`
	Author      string  `json:"author"`
	Email       string  `json:"email"`
	URL         *string `json:"url"`
	IPAddress   *string `json:"ipAddress"`
	OS          *string `json:"os"`
	Browser     *string `json:"browser"`
	ContentText string  `json:"contentText"`
	ContentHtml string  `json:"contentHtml"`
	Status      string  `json:"status"`
}

// Pagination 分页元数据
type Pagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalPage int64 `json:"totalPage"` // 总页数
}
