import { useEffect, useState } from "react";
import { Button } from "~/components/ui/button";
import { followClient, tokenStore } from "~/lib/api-client";
import { formatUser } from "~/lib/resource-name";
import {
  ConfirmActionDialog,
  type ConfirmActionType,
} from "~/components/confirm-action-dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import {
  Bell,
  ChevronDown,
  ClipboardCopy,
  Flag,
  Link2,
  MessageCircle,
  Repeat2,
  ShieldBan,
  VolumeOff,
} from "lucide-react";

export function FollowButton({
  customId,
  initialFollowing,
  followsYou = false,
  size = "sm",
  onFollowChange,
  onBlockChange,
}: {
  customId: string;
  initialFollowing: boolean;
  /** Whether this user follows you */
  followsYou?: boolean;
  /** "sm" for lists, "md" for profile page */
  size?: "sm" | "md";
  onFollowChange?: (following: boolean) => void;
  onBlockChange?: (blocked: boolean) => void;
}) {
  const [isFollowing, setIsFollowing] = useState(initialFollowing);
  const [isMuted, setIsMuted] = useState(false);
  const [isRepostMuted, setIsRepostMuted] = useState(false);
  const [isBlocked, setIsBlocked] = useState(false);
  const [pending, setPending] = useState(false);
  const [confirmAction, setConfirmAction] = useState<ConfirmActionType | null>(
    null,
  );

  useEffect(() => {
    setIsFollowing(initialFollowing);
  }, [initialFollowing]);

  const updateFollowing = (value: boolean) => {
    setIsFollowing(value);
    onFollowChange?.(value);
  };

  const applyFollow = async (next: boolean) => {
    if (!tokenStore.getAccessToken()) {
      window.location.href = `/login?returnTo=${encodeURIComponent(
        window.location.pathname,
      )}`;
      return;
    }
    const previous = isFollowing;
    updateFollowing(next);
    setPending(true);
    try {
      const name = formatUser(customId);
      if (next) {
        await followClient.followUser({ name });
      } else {
        await followClient.unfollowUser({ name });
      }
    } catch (err) {
      console.error("follow toggle failed:", err);
      updateFollowing(previous);
    } finally {
      setPending(false);
    }
  };

  const updateBlocked = (value: boolean) => {
    setIsBlocked(value);
    onBlockChange?.(value);
  };

  const btnHeight = size === "md" ? "h-9" : "h-8";
  const btnPx = size === "md" ? "px-4" : "px-3 text-xs";
  const chevronPx = size === "md" ? "px-2" : "px-1.5";
  const chevronSize = size === "md" ? "size-3.5" : "size-3";

  if (isBlocked) {
    return (
      <Button
        variant="destructive"
        className={`rounded-full ${btnHeight} ${btnPx}`}
        onClick={(e) => {
          e.preventDefault();
          updateBlocked(false);
        }}
      >
        ブロックを解除
      </Button>
    );
  }

  return (
    <>
      <div className="flex items-center">
        <Button
          variant={isFollowing ? "outline" : "default"}
          disabled={pending}
          className={`rounded-r-none rounded-l-full ${btnHeight} ${btnPx}`}
          onClick={(e) => {
            e.preventDefault();
            if (isFollowing) {
              setConfirmAction("unfollow");
            } else {
              void applyFollow(true);
            }
          }}
        >
          {isFollowing
            ? "フォロー中"
            : followsYou
              ? "フォローを返す"
              : "フォローする"}
        </Button>
        <DropdownMenu>
          <DropdownMenuTrigger
            render={
              <Button
                variant={isFollowing ? "outline" : "default"}
                className={`rounded-l-none rounded-r-full border-l-0 ${chevronPx} ${btnHeight}`}
                onClick={(e) => e.preventDefault()}
              />
            }
          >
            <ChevronDown className={chevronSize} />
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="min-w-48">
            <DropdownMenuItem
              onClick={() => navigator.clipboard.writeText(`@${customId}`)}
            >
              <ClipboardCopy className="size-4" />
              ユーザーIDをコピー
            </DropdownMenuItem>
            <DropdownMenuItem
              onClick={() =>
                navigator.clipboard.writeText(
                  `${window.location.origin}/users/${customId}`,
                )
              }
            >
              <Link2 className="size-4" />
              プロフィールURLをコピー
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <MessageCircle className="size-4" />
              メッセージを送信
            </DropdownMenuItem>
            <DropdownMenuItem>
              <Bell className="size-4" />
              投稿を通知
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={() => {
                if (isMuted) {
                  setIsMuted(false);
                } else {
                  setConfirmAction("mute");
                }
              }}
            >
              <VolumeOff className="size-4" />
              {isMuted ? "ミュートを解除" : "ミュート"}
            </DropdownMenuItem>
            <DropdownMenuItem
              onClick={() => {
                if (isRepostMuted) {
                  setIsRepostMuted(false);
                } else {
                  setConfirmAction("repost-mute");
                }
              }}
            >
              <Repeat2 className="size-4" />
              {isRepostMuted ? "再投稿のミュートを解除" : "再投稿をミュート"}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="text-destructive"
              onClick={() => {
                if (isBlocked) {
                  updateBlocked(false);
                } else {
                  setConfirmAction("block");
                }
              }}
            >
              <ShieldBan className="size-4" />
              {isBlocked ? "ブロックを解除" : "ブロック"}
            </DropdownMenuItem>
            <DropdownMenuItem className="text-destructive">
              <Flag className="size-4" />
              通報する
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <ConfirmActionDialog
        action={confirmAction}
        customId={customId}
        onConfirm={() => {
          if (confirmAction === "unfollow") void applyFollow(false);
          if (confirmAction === "mute") setIsMuted(true);
          if (confirmAction === "repost-mute") setIsRepostMuted(true);
          if (confirmAction === "block") {
            updateBlocked(true);
            void applyFollow(false);
          }
          setConfirmAction(null);
        }}
        onCancel={() => setConfirmAction(null)}
      />
    </>
  );
}
