import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "~/components/ui/dialog";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";
import { Phone, Video } from "lucide-react";

export function CreateCallDialog({
  open,
  onClose,
}: {
  open: boolean;
  onClose: () => void;
}) {
  const [name, setName] = useState("");
  const [type, setType] = useState<"audio" | "video">("audio");
  const [visibility, setVisibility] = useState("すべてのユーザー");

  const handleCreate = () => {
    // TODO: 通話作成処理
    setName("");
    setType("audio");
    setVisibility("すべてのユーザー");
    onClose();
  };

  const handleOpenChange = (nextOpen: boolean) => {
    if (!nextOpen) {
      onClose();
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent showCloseButton={false} className="max-w-80 sm:max-w-80">
        <DialogHeader>
          <DialogTitle>通話を作成</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <Input
            placeholder="通話の名前"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="h-10"
          />
          <div className="flex gap-2">
            <Button
              variant={type === "audio" ? "default" : "outline"}
              onClick={() => setType("audio")}
              className="flex-1 h-10 gap-2"
            >
              <Phone className="size-4" />
              音声通話
            </Button>
            <Button
              variant={type === "video" ? "default" : "outline"}
              onClick={() => setType("video")}
              className="flex-1 h-10 gap-2"
            >
              <Video className="size-4" />
              ビデオ通話
            </Button>
          </div>
          <div className="space-y-1.5">
            <span className="text-xs text-muted-foreground">公開範囲</span>
            <Select
              value={visibility}
              onValueChange={(v) => v && setVisibility(v)}
            >
              <SelectTrigger className="w-full h-10!">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="すべてのユーザー">
                  すべてのユーザー
                </SelectItem>
                <SelectItem value="フォロワーのみ">フォロワーのみ</SelectItem>
                <SelectItem value="相互フォローのみ">
                  相互フォローのみ
                </SelectItem>
                <SelectItem value="招待した人のみ">招待した人のみ</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
        <DialogFooter>
          <Button
            onClick={handleCreate}
            disabled={name.trim().length === 0}
            className="w-full h-10"
          >
            通話を開始する
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
