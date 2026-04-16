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
  listPublicTimeline,
  listFollowingTimeline,
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
    likeCount: p.like_count ?? 0,
  };
}

export function meta({}: Route.MetaArgs) {
  return [{ title: "Reverie" }];
}

function TimelineTab({
  fetcher,
}: {
  fetcher: (pageToken?: string) => Promise<{ posts: ApiPost[]; next_page_token?: string }>;
}) {
  const [posts, setPosts] = useState<Post[]>([]);
  const [nextPageToken, setNextPageToken] = useState<string | undefined>();
  const [loading, setLoading] = useState(false);
  const loggedIn = isLoggedIn();

  const fetch = useCallback(async (pageToken?: string) => {
    if (!loggedIn) return;
    setLoading(true);
    try {
      const res = await fetcher(pageToken);
      const newPosts = (res.posts ?? []).map(apiPostToUiPost);
      setPosts((prev) => (pageToken ? [...prev, ...newPosts] : newPosts));
      setNextPageToken(res.next_page_token || undefined);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  }, [loggedIn]);

  useEffect(() => { fetch(); }, [fetch]);

  if (!loggedIn) {
    return (
      <p className="text-center text-muted-foreground text-sm py-8">
        ログインして投稿を見る
      </p>
    );
  }

  return (
    <div>
      {posts.map((post) => (
        <PostCard key={post.id} post={post} />
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
          onClick={() => fetch(nextPageToken)}
        >
          もっと見る
        </button>
      )}
    </div>
  );
}

export default function Home() {
  const [composeMode, setComposeMode] = useState<ComposeMode | null>(null);
  const loggedIn = isLoggedIn();

  const handlePost = async (content: string) => {
    if (!loggedIn) return;
    try {
      await createPost({ text: content });
      // ページリロードでタイムライン更新（簡易実装）
      window.location.reload();
    } catch {
      // ignore
    }
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
          <TimelineTab
            fetcher={(pageToken) =>
              listFollowingTimeline({ page_size: 20, page_token: pageToken })
            }
          />
        </TabsContent>
        <TabsContent value="public">
          <CallList calls={sampleCalls} tab="public" />
          <TimelineTab
            fetcher={(pageToken) =>
              listPublicTimeline({ page_size: 20, page_token: pageToken })
            }
          />
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
