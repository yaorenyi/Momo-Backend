import { prisma } from "../orm/prisma";

export const DEFAULT_ADMIN_NAME = "momo";
export const DEFAULT_ADMIN_PASSWORD = "momo";

export async function getSetting(key: string): Promise<string | null> {
  try {
    const setting = await prisma.setting.findUnique({ where: { key } });
    return setting?.value ?? null;
  } catch {
    return null;
  }
}

export async function setSetting(key: string, value: string): Promise<void> {
  await prisma.setting.upsert({
    where: { key },
    update: { value, updated_at: new Date() },
    create: { key, value },
  });
}

export async function getAllSettings(): Promise<Record<string, string>> {
  try {
    const settings = await prisma.setting.findMany();
    const map: Record<string, string> = {};
    settings.forEach((s) => { map[s.key] = s.value; });
    return map;
  } catch {
    return {};
  }
}

export async function isDefaultAdmin(): Promise<boolean> {
  const changed = await getSetting("password_changed");
  return changed !== "true";
}

export async function checkAdminCredentials(name: string, password: string): Promise<boolean> {
  const dbName = await getSetting("admin_name");
  const dbPass = await getSetting("admin_password");

  if (dbName && dbPass) {
    return name === dbName && password === dbPass;
  }

  return name === DEFAULT_ADMIN_NAME && password === DEFAULT_ADMIN_PASSWORD;
}

export async function changeAdminPassword(name: string, password: string): Promise<void> {
  await setSetting("admin_name", name);
  await setSetting("admin_password", password);
  await setSetting("password_changed", "true");
}
