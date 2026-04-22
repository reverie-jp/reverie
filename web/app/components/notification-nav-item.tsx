import { Bell } from "lucide-react";
import { Link } from "react-router";
import { useNotifications } from "~/lib/notification-context";
import {
  Notification,
  NotificationType,
} from "~/lib/gen/notification/v1/notification_pb";
import { parseCall } from "~/lib/resource-name";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { cn } from "~/lib/utils";

function itemText(n: Notification): string {
  const actor = n.actor?.displayName || n.actor?.customId || "誰か";
  switch (n.type) {
    case NotificationType.USER_FOLLOWED:
      return `${actor} があなたをフォローしました`;
    case NotificationType.FOLLOWING_USER_CALL_STARTED:
      return `${actor} が通話を開始しました`;
    default:
      return "新しい通知";
  }
}

function itemHref(n: Notification): string | null {
  switch (n.type) {
    case NotificationType.USER_FOLLOWED:
      return n.actor ? `/@${n.actor.customId}` : null;
    case NotificationType.FOLLOWING_USER_CALL_STARTED:
      if (n.resourceName.startsWith("calls/")) {
        try {
          return `/calls/${parseCall(n.resourceName)}`;
        } catch {
          return null;
        }
      }
      return null;
    default:
      return null;
  }
}

function formatRelative(seconds: bigint | undefined): string {
  if (!seconds) return "";
  const diffMs = Date.now() - Number(seconds) * 1000;
  const mins = Math.floor(diffMs / 60_000);
  if (mins < 1) return "たった今";
  if (mins < 60) return `${mins}分前`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours}時間前`;
  const days = Math.floor(hours / 24);
  return `${days}日前`;
}

export function NotificationNavItem() {
  const { notifications, unreadCount, markAllRead, markRead } =
    useNotifications();

  const handleOpenChange = (open: boolean) => {
    if (open && unreadCount > 0) {
      void markAllRead();
    }
  };

  return (
    <DropdownMenu onOpenChange={handleOpenChange}>
      <DropdownMenuTrigger
        render={
          <button
            type="button"
            aria-label="通知"
            className="group relative flex flex-col items-center gap-1 px-3 py-1 text-muted-foreground/70 transition-colors hover:text-foreground aria-expanded:text-foreground"
          />
        }
      >
        <span className="relative grid place-items-center w-10 h-7 rounded-xl transition-all group-aria-expanded:bg-(--reverie-accent)/15 group-aria-expanded:text-(--reverie-accent)">
          <Bell className="size-4" />
          {unreadCount > 0 && (
            <span className="absolute top-0.5 right-1 min-w-3.5 h-3.5 px-1 rounded-full bg-destructive text-[9px] leading-3.5 text-center font-medium text-destructive-foreground">
              {unreadCount > 99 ? "99+" : unreadCount}
            </span>
          )}
        </span>
        <span className="text-[10px] leading-none">通知</span>
      </DropdownMenuTrigger>

      <DropdownMenuContent
        side="top"
        align="center"
        sideOffset={8}
        className="w-80 max-h-96 p-0"
      >
        <div className="flex items-center justify-between px-3 py-2 border-b border-foreground/10">
          <span className="text-sm font-medium">通知</span>
          {notifications.some((n) => !n.readTime) && (
            <button
              type="button"
              className="text-xs text-muted-foreground hover:text-foreground"
              onClick={() => void markAllRead()}
            >
              すべて既読
            </button>
          )}
        </div>
        {notifications.length === 0 ? (
          <p className="px-3 py-6 text-xs text-muted-foreground text-center">
            通知はまだありません
          </p>
        ) : (
          <ul className="py-1">
            {notifications.map((n) => {
              const href = itemHref(n);
              const unread = !n.readTime;
              const body = (
                <div
                  className={cn(
                    "flex items-start gap-2 px-3 py-2 text-sm",
                    unread && "bg-foreground/5",
                  )}
                >
                  {unread && (
                    <span className="mt-1.5 size-1.5 rounded-full bg-primary shrink-0" />
                  )}
                  <div className="min-w-0 flex-1">
                    <p className="truncate">{itemText(n)}</p>
                    <p className="text-[10px] text-muted-foreground">
                      {formatRelative(n.createTime?.seconds)}
                    </p>
                  </div>
                </div>
              );
              return (
                <li key={n.name}>
                  {href ? (
                    <Link
                      to={href}
                      onClick={() => {
                        if (unread) void markRead(n.name);
                      }}
                      className="block hover:bg-foreground/5"
                    >
                      {body}
                    </Link>
                  ) : (
                    body
                  )}
                </li>
              );
            })}
          </ul>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
