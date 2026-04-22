import { useEffect, useRef } from "react";
import { ParticipantEvent, Track, type Participant } from "livekit-client";
import { Crown, MicOff } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { cn } from "~/lib/utils";

const NUM_BARS = 6;
const BAR_MIN_HEIGHT = 0.18;

// Drives the avatar's scale + lavender glow AND the waveform bars below the
// name from a single audio analyser. One RAF loop feeds both: RMS on
// time-domain data for the avatar, log-spaced frequency bins for the bars.
function useAudioReactive(
  avatarRef: React.RefObject<HTMLDivElement | null>,
  barsRef: React.RefObject<(HTMLSpanElement | null)[]>,
  participant: Participant | null,
  active: boolean,
) {
  useEffect(() => {
    if (!participant || !active) return;
    const avatarEl = avatarRef.current;
    if (!avatarEl) return;

    let audioCtx: AudioContext | null = null;
    let analyser: AnalyserNode | null = null;
    let source: MediaStreamAudioSourceNode | null = null;
    let timeData: Uint8Array<ArrayBuffer> | null = null;
    let freqData: Uint8Array<ArrayBuffer> | null = null;
    let animationId = 0;
    let cleaned = false;

    const resetAvatar = () => {
      avatarEl.style.transform = "";
      avatarEl.style.boxShadow = "";
    };

    const resetBars = () => {
      const els = barsRef.current;
      if (!els) return;
      for (const el of els) {
        if (el) el.style.height = `${BAR_MIN_HEIGHT * 100}%`;
      }
    };

    const tryAttach = () => {
      if (cleaned || analyser) return;
      const pub = participant
        .getTrackPublications()
        .find((p) => p.kind === Track.Kind.Audio);
      const mst = pub?.track?.mediaStreamTrack;
      if (!mst || pub?.isMuted) {
        resetAvatar();
        resetBars();
        return;
      }
      audioCtx = new AudioContext();
      source = audioCtx.createMediaStreamSource(new MediaStream([mst]));
      analyser = audioCtx.createAnalyser();
      analyser.fftSize = 256;
      analyser.smoothingTimeConstant = 0.5;
      source.connect(analyser);
      timeData = new Uint8Array(new ArrayBuffer(analyser.frequencyBinCount));
      freqData = new Uint8Array(new ArrayBuffer(analyser.frequencyBinCount));

      // Log-spaced bin ranges + per-bin gain (flatter response at low freq).
      const binHz = audioCtx.sampleRate / analyser.fftSize;
      const minFreq = 150;
      const maxFreq = 5000;
      const ratio = maxFreq / minFreq;
      const barBins: [number, number][] = [];
      const barCenterHz: number[] = [];
      for (let b = 0; b < NUM_BARS; b++) {
        const lo = minFreq * Math.pow(ratio, b / NUM_BARS);
        const hi = minFreq * Math.pow(ratio, (b + 1) / NUM_BARS);
        const loBin = Math.max(0, Math.floor(lo / binHz));
        const hiBin = Math.min(freqData.length - 1, Math.ceil(hi / binHz));
        barBins.push([loBin, Math.max(loBin, hiBin)]);
        barCenterHz.push(Math.sqrt(lo * hi));
      }
      const gains = barCenterHz.map((hz) => Math.sqrt(hz / barCenterHz[0]));
      const BAR_SENSITIVITY = 0.55;
      const SCALE_GAIN = 6;

      const draw = () => {
        if (cleaned || !analyser || !timeData || !freqData) return;
        animationId = requestAnimationFrame(draw);

        // Avatar: RMS → scale + lavender glow.
        analyser.getByteTimeDomainData(timeData);
        let sum = 0;
        for (let i = 0; i < timeData.length; i++) {
          const v = (timeData[i] - 128) / 128;
          sum += v * v;
        }
        const rms = Math.sqrt(sum / timeData.length);
        const level = Math.min(1, rms * SCALE_GAIN);
        avatarEl.style.transform = `scale(${(1 + level * 0.18).toFixed(3)})`;
        if (level > 0.04) {
          const ringPx = 2 + Math.round(level * 6);
          const alpha = (0.25 + level * 0.5).toFixed(2);
          avatarEl.style.boxShadow = `0 0 0 ${ringPx}px rgb(184 164 255 / ${alpha})`;
        } else {
          avatarEl.style.boxShadow = "";
        }

        // Bars: frequency-domain → per-bar height.
        analyser.getByteFrequencyData(freqData);
        const els = barsRef.current;
        if (els) {
          for (let b = 0; b < NUM_BARS; b++) {
            const el = els[b];
            if (!el) continue;
            const [loBin, hiBin] = barBins[b];
            let s = 0;
            const n = hiBin - loBin + 1;
            for (let i = loBin; i <= hiBin; i++) s += freqData[i];
            const avg = n > 0 ? s / n / 255 : 0;
            const lvl = Math.min(1, avg * gains[b] * BAR_SENSITIVITY);
            const h = Math.max(BAR_MIN_HEIGHT, lvl);
            el.style.height = `${(h * 100).toFixed(1)}%`;
          }
        }
      };
      draw();
    };

    const teardown = () => {
      cancelAnimationFrame(animationId);
      animationId = 0;
      try {
        source?.disconnect();
      } catch {}
      source = null;
      analyser = null;
      timeData = null;
      freqData = null;
      audioCtx?.close().catch(() => undefined);
      audioCtx = null;
      resetAvatar();
      resetBars();
    };

    const onTrackChanged = () => {
      const pub = participant
        .getTrackPublications()
        .find((p) => p.kind === Track.Kind.Audio);
      const mst = pub?.track?.mediaStreamTrack;
      if (!mst || pub?.isMuted) {
        if (analyser) teardown();
        else {
          resetAvatar();
          resetBars();
        }
        return;
      }
      if (!analyser) tryAttach();
    };

    resetBars();
    tryAttach();
    participant.on(ParticipantEvent.TrackSubscribed, onTrackChanged);
    participant.on(ParticipantEvent.TrackUnsubscribed, onTrackChanged);
    participant.on(ParticipantEvent.TrackMuted, onTrackChanged);
    participant.on(ParticipantEvent.TrackUnmuted, onTrackChanged);
    participant.on(ParticipantEvent.LocalTrackPublished, onTrackChanged);
    participant.on(ParticipantEvent.LocalTrackUnpublished, onTrackChanged);

    return () => {
      cleaned = true;
      participant.off(ParticipantEvent.TrackSubscribed, onTrackChanged);
      participant.off(ParticipantEvent.TrackUnsubscribed, onTrackChanged);
      participant.off(ParticipantEvent.TrackMuted, onTrackChanged);
      participant.off(ParticipantEvent.TrackUnmuted, onTrackChanged);
      participant.off(ParticipantEvent.LocalTrackPublished, onTrackChanged);
      participant.off(ParticipantEvent.LocalTrackUnpublished, onTrackChanged);
      cancelAnimationFrame(animationId);
      try {
        source?.disconnect();
      } catch {}
      audioCtx?.close().catch(() => undefined);
      resetAvatar();
      resetBars();
    };
  }, [participant, active, avatarRef, barsRef]);
}

// Lavender theme gradient — avatar circle.
const AVATAR_BACKGROUND =
  "linear-gradient(160deg, #c9b5ff, #6b4ee0 55%, #4a2d7d)";

export function ParticipantTile({
  displayName,
  avatarUrl,
  connected,
  muted,
  isHost,
  participant,
  onClick,
}: {
  displayName: string;
  avatarUrl?: string;
  connected: boolean;
  muted: boolean;
  isHost: boolean;
  participant: Participant | null;
  onClick?: () => void;
}) {
  const avatarRef = useRef<HTMLDivElement>(null);
  const barsRef = useRef<(HTMLSpanElement | null)[]>([]);
  const speaking = connected && !muted;
  useAudioReactive(avatarRef, barsRef, participant, speaking);
  const initials = displayName.slice(0, 2).toUpperCase();

  return (
    <button
      type="button"
      onClick={onClick}
      className="flex flex-col items-center gap-1.5 p-1 text-center rounded-lg hover:bg-foreground/5 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--reverie-accent)]/60"
    >
      <div className={cn("relative", !connected && "opacity-40")}>
        <div
          ref={avatarRef}
          className="rounded-full transition-[transform,box-shadow] duration-75"
        >
          <Avatar
            className="size-14"
            style={{ background: AVATAR_BACKGROUND }}
          >
            {avatarUrl && <AvatarImage src={avatarUrl} alt={displayName} />}
            <AvatarFallback className="bg-transparent text-white font-display font-medium text-base">
              {initials}
            </AvatarFallback>
          </Avatar>
        </div>

        {connected && muted && (
          <span
            className="absolute -bottom-0.5 -right-0.5 grid place-items-center size-5 rounded-full bg-background/90 ring-1 ring-white/15"
            title="ミュート中"
          >
            <MicOff className="size-2.5 text-muted-foreground" />
          </span>
        )}

        {isHost && (
          <span
            className="absolute -top-1 -right-1 grid place-items-center size-5 rounded-full text-[10px] text-[#120a2e] font-bold ring-[1.5px] ring-[#0a081a]"
            style={{
              background: "linear-gradient(160deg, #c9b5ff, #8a6dff)",
            }}
            title="ホスト"
          >
            <Crown className="size-2.5" strokeWidth={2.5} />
          </span>
        )}
      </div>

      <span
        className={cn(
          "text-[11px] leading-tight truncate max-w-full",
          speaking
            ? "text-foreground font-medium"
            : "text-muted-foreground/80",
        )}
      >
        {displayName}
      </span>

      {/* Audio-reactive waveform. items-center so bars grow from the middle
          (up AND down) as in the design, not from a fixed baseline. */}
      <span className="h-3 flex items-center gap-[1.5px]" aria-hidden>
        {speaking &&
          Array.from({ length: NUM_BARS }).map((_, i) => (
            <span
              key={i}
              ref={(el) => {
                barsRef.current[i] = el;
              }}
              className="w-[1.5px] rounded-[1px] bg-[var(--reverie-accent)]"
              style={{
                height: `${BAR_MIN_HEIGHT * 100}%`,
                boxShadow: "0 0 6px var(--reverie-accent-glow)",
                transition: "height 60ms linear",
              }}
            />
          ))}
      </span>
    </button>
  );
}
