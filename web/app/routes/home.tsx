import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { CallVisibility, type Call } from "~/lib/gen/call/v1/call_pb";
import { callClient, tokenStore } from "~/lib/api-client";
import { parseCall } from "~/lib/resource-name";
import { Button } from "~/components/ui/button";

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

export default function Home() {
  const navigate = useNavigate();
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [visibility, setVisibility] = useState<VisibilityChoice>("OPEN");
  const [creating, setCreating] = useState(false);
  const [joinId, setJoinId] = useState("");
  const [openCalls, setOpenCalls] = useState<Call[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setIsAuthenticated(Boolean(tokenStore.getAccessToken()));
  }, []);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const res = await callClient.listPublicCalls({});
        if (!cancelled) setOpenCalls(res.calls);
      } catch (err) {
        console.error("ListPublicCalls failed:", err);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  const handleCreate = async () => {
    setCreating(true);
    setError(null);
    try {
      const res = await callClient.createCall({
        visibility: toProtoVisibility(visibility),
      });
      if (!res.call) throw new Error("call was not returned");
      navigate(`/calls/${parseCall(res.call.name)}`);
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

  return (
    <div className="w-full min-h-full flex flex-col items-center justify-center px-6 py-8">
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

        <div className="flex flex-col gap-2">
          <p className="text-sm font-medium">公開中の通話</p>
          {openCalls.length === 0 ? (
            <p className="text-xs text-muted-foreground">
              現在公開中の通話はありません
            </p>
          ) : (
            <ul className="flex flex-col gap-2">
              {openCalls.map((call) => {
                const callId = parseCall(call.name);
                return (
                  <li
                    key={call.name}
                    className="rounded-md border p-3 flex items-center justify-between"
                  >
                    <div className="min-w-0">
                      <p className="text-xs font-mono truncate">{callId}</p>
                      <p className="text-[10px] text-muted-foreground truncate">
                        host: {call.host?.displayName ?? "unknown"}
                        {call.host?.customId ? ` @${call.host.customId}` : ""}
                      </p>
                    </div>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => navigate(`/calls/${callId}`)}
                    >
                      参加
                    </Button>
                  </li>
                );
              })}
            </ul>
          )}
        </div>

        {error && (
          <p className="text-sm text-destructive text-center">{error}</p>
        )}
      </div>
    </div>
  );
}
