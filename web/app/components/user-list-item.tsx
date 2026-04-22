import { useNavigate } from "react-router";
import type { User } from "~/lib/gen/user/v1/user_pb";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { FollowButton } from "~/components/follow-button";
import { OnlineDot } from "~/components/online-dot";

export function UserListItem({ user }: { user: User }) {
  const navigate = useNavigate();
  return (
    <li
      className="flex items-start gap-3 px-4 py-3 hover:bg-muted/40 cursor-pointer border-b"
      onClick={() => navigate(`/@${user.customId}`)}
    >
      <Avatar className="size-10 shrink-0">
        <AvatarImage src={user.avatarUrl} alt={user.displayName} />
        <AvatarFallback>{user.displayName.slice(0, 2)}</AvatarFallback>
        <OnlineDot status={user.onlineStatus} />
      </Avatar>
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between gap-2">
          <div className="min-w-0">
            <p className="text-sm font-bold truncate">{user.displayName}</p>
            <p className="text-xs text-muted-foreground truncate">
              @{user.customId}
              {user.isFollowedBy && (
                <span className="ml-1.5 inline-block text-[10px] bg-muted text-muted-foreground rounded px-1.5 py-0.5 leading-none align-middle">
                  フォローされています
                </span>
              )}
            </p>
          </div>
          {!user.isMe && !user.isBlockedBy && (
            <div onClick={(e) => e.stopPropagation()}>
              <FollowButton
                customId={user.customId}
                initialFollowing={user.isFollowing}
                followsYou={user.isFollowedBy}
                size="sm"
              />
            </div>
          )}
        </div>
        {user.biography && (
          <p className="text-xs text-foreground/80 mt-1 line-clamp-2">
            {user.biography}
          </p>
        )}
      </div>
    </li>
  );
}
