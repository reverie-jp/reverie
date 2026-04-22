import { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { CallVisibility, type Call } from "~/lib/gen/call/v1/call_pb";
import { callClient, tokenStore } from "~/lib/api-client";
import { useCall } from "~/lib/call-context";
import { parseCall } from "~/lib/resource-name";
import { Button } from "~/components/ui/button";
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "~/components/ui/tabs";

type VisibilityChoice = "OPEN" | "USERS_ONLY" | "LOCKED";

const VISIBILITY_LABELS: Record<VisibilityChoice, string> = {
  OPEN: "オープン",
  USERS_ONLY: "ユーザーのみ",
  LOCKED: "非公開",
};

function toProtoVisibility(v: VisibilityChoice): CallVisibility {
  switch (v) {
    case "OPEN":
      return CallVisibility.OPEN;
    case "USERS_ONLY":
      return CallVisibility.USERS_ONLY;
    case "LOCKED":
      return CallVisibility.LOCKED;
  }
}

function CallListItem({
  call,
  onJoin,
}: {
  call: Call;
  onJoin: (id: string) => void;
}) {
  const callId = parseCall(call.name);
  return (
    <li className="rounded-md border p-3 flex items-center justify-between">
      <div className="min-w-0">
        <p className="text-xs font-mono truncate">{callId}</p>
        <p className="text-[10px] text-muted-foreground truncate">
          host: {call.host?.displayName ?? "unknown"}
          {call.host?.customId ? ` @${call.host.customId}` : ""}
        </p>
      </div>
      <Button size="sm" variant="outline" onClick={() => onJoin(callId)}>
        参加
      </Button>
    </li>
  );
}

export default function Home() {
  const navigate = useNavigate();
  const call = useCall();
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [visibility, setVisibility] = useState<VisibilityChoice>("OPEN");
  const [creating, setCreating] = useState(false);
  const [joinId, setJoinId] = useState("");
  const [publicCalls, setPublicCalls] = useState<Call[]>([]);
  const [followingCalls, setFollowingCalls] = useState<Call[]>([]);
  const [followingLoading, setFollowingLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [tab, setTab] = useState<"public" | "following">("public");

  useEffect(() => {
    setIsAuthenticated(Boolean(tokenStore.getAccessToken()));
  }, []);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const res = await callClient.listPublicCalls({});
        if (!cancelled) setPublicCalls(res.calls);
      } catch (err) {
        console.error("ListPublicCalls failed:", err);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  const fetchFollowing = useCallback(async () => {
    setFollowingLoading(true);
    try {
      const res = await callClient.listFollowingCalls({});
      setFollowingCalls(res.calls);
    } catch (err) {
      console.error("ListFollowingCalls failed:", err);
    } finally {
      setFollowingLoading(false);
    }
  }, []);

  useEffect(() => {
    if (!isAuthenticated) return;
    if (tab !== "following") return;
    void fetchFollowing();
  }, [isAuthenticated, tab, fetchFollowing]);

  const handleCreate = async () => {
    setCreating(true);
    setError(null);
    try {
      const res = await callClient.createCall({
        visibility: toProtoVisibility(visibility),
      });
      if (!res.call) throw new Error("call was not returned");
      const newCallId = parseCall(res.call.name);
      // Host joins immediately so landing on /calls/{id} skips the preview
      // screen. If a different call was previously joined, leave it first.
      if (call.callId && call.callId !== newCallId) {
        await call.leave();
      }
      const joinResult = await call.join(newCallId);
      if (!joinResult.ok) {
        throw new Error(joinResult.error);
      }
      navigate(`/calls/${newCallId}`);
    } catch (err) {
      console.error("CreateCall failed:", err);
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setCreating(false);
    }
  };

  const handleJoin = () => {
    if (!joinId.trim()) return;
    navigate(`/calls/${joinId.trim()}`);
  };

  const handleJoinCall = (id: string) => navigate(`/calls/${id}`);

  return (
    <div className="w-full min-h-full flex flex-col items-center px-6 py-8">
      <div className="w-full max-w-md flex flex-col gap-8">
        {isAuthenticated && (
          <div className="flex flex-col gap-3">
            <p className="text-sm font-medium">新しい通話を作成</p>
            <div className="flex flex-col gap-2">
              {(["OPEN", "USERS_ONLY", "LOCKED"] as VisibilityChoice[]).map(
                (v) => (
                  <label
                    key={v}
                    className="flex items-center gap-2 text-sm cursor-pointer"
                  >
                    <input
                      type="radio"
                      name="visibility"
                      value={v}
                      checked={visibility === v}
                      onChange={() => setVisibility(v)}
                    />
                    {VISIBILITY_LABELS[v]}
                  </label>
                ),
              )}
            </div>
            <Button
              onClick={handleCreate}
              disabled={creating}
              className="w-full h-11"
            >
              {creating ? "作成中..." : "作成"}
            </Button>
          </div>
        )}

        <div className="flex flex-col gap-2">
          <p className="text-sm font-medium">Room ID で参加</p>
          <input
            className="w-full h-11 px-3 rounded-md border border-input bg-background text-sm"
            placeholder="ULID"
            value={joinId}
            onChange={(e) => setJoinId(e.target.value)}
          />
          <Button
            variant="outline"
            onClick={handleJoin}
            disabled={!joinId.trim()}
            className="w-full h-11"
          >
            参加
          </Button>
        </div>

        <Tabs value={tab} onValueChange={(v) => setTab(v as typeof tab)}>
          <TabsList className="w-full">
            <TabsTrigger value="public" className="text-sm">
              公開中
            </TabsTrigger>
            {isAuthenticated && (
              <TabsTrigger value="following" className="text-sm">
                フォロー中
              </TabsTrigger>
            )}
          </TabsList>

          <TabsContent value="public" className="mt-3">
            {publicCalls.length === 0 ? (
              <p className="text-xs text-muted-foreground">
                現在公開中の通話はありません
              </p>
            ) : (
              <ul className="flex flex-col gap-2">
                {publicCalls.map((call) => (
                  <CallListItem
                    key={call.name}
                    call={call}
                    onJoin={handleJoinCall}
                  />
                ))}
              </ul>
            )}
          </TabsContent>

          {isAuthenticated && (
            <TabsContent value="following" className="mt-3">
              {followingLoading ? (
                <p className="text-xs text-muted-foreground">読み込み中...</p>
              ) : followingCalls.length === 0 ? (
                <p className="text-xs text-muted-foreground">
                  フォロー中のユーザーが参加している通話はありません
                </p>
              ) : (
                <ul className="flex flex-col gap-2">
                  {followingCalls.map((call) => (
                    <CallListItem
                      key={call.name}
                      call={call}
                      onJoin={handleJoinCall}
                    />
                  ))}
                </ul>
              )}
            </TabsContent>
          )}
        </Tabs>

        {error && (
          <p className="text-sm text-destructive text-center">{error}</p>
        )}
      </div>
    </div>
  );
}
