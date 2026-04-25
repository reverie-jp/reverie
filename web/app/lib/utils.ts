import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import type { Post as ApiPost } from "./api";
import type { Post as PostCardPost } from "~/components/post-card";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function apiPostToUiPost(apiPost: ApiPost): PostCardPost {
  const postCardPost: PostCardPost = {
    id: apiPost.id,
    author: {
      id: apiPost.author?.id ?? "",
      name: apiPost.author?.displayName ?? "unknown",
      customId: apiPost.author?.customId ?? "",
      avatarUrl: undefined,
      isFollowing: apiPost.author?.isFollowing ?? false,
    },
    content: apiPost.text,
    createdAt: new Date(apiPost.createTime),
    replyCount: apiPost.replyCount ?? 0,
    repostCount: apiPost.repostCount ?? 0,
    likeCount: apiPost.likeCount ?? 0,
    isLiked: apiPost.isLiked,
  };

  if (apiPost.repostOf) {
    postCardPost.repostOf = apiPostToUiPost(apiPost.repostOf);
  }

  // TODO: Handle replyTo if needed

  return postCardPost;
}
