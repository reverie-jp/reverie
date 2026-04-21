import { useEffect, useRef } from "react";
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
import { useCall } from "~/lib/call-context";
import { Button } from "~/components/ui/button";

function MixedVoiceMeter() {
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
          const src = audioCtx.createMediaStreamSource(
            new MediaStream([mst]),
          );
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

    let animationId = 0;
    let cleaned = false;
    const draw = () => {
      if (cleaned) return;
      animationId = requestAnimationFrame(draw);
      analyser.getByteFrequencyData(data);
      const w = canvas.width;
      const h = canvas.height;
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
        const y = h - barH;
        ctx2d.fillStyle = level > 0.05 ? "#22c55e" : "#9ca3af";
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
    <canvas ref={canvasRef} width={28} height={14} aria-hidden />
  );
}

export function CallHeader() {
  const { callId, connected, isSelfMuted, toggleSelfMic, leave } = useCall();
  const location = useLocation();
  if (!connected || !callId) return null;
  if (location.pathname === `/calls/${callId}`) return null;
  const muted = isSelfMuted();
  return (
    <div className="shrink-0 z-40 bg-green-500/10 border-b border-green-500/30 flex items-center gap-2 px-3 h-11">
      <span className="size-2 rounded-full bg-green-500 animate-pulse" />
      <Link
        to={`/calls/${callId}`}
        className="text-sm flex-1 truncate hover:underline flex items-center gap-2"
      >
        <span>通話中 — タップで戻る</span>
        <MixedVoiceMeter />
      </Link>
      <Button
        size="sm"
        variant="outline"
        onClick={() => void toggleSelfMic()}
        title={muted ? "ミュートを解除" : "ミュート"}
      >
        {muted ? <MicOff className="size-4" /> : <Mic className="size-4" />}
      </Button>
      <Button
        size="sm"
        variant="destructive"
        onClick={() => void leave()}
        title="退出"
      >
        <PhoneOff className="size-4" />
      </Button>
    </div>
  );
}
