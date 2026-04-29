import type koa from "koa";
import CommentService from "../../orm/commentService";
import { checkKey, extractToken } from "../../utils/security"

export default async (ctx: koa.Context): Promise<void> => {
  const authHeader = ctx.get("Authorization");
  const key = extractToken(authHeader);

  if (!key || !checkKey(key)) {
    ctx.status = 401;
    ctx.body = { code: 401, message: "Invalid token" };
    return;
  }

  const stats = await CommentService.getStatsOverview();

  ctx.body = {
    code: 200,
    message: "Stats fetched successfully",
    data: stats
  };
};