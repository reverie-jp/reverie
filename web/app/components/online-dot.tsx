import { OnlineStatus } from "~/lib/gen/user/v1/user_pb";
import { AvatarBadge } from "~/components/ui/avatar";
import { cn } from "~/lib/utils";

// OnlineDot always renders an AvatarBadge. Color communicates presence:
// green for ONLINE, muted gray (disabled-style) otherwise. Place as a sibling
// of <AvatarImage> / <AvatarFallback>.
export function OnlineDot({
  status,
  className,
}: {
  status: OnlineStatus | undefined;
  className?: string;
}) {
  const isOnline = status === OnlineStatus.ONLINE;
  return (
    <AvatarBadge
      className={cn(
        "text-transparent",
        isOnline
          ? "bg-green-500 dark:bg-green-500"
          : "bg-muted-foreground/40 dark:bg-muted-foreground/30",
        className,
      )}
    />
  );
}
