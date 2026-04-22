import { CallVisibility, type Call } from "~/lib/gen/call/v1/call_pb";
import {
  Avatar,
  AvatarFallback,
  AvatarGroup,
  AvatarGroupCount,
  AvatarImage,
} from "~/components/ui/avatar";
import { parseCall } from "~/lib/resource-name";
import { cn } from "~/lib/utils";

const STACK_SIZE = 6;

const VISIBILITY_LABEL: Record<CallVisibility, string> = {
  [CallVisibility.UNSPECIFIED]: "",
  [CallVisibility.OPEN]: "オープン",
  [CallVisibility.USERS_ONLY]: "ユーザーのみ",
  [CallVisibility.LOCKED]: "非公開",
};

// Lavender theme gradient for the fallback avatar (matches ParticipantTile).
const AVATAR_BACKGROUND =
  "linear-gradient(160deg, #c9b5ff, #6b4ee0 55%, #4a2d7d)";

function Waveform({ bars = 14 }: { bars?: number }) {
  return (
    <span className="reverie-wave" aria-hidden>
      {Array.from({ length: bars }).map((_, i) => (
        <span
          key={i}
          style={{
            height: `${20 + ((i * 37) % 80)}%`,
            animationDelay: `${(i * 0.08) % 1}s`,
          }}
        />
      ))}
    </span>
  );
}

export function CallCard({
  call,
  onJoin,
}: {
  call: Call;
  onJoin: (id: string) => void;
}) {
  const callId = parseCall(call.name);
  const hostName = call.host?.displayName ?? "ゲストホスト";
  const customId = call.host?.customId ?? "";
  const initials = hostName.slice(0, 2).toUpperCase();
  const visibilityLabel = VISIBILITY_LABEL[call.visibility];
  const title = call.title?.trim() || hostName;

  return (
    <button
      type="button"
      onClick={() => onJoin(callId)}
      className={cn(
        "group relative block w-full overflow-hidden rounded-2xl p-4 text-left",
        "border border-white/10 bg-linear-to-b from-white/[0.07] to-white/2",
        "backdrop-blur-xl backdrop-saturate-150 transition",
        "hover:border-white/20 hover:from-white/10 hover:to-white/4",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-(--reverie-accent)",
      )}
    >
      <span
        aria-hidden
        className="pointer-events-none absolute inset-0"
        style={{
          background:
            "radial-gradient(220px 180px at 85% 18%, rgba(184,164,255,0.2), transparent 70%)",
        }}
      />

      {/* Row 1: avatar + meta */}
      <div className="relative flex items-center gap-2.5">
        <div className="reverie-speaking-ring shrink-0">
          <Avatar
            size="lg"
            className="size-10"
            style={{ background: AVATAR_BACKGROUND }}
          >
            {call.host?.avatarUrl && (
              <AvatarImage src={call.host.avatarUrl} alt={hostName} />
            )}
            <AvatarFallback className="bg-transparent text-white font-display font-medium">
              {initials}
            </AvatarFallback>
          </Avatar>
        </div>

        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span className="reverie-live-pill">LIVE</span>
            {visibilityLabel && (
              <span className="reverie-tag">{visibilityLabel}</span>
            )}
          </div>
          <p className="mt-1 font-display text-lg leading-tight text-foreground line-clamp-2 wrap-break-word">
            {title}
          </p>
        </div>
      </div>

      {/* Row 2: host line — spans full card width, mirrors design */}
      <p className="relative mt-2.5 text-[11px] text-muted-foreground truncate">
        {hostName}
        {customId && (
          <span className="text-muted-foreground/70"> @{customId}</span>
        )}
        <span className="mx-1.5 text-muted-foreground/50">・</span>
        <span style={{ color: "var(--reverie-accent)" }}>通話中</span>
      </p>

      {/* Row 3: stacked participant avatars + waveform. Host is filtered
          out — they're already prominent on row 1 as the big avatar. */}
      {(() => {
        const others = call.activeParticipants.filter(
          (p) => !p.user || !call.host || p.user.name !== call.host.name,
        );
        return (
          <div className="relative mt-3 flex items-center justify-between gap-3">
            {others.length > 0 ? (
              <AvatarGroup>
                {others.slice(0, STACK_SIZE).map((p) => {
                  const name = p.user?.displayName || p.displayName || "?";
                  return (
                    <Avatar
                      key={p.name}
                      size="sm"
                      style={{ background: AVATAR_BACKGROUND }}
                    >
                      {p.user?.avatarUrl && (
                        <AvatarImage src={p.user.avatarUrl} alt={name} />
                      )}
                      <AvatarFallback className="bg-transparent text-white text-[10px] font-medium">
                        {name.slice(0, 2).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                  );
                })}
                {others.length > STACK_SIZE && (
                  <AvatarGroupCount>
                    +{others.length - STACK_SIZE}
                  </AvatarGroupCount>
                )}
              </AvatarGroup>
            ) : (
              <span className="text-[10px] font-mono uppercase tracking-wider text-muted-foreground/60">
                {callId.slice(0, 8).toLowerCase()}
              </span>
            )}
            <Waveform bars={14} />
          </div>
        );
      })()}
    </button>
  );
}
