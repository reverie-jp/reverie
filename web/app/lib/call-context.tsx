import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";
import {
  DisconnectReason,
  Room,
  RoomEvent,
  Track,
  type Participant,
  type RemoteParticipant,
  type RemoteTrack,
  type RemoteTrackPublication,
} from "livekit-client";
import type { JoinCallResponse } from "./gen/call/v1/call_pb";
import {
  API_BASE_URL,
  callClient,
  tokenStore,
  userClient,
} from "./api-client";
import { formatCall } from "./resource-name";

const GUEST_DISPLAY_NAME_KEY = "reverie.guest_display_name";
const REFRESH_MARGIN_MS = 60_000;
const HEARTBEAT_INTERVAL_MS = 30_000;
export const CALL_DEFAULT_VOLUME = 0.7;
export const CALL_UPDATED_EVENT = "reverie:call_updated";

function friendlyJoinError(err: unknown): string {
  const raw = err instanceof Error ? err.message : String(err);
  if (raw.includes("banned")) return "この通話から追放されています";
  if (raw.includes("ended")) return "この通話は終了しています";
  if (raw.includes("locked")) return "この通話は非公開で参加できません";
  if (raw.includes("guest")) return "この通話はログインユーザーのみ参加できます";
  if (raw.includes("already_in")) return "別の通話に参加しています";
  return raw;
}

export type ChatMessage = {
  id: string;
  identity: string;
  displayName: string;
  text: string;
  time: number;
  system?: boolean;
};

type JoinResult = { ok: true } | { ok: false; error: string };

type CallContextValue = {
  callId: string | null;
  callName: string;
  connected: boolean;
  connecting: boolean;
  joinError: string | null;
  identity: string | null;
  isAuthenticated: boolean;
  myUserName: string | null;
  chatMessages: ChatMessage[];
  volumes: Record<string, number>;
  hostMutedMe: boolean;
  hostMutedIdentities: Set<string>;
  becameHost: boolean;
  tick: number;

  join: (callId: string, guestDisplayName?: string) => Promise<JoinResult>;
  leave: () => Promise<void>;
  toggleSelfMic: () => Promise<void>;
  sendChat: (text: string) => Promise<void>;
  setVolume: (identity: string, value: number) => void;
  publishData: (payload: unknown) => Promise<void>;
  addSystemMessage: (text: string) => void;
  appendChatMessage: (msg: ChatMessage) => void;
  clearChatMessages: () => void;
  dismissBecameHost: () => void;
  markHostMuted: (identity: string, muted: boolean) => void;

  isSelfMuted: () => boolean;
  isParticipantMuted: (identity: string) => boolean;
  getRoom: () => Room | null;
  getLKParticipant: (identity: string) => Participant | null;
};

const CallContext = createContext<CallContextValue | null>(null);

export function useCall() {
  const ctx = useContext(CallContext);
  if (!ctx) throw new Error("useCall must be used within CallProvider");
  return ctx;
}

export function CallProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [myUserName, setMyUserName] = useState<string | null>(null);

  const [callId, setCallId] = useState<string | null>(null);
  const [connected, setConnected] = useState(false);
  const [connecting, setConnecting] = useState(false);
  const [joinError, setJoinError] = useState<string | null>(null);
  const [identity, setIdentity] = useState<string | null>(null);

  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
  const [volumes, setVolumes] = useState<Record<string, number>>({});
  const [hostMutedMe, setHostMutedMe] = useState(false);
  const [hostMutedIdentities, setHostMutedIdentities] = useState<Set<string>>(
    new Set(),
  );
  const [becameHost, setBecameHost] = useState(false);
  const [tick, setTick] = useState(0);

  const markHostMuted = useCallback((identity: string, muted: boolean) => {
    setHostMutedIdentities((prev) => {
      const has = prev.has(identity);
      if (muted && has) return prev;
      if (!muted && !has) return prev;
      const next = new Set(prev);
      if (muted) next.add(identity);
      else next.delete(identity);
      return next;
    });
  }, []);
  const rerender = useCallback(() => setTick((t) => t + 1), []);

  const roomRef = useRef<Room | null>(null);
  const joinResRef = useRef<JoinCallResponse | null>(null);
  const callIdRef = useRef<string | null>(null);
  const isAuthenticatedRef = useRef(false);
  const refreshTimerRef = useRef<number | null>(null);
  const heartbeatTimerRef = useRef<number | null>(null);
  const intentionalLeaveRef = useRef(false);
  const volumesRef = useRef<Record<string, number>>({});
  const hostMutedMeRef = useRef(false);
  const myUserNameRef = useRef<string | null>(null);
  const audioContainerRef = useRef<HTMLDivElement | null>(null);

  volumesRef.current = volumes;
  hostMutedMeRef.current = hostMutedMe;
  myUserNameRef.current = myUserName;
  callIdRef.current = callId;
  isAuthenticatedRef.current = isAuthenticated;

  const callName = callId ? formatCall(callId) : "";

  const addSystemMessage = useCallback((text: string) => {
    setChatMessages((prev) => [
      ...prev,
      {
        id: `system-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
        identity: "system",
        displayName: "システム",
        text,
        time: Date.now(),
        system: true,
      },
    ]);
  }, []);

  const appendChatMessage = useCallback((msg: ChatMessage) => {
    setChatMessages((prev) => [...prev, msg]);
  }, []);

  const clearChatMessages = useCallback(() => setChatMessages([]), []);

  const publishData = useCallback(async (payload: unknown) => {
    const room = roomRef.current;
    if (!room) return;
    try {
      const bytes = new TextEncoder().encode(JSON.stringify(payload));
      await room.localParticipant.publishData(bytes, { reliable: true });
    } catch (err) {
      console.warn("publishData failed:", err);
    }
  }, []);

  const getLKParticipant = useCallback(
    (id: string): Participant | null => {
      const room = roomRef.current;
      if (!room) return null;
      if (room.localParticipant.identity === id) return room.localParticipant;
      return room.remoteParticipants.get(id) ?? null;
    },
    [],
  );

  const isParticipantMuted = useCallback(
    (id: string): boolean => {
      const p = getLKParticipant(id);
      if (!p) return false;
      const audio = p
        .getTrackPublications()
        .filter((pub) => pub.kind === Track.Kind.Audio);
      if (audio.length === 0) return false;
      return audio.every((pub) => pub.isMuted);
    },
    [getLKParticipant],
  );

  const isSelfMuted = useCallback((): boolean => {
    const room = roomRef.current;
    if (!room || !connected) return false;
    return isParticipantMuted(room.localParticipant.identity);
  }, [connected, isParticipantMuted]);

  const clearRefreshTimer = () => {
    if (refreshTimerRef.current !== null) {
      window.clearTimeout(refreshTimerRef.current);
      refreshTimerRef.current = null;
    }
  };

  const clearHeartbeatTimer = () => {
    if (heartbeatTimerRef.current !== null) {
      window.clearInterval(heartbeatTimerRef.current);
      heartbeatTimerRef.current = null;
    }
  };

  const startHeartbeat = () => {
    clearHeartbeatTimer();
    heartbeatTimerRef.current = window.setInterval(() => {
      const cid = callIdRef.current;
      const res = joinResRef.current;
      if (!cid || !res) return;
      void callClient
        .heartbeatCall({
          name: formatCall(cid),
          guestIdentity: isAuthenticatedRef.current ? "" : res.identity,
        })
        .catch((err) => console.warn("HeartbeatCall failed:", err));
    }, HEARTBEAT_INTERVAL_MS);
  };

  const sendLeaveBeacon = () => {
    const cid = callIdRef.current;
    const res = joinResRef.current;
    if (!cid || !res) return;
    const token = tokenStore.getAccessToken();
    const body = JSON.stringify({
      name: formatCall(cid),
      guestIdentity: isAuthenticatedRef.current ? "" : res.identity,
    });
    try {
      void fetch(`${API_BASE_URL}/call.v1.CallService/LeaveCall`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        },
        body,
        keepalive: true,
      }).catch(() => undefined);
    } catch {
      // ignore
    }
  };

  const scheduleTokenRefresh = () => {
    clearRefreshTimer();
    if (!isAuthenticatedRef.current) return;
    const res = joinResRef.current;
    if (!res?.expireTime) return;
    const expireMs = Number(res.expireTime.seconds) * 1000;
    const delay = expireMs - Date.now() - REFRESH_MARGIN_MS;
    if (delay <= 0) return;
    refreshTimerRef.current = window.setTimeout(() => {
      void refreshToken();
    }, delay);
  };

  const refreshToken = async () => {
    const cid = callIdRef.current;
    if (!cid) return;
    try {
      const res = await callClient.joinCall({
        name: formatCall(cid),
        guestDisplayName: "",
      });
      if (res.identity !== joinResRef.current?.identity) {
        console.warn("Token refresh returned different identity; ignoring");
        return;
      }
      joinResRef.current = res;
      scheduleTokenRefresh();
    } catch (err) {
      console.error("Token refresh failed:", err);
    }
  };

  const applyVolumeToParticipant = (id: string, value: number) => {
    const room = roomRef.current;
    if (!room) return;
    const p = room.remoteParticipants.get(id);
    if (!p) return;
    p.getTrackPublications()
      .filter((pub) => pub.kind === Track.Kind.Audio)
      .forEach((pub) => {
        const track = pub.track;
        if (track && "setVolume" in track) {
          (track as { setVolume: (v: number) => void }).setVolume(value);
        }
      });
  };

  const setVolume = (id: string, value: number) => {
    setVolumes((prev) => ({ ...prev, [id]: value }));
    applyVolumeToParticipant(id, value);
  };

  const attachRoomListeners = (room: Room) => {
    room
      .on(RoomEvent.ParticipantConnected, (p: RemoteParticipant) => {
        rerender();
        addSystemMessage(
          `${p.name || p.identity} が通話に参加しました`,
        );
      })
      .on(RoomEvent.ParticipantDisconnected, (p: RemoteParticipant) => {
        rerender();
        addSystemMessage(`${p.name || p.identity} が退出しました`);
      })
      .on(
        RoomEvent.TrackSubscribed,
        (track: RemoteTrack, _pub, participant: RemoteParticipant) => {
          if (track.kind === Track.Kind.Audio && audioContainerRef.current) {
            audioContainerRef.current.appendChild(track.attach());
            if ("setVolume" in track) {
              const saved = volumesRef.current[participant.identity];
              const v =
                typeof saved === "number" ? saved : CALL_DEFAULT_VOLUME;
              (track as { setVolume: (v: number) => void }).setVolume(v);
            }
          }
        },
      )
      .on(RoomEvent.TrackUnsubscribed, (track: RemoteTrack) => {
        track.detach().forEach((el) => el.remove());
      })
      .on(RoomEvent.TrackMuted, () => rerender())
      .on(RoomEvent.TrackUnmuted, () => rerender())
      .on(
        RoomEvent.DataReceived,
        (payload: Uint8Array, participant?: RemoteParticipant) => {
          try {
            const msg = JSON.parse(new TextDecoder().decode(payload));
            if (msg?.type === "call_updated") {
              window.dispatchEvent(new CustomEvent(CALL_UPDATED_EVENT));
              return;
            }
            if (msg?.type === "host_muted" || msg?.type === "host_unmuted") {
              const me = roomRef.current?.localParticipant.identity;
              if (me && msg.targetIdentity === me) {
                const v = msg.type === "host_muted";
                setHostMutedMe(v);
                hostMutedMeRef.current = v;
              }
              if (typeof msg.targetIdentity === "string") {
                markHostMuted(
                  msg.targetIdentity,
                  msg.type === "host_muted",
                );
              }
              const targetName =
                typeof msg.targetDisplayName === "string"
                  ? msg.targetDisplayName
                  : "参加者";
              addSystemMessage(
                msg.type === "host_muted"
                  ? `ホストが ${targetName} をミュートしました`
                  : `ホストが ${targetName} のミュートを解除しました`,
              );
              return;
            }
            if (msg?.type === "unmute_request") {
              const requester =
                typeof msg.requesterDisplayName === "string"
                  ? msg.requesterDisplayName
                  : "参加者";
              addSystemMessage(
                `${requester} がミュートの解除を求めています`,
              );
              return;
            }
            if (msg?.type === "host_kicked") {
              const target =
                typeof msg.targetDisplayName === "string"
                  ? msg.targetDisplayName
                  : "参加者";
              addSystemMessage(`ホストが ${target} をキックしました`);
              return;
            }
            if (msg?.type === "host_banned") {
              const target =
                typeof msg.targetDisplayName === "string"
                  ? msg.targetDisplayName
                  : "参加者";
              addSystemMessage(`ホストが ${target} を追放しました`);
              return;
            }
            if (msg?.type === "host_transferred") {
              const newHostDisplayName =
                typeof msg.newHostDisplayName === "string"
                  ? msg.newHostDisplayName
                  : "参加者";
              const newHostName =
                typeof msg.newHostName === "string" ? msg.newHostName : "";
              addSystemMessage(
                `${newHostDisplayName} が新しいホストになりました`,
              );
              if (
                myUserNameRef.current &&
                myUserNameRef.current === newHostName
              ) {
                setBecameHost(true);
              }
              window.dispatchEvent(new CustomEvent(CALL_UPDATED_EVENT));
              return;
            }
            if (
              msg?.type === "chat_message" &&
              typeof msg.text === "string" &&
              participant
            ) {
              appendChatMessage({
                id: `${participant.identity}-${Date.now()}-${Math.random()
                  .toString(36)
                  .slice(2, 8)}`,
                identity: participant.identity,
                displayName: participant.name || participant.identity,
                text: msg.text,
                time: Date.now(),
              });
            }
          } catch {
            // ignore malformed
          }
        },
      )
      .on(RoomEvent.Disconnected, (reason?: DisconnectReason) => {
        void handleDisconnected(reason);
      });
  };

  const handleDisconnected = async (reason?: DisconnectReason) => {
    clearRefreshTimer();
    clearHeartbeatTimer();
    if (intentionalLeaveRef.current) {
      intentionalLeaveRef.current = false;
      roomRef.current = null;
      joinResRef.current = null;
      setConnected(false);
      setIdentity(null);
      setCallId(null);
      return;
    }
    if (reason === DisconnectReason.PARTICIPANT_REMOVED) {
      // Keep callId so the route-level gating still matches and the
      // "removed by host" notice shows on the call page until the user
      // navigates away.
      roomRef.current = null;
      joinResRef.current = null;
      setConnected(false);
      setIdentity(null);
      setJoinError("ホストにより通話から退出させられました");
      window.dispatchEvent(new CustomEvent(CALL_UPDATED_EVENT));
      return;
    }
    if (reason === DisconnectReason.ROOM_DELETED) {
      // Call was ended by the host. Clean up and surface the notice.
      // Server has also set end_time, so any rejoin attempt would fail.
      roomRef.current = null;
      joinResRef.current = null;
      setConnected(false);
      setIdentity(null);
      setJoinError("通話は終了しました");
      setChatMessages([]);
      setVolumes({});
      setHostMutedMe(false);
      hostMutedMeRef.current = false;
      window.dispatchEvent(new CustomEvent(CALL_UPDATED_EVENT));
      return;
    }
    if (!isAuthenticatedRef.current) {
      setConnected(false);
      setJoinError("接続が切れました。再参加してください");
      return;
    }
    const cid = callIdRef.current;
    if (!cid) return;
    try {
      const fresh = await callClient.joinCall({
        name: formatCall(cid),
        guestDisplayName: "",
      });
      joinResRef.current = fresh;
      const newRoom = new Room();
      attachRoomListeners(newRoom);
      await newRoom.connect(fresh.url, fresh.accessToken);
      await newRoom.localParticipant.setMicrophoneEnabled(true);
      roomRef.current = newRoom;
      scheduleTokenRefresh();
      startHeartbeat();
    } catch (err) {
      console.error("Auto-reconnect failed:", err);
      setConnected(false);
      setJoinError("再接続に失敗しました。再参加してください");
    }
  };

  const join = async (
    id: string,
    guestDisplayName?: string,
  ): Promise<JoinResult> => {
    if (roomRef.current) {
      return { ok: false, error: "already in a call" };
    }
    setConnecting(true);
    setJoinError(null);
    setCallId(id);
    callIdRef.current = id;
    try {
      if (!isAuthenticatedRef.current && guestDisplayName) {
        localStorage.setItem(GUEST_DISPLAY_NAME_KEY, guestDisplayName);
      }
      const res = await callClient.joinCall({
        name: formatCall(id),
        guestDisplayName: isAuthenticatedRef.current
          ? ""
          : guestDisplayName || "",
      });
      joinResRef.current = res;
      setIdentity(res.identity);

      const room = new Room();
      attachRoomListeners(room);
      await room.connect(res.url, res.accessToken);
      await room.localParticipant.setMicrophoneEnabled(true);

      roomRef.current = room;
      setConnected(true);
      setHostMutedMe(false);
      hostMutedMeRef.current = false;
      setChatMessages([]);
      setVolumes({});
      scheduleTokenRefresh();
      startHeartbeat();
      return { ok: true };
    } catch (err) {
      console.error("JoinCall failed:", err);
      const msg = friendlyJoinError(err);
      setJoinError(msg);
      // Keep callId so the route's joinError gating (call.callId === callId)
      // continues to surface the error. connected stays false and roomRef is
      // null, so the UI remains in the pre-join state and the user can retry.
      setConnected(false);
      joinResRef.current = null;
      setIdentity(null);
      return { ok: false, error: msg };
    } finally {
      setConnecting(false);
    }
  };

  const leave = async () => {
    intentionalLeaveRef.current = true;
    clearRefreshTimer();
    clearHeartbeatTimer();
    const cid = callIdRef.current;
    const res = joinResRef.current;
    if (cid && res) {
      try {
        await callClient.leaveCall({
          name: formatCall(cid),
          guestIdentity: isAuthenticatedRef.current ? "" : res.identity,
        });
      } catch (err) {
        console.warn("LeaveCall failed:", err);
      }
    }
    await roomRef.current?.disconnect();
    roomRef.current = null;
    joinResRef.current = null;
    setConnected(false);
    setCallId(null);
    callIdRef.current = null;
    setIdentity(null);
    setChatMessages([]);
    setVolumes({});
    setHostMutedMe(false);
    hostMutedMeRef.current = false;
    setHostMutedIdentities(new Set());
  };

  const toggleSelfMic = async () => {
    const room = roomRef.current;
    if (!room) return;
    const wantUnmute = isParticipantMuted(room.localParticipant.identity);
    if (wantUnmute && hostMutedMeRef.current) {
      const myName =
        room.localParticipant.name || room.localParticipant.identity;
      addSystemMessage(`${myName} がミュートの解除を求めています`);
      await publishData({
        type: "unmute_request",
        requesterDisplayName: myName,
      });
      return;
    }
    try {
      await room.localParticipant.setMicrophoneEnabled(wantUnmute);
    } catch (err) {
      console.warn("toggle mic failed:", err);
    }
  };

  const sendChat = async (text: string) => {
    const room = roomRef.current;
    if (!room) return;
    const trimmed = text.trim();
    if (!trimmed) return;
    await publishData({ type: "chat_message", text: trimmed });
    const local = room.localParticipant;
    appendChatMessage({
      id: `${local.identity}-${Date.now()}-${Math.random()
        .toString(36)
        .slice(2, 8)}`,
      identity: local.identity,
      displayName: local.name || local.identity,
      text: trimmed,
      time: Date.now(),
    });
  };

  useEffect(() => {
    const authed = Boolean(tokenStore.getAccessToken());
    setIsAuthenticated(authed);
    isAuthenticatedRef.current = authed;
    if (authed) {
      userClient
        .getMyUser({})
        .then((res) => {
          if (res.user) setMyUserName(res.user.name);
        })
        .catch((err) => console.warn("GetMyUser failed:", err));
    }
  }, []);

  useEffect(() => {
    const onPageHide = () => {
      sendLeaveBeacon();
    };
    window.addEventListener("pagehide", onPageHide);
    return () => {
      window.removeEventListener("pagehide", onPageHide);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);


  const value: CallContextValue = {
    callId,
    callName,
    connected,
    connecting,
    joinError,
    identity,
    isAuthenticated,
    myUserName,
    chatMessages,
    volumes,
    hostMutedMe,
    hostMutedIdentities,
    becameHost,
    tick,
    join,
    leave,
    toggleSelfMic,
    sendChat,
    setVolume,
    publishData,
    addSystemMessage,
    appendChatMessage,
    clearChatMessages,
    dismissBecameHost: () => setBecameHost(false),
    markHostMuted,
    isSelfMuted,
    isParticipantMuted,
    getRoom: () => roomRef.current,
    getLKParticipant,
  };

  return (
    <CallContext.Provider value={value}>
      {children}
      <div ref={audioContainerRef} className="hidden" aria-hidden />
    </CallContext.Provider>
  );
}
