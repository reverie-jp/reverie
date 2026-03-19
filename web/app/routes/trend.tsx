import { useState } from "react";
import { Link, useParams } from "react-router";
import { ArrowLeft } from "lucide-react";
import { BottomNav } from "~/components/bottom-nav";
import { PostCard, type Post } from "~/components/post-card";
import {
  ComposePostDialog,
  type ComposeMode,
} from "~/components/compose-post-dialog";

const sampleTrendPosts: Record<string, Post[]> = {
  "React Router v7": [
    {
      id: "tr1",
      author: { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
      content:
        "React Router v7のローダーがかなり便利になってる。データフェッチの書き方がシンプルになった。",
      createdAt: new Date(Date.now() - 30 * 60_000),
      replyCount: 8,
      repostCount: 15,
      likeCount: 42,
    },
    {
      id: "tr2",
      author: { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
      content:
        "React Router v7、ファイルベースルーティングが最高。Next.jsから移行したけど全然いい。",
      createdAt: new Date(Date.now() - 2 * 3_600_000),
      replyCount: 3,
      repostCount: 7,
      likeCount: 25,
    },
    {
      id: "tr3",
      author: { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
      content:
        "React Router v7のチュートリアルやっと終わった。SSRの設定が思ったより楽だった。",
      createdAt: new Date(Date.now() - 5 * 3_600_000),
      replyCount: 1,
      repostCount: 2,
      likeCount: 14,
    },
    {
      id: "tr4",
      author: { name: "井上翔", customId: "sho_inoue", avatarUrl: "" },
      content:
        "React Router v7でSPAモードとSSRモードを切り替えられるの便利すぎる。プロジェクトの要件に合わせて選べる。",
      createdAt: new Date(Date.now() - 1 * 86_400_000),
      replyCount: 5,
      repostCount: 10,
      likeCount: 38,
    },
    {
      id: "tr5",
      author: { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
      content:
        "React Router v7、型安全なルーティングが嬉しい。パラメータの型が自動で付くのは開発体験が良い。",
      createdAt: new Date(Date.now() - 2 * 86_400_000),
      replyCount: 4,
      repostCount: 6,
      likeCount: 30,
    },
  ],
};

const defaultPosts: Post[] = [
  {
    id: "td1",
    author: { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
    content: "このトピックについて語りたい！みんなはどう思う？",
    createdAt: new Date(Date.now() - 1 * 3_600_000),
    replyCount: 2,
    repostCount: 1,
    likeCount: 8,
  },
  {
    id: "td2",
    author: { name: "中村悠", customId: "yu_nkmr", avatarUrl: "" },
    content: "最近この話題をよく見かける。盛り上がってるね。",
    createdAt: new Date(Date.now() - 3 * 3_600_000),
    replyCount: 0,
    repostCount: 3,
    likeCount: 12,
  },
  {
    id: "td3",
    author: { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
    content: "面白いトレンドだなー。もっと詳しく知りたい。",
    createdAt: new Date(Date.now() - 6 * 3_600_000),
    replyCount: 1,
    repostCount: 0,
    likeCount: 5,
  },
];

export default function TrendPage() {
  const { keyword } = useParams();
  const decodedKeyword = decodeURIComponent(keyword ?? "");
  const posts = sampleTrendPosts[decodedKeyword] ?? defaultPosts;
  const [composeMode, setComposeMode] = useState<ComposeMode | null>(null);

  return (
    <div className="w-full min-h-full flex flex-col">
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center px-4 h-14 gap-3">
          <button
            onClick={() => history.back()}
            className="text-muted-foreground hover:text-foreground transition-colors"
          >
            <ArrowLeft className="size-5" />
          </button>
          <div className="min-w-0">
            <h1 className="font-bold text-sm truncate">{decodedKeyword}</h1>
            <p className="text-xs text-muted-foreground">
              {posts.length}件の投稿
            </p>
          </div>
        </div>
      </div>

      <div className="flex-1">
        {posts.map((post) => (
          <PostCard
            key={post.id}
            post={post}
            onReply={(p) => setComposeMode({ type: "reply", post: p })}
            onRepost={(p) => setComposeMode({ type: "repost", post: p })}
          />
        ))}
      </div>

      <ComposePostDialog
        open={composeMode !== null}
        onClose={() => setComposeMode(null)}
        onPost={() => setComposeMode(null)}
        mode={composeMode ?? undefined}
      />
      <BottomNav />
    </div>
  );
}
