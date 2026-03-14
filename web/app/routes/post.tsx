import { useState } from "react";
import { Link, useParams } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { BottomNav } from "~/components/bottom-nav";
import { PostCard, formatRelativeTime, type Post } from "~/components/post-card";
import {
  ComposePostDialog,
  type ComposeMode,
} from "~/components/compose-post-dialog";
import { useNavigate } from "react-router";
import { ArrowLeft, Heart, MessageCircle, Repeat2 } from "lucide-react";

const allPosts: Post[] = [
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

const parentPost = allPosts[0]; // id: "1"

const sampleReplies: Record<string, Post[]> = {
  "1": [
    {
      id: "r1",
      author: { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
      content: "いいですね！天気が良い日は外に出たくなりますよね。",
      createdAt: new Date(Date.now() - 2 * 60_000),
      replyCount: 0,
      repostCount: 0,
      likeCount: 2,
      replyTo: parentPost,
    },
    {
      id: "r2",
      author: { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
      content: "どこを散歩したんですか？最近いい散歩コースを探してます。",
      createdAt: new Date(Date.now() - 1 * 60_000),
      replyCount: 1,
      repostCount: 0,
      likeCount: 1,
      replyTo: parentPost,
    },
  ],
};

allPosts.push({
  id: "rp1",
  author: { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
  content: "ほんとにいい天気だった！",
  createdAt: new Date(Date.now() - 30 * 60_000),
  replyCount: 0,
  repostCount: 0,
  likeCount: 3,
  repostOf: allPosts[0],
});

const allPostsMap = new Map<string, Post>();
for (const p of allPosts) allPostsMap.set(p.id, p);
for (const replies of Object.values(sampleReplies)) {
  for (const r of replies) allPostsMap.set(r.id, r);
}

function formatFullDate(date: Date): string {
  return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日 ${date.getHours()}:${String(date.getMinutes()).padStart(2, "0")}`;
}

export default function PostDetailRoute() {
  const { id } = useParams();
  return <PostDetail key={id} id={id!} />;
}

function PostDetail({ id }: { id: string }) {
  const navigate = useNavigate();
  const [liked, setLiked] = useState(false);
  const [composeMode, setComposeMode] = useState<ComposeMode | null>(null);
  const post = allPostsMap.get(id) ?? allPosts[0];
  const [replies, setReplies] = useState<Post[]>(sampleReplies[id] ?? []);

  const handlePost = (content: string) => {
    const newPost: Post = {
      id: `reply-${Date.now()}`,
      author: { name: "自分", customId: "me", avatarUrl: "" },
      content,
      createdAt: new Date(),
      replyCount: 0,
      repostCount: 0,
      likeCount: 0,
    };
    setReplies((prev) => [newPost, ...prev]);
  };

  const handleReply = (target: Post) => {
    setComposeMode({ type: "reply", post: target });
  };

  const handleRepost = (target: Post) => {
    setComposeMode({ type: "repost", post: target });
  };

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
            <span className="font-bold">
              {post.likeCount + (liked ? 1 : 0)}
            </span>
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
            onClick={() => setLiked(!liked)}
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
              onReply={handleReply}
              onRepost={handleRepost}
            />
          ))}
        </div>
      </div>

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
