/**
 * 辅助函数：生成 Gravatar 头像地址 (MD5 算法)
 */
export const getCravatar = async (email: string): Promise<string> => {
  const cleanEmail = email.trim().toLowerCase();
  
  // 1. 计算 MD5 Hash (Web Crypto API)
  const msgUint8 = new TextEncoder().encode(cleanEmail);
  const hashBuffer = await crypto.subtle.digest('MD5', msgUint8);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');

  const cravatarUrl = `https://cravatar.cn/avatar/${hashHex}?s=200&d=retro`;

  try {
    // 2. 发起 HEAD 请求，设置较短的超时以保证 Worker 性能
    const response = await fetch(cravatarUrl, { 
      method: 'HEAD',
      // Cloudflare Workers 特有配置，防止由于缓存导致头信息判断不准
      cf: { cacheTtl: 3600 } 
    });

    // 3. 检查 Cravatar 特有的标志位
    const isDefault = response.headers.get('avatar-from') === 'default';

    if (isDefault) {
      // 4. 返回 DiceBear 的风格化头像 (类似 Boring Avatars)
      // 这里的 'shapes' 风格非常接近艺术色块感，'thumbs' 则更有设计感
      return `https://avatar.motues.top/avatar?variant=beam&name=${hashHex}&colors=FFADAD,FFD6A5,FDFFB6,FF9900,AABBCC`;
    }

    return cravatarUrl;
  } catch (e) {
    // 如果网络波动，直接返回原始地址兜底
    return cravatarUrl;
  }
};