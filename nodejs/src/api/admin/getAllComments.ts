import type koa from "koa";
import CommentService from "../../orm/commentService";
import { Comment } from "../../type/prisma";
import { getResponseCommentAdmin } from "../../utils/content";
import { checkAdmin, checkKey, extractToken } from "../../utils/security"
import { getQueryNumber, getQueryBoolean, getQueryString } from "../../utils/url";

export default async (ctx: koa.Context, next: koa.Next): Promise<void> => {

  const page = getQueryNumber(ctx.query.page as string, 1);
  // const key = getQueryString(ctx.query.key as string, "");
  const authHeader = ctx.get("Authorization");
  const key = extractToken(authHeader);

  if(!key || !checkKey(key)) {
    ctx.status = 401;
    ctx.body = { 
      code: 400,
      message: "Invalid key" 
    };
    return;
  }

  // 获取所有评论
  const comments: Comment[] = await CommentService.getAllComments();
  
  const groupedComments = await getResponseCommentAdmin(comments, page);

  // console.log(groupedComments);
  
  ctx.body = groupedComments;
};