import { useState, useEffect, useCallback, useRef } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { BottomNav } from "~/components/bottom-nav";
import { PostCard, type Post } from "~/components/post-card";
import { CallList, type Call } from "~/components/call-list";
import { ComposeFab } from "~/components/compose-fab";
import {
  ComposePostDialog,
  type ComposeMode,
  type PostOptions,
} from "~/components/compose-post-dialog";
import type { Route } from "./+types/home";
import {
  listPublicTimeline,
  listFollowingTimeline,
  createPost,
  likePost,
  unlikePost,
  deletePost,
  isLoggedIn,
  type Post as ApiPost,
} from "~/lib/api";
import { useCurrentUser } from "~/lib/use-current-user";
import { apiPostToUiPost } from "~/lib/utils";

const sampleCalls: Call[] = [
  {
    id: "c1",
    name: "雑談部屋",
    type: "audio",
    host: "田中太郎",
    participants: [{ name: "田中太郎", customId: "tanaka", avatarUrl: "" }],
  },
  {
    id: "c2",
    name: "デザインレビュー",
    type: "video",
    host: "佐藤花子",
    participants: [
      { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
      { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
    ],
  },
];



export function meta({}: Route.MetaArgs) {
  return [{ title: "Reverie" }];
}

function TimelineTab({
  fetcher,
  refreshKey,
  currentUserId,
  onReply,
  onRepost,
  onLike,
  onUnlike,
  onDelete,
}: {
  fetcher: (pageToken?: string) => Promise<{ posts: ApiPost[]; nextPageToken?: string }>;
  refreshKey: number;
  currentUserId?: string;
  onReply?: (post: Post) => void;
  onRepost?: (post: Post) => void;
  onLike?: (postId: string) => void;
  onUnlike?: (postId: string) => void;
  onDelete?: (postId: string) => void;
}) {
  const [posts, setPosts] = useState<Post[]>([]);
  const [nextPageToken, setNextPageToken] = useState<string | undefined>();
  const [loading, setLoading] = useState(false);

  // fetcher を ref に保持して load を安定させる（毎 render で新関数が生まれても OK）
  const fetcherRef = useRef(fetcher);
  fetcherRef.current = fetcher;

  const load = useCallback(async (pageToken?: string) => {
    if (typeof window === "undefined" || !isLoggedIn()) return;
    setLoading(true);
    try {
      const res = await fetcherRef.current(pageToken);
      const newPosts = (res.posts ?? []).map(apiPostToUiPost);
      setPosts((prev) => (pageToken ? [...prev, ...newPosts] : newPosts));
      setNextPageToken(res.nextPageToken || undefined);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  }, []); // deps なし → 関数が安定し無限ループしない

  useEffect(() => {
    setPosts([]);
    setNextPageToken(undefined);
    load();
  }, [load, refreshKey]); // refreshKey が変わった時だけ再 fetch

  const handleDelete = useCallback((postId: string) => {
    setPosts((prev) => prev.filter((p) => p.id !== postId));
    onDelete?.(postId);
  }, [onDelete]);

  // SSR 中は何も出さない（hydration mismatch 回避）
  if (typeof window === "undefined") return null;

  if (!isLoggedIn()) {
    return (
      <p className="text-center text-muted-foreground text-sm py-8">
        ログインして投稿を見る
      </p>
    );
  }

  return (
    <div>
      {posts.map((post) => (
        <PostCard
          key={post.id}
          post={post}
          currentUserId={currentUserId}
          onReply={onReply}
          onRepost={onRepost}
          onLike={onLike}
          onUnlike={onUnlike}
          onDelete={handleDelete}
        />
      ))}
      {loading && (
        <p className="text-center text-muted-foreground text-sm py-8">読み込み中...</p>
      )}
      {!loading && posts.length === 0 && (
        <p className="text-center text-muted-foreground text-sm py-8">
          まだ投稿がありません
        </p>
      )}
      {nextPageToken && !loading && (
        <button
          className="w-full py-4 text-sm text-primary"
          onClick={() => load(nextPageToken)}
        >
          もっと見る
        </button>
      )}
    </div>
  );
}

export default function Home() {
  const [composeMode, setComposeMode] = useState<ComposeMode | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);
  // SSR では isLoggedIn() が常に false → useEffect で再評価してハイドレーション後に正しい値にする
  const [loggedIn, setLoggedIn] = useState(false);
  useEffect(() => { setLoggedIn(isLoggedIn()); }, []);
  const currentUser = useCurrentUser();

  const handlePost = async (content: string, options?: PostOptions) => {
    if (!loggedIn) return;
    try {
      await createPost({
        text: content,
        replyToId: options?.replyToId,
        repostId: options?.repostId,
      });
      setRefreshKey((k) => k + 1);
    } catch (e: any) {
      console.error("[createPost] failed:", e);
      alert(`投稿に失敗しました: ${e?.message ?? "不明なエラー"}`);
    }
  };

  const handleLike = useCallback(async (postId: string) => {
    try { await likePost(postId); } catch {}
  }, []);

  const handleUnlike = useCallback(async (postId: string) => {
    try { await unlikePost(postId); } catch {}
  }, []);

  const handleDelete = useCallback(async (postId: string) => {
    try { await deletePost(postId); } catch {}
  }, []);

  const handleReply = useCallback((post: Post) => {
    setComposeMode({ type: "reply", post });
  }, []);

  const handleRepost = useCallback((post: Post) => {
    setComposeMode({ type: "repost", post });
  }, []);

  const currentUserForUi = currentUser
    ? { name: currentUser.displayName, avatarUrl: undefined }
    : undefined;

  return (
    <div className="w-full min-h-full flex flex-col">
      <Tabs defaultValue="public" className="gap-0 flex-1">
        <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
          <TabsList variant="line" className="w-full h-14">
            <TabsTrigger value="following">フォロー中</TabsTrigger>
            <TabsTrigger value="public">オープン</TabsTrigger>
          </TabsList>
        </div>
        <TabsContent value="following">
          <CallList calls={sampleCalls} />
          <TimelineTab
            refreshKey={refreshKey}
            fetcher={(pageToken) =>
              listFollowingTimeline({ pageSize: 20, pageToken: pageToken })
            }
            currentUserId={currentUser?.id}
            onReply={handleReply}
            onRepost={handleRepost}
            onLike={handleLike}
            onUnlike={handleUnlike}
            onDelete={handleDelete}
          />
        </TabsContent>
        <TabsContent value="public">
          <CallList calls={sampleCalls} tab="public" />
          <TimelineTab
            refreshKey={refreshKey}
            fetcher={(pageToken) =>
              listPublicTimeline({ pageSize: 20, pageToken: pageToken })
            }
            currentUserId={currentUser?.id}
            onReply={handleReply}
            onRepost={handleRepost}
            onLike={handleLike}
            onUnlike={handleUnlike}
            onDelete={handleDelete}
          />
        </TabsContent>
      </Tabs>
      {loggedIn && (
        <>
          <ComposeFab onPost={handlePost} currentUser={currentUserForUi} />
          <ComposePostDialog
            open={composeMode !== null}
            onClose={() => setComposeMode(null)}
            onPost={handlePost}
            mode={composeMode ?? undefined}
            currentUser={currentUserForUi}
          />
        </>
      )}
      <BottomNav />
    </div>
  );
}
