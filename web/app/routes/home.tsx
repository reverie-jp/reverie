import { useState, useEffect, useCallback } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { BottomNav } from "~/components/bottom-nav";
import { PostCard, type Post } from "~/components/post-card";
import { CallList, type Call } from "~/components/call-list";
import { ComposeFab } from "~/components/compose-fab";
import {
  ComposePostDialog,
  type ComposeMode,
} from "~/components/compose-post-dialog";
import type { Route } from "./+types/home";
import {
  listTimeline,
  createPost,
  isLoggedIn,
  type Post as ApiPost,
} from "~/lib/api";

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

function apiPostToUiPost(p: ApiPost): Post {
  return {
    id: p.id,
    author: {
      name: p.author?.display_name ?? "unknown",
      customId: p.author?.custom_id ?? "",
      avatarUrl: "",
    },
    content: p.text,
    createdAt: new Date(p.create_time),
    replyCount: p.reply_count ?? 0,
    repostCount: p.repost_count ?? 0,
    likeCount: p.favorite_count ?? 0,
  };
}

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Reverie" },
    { name: "description", content: "Reverie SNS" },
  ];
}

export default function Home() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [nextCursor, setNextCursor] = useState<string | undefined>();
  const [loading, setLoading] = useState(false);
  const [composeMode, setComposeMode] = useState<ComposeMode | null>(null);
  const loggedIn = isLoggedIn();

  const fetchTimeline = useCallback(async (cursor?: string) => {
    if (!loggedIn) return;
    setLoading(true);
    try {
      const res = await listTimeline({ limit: 20, cursor });
      const newPosts = (res.posts ?? []).map(apiPostToUiPost);
      setPosts((prev) => cursor ? [...prev, ...newPosts] : newPosts);
      setNextCursor(res.next_cursor);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  }, [loggedIn]);

  useEffect(() => {
    fetchTimeline();
  }, [fetchTimeline]);

  const handlePost = async (content: string) => {
    if (!loggedIn) return;
    try {
      const res = await createPost({ text: content });
      const newPost = apiPostToUiPost(res.post);
      setPosts((prev) => [newPost, ...prev]);
    } catch {
      // fallback: show locally without API
    }
  };

  const handleReply = (post: Post) => {
    setComposeMode({ type: "reply", post });
  };

  const handleRepost = (post: Post) => {
    setComposeMode({ type: "repost", post });
  };

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
          <div>
            {posts.map((post) => (
              <PostCard
                key={post.id}
                post={post}
                onReply={handleReply}
                onRepost={handleRepost}
              />
            ))}
            {!loggedIn && (
              <p className="text-center text-muted-foreground text-sm py-8">
                ログインして投稿を見る
              </p>
            )}
          </div>
        </TabsContent>
        <TabsContent value="public">
          <CallList calls={sampleCalls} tab="public" />
          <div>
            {posts.map((post) => (
              <PostCard
                key={post.id}
                post={post}
                onReply={handleReply}
                onRepost={handleRepost}
              />
            ))}
            {loading && (
              <p className="text-center text-muted-foreground text-sm py-8">
                読み込み中...
              </p>
            )}
            {!loading && posts.length === 0 && loggedIn && (
              <p className="text-center text-muted-foreground text-sm py-8">
                まだ投稿がありません
              </p>
            )}
            {!loggedIn && (
              <p className="text-center text-muted-foreground text-sm py-8">
                ログインして投稿を見る
              </p>
            )}
            {nextCursor && !loading && (
              <button
                className="w-full py-4 text-sm text-primary"
                onClick={() => fetchTimeline(nextCursor)}
              >
                もっと見る
              </button>
            )}
          </div>
        </TabsContent>
      </Tabs>
      {loggedIn && (
        <>
          <ComposeFab onPost={handlePost} />
          <ComposePostDialog
            open={composeMode !== null}
            onClose={() => setComposeMode(null)}
            onPost={handlePost}
            mode={composeMode ?? undefined}
          />
        </>
      )}
      <BottomNav />
    </div>
  );
}
