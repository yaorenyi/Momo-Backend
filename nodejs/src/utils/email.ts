import { Resend } from 'resend';
import LogService from "../utils/log";

// 单例实例
let resendInstance: Resend | null = null;
let isInitialized = false;

// 懒加载初始化函数
function getResendClient(): Resend | null {
  if (isInitialized) {
    return resendInstance;
  }
  
  const apiKey = process.env.RESEND_API_KEY;
  
  if (!apiKey) {
    LogService.warn('RESEND_API_KEY 未配置，邮件服务不可用');
    isInitialized = true;
    return null;
  }
  
  resendInstance = new Resend(apiKey);
  isInitialized = true;
  LogService.info('Resend 客户端初始化成功');
  return resendInstance;
}

// 服务可用性检查
export function isEmailServiceAvailable(): boolean {
  return getResendClient() !== null;
}

export async function sendCommentReplyNotification({
  toEmail,
  toName,
  postTitle,
  parentComment,
  replyAuthor,
  replyContent,
  postUrl,
}: {
  toEmail: string;
  toName: string;
  postTitle: string;
  parentComment: string;
  replyAuthor: string;
  replyContent: string;
  postUrl: string;
}) {
  const resend = getResendClient();
  
  if (!resend) {
    LogService.warn('邮件发送跳过：服务未初始化');
    return null;
  }

  const { data, error } = await resend.emails.send({
    from: `评论通知 ${process.env.RESEND_FROM_EMAIL}`,
    to: toEmail,
    subject: `你在 blog.motues.top 上的评论有了新回复`,
    html: `
      <p>Hi ${toName}，</p>
      <p>${replyAuthor} 回复了你在 <b>${postTitle}</b> 中的评论：</p>
      <blockquote style="margin: 10px 0; padding-left: 10px; border-left: 3px solid #ccc;">
        ${parentComment}
      </blockquote>
      <p>回复内容：</p>
      <blockquote style="margin: 10px 0; padding-left: 10px; border-left: 3px solid #007acc;">
        ${replyContent}
      </blockquote>
      <p>
        <a href="${postUrl}" 
           style="background: #007acc; color: white; padding: 8px 16px; text-decoration: none; border-radius: 4px;">
          查看回复
        </a>
      </p>
      <hr>
      <p><small>此邮件由系统自动发送，请勿直接回复。</small></p>
    `,
  });

  if (error) {
    LogService.error('邮件发送失败:', error);
    throw new Error('Failed to send email');
  }

  LogService.info('邮件已发送:', data.id);
  return data;
}

export async function sendCommentNotification({
  postTitle,
  postUrl,
  commentAuthor,
  commentContent,
}: {
  postTitle: string;
  postUrl: string;
  commentAuthor: string;
  commentContent: string;
}) {
  const resend = getResendClient();
  
  if (!resend) {
    LogService.warn('邮件发送跳过：服务未初始化');
    return null;
  }
  
  const { data, error } = await resend.emails.send({
    from: `评论通知 ${process.env.RESEND_FROM_EMAIL}`,
    to: process.env.EMAIL_ADDRESS as string,
    subject: `你在 blog.motues.top 上有新的评论`,
    html: `
      <p>${commentAuthor} 评论了你的文章 <b>${postTitle}</b>：</p>
      <p>回复内容：</p>
      <blockquote style="margin: 10px 0; padding-left: 10px; border-left: 3px solid #ccc;">
        ${commentContent}
      </blockquote>
      <p>
        <a href="${postUrl}" 
           style="background: #007acc; color: white; padding: 8px 16px; text-decoration: none; border-radius: 4px;">
          查看评论
        </a>
      </p>
      <hr>
      <p><small>此邮件由系统自动发送，请勿直接回复。</small></p>
    `,
  });

  if (error) {
    LogService.error('邮件发送失败:', error);
    throw new Error('Failed to send email');
  }

  LogService.info('邮件已发送:', data.id);
  return data;
}