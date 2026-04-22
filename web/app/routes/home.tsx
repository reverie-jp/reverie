import { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { type Call } from "~/lib/gen/call/v1/call_pb";
import { callClient, tokenStore } from "~/lib/api-client";
import { Button } from "~/components/ui/button";
import { CallCard } from "~/components/call-card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";

function CallList({
  calls,
  onJoin,
  empty,
}: {
  calls: Call[];
  onJoin: (id: string) => void;
  empty: string;
}) {
  if (calls.length === 0) {
    return <p className="text-xs text-muted-foreground px-1 py-6">{empty}</p>;
  }
  return (
    <ul className="flex flex-col gap-3">
      {calls.map((c) => (
        <li key={c.name}>
          <CallCard call={c} onJoin={onJoin} />
        </li>
      ))}
    </ul>
  );
}

function GuestFollowingPlaceholder({
  peekCalls,
  onLogin,
}: {
  peekCalls: Call[];
  onLogin: () => void;
}) {
  // Blur the current public feed so the tab has visual weight even for
  // guests — mirrors the design's "peek under frosted glass" treatment.
  const peek = peekCalls.slice(0, 4);
  return (
    <div className="relative min-h-90">
      <div
        aria-hidden
        className="pointer-events-none select-none flex flex-col gap-3 blur-[10px] opacity-70"
      >
        {peek.length > 0
          ? peek.map((c) => (
              <CallCard key={c.name} call={c} onJoin={() => {}} />
            ))
          : Array.from({ length: 3 }).map((_, i) => (
              <div
                key={i}
                className="h-28 rounded-2xl border border-white/10 bg-white/5"
              />
            ))}
      </div>

      <div className="absolute inset-x-0 top-6 flex justify-center px-2">
        <div
          className="w-full max-w-sm rounded-2xl border border-(--reverie-accent)/30 p-5 backdrop-blur-xl"
          style={{
            background:
              "linear-gradient(180deg, rgba(123,92,255,0.18), rgba(255,255,255,0.05))",
            boxShadow: "0 20px 60px rgba(0,0,0,0.4)",
          }}
        >
          <p className="font-display text-xl leading-tight">
            フォロー中を見るには
          </p>
          <p className="mt-1.5 text-[13px] text-muted-foreground leading-relaxed">
            アカウントを作ると、気になる人の通話が夢のログのように並びます。
          </p>
          <Button
            size="lg"
            onClick={onLogin}
            className="mt-4 w-full h-11 rounded-full font-medium text-[#120a2e]"
            style={{
              background:
                "linear-gradient(180deg, #c9b5ff, var(--reverie-accent))",
              boxShadow:
                "0 0 24px var(--reverie-accent-glow), inset 0 1px 0 rgba(255,255,255,0.5)",
            }}
          >
            ログインして続ける
          </Button>
        </div>
      </div>
    </div>
  );
}

export default function Home() {
  const navigate = useNavigate();
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [publicCalls, setPublicCalls] = useState<Call[]>([]);
  const [followingCalls, setFollowingCalls] = useState<Call[]>([]);
  const [followingLoading, setFollowingLoading] = useState(false);
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

  const handleJoinCall = (id: string) => navigate(`/calls/${id}`);

  const currentCount =
    tab === "following" ? followingCalls.length : publicCalls.length;

  return (
    <div className="w-full min-h-full px-6 pt-10 pb-16 flex flex-col items-center">
      <div className="w-full max-w-xl flex flex-col gap-8">
        {/* Hero */}
        <header className="flex flex-col gap-2">
          <div className="flex items-baseline justify-between gap-3">
            <h1 className="font-display text-4xl leading-none tracking-tight">
              Reverie
            </h1>
            <span className="text-[11px] font-medium uppercase tracking-widest text-muted-foreground/80">
              {currentCount} rooms
            </span>
          </div>
        </header>

        <Tabs
          value={tab}
          onValueChange={(v) => setTab(v as typeof tab)}
          className="gap-6"
        >
          <TabsList className="w-full rounded-full bg-white/5 border border-white/10 p-1 h-10">
            <TabsTrigger
              value="following"
              className="rounded-full text-sm h-8 data-active:bg-white/10 data-active:text-foreground"
            >
              フォロー中
            </TabsTrigger>
            <TabsTrigger
              value="public"
              className="rounded-full text-sm h-8 data-active:bg-white/10 data-active:text-foreground"
            >
              オープン
            </TabsTrigger>
          </TabsList>

          <TabsContent value="public" className="mt-0">
            <CallList
              calls={publicCalls}
              onJoin={handleJoinCall}
              empty="現在公開中の通話はありません"
            />
          </TabsContent>

          <TabsContent value="following" className="mt-0">
            {isAuthenticated ? (
              followingLoading ? (
                <p className="text-xs text-muted-foreground px-1 py-6">
                  読み込み中...
                </p>
              ) : (
                <CallList
                  calls={followingCalls}
                  onJoin={handleJoinCall}
                  empty="フォロー中のユーザーが参加している通話はありません"
                />
              )
            ) : (
              <GuestFollowingPlaceholder
                peekCalls={publicCalls}
                onLogin={() => navigate("/login")}
              />
            )}
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
