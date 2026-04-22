import { useCallback, useEffect, useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router";
import { ArrowLeft } from "lucide-react";
import { followClient, userClient } from "~/lib/api-client";
import { formatUser } from "~/lib/resource-name";
import type { User } from "~/lib/gen/user/v1/user_pb";
import { Button } from "~/components/ui/button";
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "~/components/ui/tabs";
import { UserListItem } from "~/components/user-list-item";

type Tab = "following" | "followers";

export default function UserConnectionsRoute() {
  const params = useParams<{ handle: string }>();
  const handle = params.handle ?? "";
  if (!handle.startsWith("@") || handle.length < 2) {
    throw new Response("Not Found", { status: 404 });
  }
  const customId = handle.slice(1);
  const location = useLocation();
  const navigate = useNavigate();

  const currentTab: Tab = location.pathname.endsWith("/followers")
    ? "followers"
    : "following";

  const [profile, setProfile] = useState<User | null>(null);
  const [following, setFollowing] = useState<User[] | null>(null);
  const [followers, setFollowers] = useState<User[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const res = await userClient.getUser({ name: formatUser(customId) });
        if (!cancelled) setProfile(res.user ?? null);
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : String(err));
        }
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [customId]);

  const fetchList = useCallback(
    async (which: Tab) => {
      setLoading(true);
      setError(null);
      try {
        const name = formatUser(customId);
        if (which === "following") {
          const res = await followClient.listFollowingUsers({ name });
          setFollowing(res.users);
        } else {
          const res = await followClient.listUserFollowers({ name });
          setFollowers(res.users);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : String(err));
      } finally {
        setLoading(false);
      }
    },
    [customId],
  );

  useEffect(() => {
    if (currentTab === "following" && following === null) {
      void fetchList("following");
    }
    if (currentTab === "followers" && followers === null) {
      void fetchList("followers");
    }
  }, [currentTab, following, followers, fetchList]);

  const handleTabChange = (value: string) => {
    const next = value as Tab;
    if (next === currentTab) return;
    navigate(`/${handle}/${next}`, { replace: true });
  };

  return (
    <div className="w-full min-h-full flex flex-col">
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-14">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <div className="min-w-0 flex-1">
            <h1 className="text-base font-bold truncate">
              {profile?.displayName ?? customId}
            </h1>
            <p className="text-xs text-muted-foreground truncate">
              @{customId}
            </p>
          </div>
        </div>
      </div>

      <Tabs value={currentTab} onValueChange={handleTabChange}>
        <TabsList className="w-full rounded-none border-b bg-background p-0 h-11">
          <TabsTrigger value="following" className="text-sm h-full">
            フォロー中
          </TabsTrigger>
          <TabsTrigger value="followers" className="text-sm h-full">
            フォロワー
          </TabsTrigger>
        </TabsList>

        <TabsContent value="following" className="mt-0">
          <UserList
            users={following}
            loading={loading && currentTab === "following"}
            error={currentTab === "following" ? error : null}
            emptyMessage="フォロー中のユーザーはいません"
          />
        </TabsContent>
        <TabsContent value="followers" className="mt-0">
          <UserList
            users={followers}
            loading={loading && currentTab === "followers"}
            error={currentTab === "followers" ? error : null}
            emptyMessage="フォロワーはいません"
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}

function UserList({
  users,
  loading,
  error,
  emptyMessage,
}: {
  users: User[] | null;
  loading: boolean;
  error: string | null;
  emptyMessage: string;
}) {
  if (loading && users === null) {
    return (
      <p className="px-6 py-8 text-center text-sm text-muted-foreground">
        読み込み中...
      </p>
    );
  }
  if (error) {
    return (
      <p className="px-6 py-8 text-center text-sm text-destructive">{error}</p>
    );
  }
  if (!users || users.length === 0) {
    return (
      <p className="px-6 py-8 text-center text-sm text-muted-foreground">
        {emptyMessage}
      </p>
    );
  }
  return (
    <ul className="flex flex-col">
      {users.map((u) => (
        <UserListItem key={u.name} user={u} />
      ))}
    </ul>
  );
}
