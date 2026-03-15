import { useState } from "react";
import { Link } from "react-router";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "~/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogAction,
  AlertDialogCancel,
} from "~/components/ui/alert-dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { Ellipsis, Flag, ShieldBan, Phone, Video } from "lucide-react";
import type { Call } from "~/components/call-list";
import { useCall } from "~/components/call-context";
import { useAnyCallActive } from "~/components/use-any-call-active";

export function JoinCallDialog({
  call,
  onClose,
}: {
  call: Call | null;
  onClose: () => void;
}) {
  const { joinCall } = useCall();
  const { isInCall, currentCallName, endCurrentCall } = useAnyCallActive();
  const [showConfirm, setShowConfirm] = useState(false);
  const TypeIcon = call?.type === "video" ? Video : Phone;

  const handleJoin = () => {
    if (!call) return;
    if (isInCall) {
      setShowConfirm(true);
      return;
    }
    joinCall(call);
    onClose();
  };

  const handleConfirmSwitch = () => {
    if (!call) return;
    endCurrentCall();
    joinCall(call);
    setShowConfirm(false);
    onClose();
  };

  return (
    <>
      <Dialog
        open={call !== null}
        onOpenChange={(open) => {
          if (!open) onClose();
        }}
      >
        <DialogContent showCloseButton={false} className="max-w-72 sm:max-w-72">
          <DialogHeader>
            <div className="flex items-start justify-between">
              <DialogTitle>{call?.name}</DialogTitle>
              <DropdownMenu>
                <DropdownMenuTrigger
                  render={<Button variant="ghost" size="icon-sm" className="-mt-1 -mr-1" />}
                >
                  <Ellipsis className="size-4" />
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem>
                    <Flag className="size-4" />
                    通報する
                  </DropdownMenuItem>
                  <DropdownMenuItem className="text-destructive">
                    <ShieldBan className="size-4" />
                    ブロックする
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </DialogHeader>
          {call && (
            <div className="grid grid-cols-5 gap-2 justify-items-center">
              {call.participants.map((p, i) => (
                <Link
                  key={i}
                  to={`/users/${p.customId}`}
                  onClick={onClose}
                  title={p.name}
                >
                  <Avatar className="size-10 ring-2 ring-background hover:opacity-80 transition-opacity">
                    <AvatarImage src={p.avatarUrl} alt={p.name} />
                    <AvatarFallback>{p.name.slice(0, 2)}</AvatarFallback>
                  </Avatar>
                </Link>
              ))}
            </div>
          )}
          <div className="flex items-center justify-between">
            <DialogDescription>{call?.host}の通話</DialogDescription>
            <span className="text-xs text-muted-foreground whitespace-nowrap">
              {call?.participants.length}人が参加中
            </span>
          </div>
          <DialogFooter>
            <Button
              className="w-full h-10 gap-2"
              onClick={handleJoin}
            >
              <TypeIcon className="size-4" />
              この通話に参加する
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <AlertDialog
        open={showConfirm}
        onOpenChange={(open) => {
          if (!open) setShowConfirm(false);
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>通話を切り替えますか？</AlertDialogTitle>
            <AlertDialogDescription>
              「{currentCallName}」を終了して、「{call?.name}」に参加しますか？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>キャンセル</AlertDialogCancel>
            <AlertDialogAction onClick={handleConfirmSwitch}>
              切り替える
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
