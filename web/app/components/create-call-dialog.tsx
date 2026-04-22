import { useState } from "react";
import { useNavigate } from "react-router";
import { Mic, X } from "lucide-react";
import { CallVisibility } from "~/lib/gen/call/v1/call_pb";
import { callClient } from "~/lib/api-client";
import { useCall } from "~/lib/call-context";
import { parseCall } from "~/lib/resource-name";
import { Button } from "~/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogTitle,
} from "~/components/ui/dialog";
import { cn } from "~/lib/utils";

type VisibilityChoice = "OPEN" | "USERS_ONLY" | "LOCKED";
const TITLE_MAX = 60;

const VISIBILITY_ITEMS: Array<{
  value: VisibilityChoice;
  label: string;
  desc: string;
}> = [
  {
    value: "OPEN",
    label: "オープン",
    desc: "タイムラインに載り、誰でも招待なしで参加できます。",
  },
  {
    value: "USERS_ONLY",
    label: "ユーザーのみ",
    desc: "アカウントを持つユーザーだけが参加できます。",
  },
  {
    value: "LOCKED",
    label: "招待のみ",
    desc: "URL を共有しても他人は参加できません。",
  },
];

const TITLE_SUGGESTIONS = [
  "☕ 雑談",
  "🎮 ゲーム",
  "📚 勉強",
  "🎵 音楽",
  "💬 相談",
];

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

export function CreateCallDialog({
  open,
  onOpenChange,
}: {
  open: boolean;
  onOpenChange: (o: boolean) => void;
}) {
  const navigate = useNavigate();
  const call = useCall();
  const [title, setTitle] = useState("");
  const [visibility, setVisibility] = useState<VisibilityChoice>("OPEN");
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleOpenChange = (next: boolean) => {
    onOpenChange(next);
    if (next) setError(null);
  };

  const handleCreate = async () => {
    setCreating(true);
    setError(null);
    try {
      const res = await callClient.createCall({
        visibility: toProtoVisibility(visibility),
        title: title.trim(),
      });
      if (!res.call) throw new Error("call was not returned");
      const newCallId = parseCall(res.call.name);
      if (call.callId && call.callId !== newCallId) {
        await call.leave();
      }
      const joinResult = await call.join(newCallId);
      if (!joinResult.ok) {
        throw new Error(joinResult.error);
      }
      onOpenChange(false);
      navigate(`/calls/${newCallId}`);
    } catch (err) {
      console.error("CreateCall failed:", err);
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setCreating(false);
    }
  };

  const visibilityDesc =
    VISIBILITY_ITEMS.find((v) => v.value === visibility)?.desc ?? "";

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent
        showCloseButton={false}
        className={cn(
          // Full-screen sheet: override shadcn's centered-modal positioning.
          "fixed inset-0 top-0 left-0 translate-x-0 translate-y-0",
          "w-screen h-dvh max-w-none rounded-none",
          "gap-0 p-0 ring-0 border-0",
          "flex flex-col overflow-hidden",
          // Aurora-tinted background so the sheet feels like the rest of
          // the app, not a bright white modal.
          "bg-linear-to-b from-[#0f0a28] via-[#07061a] to-[#04020f]",
        )}
      >
        <span
          aria-hidden
          className="pointer-events-none absolute inset-0"
          style={{
            background:
              "radial-gradient(600px 500px at 20% 0%, rgba(123,92,255,0.25), transparent 55%), radial-gradient(500px 400px at 100% 100%, rgba(255,107,157,0.12), transparent 60%)",
          }}
        />

        {/* Header */}
        <div className="relative flex items-center gap-3 px-5 pt-4 pb-3 shrink-0">
          <button
            type="button"
            onClick={() => handleOpenChange(false)}
            aria-label="閉じる"
            className="grid place-items-center size-8 rounded-lg border border-white/10 bg-white/4 text-muted-foreground hover:text-foreground hover:bg-white/8 transition-colors"
          >
            <X className="size-4" />
          </button>
          <DialogTitle className="font-display text-xl leading-none">
            新しい通話
          </DialogTitle>
        </div>

        {/* Scrollable form body */}
        <div className="relative flex-1 min-h-0 overflow-y-auto">
          <div className="mx-auto w-full max-w-xl flex flex-col gap-6 px-5 pt-4 pb-32">
            {/* Title */}
            <div className="flex flex-col gap-2">
              <span className="text-[10px] font-medium uppercase tracking-widest text-muted-foreground/80">
                タイトル
              </span>
              <div className="rounded-2xl border border-white/10 bg-white/5 px-4 py-5 min-h-19 focus-within:border-(--reverie-accent)/50 transition-colors">
                <input
                  autoFocus
                  maxLength={TITLE_MAX}
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="深夜の雑談ラジオ"
                  className="w-full bg-transparent font-display text-2xl leading-tight text-foreground placeholder:text-muted-foreground/50 focus:outline-none"
                />
              </div>
              <div className="flex flex-wrap gap-1.5">
                {TITLE_SUGGESTIONS.map((t) => (
                  <button
                    key={t}
                    type="button"
                    onClick={() => setTitle(t)}
                    className="reverie-tag hover:bg-white/10 transition-colors"
                  >
                    {t}
                  </button>
                ))}
              </div>
            </div>

            {/* Visibility */}
            <div className="flex flex-col gap-2">
              <span className="text-[10px] font-medium uppercase tracking-widest text-muted-foreground/80">
                公開範囲
              </span>
              <div className="flex rounded-full border border-white/10 bg-white/4 p-1">
                {VISIBILITY_ITEMS.map((v) => {
                  const selected = visibility === v.value;
                  return (
                    <button
                      key={v.value}
                      type="button"
                      onClick={() => setVisibility(v.value)}
                      className={cn(
                        "flex-1 rounded-full py-2 text-xs font-medium transition-colors",
                        selected
                          ? "bg-white/10 text-foreground shadow-[inset_0_0_0_1px_rgba(255,255,255,0.18)]"
                          : "text-muted-foreground hover:text-foreground",
                      )}
                    >
                      {v.label}
                    </button>
                  );
                })}
              </div>
              <p className="text-[11px] leading-relaxed text-muted-foreground">
                {visibilityDesc}
              </p>
            </div>

            {error && <p className="text-xs text-destructive">{error}</p>}
          </div>
        </div>

        {/* Floating primary action — sits above scroll. Backdrop fade blends
            the button into the content without a hard edge. */}
        <div className="pointer-events-none absolute inset-x-0 bottom-0 pb-[calc(env(safe-area-inset-bottom)+20px)] pt-10 bg-linear-to-t from-[#04020f] via-[#04020f]/80 to-transparent">
          <div className="pointer-events-auto mx-auto w-full max-w-xl px-5">
            <Button
              size="lg"
              disabled={creating}
              onClick={handleCreate}
              className="w-full h-12 rounded-full gap-2 font-medium text-[#120a2e]"
              style={{
                background:
                  "linear-gradient(180deg, #c9b5ff, var(--reverie-accent))",
                boxShadow:
                  "0 0 24px var(--reverie-accent-glow), inset 0 1px 0 rgba(255,255,255,0.5)",
              }}
            >
              <Mic className="size-4" strokeWidth={2.5} />
              {creating ? "開始中..." : "開始する"}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
