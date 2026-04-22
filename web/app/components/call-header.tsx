import { useEffect, useRef, useState } from "react";
import { Link, useLocation } from "react-router";
import { Mic, MicOff, PhoneOff } from "lucide-react";
import {
  RoomEvent,
  Track,
  type Participant,
  type RemoteParticipant,
  type RemoteTrack,
  type RemoteTrackPublication,
} from "livekit-client";
import { callClient } from "~/lib/api-client";
import { CALL_UPDATED_EVENT, useCall } from "~/lib/call-context";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { cn } from "~/lib/utils";
import { formatCall } from "~/lib/resource-name";
import type { Call } from "~/lib/gen/call/v1/call_pb";

const AVATAR_BACKGROUND =
  "linear-gradient(160deg, #c9b5ff, #6b4ee0 55%, #4a2d7d)";

// Mixed-audio FFT → 5 bars drawn from center (up+down) in lavender.
// One canvas, one RAF loop, subscribes to every participant in the room.
function MiniVoiceMeter() {
  const { connected, getRoom } = useCall();
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    if (!connected) return;
    const room = getRoom();
    if (!room) return;
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx2d = canvas.getContext("2d");
    if (!ctx2d) return;

    const audioCtx = new AudioContext();
    const analyser = audioCtx.createAnalyser();
    analyser.fftSize = 256;
    analyser.smoothingTimeConstant = 0.6;
    const data = new Uint8Array(new ArrayBuffer(analyser.frequencyBinCount));

    const sources = new Map<string, MediaStreamAudioSourceNode>();

    const attachParticipant = (p: Participant) => {
      const audio = p
        .getTrackPublications()
        .filter((pub) => pub.kind === Track.Kind.Audio);
      for (const pub of audio) {
        const mst = pub.track?.mediaStreamTrack;
        if (!mst) continue;
        const key = `${p.identity}:${pub.trackSid}`;
        if (sources.has(key)) continue;
        try {
          const src = audioCtx.createMediaStreamSource(new MediaStream([mst]));
          src.connect(analyser);
          sources.set(key, src);
        } catch (err) {
          console.warn("attach source failed:", err);
        }
      }
    };

    const detachByPrefix = (prefix: string) => {
      for (const key of Array.from(sources.keys())) {
        if (key.startsWith(prefix)) {
          const src = sources.get(key);
          try {
            src?.disconnect();
          } catch {}
          sources.delete(key);
        }
      }
    };

    attachParticipant(room.localParticipant);
    room.remoteParticipants.forEach((p) => attachParticipant(p));

    const onSubscribed = (
      _track: RemoteTrack,
      _pub: RemoteTrackPublication,
      participant: RemoteParticipant,
    ) => attachParticipant(participant);
    const onUnsubscribed = (
      _track: RemoteTrack,
      pub: RemoteTrackPublication,
      participant: RemoteParticipant,
    ) => {
      const key = `${participant.identity}:${pub.trackSid}`;
      const src = sources.get(key);
      try {
        src?.disconnect();
      } catch {}
      sources.delete(key);
    };
    const onParticipantDisconnected = (participant: RemoteParticipant) => {
      detachByPrefix(`${participant.identity}:`);
    };
    const onLocalPublished = () => attachParticipant(room.localParticipant);
    const onLocalUnpublished = () => {
      detachByPrefix(`${room.localParticipant.identity}:`);
    };

    room.on(RoomEvent.TrackSubscribed, onSubscribed);
    room.on(RoomEvent.TrackUnsubscribed, onUnsubscribed);
    room.on(RoomEvent.ParticipantDisconnected, onParticipantDisconnected);
    room.on(RoomEvent.LocalTrackPublished, onLocalPublished);
    room.on(RoomEvent.LocalTrackUnpublished, onLocalUnpublished);

    const NUM_BARS = 5;
    const BAR_GAP = 2;
    const MIN_BAR_HEIGHT = 2;
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
      const hiBin = Math.min(data.length - 1, Math.ceil(hi / binHz));
      barBins.push([loBin, Math.max(loBin, hiBin)]);
      barCenterHz.push(Math.sqrt(lo * hi));
    }
    const gains = barCenterHz.map((hz) => Math.sqrt(hz / barCenterHz[0]));
    const SENSITIVITY = 0.55;

    const dpr = window.devicePixelRatio || 1;
    canvas.width = 28 * dpr;
    canvas.height = 16 * dpr;
    ctx2d.scale(dpr, dpr);

    let animationId = 0;
    let cleaned = false;
    const draw = () => {
      if (cleaned) return;
      animationId = requestAnimationFrame(draw);
      analyser.getByteFrequencyData(data);
      const w = 28;
      const h = 16;
      ctx2d.clearRect(0, 0, w, h);
      const barWidth = (w - BAR_GAP * (NUM_BARS - 1)) / NUM_BARS;
      for (let b = 0; b < NUM_BARS; b++) {
        const [loBin, hiBin] = barBins[b];
        let sum = 0;
        const n = hiBin - loBin + 1;
        for (let i = loBin; i <= hiBin; i++) sum += data[i];
        const avg = n > 0 ? sum / n / 255 : 0;
        const level = Math.min(1, avg * gains[b] * SENSITIVITY);
        const barH = Math.max(MIN_BAR_HEIGHT, level * h);
        const x = b * (barWidth + BAR_GAP);
        const y = (h - barH) / 2; // grow from center
        ctx2d.fillStyle = `rgba(184, 164, 255, ${(0.45 + level * 0.55).toFixed(2)})`;
        ctx2d.fillRect(x, y, barWidth, barH);
      }
    };
    draw();

    return () => {
      cleaned = true;
      cancelAnimationFrame(animationId);
      room.off(RoomEvent.TrackSubscribed, onSubscribed);
      room.off(RoomEvent.TrackUnsubscribed, onUnsubscribed);
      room.off(RoomEvent.ParticipantDisconnected, onParticipantDisconnected);
      room.off(RoomEvent.LocalTrackPublished, onLocalPublished);
      room.off(RoomEvent.LocalTrackUnpublished, onLocalUnpublished);
      sources.forEach((src) => {
        try {
          src.disconnect();
        } catch {}
      });
      sources.clear();
      audioCtx.close().catch(() => undefined);
    };
  }, [connected, getRoom]);

  return (
    <canvas ref={canvasRef} style={{ width: 28, height: 16 }} aria-hidden />
  );
}

export function CallHeader() {
  const {
    callId,
    connected,
    isSelfMuted,
    toggleSelfMic,
    leave,
    getRoom,
    participantTick,
  } = useCall();
  const location = useLocation();
  const [callInfo, setCallInfo] = useState<Call | null>(null);

  // Fetch call metadata (title, host) once when the mini bar appears.
  useEffect(() => {
    if (!connected || !callId) {
      setCallInfo(null);
      return;
    }
    let cancelled = false;
    const fetchInfo = async () => {
      try {
        const res = await callClient.getCall({ name: formatCall(callId) });
        if (!cancelled) setCallInfo(res.call ?? null);
      } catch (err) {
        console.warn("CallHeader GetCall failed:", err);
      }
    };
    void fetchInfo();
    const onUpdated = () => void fetchInfo();
    window.addEventListener(CALL_UPDATED_EVENT, onUpdated);
    return () => {
      cancelled = true;
      window.removeEventListener(CALL_UPDATED_EVENT, onUpdated);
    };
  }, [callId, connected]);

  if (!connected || !callId) return null;
  if (location.pathname === `/calls/${callId}`) return null;

  const muted = isSelfMuted();
  const room = getRoom();
  // Recomputed each render (participantTick drives re-render on join/leave).
  void participantTick;
  const participantCount = room ? room.remoteParticipants.size + 1 : 1;
  const hostName = callInfo?.host?.displayName ?? "";
  const hostAvatar = callInfo?.host?.avatarUrl;
  const title = callInfo?.title?.trim() || hostName || "通話中";
  const initials = (hostName || "?").slice(0, 2).toUpperCase();

  return (
    <div className="sticky top-4 z-40 px-2.5 pb-1.5 pointer-events-none">
      <div
        className={cn(
          "pointer-events-auto flex items-center gap-2.5 rounded-2xl px-2 py-1.5",
          "border border-(--reverie-accent)/25",
        )}
        style={{
          background:
            "linear-gradient(180deg, rgba(123,92,255,0.18), rgba(184,164,255,0.06))",
          backdropFilter: "blur(20px) saturate(140%)",
          WebkitBackdropFilter: "blur(20px) saturate(140%)",
          boxShadow:
            "0 8px 24px rgba(0,0,0,0.3), 0 0 16px rgba(184,164,255,0.15)",
        }}
      >
        <Link
          to={`/calls/${callId}`}
          className="flex flex-1 min-w-0 items-center gap-2.5"
          aria-label="通話に戻る"
        >
          <div className="reverie-speaking-ring shrink-0">
            <Avatar
              className="size-8"
              style={{ background: AVATAR_BACKGROUND }}
            >
              {hostAvatar && <AvatarImage src={hostAvatar} alt={hostName} />}
              <AvatarFallback className="bg-transparent text-white text-[10px] font-display font-medium">
                {initials}
              </AvatarFallback>
            </Avatar>
          </div>

          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-1.5">
              <span className="reverie-live-pill" style={{ fontSize: 9 }}>
                LIVE
              </span>
              <span className="text-[10px] text-muted-foreground/80">
                {participantCount}人
              </span>
            </div>
            <p className="font-display text-[13px] leading-tight truncate mt-0.5">
              {title}
            </p>
          </div>

          <MiniVoiceMeter />
        </Link>

        <button
          type="button"
          onClick={() => void toggleSelfMic()}
          title={muted ? "ミュートを解除" : "ミュート"}
          aria-label={muted ? "ミュートを解除" : "ミュート"}
          className={cn(
            "grid place-items-center size-8 rounded-xl shrink-0 transition-colors",
            "border border-white/20 bg-white/12 hover:bg-white/18",
            muted && "text-(--reverie-live)",
          )}
        >
          {muted ? <MicOff className="size-4" /> : <Mic className="size-4" />}
        </button>

        <button
          type="button"
          onClick={() => void leave()}
          title="退出"
          aria-label="退出"
          className="grid place-items-center size-8 rounded-xl shrink-0 text-white font-bold transition-colors"
          style={{
            background: "rgba(255,90,120,0.85)",
            boxShadow: "0 0 12px rgba(255,90,120,0.5)",
          }}
        >
          <PhoneOff className="size-4" />
        </button>
      </div>
    </div>
  );
}
