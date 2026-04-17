import { useState, useEffect } from "react";
import { useParams } from "react-router";
import { Button } from "~/components/ui/button";
import { BottomNav } from "~/components/bottom-nav";
import { PostCard, type Post } from "~/components/post-card";
import {
  ComposePostDialog,
  type ComposeMode,
} from "~/components/compose-post-dialog";
import { ArrowLeft } from "lucide-react";
import { listPostReposts } from "~/lib/api";
import { apiPostToUiPost } from "~/lib/utils";

export default function PostReposts() {
  const { id } = useParams();
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [composeMode, setComposeMode] = useState<ComposeMode | null>(null);

  useEffect(() => {
    if (!id) return;
    listPostReposts(id)
      .then((res) => setPosts(res.posts.map(apiPostToUiPost)))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [id]);

  const handleReply = (post: Post) => {
    setComposeMode({ type: "reply", post });
  };

  const handleRepost = (post: Post) => {
    setComposeMode({ type: "repost", post });
  };

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
        {loading && (
          <div className="flex justify-center items-center py-12 text-muted-foreground text-sm">
            読み込み中...
          </div>
        )}
        {!loading && posts.length === 0 && (
          <div className="flex justify-center items-center py-12 text-muted-foreground text-sm">
            再投稿はまだありません
          </div>
        )}
        {posts.map((post) => (
          <PostCard key={post.id} post={post} onReply={handleReply} onRepost={handleRepost} />
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
