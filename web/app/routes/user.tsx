import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { FollowButton } from "~/components/follow-button";
import {
  ArrowLeft,
  CalendarDays,
  Link2,
  MapPin,
  ShieldBan,
} from "lucide-react";
import { callClient, userClient } from "~/lib/api-client";
import { formatUser, parseCall } from "~/lib/resource-name";
import { OnlineStatus, type User } from "~/lib/gen/user/v1/user_pb";
import type { Call } from "~/lib/gen/call/v1/call_pb";

export default function UserRoute() {
  const params = useParams<{ handle: string }>();
  const handle = params.handle ?? "";
  // `/:handle` catches any top-level path, so anything that doesn't look like
  // `@<custom_id>` should surface as a proper 404 via the root ErrorBoundary.
  if (!handle.startsWith("@") || handle.length < 2) {
    throw new Response("Not Found", { status: 404 });
  }
  const customId = handle.slice(1);
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const [participatingCall, setParticipatingCall] = useState<Call | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isBlocked, setIsBlocked] = useState(false);
  const [showBlockedProfile, setShowBlockedProfile] = useState(false);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(null);

    const fetchProfile = async () => {
      try {
        const userName = formatUser(customId);
        const [userRes, callRes] = await Promise.all([
          userClient.getUser({ name: userName }),
          callClient
            .getUserParticipatingCall({ name: userName })
            .catch((err) => {
              console.error("GetUserParticipatingCall failed:", err);
              return { call: null } as { call: Call | null };
            }),
        ]);
        if (cancelled) return;
        if (!userRes.user) {
          setError("ユーザーが見つかりません");
          return;
        }
        setUser(userRes.user);
        setParticipatingCall(callRes.call ?? null);
      } catch (err) {
        if (cancelled) return;
        console.error("GetUser failed:", err);
        setError(err instanceof Error ? err.message : String(err));
      } finally {
        if (!cancelled) setLoading(false);
      }
    };

    fetchProfile();
    return () => {
      cancelled = true;
    };
  }, [customId]);

  if (loading) {
    return (
      <div className="w-full min-h-full flex items-center justify-center">
        <p className="text-sm text-muted-foreground">読み込み中...</p>
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="w-full min-h-full flex flex-col items-center justify-center px-6 gap-4">
        <p className="text-sm text-destructive">
          {error ?? "ユーザーを取得できませんでした"}
        </p>
        <Button
          variant="outline"
          size="sm"
          onClick={() => navigate("/login", { replace: true })}
        >
          ログイン画面へ
        </Button>
      </div>
    );
  }

  const joinedAt = user.createTime?.toDate() ?? new Date();

  return (
    <div className="w-full min-h-full flex flex-col">
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-14">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <div className="min-w-0 flex-1">
            <h1 className="text-base font-bold truncate">{user.displayName}</h1>
          </div>
        </div>
      </div>

      <div className="w-full h-32 bg-muted" />

      <div className="px-4 relative">
        <div className="absolute top-3 right-4">
          {!user.isMe && !user.isBlockedBy && (
            <FollowButton
              customId={user.customId}
              initialFollowing={user.isFollowing}
              followsYou={user.isFollowedBy}
              size="md"
              onBlockChange={setIsBlocked}
            />
          )}
        </div>

        <div className="-mt-12 relative w-fit">
          <Avatar className="size-20 ring-4 ring-background">
            {!user.isBlockedBy && (
              <AvatarImage src={user.avatarUrl} alt={user.displayName} />
            )}
            <AvatarFallback className="text-2xl">
              {user.displayName.slice(0, 2)}
            </AvatarFallback>
          </Avatar>
          <span
            className={`absolute bottom-0.5 right-0.5 size-4 rounded-full ring-[3px] ring-background ${
              user.onlineStatus === OnlineStatus.ONLINE
                ? "bg-green-500"
                : "bg-muted-foreground"
            }`}
          />
        </div>

        <div className="mt-3">
          <div className="flex items-center gap-2">
            <p className="text-lg font-bold">{user.displayName}</p>
            {user.isFollowedBy && (
              <span className="text-[10px] bg-muted text-muted-foreground rounded px-1.5 py-0.5 leading-none">
                フォローされています
              </span>
            )}
          </div>
          <p className="text-sm text-muted-foreground">@{user.customId}</p>
        </div>

        {!user.isBlockedBy && participatingCall && (
          <button
            type="button"
            onClick={() =>
              navigate(`/calls/${parseCall(participatingCall.name)}`)
            }
            className="mt-3 w-full flex items-center justify-between rounded-md border px-3 py-2 text-sm hover:bg-muted"
          >
            <span className="flex items-center gap-2">
              <span className="size-2 rounded-full bg-green-500 animate-pulse" />
              通話中
            </span>
            <span className="text-xs text-muted-foreground font-mono truncate max-w-[50%]">
              {parseCall(participatingCall.name)}
            </span>
          </button>
        )}

        {!user.isBlockedBy && (
          <>
            {user.biography && (
              <p className="mt-2 text-sm whitespace-pre-wrap">
                {user.biography}
              </p>
            )}

            <div className="flex flex-wrap gap-x-4 gap-y-1 mt-2 text-sm text-muted-foreground">
              {user.location && (
                <span className="flex items-center gap-1">
                  <MapPin className="size-3.5" />
                  {user.location}
                </span>
              )}
              {user.website && (
                <a
                  href={
                    user.website.startsWith("http")
                      ? user.website
                      : `https://${user.website}`
                  }
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-1 hover:underline"
                >
                  <Link2 className="size-3.5" />
                  <span className="text-primary">{user.website}</span>
                </a>
              )}
              <span className="flex items-center gap-1">
                <CalendarDays className="size-3.5" />
                {joinedAt.getFullYear()}年{joinedAt.getMonth() + 1}月に登録
              </span>
            </div>

            <div className="flex gap-4 mt-3 text-sm">
              <button
                type="button"
                onClick={() => navigate(`/@${user.customId}/following`)}
                className="hover:underline"
              >
                <span className="font-bold">{user.followingCount}</span>{" "}
                <span className="text-muted-foreground">フォロー中</span>
              </button>
              <button
                type="button"
                onClick={() => navigate(`/@${user.customId}/followers`)}
                className="hover:underline"
              >
                <span className="font-bold">{user.followerCount}</span>{" "}
                <span className="text-muted-foreground">フォロワー</span>
              </button>
            </div>
          </>
        )}
      </div>

      {user.isBlockedBy ? (
        <div className="flex-1 flex flex-col items-center justify-center px-6 py-16 text-center">
          <ShieldBan className="size-12 text-muted-foreground/30 mb-4" />
          <p className="text-lg font-bold">あなたをブロックしました</p>
          <p className="text-sm text-muted-foreground mt-2">
            @{user.customId}{" "}
            さんにブロックされているため、プロフィールを閲覧できません。
          </p>
        </div>
      ) : isBlocked && !showBlockedProfile ? (
        <div className="flex-1 flex flex-col items-center justify-center px-6 py-16 text-center">
          <ShieldBan className="size-12 text-muted-foreground/30 mb-4" />
          <p className="text-lg font-bold">
            @{user.customId} さんをブロックしています
          </p>
          <Button
            variant="outline"
            size="sm"
            className="mt-4 rounded-full"
            onClick={() => setShowBlockedProfile(true)}
          >
            プロフィールを表示
          </Button>
        </div>
      ) : null}
    </div>
  );
}
