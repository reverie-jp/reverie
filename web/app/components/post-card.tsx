import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import {
  MessageCircle,
  Repeat2,
  Heart,
  Ellipsis,
} from "lucide-react";

export interface Post {
  id: string;
  author: {
    name: string;
    customId: string;
    avatarUrl?: string;
  };
  content: string;
  createdAt: Date;
  replyCount: number;
  repostCount: number;
  likeCount: number;
}

function formatRelativeTime(date: Date): string {
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMinutes = Math.floor(diffMs / 60_000);
  const diffHours = Math.floor(diffMs / 3_600_000);
  const diffDays = Math.floor(diffMs / 86_400_000);

  if (diffMinutes < 1) return "たった今";
  if (diffMinutes < 60) return `${diffMinutes}分前`;
  if (diffHours < 24) return `${diffHours}時間前`;
  if (diffDays < 4) return `${diffDays}日前`;

  return `${date.getMonth() + 1}月${date.getDate()}日`;
}

export function PostCard({ post }: { post: Post }) {
  const initials = post.author.name.slice(0, 2);

  return (
    <article className="flex gap-3 px-4 py-3 border-b">
      <Avatar className="size-10 shrink-0">
        <AvatarImage src={post.author.avatarUrl} alt={post.author.name} />
        <AvatarFallback>{initials}</AvatarFallback>
      </Avatar>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-1.5 text-sm">
          <span className="font-bold truncate">{post.author.name}</span>
          <span className="text-muted-foreground truncate">
            @{post.author.customId}
          </span>
          <span className="text-muted-foreground">·</span>
          <span className="text-muted-foreground whitespace-nowrap">
            {formatRelativeTime(post.createdAt)}
          </span>
        </div>
        <p className="mt-1 text-sm whitespace-pre-wrap wrap-break-word">
          {post.content}
        </p>
        <div className="flex items-center justify-between mt-2 max-w-xs">
          <button className="flex items-center gap-1.5 text-muted-foreground hover:text-blue-400 transition-colors group">
            <MessageCircle className="size-4" />
            {post.replyCount > 0 && (
              <span className="text-xs">{post.replyCount}</span>
            )}
          </button>
          <button className="flex items-center gap-1.5 text-muted-foreground hover:text-green-400 transition-colors group">
            <Repeat2 className="size-4" />
            {post.repostCount > 0 && (
              <span className="text-xs">{post.repostCount}</span>
            )}
          </button>
          <button className="flex items-center gap-1.5 text-muted-foreground hover:text-pink-400 transition-colors group">
            <Heart className="size-4" />
            {post.likeCount > 0 && (
              <span className="text-xs">{post.likeCount}</span>
            )}
          </button>
          <button className="text-muted-foreground hover:text-foreground transition-colors">
            <Ellipsis className="size-4" />
          </button>
        </div>
      </div>
    </article>
  );
}
