import { useState } from "react";
import { Dialog, DialogContent } from "~/components/ui/dialog";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { Textarea } from "~/components/ui/textarea";
import { Image, BarChart3, Send } from "lucide-react";

export function ComposePostDialog({
  open,
  onClose,
  onPost,
}: {
  open: boolean;
  onClose: () => void;
  onPost?: (content: string) => void;
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

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent showCloseButton={false} className="max-w-86 md:max-w-md">
        <div className="flex gap-3">
          <Avatar className="size-10 shrink-0">
            <AvatarImage src="" alt="自分" />
            <AvatarFallback>自分</AvatarFallback>
          </Avatar>
          <Textarea
            placeholder="いまどうしてる？"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            className="min-h-24 border-none bg-transparent dark:bg-transparent px-0 focus-visible:border-transparent focus-visible:ring-0 break-all"
          />
        </div>
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
            投稿する
            <Send className="size-3.5" />
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
