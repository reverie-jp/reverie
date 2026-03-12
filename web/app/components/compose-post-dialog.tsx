import { useState } from "react";
import { Dialog, DialogContent } from "~/components/ui/dialog";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { Textarea } from "~/components/ui/textarea";
import { Image, BarChart3, Send } from "lucide-react";
import type { Post } from "~/components/post-card";

export type ComposeMode =
  | { type: "new" }
  | { type: "reply"; post: Post }
  | { type: "repost"; post: Post };

export function ComposePostDialog({
  open,
  onClose,
  onPost,
  mode = { type: "new" },
}: {
  open: boolean;
  onClose: () => void;
  onPost?: (content: string) => void;
  mode?: ComposeMode;
}) {
  const [content, setContent] = useState("");

  const handlePost = () => {
    onPost?.(content);
    setContent("");
    onClose();
  };

  const handleDraft = () => {
    // TODO: 下書き保存処理
    setContent("");
    onClose();
  };

  const handleOpenChange = (nextOpen: boolean) => {
    if (!nextOpen) {
      onClose();
    }
  };

  const placeholder =
    mode.type === "reply"
      ? `@${mode.post.author.customId} への返信`
      : mode.type === "repost"
        ? "コメントを追加..."
        : "いまどうしてる？";

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent showCloseButton={false} className="max-w-86 md:max-w-md">
        {mode.type === "reply" ? (
          <div className="flex gap-3">
            <div className="flex flex-col items-center shrink-0">
              <Avatar className="size-8">
                <AvatarImage
                  src={mode.post.author.avatarUrl}
                  alt={mode.post.author.name}
                />
                <AvatarFallback>
                  {mode.post.author.name.slice(0, 2)}
                </AvatarFallback>
              </Avatar>
              <div className="w-px flex-1 bg-border my-1" />
              <Avatar className="size-8">
                <AvatarImage src="" alt="自分" />
                <AvatarFallback className="text-xs">自分</AvatarFallback>
              </Avatar>
            </div>
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-1.5 text-xs">
                <span className="font-bold truncate">
                  {mode.post.author.name}
                </span>
                <span className="text-muted-foreground truncate">
                  @{mode.post.author.customId}
                </span>
              </div>
              <p className="mt-0.5 text-xs text-muted-foreground line-clamp-3 whitespace-pre-wrap wrap-break-word">
                {mode.post.content}
              </p>
              <Textarea
                placeholder={placeholder}
                value={content}
                onChange={(e) => setContent(e.target.value)}
                className="min-h-24 mt-2 border-none bg-transparent dark:bg-transparent px-0 focus-visible:border-transparent focus-visible:ring-0 break-all"
              />
            </div>
          </div>
        ) : (
          <div className="flex gap-3">
            <Avatar className="size-10 shrink-0">
              <AvatarImage src="" alt="自分" />
              <AvatarFallback>自分</AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <Textarea
                placeholder={placeholder}
                value={content}
                onChange={(e) => setContent(e.target.value)}
                className={`${mode.type === "repost" ? "min-h-12" : "min-h-24"} border-none bg-transparent dark:bg-transparent px-0 focus-visible:border-transparent focus-visible:ring-0 break-all`}
              />
            </div>
          </div>
        )}

        {mode.type === "repost" && (
          <div className="rounded-lg border p-3">
            <div className="flex items-center gap-2">
              <Avatar className="size-5 shrink-0">
                <AvatarImage
                  src={mode.post.author.avatarUrl}
                  alt={mode.post.author.name}
                />
                <AvatarFallback className="text-[8px]">
                  {mode.post.author.name.slice(0, 2)}
                </AvatarFallback>
              </Avatar>
              <span className="text-xs font-bold truncate">
                {mode.post.author.name}
              </span>
              <span className="text-xs text-muted-foreground truncate">
                @{mode.post.author.customId}
              </span>
            </div>
            <p className="mt-1.5 text-xs text-muted-foreground line-clamp-3 whitespace-pre-wrap wrap-break-word">
              {mode.post.content}
            </p>
          </div>
        )}

        <div className="flex items-center border-t pt-3">
          <Button variant="ghost" size="icon" className="text-muted-foreground">
            <Image className="size-4" />
          </Button>
          <Button variant="ghost" size="icon" className="text-muted-foreground">
            <BarChart3 className="size-4" />
          </Button>
          <div className="flex-1" />
          <Button
            variant="link"
            onClick={handleDraft}
            className="text-muted-foreground px-2 mr-1"
          >
            下書き
          </Button>
          <Button
            onClick={handlePost}
            disabled={content.trim().length === 0}
            className="h-8 px-3 gap-1.5"
          >
            {mode.type === "reply"
              ? "返信する"
              : mode.type === "repost"
                ? "再投稿する"
                : "投稿する"}
            <Send className="size-3.5" />
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
