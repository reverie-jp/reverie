import { Button } from "~/components/ui/button";
import { BottomNav } from "~/components/bottom-nav";
import { PostCard, type Post } from "~/components/post-card";
import { ArrowLeft } from "lucide-react";

const originalPost: Post = {
  id: "1",
  author: { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
  content: "今日はとてもいい天気ですね。散歩に行ってきました！",
  createdAt: new Date(Date.now() - 3 * 60_000),
  replyCount: 2,
  repostCount: 1,
  likeCount: 5,
};

const repostPosts: Post[] = [
  {
    id: "rp1",
    author: { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
    content: "ほんとにいい天気だった！",
    createdAt: new Date(Date.now() - 30 * 60_000),
    replyCount: 0,
    repostCount: 0,
    likeCount: 3,
    repostOf: originalPost,
  },
];

export default function PostReposts() {
  return (
    <div className="w-full min-h-full flex flex-col">
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-14">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <h1 className="text-base font-bold">再投稿</h1>
        </div>
      </div>

      <div className="flex-1">
        {repostPosts.map((post) => (
          <PostCard key={post.id} post={post} />
        ))}
      </div>

      <BottomNav />
    </div>
  );
}
