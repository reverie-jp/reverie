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

export type ConfirmActionType =
  | "unfollow"
  | "mute"
  | "repost-mute"
  | "block";

const actionConfig: Record<
  ConfirmActionType,
  {
    title: (id: string) => string;
    description: string;
    action: string;
    destructive?: boolean;
  }
> = {
  unfollow: {
    title: (id) => `@${id}のフォローを解除しますか？`,
    description:
      "フォローを解除すると、このユーザーの投稿がタイムラインに表示されなくなります。",
    action: "フォロー解除",
  },
  mute: {
    title: (id) => `@${id}をミュートしますか？`,
    description:
      "ミュートすると、このユーザーの投稿や通知が非表示になります。相手には通知されません。",
    action: "ミュート",
  },
  "repost-mute": {
    title: (id) => `@${id}の再投稿をミュートしますか？`,
    description:
      "このユーザーの再投稿がタイムラインに表示されなくなります。相手には通知されません。",
    action: "ミュート",
  },
  block: {
    title: (id) => `@${id}をブロックしますか？`,
    description:
      "ブロックすると、このユーザーはあなたのプロフィールや投稿を閲覧できなくなります。",
    action: "ブロック",
    destructive: true,
  },
};

export function ConfirmActionDialog({
  action,
  customId,
  onConfirm,
  onCancel,
}: {
  action: ConfirmActionType | null;
  customId: string;
  onConfirm: () => void;
  onCancel: () => void;
}) {
  const config = action ? actionConfig[action] : null;

  return (
    <AlertDialog
      open={action !== null}
      onOpenChange={(open) => {
        if (!open) onCancel();
      }}
    >
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{config?.title(customId)}</AlertDialogTitle>
          <AlertDialogDescription>
            {config?.description}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>キャンセル</AlertDialogCancel>
          <AlertDialogAction
            className={
              config?.destructive
                ? "bg-destructive text-destructive-foreground hover:bg-destructive/90"
                : ""
            }
            onClick={onConfirm}
          >
            {config?.action}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
