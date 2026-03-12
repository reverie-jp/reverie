import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "~/components/ui/dialog";
import { Button } from "~/components/ui/button";
import type { Call } from "~/components/call-list";

export function JoinCallDialog({
  call,
  onClose,
}: {
  call: Call | null;
  onClose: () => void;
}) {
  return (
    <Dialog
      open={call !== null}
      onOpenChange={(open) => {
        if (!open) onClose();
      }}
    >
      <DialogContent showCloseButton={false}>
        <DialogHeader className="items-center text-center">
          <DialogTitle>この通話に参加しますか？</DialogTitle>
          <DialogDescription>
            「{call?.name}」に
            {call?.type === "video" ? "ビデオ" : "音声"}
            通話で参加します。
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="flex-row">
          <Button
            variant="outline"
            className="flex-1 h-10"
            onClick={onClose}
          >
            キャンセル
          </Button>
          <Button className="flex-1 h-10">参加</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
