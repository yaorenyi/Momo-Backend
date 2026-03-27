export type Bindings = {
  MOMO_DB: D1Database
  MOMO_AUTH_KV: KVNamespace;
  ALLOW_ORIGIN: string
  SITE_NAME?: string
  ADMIN_NAME: string
  ADMIN_PASSWORD: string
  ADMIN_EMAIL?: string
  SMTP_HOST?: string
  SMTP_PORT?: string
  EMAIL_USER?: string
  EMAIL_PASSWORD?: string
  EMAIL_SECURE?: string
}