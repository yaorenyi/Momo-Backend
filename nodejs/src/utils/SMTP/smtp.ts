import * as nodemailer from 'nodemailer';

// 定义 SMTP 服务器配置接口
export interface SmtpConfig {
  host: string;
  port: number;
  secure: boolean; // true 对应 465 端口, false 对应其他端口 (如 587)
  auth: {
    user: string;
    pass: string; // 注意：如果是 QQ、网易或 Gmail，通常需要使用授权码或应用专用密码，而不是登录密码
  };
}

// 定义邮件内容选项接口
export interface EmailOptions {
  from: string;            // 发件人地址 (如: "我的应用 <noreply@example.com>")
  to: string | string[];   // 收件人地址，可以是单个字符串或数组
  subject: string;         // 邮件主题
  text?: string;           // 纯文本内容
  html?: string;           // HTML 格式内容
  attachments?: nodemailer.SendMailOptions['attachments']; // 附件 (可选)
}

/**
 * 发送 SMTP 邮件的核心函数
 * * @param config SMTP 服务器配置
 * @param options 邮件发送选项
 * @returns 返回一个 Promise，成功时解析为发送结果对象
 */
export async function sendEmail(config: SmtpConfig, options: EmailOptions): Promise<nodemailer.SentMessageInfo> {
  try {
    // 1. 创建 Nodemailer 传输器实例
    const transporter = nodemailer.createTransport({
      host: config.host,
      port: config.port,
      secure: config.secure,
      auth: {
        user: config.auth.user,
        pass: config.auth.pass,
      },
    });

    // 2. 验证连接配置是否正确（可选，但推荐用于提前发现配置错误）
    await transporter.verify();

    // 3. 发送邮件
    const info = await transporter.sendMail({
      from: options.from,
      to: options.to,
      subject: options.subject,
      text: options.text,
      html: options.html,
      attachments: options.attachments,
    });

    console.log(`✅ 邮件发送成功! 消息 ID: ${info.messageId}`);
    return info;

  } catch (error) {
    console.error('❌ 邮件发送失败:', error);
    throw error; // 将错误向上抛出，由调用方决定如何处理
  }
}