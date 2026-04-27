import CommentService   from "../orm/commentService"
import crypto from "crypto";

// 存储临时密钥的Map，包含密钥和过期时间
const tempKeys = new Map<string, { key: string; expiresAt: number }>();

export async function canPostComment(ip: string): Promise<boolean> {
    const lastComment = await CommentService.getlastCommentByIP(ip);
    // 如果没有找到评论或返回为空，则允许发布
    if (!lastComment || lastComment.length === 0) return true;
    
    // 确保有评论后再检查时间
    const comment = lastComment[0];
    if (!comment || !comment.pub_date) return true; // 防止 comment 或 pub_date 不存在的情况
    
    return Date.now() - comment.pub_date.getTime() > 60 * 1000;
}

// 检查内容，删除 XSS 攻击脚本
export function checkContent(content: string): string {
    return content
        // Remove script/style blocks and their content
        .replace(/<(?:script|style)[\s\S]*?<\/(?:script|style)>/gi, '')
        // Remove event handler attributes (onclick, onerror, onload, etc.)
        .replace(/\s+on\w+\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]+)/gi, '')
        // Remove javascript: and vbscript: links in href/src/action
        .replace(/(?:href|src|action|formaction)\s*=\s*"(?:javascript|vbscript):[^"]*"/gi, '')
        .replace(/(?:href|src|action|formaction)\s*=\s*'(?:javascript|vbscript):[^']*'/gi, '')
        // Remove standalone javascript: and vbscript: protocol
        .replace(/^(?:javascript|vbscript):\s*/gi, '')
        // Remove dangerous embedding tags
        .replace(/<\/?(?:iframe|object|embed|frame|meta|link|base|form|input)\b[^>]*>/gi, '');
}

// 生成基于用户名和时间的临时密钥
export function generateTempKey(username: string): string {
    const now = Date.now();
    const data = `${username}-${now}-${process.env.ADMIN_PASSWORD}`;
    const hash = crypto.createHash('sha256').update(data).digest('hex');
    const expiresAt = now + 20 * 60 * 1000; // 5分钟后过期
    
    tempKeys.set(username, { key: hash, expiresAt });
    
    return hash;
}

// 检查密钥是否有效
export function checkKey(key: string): boolean {
    
    // 检查是否是有效的临时密钥
    const now = Date.now();
    const entries = Array.from(tempKeys.entries());
    for (let i = 0; i < entries.length; i++) {
        const [username, tempKey] = entries[i];
        // 检查密钥是否匹配且未过期
        if (tempKey.key === key && tempKey.expiresAt > now) {
            return true;
        }
        // 清除过期密钥
        if (tempKey.expiresAt <= now) {
            tempKeys.delete(username);
        }
    }
    
    return false;
}

/**
 * 从 Authorization header 中提取 token
 * 支持格式：Bearer <token> 或直接返回 token
 */
export function extractToken(authHeader: string): string {
    if (!authHeader) return "";
    
    // 如果是 "Bearer <token>" 格式，提取 token
    if (authHeader.startsWith("Bearer ")) {
        return authHeader.substring(7);
    }
    
    // 否则直接返回 header 值
    return authHeader;
}

export function checkAdmin(name: string, password: string): boolean { 
    return name === process.env.ADMIN_NAME && password === process.env.ADMIN_PASSWORD;
}