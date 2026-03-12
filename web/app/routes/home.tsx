import { useState } from "react";
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
  {
    id: "c3",
    name: "開発ミーティング",
    type: "video",
    host: "鈴木一郎",
    participants: [
      { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
      { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
      { name: "中村悠", customId: "yu_nkmr", avatarUrl: "" },
      { name: "森田陽介", customId: "yosuke_m", avatarUrl: "" },
      { name: "川口真理", customId: "mari_k", avatarUrl: "" },
    ],
  },
  {
    id: "c4",
    name: "チーム定例",
    type: "audio",
    host: "小林あおい",
    participants: [
      { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
      { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
      { name: "伊藤さくら", customId: "sakura_ito", avatarUrl: "" },
      { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
    ],
  },
  {
    id: "c5",
    name: "全体ミーティング",
    type: "video",
    host: "田中太郎",
    participants: [
      { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
      { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
      { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
      { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
      { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
      { name: "中村悠", customId: "yu_nkmr", avatarUrl: "" },
      { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
      { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
      { name: "伊藤さくら", customId: "sakura_ito", avatarUrl: "" },
      { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
      { name: "松本りな", customId: "rina_m", avatarUrl: "" },
      { name: "井上翔", customId: "sho_inoue", avatarUrl: "" },
      { name: "森田陽介", customId: "yosuke_m", avatarUrl: "" },
      { name: "川口真理", customId: "mari_k", avatarUrl: "" },
      { name: "藤原誠", customId: "makoto_f", avatarUrl: "" },
    ],
  },
  {
    id: "c6",
    name: "作業通話",
    type: "audio",
    host: "松本りな",
    participants: [
      { name: "松本りな", customId: "rina_m", avatarUrl: "" },
      { name: "井上翔", customId: "sho_inoue", avatarUrl: "" },
    ],
  },
  {
    id: "c7",
    name: "朝会",
    type: "video",
    host: "森田陽介",
    participants: [
      { name: "森田陽介", customId: "yosuke_m", avatarUrl: "" },
      { name: "川口真理", customId: "mari_k", avatarUrl: "" },
      { name: "藤原誠", customId: "makoto_f", avatarUrl: "" },
      { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
      { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
      { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
      { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
      { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
      { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
    ],
  },
];

const followingPosts: Post[] = [
  {
    id: "1",
    author: { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
    content: "今日はとてもいい天気ですね。散歩に行ってきました！",
    createdAt: new Date(Date.now() - 3 * 60_000),
    replyCount: 2,
    repostCount: 1,
    likeCount: 5,
  },
  {
    id: "2",
    author: { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
    content:
      "新しいカフェを見つけました。コーヒーがとても美味しかったです。おすすめのメニューはカフェラテです。",
    createdAt: new Date(Date.now() - 2 * 3_600_000),
    replyCount: 0,
    repostCount: 3,
    likeCount: 12,
  },
  {
    id: "3",
    author: { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
    content:
      "React Routerの新しいバージョンを試してみたけど、かなり使いやすくなってる。特にローダーの仕組みが良い。",
    createdAt: new Date(Date.now() - 1 * 86_400_000),
    replyCount: 8,
    repostCount: 15,
    likeCount: 42,
  },
  {
    id: "4",
    author: { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
    content: "週末に映画を観に行きました。ストーリーが素晴らしかった！",
    createdAt: new Date(Date.now() - 3 * 86_400_000),
    replyCount: 1,
    repostCount: 0,
    likeCount: 8,
  },
  {
    id: "5",
    author: { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
    content:
      "プログラミングの勉強を始めて半年。少しずつ書けるようになってきた気がする。",
    createdAt: new Date(Date.now() - 5 * 86_400_000),
    replyCount: 3,
    repostCount: 2,
    likeCount: 20,
  },
  {
    id: "6",
    author: { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
    content:
      "朝のランニングを始めて1ヶ月。体が軽くなった気がする。続けることが大事ですね。",
    createdAt: new Date(Date.now() - 6 * 3_600_000),
    replyCount: 4,
    repostCount: 1,
    likeCount: 18,
  },
  {
    id: "7",
    author: { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
    content:
      "TypeScriptの型パズル、難しいけど楽しい。最近はConditional Typesにハマってます。",
    createdAt: new Date(Date.now() - 2 * 86_400_000),
    replyCount: 6,
    repostCount: 10,
    likeCount: 35,
  },
];

const publicPosts: Post[] = [
  {
    id: "p1",
    author: { name: "中村悠", customId: "yu_nkmr", avatarUrl: "" },
    content: "初めての投稿です！よろしくお願いします。",
    createdAt: new Date(Date.now() - 10 * 60_000),
    replyCount: 1,
    repostCount: 0,
    likeCount: 3,
  },
  {
    id: "p2",
    author: { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
    content:
      "今日の夕焼けが本当にきれいだった。写真では伝わらないくらい。自然の美しさに感動する日々。",
    createdAt: new Date(Date.now() - 45 * 60_000),
    replyCount: 2,
    repostCount: 5,
    likeCount: 28,
  },
  {
    id: "p3",
    author: { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
    content:
      "新しいキーボードを買いました。打鍵感が最高すぎて仕事が捗る。静電容量無接点方式、一度使うと戻れない。",
    createdAt: new Date(Date.now() - 5 * 3_600_000),
    replyCount: 12,
    repostCount: 8,
    likeCount: 56,
  },
  {
    id: "p4",
    author: { name: "伊藤さくら", customId: "sakura_ito", avatarUrl: "" },
    content: "読書記録：今月は5冊読めた。来月はもう少しペースを上げたい。",
    createdAt: new Date(Date.now() - 1 * 86_400_000),
    replyCount: 0,
    repostCount: 2,
    likeCount: 9,
  },
  {
    id: "p5",
    author: { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
    content:
      "Rustでウェブサーバーを書いてみた。所有権の概念、最初は戸惑ったけどコンパイラに怒られながら学ぶのが逆に楽しい。",
    createdAt: new Date(Date.now() - 3 * 86_400_000),
    replyCount: 15,
    repostCount: 20,
    likeCount: 78,
  },
  {
    id: "p6",
    author: { name: "松本りな", customId: "rina_m", avatarUrl: "" },
    content: "引っ越し完了！新しい街を探索するのが楽しみ。",
    createdAt: new Date(Date.now() - 8 * 3_600_000),
    replyCount: 3,
    repostCount: 0,
    likeCount: 15,
  },
  {
    id: "p7",
    author: { name: "井上翔", customId: "sho_inoue", avatarUrl: "" },
    content:
      "デザインシステムを一から構築中。コンポーネントの粒度をどこまで細かくするか、チームで議論が白熱してる。",
    createdAt: new Date(Date.now() - 4 * 86_400_000),
    replyCount: 7,
    repostCount: 12,
    likeCount: 44,
  },
];

export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Welcome to React Router!" },
  ];
}

export default function Home() {
  const [myPosts, setMyPosts] = useState<Post[]>([]);
  const [composeMode, setComposeMode] = useState<ComposeMode | null>(null);

  const handlePost = (content: string) => {
    const newPost: Post = {
      id: `my-${Date.now()}`,
      author: { name: "自分", customId: "me", avatarUrl: "" },
      content,
      createdAt: new Date(),
      replyCount: 0,
      repostCount: 0,
      likeCount: 0,
    };
    setMyPosts((prev) => [newPost, ...prev]);
  };

  const handleReply = (post: Post) => {
    setComposeMode({ type: "reply", post });
  };

  const handleRepost = (post: Post) => {
    setComposeMode({ type: "repost", post });
  };

  return (
    <div className="w-full min-h-full flex flex-col">
      <Tabs defaultValue="following" className="gap-0 flex-1">
        <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
          <TabsList variant="line" className="w-full h-14">
            <TabsTrigger value="following">フォロー中</TabsTrigger>
            <TabsTrigger value="public">オープン</TabsTrigger>
          </TabsList>
        </div>
        <TabsContent value="following">
          <CallList calls={sampleCalls} />
          <div>
            {[...myPosts, ...followingPosts].map((post) => (
              <PostCard
                key={post.id}
                post={post}
                onReply={handleReply}
                onRepost={handleRepost}
              />
            ))}
          </div>
        </TabsContent>
        <TabsContent value="public">
          <CallList calls={sampleCalls} tab="public" />
          <div>
            {[...myPosts, ...publicPosts].map((post) => (
              <PostCard
                key={post.id}
                post={post}
                onReply={handleReply}
                onRepost={handleRepost}
              />
            ))}
          </div>
        </TabsContent>
      </Tabs>
      <ComposeFab onPost={handlePost} />
      <ComposePostDialog
        open={composeMode !== null}
        onClose={() => setComposeMode(null)}
        onPost={handlePost}
        mode={composeMode ?? undefined}
      />
      <BottomNav />
    </div>
  );
}
