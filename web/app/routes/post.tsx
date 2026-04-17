import { useState, useEffect } from "react";
import { Link, useParams } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { BottomNav } from "~/components/bottom-nav";
import { PostCard, formatRelativeTime, type Post } from "~/components/post-card";
import {
  ComposePostDialog,
  type ComposeMode,
  type PostOptions,
} from "~/components/compose-post-dialog";
import { useNavigate } from "react-router";
import { ArrowLeft, Heart, MessageCircle, Repeat2 } from "lucide-react";
import {
  getPost,
  createPost,
  likePost,
  unlikePost,
  listPostReplies,
} from "~/lib/api";
import { useCurrentUser } from "~/lib/use-current-user";
import { apiPostToUiPost } from "~/lib/utils";

function formatFullDate(date: Date): string {
  return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日 ${date.getHours()}:${String(date.getMinutes()).padStart(2, "0")}`;
}

export default function PostDetailRoute() {
  const { id } = useParams();
  return <PostDetail key={id} id={id!} />;
}

function PostDetail({ id }: { id: string }) {
  const navigate = useNavigate();
  const currentUser = useCurrentUser();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [liked, setLiked] = useState(false);
  const [likeCount, setLikeCount] = useState(0);
  const [replies, setReplies] = useState<Post[]>([]);
  const [composeMode, setComposeMode] = useState<ComposeMode | null>(null);

  useEffect(() => {
    setLoading(true);
    Promise.all([
      getPost(id),
      listPostReplies(id, { pageSize: 20 }),
    ])
      .then(([{ post: p }, { posts: rs }]) => {
        const ui = apiPostToUiPost(p);
        setPost(ui);
        setLiked(p.isLiked ?? false);
        setLikeCount(p.likeCount ?? 0);
        setReplies((rs ?? []).map(apiPostToUiPost));
      })
      .catch(() => setPost(null))
      .finally(() => setLoading(false));
  }, [id]);

  const handleLikeToggle = () => {
    if (!post) return;
    if (liked) {
      setLiked(false);
      setLikeCount((c) => Math.max(0, c - 1));
      unlikePost(post.id).catch(() => { setLiked(true); setLikeCount((c) => c + 1); });
    } else {
      setLiked(true);
      setLikeCount((c) => c + 1);
      likePost(post.id).catch(() => { setLiked(false); setLikeCount((c) => Math.max(0, c - 1)); });
    }
  };

  const handlePost = async (content: string, options?: PostOptions) => {
    if (!post) return;
    try {
      const { post: newApiPost } = await createPost({
        text: content,
        replyToId: options?.replyToId ?? post.id,
        repostId: options?.repostId,
      });
      const newUi = apiPostToUiPost(newApiPost);
      setReplies((prev) => [newUi, ...prev]);
      setPost((p) => p ? { ...p, replyCount: p.replyCount + 1 } : p);
    } catch (e: any) {
      alert(`投稿に失敗しました: ${e?.message ?? "不明なエラー"}`);
    }
  };

  const currentUserForUi = currentUser
    ? { name: currentUser.displayName, avatarUrl: undefined }
    : undefined;

  if (loading) {
    return (
      <div className="w-full min-h-full flex flex-col">
        <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
          <div className="flex items-center gap-3 px-4 h-14">
            <Button variant="ghost" size="icon" onClick={() => history.back()}>
              <ArrowLeft className="size-5" />
            </Button>
            <h1 className="text-base font-bold">投稿</h1>
          </div>
        </div>
        <p className="text-center text-muted-foreground text-sm py-8">読み込み中...</p>
        <BottomNav />
      </div>
    );
  }

  if (!post) {
    return (
      <div className="w-full min-h-full flex flex-col">
        <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
          <div className="flex items-center gap-3 px-4 h-14">
            <Button variant="ghost" size="icon" onClick={() => history.back()}>
              <ArrowLeft className="size-5" />
            </Button>
            <h1 className="text-base font-bold">投稿</h1>
          </div>
        </div>
        <p className="text-center text-muted-foreground text-sm py-8">投稿が見つかりません</p>
        <BottomNav />
      </div>
    );
  }

  return (
    <div className="w-full min-h-full flex flex-col">
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-14">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <h1 className="text-base font-bold">投稿</h1>
        </div>
      </div>

      <div className="flex-1">
        {post.replyTo && (
          <div
            className="flex gap-3 px-4 pt-4 cursor-pointer hover:bg-muted/30 transition-colors"
            onClick={() => navigate(`/posts/${post.replyTo!.id}`)}
          >
            <div className="flex flex-col items-center shrink-0">
              <Link
                to={`/users/${post.replyTo.author.customId}`}
                onClick={(e) => e.stopPropagation()}
              >
                <Avatar className="size-10">
                  <AvatarImage
                    src={post.replyTo.author.avatarUrl}
                    alt={post.replyTo.author.name}
                  />
                  <AvatarFallback>
                    {post.replyTo.author.name.slice(0, 2)}
                  </AvatarFallback>
                </Avatar>
              </Link>
              <div className="w-0.5 flex-1 bg-border mt-1" />
            </div>
            <div className="flex-1 min-w-0 pb-3">
              <div className="flex items-center gap-1.5 text-sm">
                <span className="font-bold truncate">
                  {post.replyTo.author.name}
                </span>
                <span className="text-muted-foreground truncate">
                  @{post.replyTo.author.customId}
                </span>
                <span className="text-muted-foreground">·</span>
                <span className="text-muted-foreground whitespace-nowrap text-xs">
                  {formatRelativeTime(post.replyTo.createdAt)}
                </span>
              </div>
              <p className="mt-1 text-sm whitespace-pre-wrap wrap-break-word">
                {post.replyTo.content}
              </p>
              <div className="flex items-center justify-between mt-2 text-xs text-muted-foreground max-w-48">
                <span className="flex items-center gap-1">
                  <MessageCircle className="size-3.5" />
                  {post.replyTo.replyCount}
                </span>
                <span className="flex items-center gap-1">
                  <Repeat2 className="size-3.5" />
                  {post.replyTo.repostCount}
                </span>
                <span className="flex items-center gap-1">
                  <Heart className="size-3.5" />
                  {post.replyTo.likeCount}
                </span>
              </div>
            </div>
          </div>
        )}
        <div className={`px-4 ${post.replyTo ? "pt-1 pb-4" : "py-4"}`}>
          <div className="flex items-center gap-3">
            <Link to={`/users/${post.author.customId}`}>
              <Avatar className="size-10">
                <AvatarImage
                  src={post.author.avatarUrl}
                  alt={post.author.name}
                />
                <AvatarFallback>{post.author.name.slice(0, 2)}</AvatarFallback>
              </Avatar>
            </Link>
            <div>
              <p className="font-bold text-sm">{post.author.name}</p>
              <p className="text-sm text-muted-foreground">
                @{post.author.customId}
              </p>
            </div>
          </div>
          <p className="mt-4 text-base whitespace-pre-wrap wrap-break-word">
            {post.content}
          </p>
          {post.repostOf && (
            <div
              className="mt-3 border rounded-lg p-3 cursor-pointer hover:bg-muted/30 transition-colors"
              onClick={() => navigate(`/posts/${post.repostOf!.id}`)}
            >
              <div className="flex items-center gap-1.5 text-xs">
                <Avatar className="size-5">
                  <AvatarImage
                    src={post.repostOf.author.avatarUrl}
                    alt={post.repostOf.author.name}
                  />
                  <AvatarFallback className="text-[8px]">
                    {post.repostOf.author.name.slice(0, 1)}
                  </AvatarFallback>
                </Avatar>
                <span className="font-bold truncate">
                  {post.repostOf.author.name}
                </span>
                <span className="text-muted-foreground truncate">
                  @{post.repostOf.author.customId}
                </span>
              </div>
              <p className="mt-1 text-sm whitespace-pre-wrap wrap-break-word">
                {post.repostOf.content}
              </p>
            </div>
          )}
          <p className="mt-3 text-xs text-muted-foreground">
            {formatFullDate(post.createdAt)}
          </p>
        </div>

        <div className="flex items-center gap-6 px-4 py-3 border-y text-sm">
          <Link
            to={`/posts/${id}/reposts`}
            className="flex items-center gap-1.5 hover:underline"
          >
            <Repeat2 className="size-4 text-muted-foreground" />
            <span className="font-bold">{post.repostCount}</span>
            <span className="text-muted-foreground">再投稿</span>
          </Link>
          <Link
            to={`/posts/${id}/likes`}
            className="flex items-center gap-1.5 hover:underline"
          >
            <Heart className="size-4 text-muted-foreground" />
            <span className="font-bold">{likeCount}</span>
            <span className="text-muted-foreground">いいね</span>
          </Link>
        </div>

        <div className="flex items-center justify-around py-2 border-b">
          <button
            onClick={() => setComposeMode({ type: "reply", post })}
            className="flex items-center gap-1.5 text-muted-foreground hover:text-blue-400 transition-colors p-2"
          >
            <MessageCircle className="size-5" />
          </button>
          <button
            onClick={() => setComposeMode({ type: "repost", post })}
            className="flex items-center gap-1.5 text-muted-foreground hover:text-green-400 transition-colors p-2"
          >
            <Repeat2 className="size-5" />
          </button>
          <button
            onClick={handleLikeToggle}
            className={`flex items-center gap-1.5 transition-colors p-2 ${liked ? "text-pink-400" : "text-muted-foreground hover:text-pink-400"}`}
          >
            <Heart className={`size-5 ${liked ? "fill-current" : ""}`} />
          </button>
        </div>

        <div>
          {replies.map((reply) => (
            <PostCard
              key={reply.id}
              post={reply}
              currentUserId={currentUser?.id}
              onReply={(p) => setComposeMode({ type: "reply", post: p })}
              onRepost={(p) => setComposeMode({ type: "repost", post: p })}
            />
          ))}
          {replies.length === 0 && (
            <p className="text-center text-muted-foreground text-sm py-8">
              返信はまだありません
            </p>
          )}
        </div>
      </div>

      <ComposePostDialog
        open={composeMode !== null}
        onClose={() => setComposeMode(null)}
        onPost={handlePost}
        mode={composeMode ?? undefined}
        currentUser={currentUserForUi}
      />
      <BottomNav />
    </div>
  );
}
