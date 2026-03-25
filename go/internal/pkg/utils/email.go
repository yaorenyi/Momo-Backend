package utils

import (
	"fmt"
	"momo-backend-go/internal/config"
	"sync"

	"github.com/resend/resend-go/v3"
)

// 定义结构体
type EmailService struct {
	client    *resend.Client
	fromEmail string
	adminMail string
}

var (
	instance *EmailService
	once     sync.Once
)

// GetService 相当于 TS 中的单例/懒加载初始化
func GetService() *EmailService {
	once.Do(func() {
		cfg := config.GlobalConfig
		if cfg == nil {
			fmt.Println("错误：配置未初始化，邮件服务不可用")
			return
		}

		if cfg.ResendAPIKey == "" || cfg.ResendAPIKey == "re_xxxxxx" {
			fmt.Println("警告：ResendAPIKey 未配置，邮件服务不可用")
			return
		}

		client := resend.NewClient(cfg.ResendAPIKey)
		instance = &EmailService{
			client:    client,
			fromEmail: cfg.ResendFromEmail,
			adminMail: cfg.EmailAddress,
		}
		fmt.Println("Resend 客户端初始化成功")
	})
	return instance
}

// IsAvailable 服务可用性检查
func (s *EmailService) IsAvailable() bool {
	return s != nil && s.client != nil
}

// SendCommentReplyNotification 发送评论回复通知
func (s *EmailService) SendCommentReplyNotification(
	toEmail, toName, postTitle, parentComment, replyAuthor, replyContent, postUrl string,
) (*resend.SendEmailResponse, error) {
	if !s.IsAvailable() {
		return nil, fmt.Errorf("邮件服务未初始化")
	}

	htmlContent := fmt.Sprintf(`
      <div style="background-color: #f4f7f9; padding: 20px 10px; font-family: sans-serif;">
        <div style="max-width: 600px; margin: 0 auto; background: #ffffff; border-radius: 8px; border: 1px solid #e1e4e8;">
          <div style="padding: 30px;">
            <h2 style="color: #333; font-size: 18px;">Hi %s，</h2>
            <p style="color: #555;"><strong>%s</strong> 回复了你在 《%s》 中的评论：</p>
            <div style="margin: 20px 0; padding: 12px; border-left: 4px solid #dfe3e8; background: #fcfcfc; color: #777; font-style: italic;">
              %s
            </div>
            <p style="font-weight: bold;">最新回复：</p>
            <div style="padding: 16px; border-radius: 6px; background: #f0f7ff; border-left: 4px solid #007acc;">
              %s
            </div>
            <div style="text-align: center; margin-top: 30px;">
              <a href="%s" style="background: #007acc; color: #fff; padding: 12px 28px; text-decoration: none; border-radius: 6px;">点击查看回复</a>
            </div>
          </div>
        </div>
      </div>`, toName, replyAuthor, postTitle, parentComment, replyContent, postUrl)

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("评论通知 %s", s.fromEmail),
		To:      []string{toEmail},
		Subject: "你在 blog.motues.top 上的评论有了新回复",
		Html:    htmlContent,
	}

	return s.client.Emails.Send(params)
}

// SendCommentNotification 发送给站长的评论通知
func (s *EmailService) SendCommentNotification(
	postTitle, postUrl, commentAuthor, commentContent string,
) (*resend.SendEmailResponse, error) {
	if !s.IsAvailable() {
		return nil, fmt.Errorf("邮件服务未初始化")
	}

	htmlContent := fmt.Sprintf(`
      <div style="background-color: #f6f8fa; padding: 40px 20px; font-family: sans-serif;">
        <div style="max-width: 600px; margin: 0 auto; background: #fff; border-radius: 12px; box-shadow: 0 8px 24px rgba(0,0,0,0.05);">
          <div style="height: 4px; background: linear-gradient(90deg, #007acc, #00c6ff);"></div>
          <div style="padding: 32px;">
            <h2>有人在你的文章下发表了评论</h2>
            <p><strong style="color: #007acc;">%s</strong> 评论了你的文章 <b>《%s》</b>：</p>
            <div style="background: #f9fafb; padding: 20px; border: 1px dashed #e1e4e8; margin-bottom: 32px;">
              %s
            </div>
            <div style="text-align: center;">
              <a href="%s" style="background: #007acc; color: #fff; padding: 12px 32px; text-decoration: none; border-radius: 50px;">立即前往查看</a>
            </div>
          </div>
        </div>
      </div>`, commentAuthor, postTitle, commentContent, postUrl)

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("评论通知 %s", s.fromEmail),
		To:      []string{s.adminMail},
		Subject: "你在 blog.motues.top 上有新的评论",
		Html:    htmlContent,
	}

	return s.client.Emails.Send(params)
}
