import { useState, useEffect } from "react";
import { Link, useSearchParams } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { ConfirmActionDialog } from "~/components/confirm-action-dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { BottomNav } from "~/components/bottom-nav";
import { ArrowLeft } from "lucide-react";
import { listFollowing, listFollowers, followUser, unfollowUser, type User } from "~/lib/api";
import type { Route } from "./+types/user-connections";

function UserListItem({ user: initial }: { user: User }) {
  const [user, setUser] = useState(initial);
  const [showUnfollowConfirm, setShowUnfollowConfirm] = useState(false);

  const handleFollow = async () => {
    try {
      const res = await followUser(user.id);
      setUser(res.user);
    } catch {}
  };

  const handleUnfollow = async () => {
    try {
      const res = await unfollowUser(user.id);
      setUser(res.user);
      setShowUnfollowConfirm(false);
    } catch {}
  };

  return (
    <>
      <div className="flex items-start gap-3 px-4 py-3 border-b">
        <Link to={`/users/${user.customId}`} className="shrink-0">
          <Avatar className="size-10">
            <AvatarFallback>{user.displayName.slice(0, 2)}</AvatarFallback>
          </Avatar>
        </Link>
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between gap-2">
            <Link to={`/users/${user.customId}`} className="min-w-0">
              <div className="flex items-center gap-1.5">
                <p className="text-sm font-bold truncate">{user.displayName}</p>
                {user.isFollowedBy && (
                  <span className="shrink-0 text-[10px] bg-muted text-muted-foreground rounded px-1 py-0.5 leading-none">
                    フォローされています
                  </span>
                )}
              </div>
              <p className="text-xs text-muted-foreground">@{user.customId}</p>
            </Link>
            {!user.isMe && (
              <Button
                variant={user.isFollowing ? "outline" : "default"}
                className="rounded-full h-8 px-3 shrink-0 text-xs"
                onClick={() => {
                  if (user.isFollowing) {
                    setShowUnfollowConfirm(true);
                  } else {
                    handleFollow();
                  }
                }}
              >
                {user.isFollowing ? "フォロー中" : "フォローする"}
              </Button>
            )}
          </div>
          {user.biography && (
            <p className="mt-1 text-sm text-muted-foreground line-clamp-2">
              {user.biography}
            </p>
          )}
        </div>
      </div>

      <ConfirmActionDialog
        action={showUnfollowConfirm ? "unfollow" : null}
        customId={user.customId}
        onConfirm={handleUnfollow}
        onCancel={() => setShowUnfollowConfirm(false)}
      />
    </>
  );
}

export default function UserConnections({ params }: Route.ComponentProps) {
  const [searchParams] = useSearchParams();
  const tab = searchParams.get("tab");
  const defaultTab =
    tab === "followers" ? "followers" : tab === "mutual" ? "mutual" : "following";

  const [following, setFollowing] = useState<User[]>([]);
  const [followers, setFollowers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [displayName, setDisplayName] = useState(params.id);

  useEffect(() => {
    if (!params.id) return;
    Promise.all([
      listFollowing(params.id, { pageSize: 100 }),
      listFollowers(params.id, { pageSize: 100 }),
    ])
      .then(([fwingRes, fwersRes]) => {
        setFollowing(fwingRes.users ?? []);
        setFollowers(fwersRes.users ?? []);
        const first = fwingRes.users?.[0] ?? fwersRes.users?.[0];
        if (first) setDisplayName(first.displayName);
      })
      .finally(() => setLoading(false));
  }, [params.id]);

  return (
    <div className="w-full min-h-full flex flex-col">
      {/* Header */}
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-14">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <div className="min-w-0">
            <h1 className="text-base font-bold truncate">{displayName}</h1>
            <p className="text-xs text-muted-foreground">@{params.id}</p>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <Tabs defaultValue={defaultTab} className="gap-0 flex-1">
        <div className="sticky top-14 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
          <TabsList variant="line" className="w-full h-12">
            <TabsTrigger value="following">フォロー中</TabsTrigger>
            <TabsTrigger value="followers">フォロワー</TabsTrigger>
          </TabsList>
        </div>

        <TabsContent value="following">
          {loading ? (
            <div className="py-12 text-center text-muted-foreground text-sm">読み込み中...</div>
          ) : following.length === 0 ? (
            <div className="py-12 text-center text-muted-foreground text-sm">
              フォロー中のユーザーはいません
            </div>
          ) : (
            following.map((user) => <UserListItem key={user.id} user={user} />)
          )}
        </TabsContent>

        <TabsContent value="followers">
          {loading ? (
            <div className="py-12 text-center text-muted-foreground text-sm">読み込み中...</div>
          ) : followers.length === 0 ? (
            <div className="py-12 text-center text-muted-foreground text-sm">
              フォロワーはいません
            </div>
          ) : (
            followers.map((user) => <UserListItem key={user.id} user={user} />)
          )}
        </TabsContent>
      </Tabs>

      <BottomNav />
    </div>
  );
}
