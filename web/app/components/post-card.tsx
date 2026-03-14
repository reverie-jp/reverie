import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import {
  MessageCircle,
  Repeat2,
  Heart,
  Ellipsis,
  ShieldBan,
  UserMinus,
  Flag,
  Link2,
  ClipboardCopy,
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
  repostOf?: Post;
  replyTo?: Post;
}

export function formatRelativeTime(date: Date): string {
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

export function PostCard({
  post,
  onReply,
  onRepost,
}: {
  post: Post;
  onReply?: (post: Post) => void;
  onRepost?: (post: Post) => void;
}) {
  const initials = post.author.name.slice(0, 2);
  const [liked, setLiked] = useState(false);
  const navigate = useNavigate();

  return (
    <article
      className="flex gap-3 px-4 py-3 border-b cursor-pointer hover:bg-muted/30 transition-colors"
      onClick={() => navigate(`/posts/${post.id}`)}
    >
      <Link
        to={`/users/${post.author.customId}`}
        className="shrink-0"
        onClick={(e) => e.stopPropagation()}
      >
        <Avatar className="size-10">
          <AvatarImage src={post.author.avatarUrl} alt={post.author.name} />
          <AvatarFallback>{initials}</AvatarFallback>
        </Avatar>
      </Link>
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
        {post.repostOf && (
          <div
            className="mt-2 border rounded-lg p-3 hover:bg-muted/30 transition-colors"
            onClick={(e) => {
              e.stopPropagation();
              navigate(`/posts/${post.repostOf!.id}`);
            }}
          >
            <div className="flex items-center gap-1.5 text-xs">
              <Avatar className="size-4">
                <AvatarImage src={post.repostOf.author.avatarUrl} alt={post.repostOf.author.name} />
                <AvatarFallback className="text-[8px]">{post.repostOf.author.name.slice(0, 1)}</AvatarFallback>
              </Avatar>
              <span className="font-bold truncate">{post.repostOf.author.name}</span>
              <span className="text-muted-foreground truncate">@{post.repostOf.author.customId}</span>
            </div>
            <p className="mt-1 text-sm whitespace-pre-wrap wrap-break-word line-clamp-3">
              {post.repostOf.content}
            </p>
          </div>
        )}
        <div className="flex items-center justify-between mt-2 max-w-xs" onClick={(e) => e.stopPropagation()}>
          <button
            onClick={() => onReply?.(post)}
            className="flex items-center gap-1.5 text-muted-foreground hover:text-blue-400 transition-colors group"
          >
            <MessageCircle className="size-4" />
            {post.replyCount > 0 && (
              <span className="text-xs">{post.replyCount}</span>
            )}
          </button>
          <button
            onClick={() => onRepost?.(post)}
            className="flex items-center gap-1.5 text-muted-foreground hover:text-green-400 transition-colors group"
          >
            <Repeat2 className="size-4" />
            {post.repostCount > 0 && (
              <span className="text-xs">{post.repostCount}</span>
            )}
          </button>
          <button
            onClick={() => setLiked(!liked)}
            className={`flex items-center gap-1.5 transition-colors group ${liked ? "text-pink-400" : "text-muted-foreground hover:text-pink-400"}`}
          >
            <Heart className={`size-4 ${liked ? "fill-current" : ""}`} />
            {(post.likeCount > 0 || liked) && (
              <span className="text-xs">
                {post.likeCount + (liked ? 1 : 0)}
              </span>
            )}
          </button>
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <button className="text-muted-foreground hover:text-foreground transition-colors" />
              }
            >
              <Ellipsis className="size-4" />
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="min-w-40">
              <DropdownMenuItem
                onClick={() =>
                  navigator.clipboard.writeText(
                    `${window.location.origin}/posts/${post.id}`
                  )
                }
              >
                <Link2 className="size-4" />
                URLをコピー
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => navigator.clipboard.writeText(post.content)}
              >
                <ClipboardCopy className="size-4" />
                テキストをコピー
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <UserMinus className="size-4" />
                フォローをやめる
              </DropdownMenuItem>
              <DropdownMenuItem>
                <Flag className="size-4" />
                通報する
              </DropdownMenuItem>
              <DropdownMenuItem className="text-destructive">
                <ShieldBan className="size-4" />
                ブロック
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </article>
  );
}
