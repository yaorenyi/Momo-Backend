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
    subject: `你在 ${process.env.SITE_NAME ?? 'Momo Blog'} 上的评论有了新回复`,
    html: `
      <div style="background-color: #f4f7f9; padding: 20px 10px; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;">
  <div style="max-width: 600px; margin: 0 auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.05); border: 1px solid #e1e4e8;">
    
    <div style="padding: 30px;">
      <h2 style="margin-top: 0; color: #333; font-size: 18px;">Hi ${toName}，</h2>
      <p style="color: #555; line-height: 1.6;">
        <strong>${replyAuthor}</strong> 回复了你在 <span style="color: #007acc;">《${postTitle}》</span> 中的评论：
      </p>

      <div style="margin: 20px 0; padding: 12px 16px; border-left: 4px solid #dfe3e8; background-color: #fcfcfc; color: #777; font-style: italic; font-size: 14px;">
        ${parentComment}
      </div>

      <p style="color: #333; font-weight: bold; margin-bottom: 8px;">最新回复：</p>
      
      <div style="margin-bottom: 30px; padding: 16px; border-radius: 6px; background-color: #f0f7ff; border-left: 4px solid #007acc; color: #2c3e50; line-height: 1.6;">
        ${replyContent}
      </div>

      <div style="text-align: center;">
        <a href="${postUrl}" 
           style="display: inline-block; background-color: #007acc; color: #ffffff; padding: 12px 28px; text-decoration: none; border-radius: 6px; font-weight: 600; font-size: 15px; box-shadow: 0 2px 5px rgba(0,122,204,0.2);">
          点击查看回复
        </a>
      </div>
    </div>

    <div style="background-color: #fafbfc; padding: 15px 30px; border-top: 1px solid #eeeeee; text-align: center;">
      <p style="margin: 0; font-size: 12px; color: #999;">
        此邮件由系统自动发送，请勿直接回复。<br>
      </p>
    </div>
  </div>
</div>
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
    subject: `你在 ${process.env.SITE_NAME ?? 'Momo Blog'} 上有新的评论`,
    html: `
      <div style="background-color: #f6f8fa; padding: 40px 20px; min-height: 100%; font-family: 'PingFang SC', 'Microsoft YaHei', Helvetica, Arial, sans-serif;">
  <div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 12px; box-shadow: 0 8px 24px rgba(0,0,0,0.05); overflow: hidden;">
    
    <div style="height: 4px; background: linear-gradient(90deg, #007acc, #00c6ff);"></div>

    <div style="padding: 32px;">
      <h2 style="margin: 0 0 16px 0; color: #1a1a1a; font-size: 20px; line-height: 1.4;">
        有人在你的文章下发表了评论
      </h2>
      
      <p style="color: #555; font-size: 15px; margin-bottom: 24px; line-height: 1.6;">
        <strong style="color: #007acc;">${commentAuthor}</strong> 评论了你的文章 
        <b style="color: #1a1a1a;">《${postTitle}》</b>：
      </p>

      <div style="background-color: #f9fafb; border-radius: 8px; padding: 20px; border: 1px dashed #e1e4e8; margin-bottom: 32px;">
        <div style="color: #444; font-size: 15px; line-height: 1.8; word-break: break-all;">
          ${commentContent}
        </div>
      </div>

      <div style="text-align: center;">
        <a href="${postUrl}" 
           style="display: inline-block; background-color: #007acc; color: #ffffff; padding: 12px 32px; text-decoration: none; border-radius: 50px; font-weight: 500; font-size: 15px; transition: all 0.3s ease;">
          立即前往查看
        </a>
      </div>
    </div>

    <div style="background-color: #fafbfc; padding: 20px; text-align: center; border-top: 1px solid #f0f0f0;">
      <p style="margin: 0; font-size: 13px; color: #99aab5; line-height: 1.5;">
        此邮件由系统自动发送，请勿直接回复。
      </p>
    </div>
  </div>
</div>
    `,
  });

  if (error) {
    LogService.error('邮件发送失败:', error);
    throw new Error('Failed to send email');
  }

  LogService.info('邮件已发送:', data.id);
  return data;
}