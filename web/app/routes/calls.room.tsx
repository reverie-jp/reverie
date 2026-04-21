import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router";
import {
  Room,
  RoomEvent,
  Track,
  type RemoteTrack,
  type RemoteTrackPublication,
  type RemoteParticipant,
} from "livekit-client";
import type { JoinCallResponse } from "~/lib/gen/call/v1/call_pb";
import { callClient, tokenStore } from "~/lib/api-client";
import { Button } from "~/components/ui/button";

const GUEST_DISPLAY_NAME_KEY = "reverie.guest_display_name";
const REFRESH_MARGIN_MS = 60_000;

export default function CallRoomRoute() {
  const { roomId } = useParams<{ roomId: string }>();
  const isAuthenticated = Boolean(tokenStore.getAccessToken());
  const [guestDisplayName, setGuestDisplayName] = useState(
    () => localStorage.getItem(GUEST_DISPLAY_NAME_KEY) ?? "",
  );
  const [connected, setConnected] = useState(false);
  const [connecting, setConnecting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [, setTick] = useState(0);
  const roomRef = useRef<Room | null>(null);
  const audioContainerRef = useRef<HTMLDivElement | null>(null);
  const latestJoinResRef = useRef<JoinCallResponse | null>(null);
  const refreshTimerRef = useRef<number | null>(null);
  const intentionalLeaveRef = useRef(false);

  const rerender = () => setTick((t) => t + 1);

  const clearRefreshTimer = () => {
    if (refreshTimerRef.current !== null) {
      window.clearTimeout(refreshTimerRef.current);
      refreshTimerRef.current = null;
    }
  };

  const scheduleTokenRefresh = () => {
    clearRefreshTimer();
    // Guests cannot refresh: a new JoinCall issues a different guest identity.
    if (!isAuthenticated) return;
    const res = latestJoinResRef.current;
    if (!res?.expireTime) return;
    const expireMs = Number(res.expireTime.seconds) * 1000;
    const delay = expireMs - Date.now() - REFRESH_MARGIN_MS;
    if (delay <= 0) return;
    refreshTimerRef.current = window.setTimeout(() => {
      void refreshToken();
    }, delay);
  };

  const refreshToken = async () => {
    if (!roomId) return;
    try {
      const res = await callClient.joinCall({ roomId, guestDisplayName: "" });
      if (res.identity !== latestJoinResRef.current?.identity) {
        console.warn("Token refresh returned a different identity; ignoring");
        return;
      }
      latestJoinResRef.current = res;
      console.log(
        "Call token refreshed; next expiry",
        res.expireTime?.seconds,
      );
      scheduleTokenRefresh();
    } catch (err) {
      console.error("Token refresh failed:", err);
    }
  };

  const attachRoomListeners = (room: Room) => {
    room
      .on(RoomEvent.ParticipantConnected, rerender)
      .on(RoomEvent.ParticipantDisconnected, rerender)
      .on(
        RoomEvent.TrackSubscribed,
        (
          track: RemoteTrack,
          _pub: RemoteTrackPublication,
          _participant: RemoteParticipant,
        ) => {
          if (track.kind === Track.Kind.Audio && audioContainerRef.current) {
            audioContainerRef.current.appendChild(track.attach());
          }
        },
      )
      .on(RoomEvent.TrackUnsubscribed, (track: RemoteTrack) => {
        track.detach().forEach((el) => el.remove());
      })
      .on(RoomEvent.Disconnected, () => {
        void handleDisconnected();
      });
  };

  const handleDisconnected = async () => {
    clearRefreshTimer();
    if (intentionalLeaveRef.current) {
      intentionalLeaveRef.current = false;
      return;
    }
    if (!isAuthenticated) {
      setConnected(false);
      setError("接続が切れました。再参加してください");
      rerender();
      return;
    }
    try {
      const fresh = await callClient.joinCall({
        roomId: roomId!,
        guestDisplayName: "",
      });
      latestJoinResRef.current = fresh;
      const newRoom = new Room();
      attachRoomListeners(newRoom);
      await newRoom.connect(fresh.url, fresh.accessToken);
      await newRoom.localParticipant.setMicrophoneEnabled(true);
      roomRef.current = newRoom;
      scheduleTokenRefresh();
      console.log("Auto-reconnected with a fresh token");
      rerender();
    } catch (err) {
      console.error("Auto-reconnect failed:", err);
      setConnected(false);
      setError("再接続に失敗しました。再参加してください");
      rerender();
    }
  };

  const handleJoin = async () => {
    if (!roomId) return;
    if (!isAuthenticated && !guestDisplayName.trim()) return;
    setConnecting(true);
    setError(null);
    try {
      if (!isAuthenticated) {
        localStorage.setItem(GUEST_DISPLAY_NAME_KEY, guestDisplayName);
      }
      const res = await callClient.joinCall({
        roomId,
        guestDisplayName: isAuthenticated ? "" : guestDisplayName,
      });
      latestJoinResRef.current = res;

      const room = new Room();
      attachRoomListeners(room);
      await room.connect(res.url, res.accessToken);
      await room.localParticipant.setMicrophoneEnabled(true);

      roomRef.current = room;
      setConnected(true);
      scheduleTokenRefresh();
      rerender();
    } catch (err) {
      console.error("JoinCall failed:", err);
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setConnecting(false);
    }
  };

  const handleLeave = async () => {
    intentionalLeaveRef.current = true;
    clearRefreshTimer();
    await roomRef.current?.disconnect();
    roomRef.current = null;
    latestJoinResRef.current = null;
    setConnected(false);
    rerender();
  };

  useEffect(() => {
    return () => {
      intentionalLeaveRef.current = true;
      clearRefreshTimer();
      roomRef.current?.disconnect();
    };
  }, []);

  const room = roomRef.current;
  const remoteParticipants = room
    ? Array.from(room.remoteParticipants.values())
    : [];

  return (
    <div className="w-full min-h-full flex flex-col items-center justify-center px-6 py-8">
      <div className="w-full max-w-md flex flex-col gap-6">
        <div className="text-center">
          <p className="text-xs text-muted-foreground">Room ID</p>
          <p className="text-sm font-mono break-all">{roomId}</p>
        </div>

        {!connected ? (
          <div className="flex flex-col gap-3">
            {!isAuthenticated && (
              <input
                className="w-full h-11 px-3 rounded-md border border-input bg-background text-sm"
                placeholder="表示名（ゲスト）"
                value={guestDisplayName}
                onChange={(e) => setGuestDisplayName(e.target.value)}
              />
            )}
            <Button
              onClick={handleJoin}
              disabled={
                connecting || (!isAuthenticated && !guestDisplayName.trim())
              }
              className="w-full h-11"
            >
              {connecting ? "参加中..." : "通話に参加"}
            </Button>
          </div>
        ) : (
          <div className="flex flex-col gap-4">
            <div className="rounded-md border p-4">
              <p className="text-xs text-muted-foreground mb-2">参加者</p>
              <ul className="flex flex-col gap-1 text-sm">
                <li>
                  {room?.localParticipant.name ||
                    room?.localParticipant.identity}{" "}
                  <span className="text-muted-foreground">(自分)</span>
                </li>
                {remoteParticipants.map((p) => (
                  <li key={p.identity}>{p.name || p.identity}</li>
                ))}
              </ul>
            </div>
            <Button
              variant="destructive"
              onClick={handleLeave}
              className="w-full h-11"
            >
              退出
            </Button>
          </div>
        )}

        {error && (
          <p className="text-sm text-destructive text-center">{error}</p>
        )}
      </div>
      <div ref={audioContainerRef} className="hidden" aria-hidden />
    </div>
  );
}
