interface Comment {
  id: number;
  author: string;
  url?: string;
  avatar?: string;
  contentText: string;
  contentHtml: string;
  pubDate: string;
  parentId: number | null;
}

interface Pagination {
  page: number;
  limit: number;
  totalPage: number;
}

interface CommentsResponse {
  code: number;
  data: {
    comments: Comment[];
    pagination: Pagination;
  }
}

interface NestedComment {
  id: number;
  author: string;
  avatar?: string;
  contentText: string;
  contentHtml: string;
  pubDate: string;
  replies: NestedComment[];
}

interface NestedCommentsResponse {
  code: number;
  message: string;
  data: {
    comments: NestedComment[];
    pagination: Pagination;
  }
}

interface CommentAdmin {
  id: number;
  pubDate: string;
  postSlug: string;
  author: string;
  email: string;
  url?: string;
  ipAddress: string;
  os: string;
  browser: string;
  contentText: string;
  contentHtml: string;
  status: string;
}
interface CommentAdminResponse {
  code: number;
  message: string;
  data: {
    comments: CommentAdmin[];
    pagination: Pagination;
  }
}

export type { CommentsResponse, NestedCommentsResponse, NestedComment, Comment, CommentAdmin, CommentAdminResponse };