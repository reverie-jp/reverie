import { useEffect, useRef, useState } from "react";
import { Link, useParams } from "react-router";
import {
  Room,
  RoomEvent,
  Track,
  type RemoteTrack,
  type RemoteTrackPublication,
  type RemoteParticipant,
} from "livekit-client";
import type {
  CallParticipant,
  GetCallResponse,
  JoinCallResponse,
} from "~/lib/gen/call/v1/call_pb";
import { CallVisibility } from "~/lib/gen/call/v1/call_pb";
import {
  API_BASE_URL,
  callClient,
  tokenStore,
  userClient,
} from "~/lib/api-client";
import { formatCall } from "~/lib/resource-name";
import { Button } from "~/components/ui/button";

const GUEST_DISPLAY_NAME_KEY = "reverie.guest_display_name";
const REFRESH_MARGIN_MS = 60_000;
const HEARTBEAT_INTERVAL_MS = 30_000;

const VISIBILITY_LABELS: Record<CallVisibility, string> = {
  [CallVisibility.UNSPECIFIED]: "不明",
  [CallVisibility.OPEN]: "オープン",
  [CallVisibility.USERS_ONLY]: "ユーザーのみ",
  [CallVisibility.LOCKED]: "非公開",
};

const UPDATABLE_VISIBILITIES: CallVisibility[] = [
  CallVisibility.OPEN,
  CallVisibility.USERS_ONLY,
  CallVisibility.LOCKED,
];

type ChatMessage = {
  id: string;
  identity: string;
  displayName: string;
  text: string;
};

function makeMessageID(identity: string): string {
  return `${identity}-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`;
}

export default function CallRoomRoute() {
  const { callId } = useParams<{ callId: string }>();
  const callName = callId ? formatCall(callId) : "";
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [myUserName, setMyUserName] = useState<string | null>(null);
  const [guestDisplayName, setGuestDisplayName] = useState("");
  const [callInfo, setCallInfo] = useState<GetCallResponse | null>(null);
  const [loadingCall, setLoadingCall] = useState(true);
  const [callLoadError, setCallLoadError] = useState<string | null>(null);
  const [connected, setConnected] = useState(false);
  const [connecting, setConnecting] = useState(false);
  const [joinError, setJoinError] = useState<string | null>(null);
  const [updatingVisibility, setUpdatingVisibility] = useState(false);
  const [inviteCopied, setInviteCopied] = useState(false);
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);
  const [chatInput, setChatInput] = useState("");
  const [, setTick] = useState(0);
  const chatScrollRef = useRef<HTMLDivElement | null>(null);
  const roomRef = useRef<Room | null>(null);
  const audioContainerRef = useRef<HTMLDivElement | null>(null);
  const latestJoinResRef = useRef<JoinCallResponse | null>(null);
  const refreshTimerRef = useRef<number | null>(null);
  const heartbeatTimerRef = useRef<number | null>(null);
  const intentionalLeaveRef = useRef(false);

  const rerender = () => setTick((t) => t + 1);

  const fetchCallInfo = async (opts?: { silent?: boolean }) => {
    if (!callName) return;
    const silent = opts?.silent ?? false;
    if (!silent) {
      setLoadingCall(true);
      setCallLoadError(null);
    }
    try {
      const res = await callClient.getCall({
        name: callName,
        guestIdentity: latestJoinResRef.current?.identity ?? "",
      });
      setCallInfo(res);
    } catch (err) {
      console.error("GetCall failed:", err);
      if (!silent) {
        setCallLoadError(err instanceof Error ? err.message : String(err));
      }
    } finally {
      if (!silent) setLoadingCall(false);
    }
  };

  useEffect(() => {
    const authed = Boolean(tokenStore.getAccessToken());
    setIsAuthenticated(authed);
    setGuestDisplayName(localStorage.getItem(GUEST_DISPLAY_NAME_KEY) ?? "");
    if (authed) {
      userClient
        .getMyUser({})
        .then((res) => {
          if (res.user) setMyUserName(res.user.name);
        })
        .catch((err) => {
          console.warn("GetMyUser failed:", err);
        });
    }
  }, []);

  const sendChat = async () => {
    const room = roomRef.current;
    const text = chatInput.trim();
    if (!room || !text) return;
    try {
      const payload = new TextEncoder().encode(
        JSON.stringify({ type: "chat_message", text }),
      );
      await room.localParticipant.publishData(payload, { reliable: true });
      const local = room.localParticipant;
      setChatMessages((prev) => [
        ...prev,
        {
          id: makeMessageID(local.identity),
          identity: local.identity,
          displayName: local.name || local.identity,
          text,
        },
      ]);
      setChatInput("");
    } catch (err) {
      console.warn("publishData (chat) failed:", err);
    }
  };

  useEffect(() => {
    if (chatScrollRef.current) {
      chatScrollRef.current.scrollTop = chatScrollRef.current.scrollHeight;
    }
  }, [chatMessages]);

  const copyInviteLink = async () => {
    if (!callId) return;
    const url = `${window.location.origin}/calls/${callId}`;
    try {
      await navigator.clipboard.writeText(url);
      setInviteCopied(true);
      window.setTimeout(() => setInviteCopied(false), 1500);
    } catch (err) {
      console.warn("clipboard.writeText failed:", err);
    }
  };

  const handleUpdateVisibility = async (visibility: CallVisibility) => {
    if (!callName) return;
    setUpdatingVisibility(true);
    try {
      await callClient.updateCall({
        call: { name: callName, visibility },
        updateMask: { paths: ["visibility"] },
      });
      await fetchCallInfo({ silent: true });
      // Push to other connected participants so they refetch.
      const room = roomRef.current;
      if (room) {
        try {
          const payload = new TextEncoder().encode(
            JSON.stringify({ type: "call_updated" }),
          );
          await room.localParticipant.publishData(payload, { reliable: true });
        } catch (err) {
          console.warn("publishData failed:", err);
        }
      }
    } catch (err) {
      console.error("UpdateCall failed:", err);
    } finally {
      setUpdatingVisibility(false);
    }
  };

  useEffect(() => {
    void fetchCallInfo();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [callName]);

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
      const res = latestJoinResRef.current;
      if (!callName || !res) return;
      void callClient
        .heartbeatCall({
          name: callName,
          guestIdentity: isAuthenticated ? "" : res.identity,
        })
        .catch((err) => {
          console.warn("HeartbeatCall failed:", err);
        });
    }, HEARTBEAT_INTERVAL_MS);
  };

  // Fire-and-forget LeaveCall that survives page unload. Regular fetch from the
  // Connect client would be aborted by the browser; keepalive lets the request
  // flush before the tab dies.
  const sendLeaveBeacon = () => {
    const res = latestJoinResRef.current;
    if (!callName || !res) return;
    const token = tokenStore.getAccessToken();
    const body = JSON.stringify({
      name: callName,
      guestIdentity: isAuthenticated ? "" : res.identity,
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
    if (!callName) return;
    try {
      const res = await callClient.joinCall({
        name: callName,
        guestDisplayName: "",
      });
      if (res.identity !== latestJoinResRef.current?.identity) {
        console.warn("Token refresh returned different identity; ignoring");
        return;
      }
      latestJoinResRef.current = res;
      scheduleTokenRefresh();
    } catch (err) {
      console.error("Token refresh failed:", err);
    }
  };

  const attachRoomListeners = (room: Room) => {
    room
      .on(RoomEvent.ParticipantConnected, () => {
        rerender();
        void fetchCallInfo({ silent: true });
      })
      .on(RoomEvent.ParticipantDisconnected, () => {
        rerender();
        void fetchCallInfo({ silent: true });
      })
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
      .on(
        RoomEvent.DataReceived,
        (payload: Uint8Array, participant?: RemoteParticipant) => {
          try {
            const msg = JSON.parse(new TextDecoder().decode(payload));
            if (msg?.type === "call_updated") {
              void fetchCallInfo({ silent: true });
              return;
            }
            if (
              msg?.type === "chat_message" &&
              typeof msg.text === "string" &&
              participant
            ) {
              setChatMessages((prev) => [
                ...prev,
                {
                  id: makeMessageID(participant.identity),
                  identity: participant.identity,
                  displayName: participant.name || participant.identity,
                  text: msg.text,
                },
              ]);
            }
          } catch {
            // ignore malformed payloads
          }
        },
      )
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
      setJoinError("接続が切れました。再参加してください");
      void fetchCallInfo();
      rerender();
      return;
    }
    try {
      const fresh = await callClient.joinCall({
        name: callName,
        guestDisplayName: "",
      });
      latestJoinResRef.current = fresh;
      const newRoom = new Room();
      attachRoomListeners(newRoom);
      await newRoom.connect(fresh.url, fresh.accessToken);
      await newRoom.localParticipant.setMicrophoneEnabled(true);
      roomRef.current = newRoom;
      scheduleTokenRefresh();
      rerender();
    } catch (err) {
      console.error("Auto-reconnect failed:", err);
      setConnected(false);
      setJoinError("再接続に失敗しました。再参加してください");
      void fetchCallInfo();
      rerender();
    }
  };

  const handleJoin = async () => {
    if (!callName) return;
    if (!isAuthenticated && !guestDisplayName.trim()) return;
    setConnecting(true);
    setJoinError(null);
    try {
      if (!isAuthenticated) {
        localStorage.setItem(GUEST_DISPLAY_NAME_KEY, guestDisplayName);
      }
      const res = await callClient.joinCall({
        name: callName,
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
      startHeartbeat();
      void fetchCallInfo({ silent: true });
      rerender();
    } catch (err) {
      console.error("JoinCall failed:", err);
      setJoinError(err instanceof Error ? err.message : String(err));
    } finally {
      setConnecting(false);
    }
  };

  const handleLeave = async () => {
    intentionalLeaveRef.current = true;
    clearRefreshTimer();
    clearHeartbeatTimer();
    const res = latestJoinResRef.current;
    if (callName && res) {
      try {
        await callClient.leaveCall({
          name: callName,
          guestIdentity: isAuthenticated ? "" : res.identity,
        });
      } catch (err) {
        console.warn("LeaveCall failed:", err);
      }
    }
    await roomRef.current?.disconnect();
    roomRef.current = null;
    latestJoinResRef.current = null;
    setConnected(false);
    void fetchCallInfo({ silent: true });
    rerender();
  };

  useEffect(() => {
    const onPageHide = () => {
      sendLeaveBeacon();
    };
    window.addEventListener("pagehide", onPageHide);
    return () => {
      window.removeEventListener("pagehide", onPageHide);
      intentionalLeaveRef.current = true;
      clearRefreshTimer();
      clearHeartbeatTimer();
      sendLeaveBeacon();
      roomRef.current?.disconnect();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (loadingCall) {
    return (
      <div className="w-full min-h-full flex items-center justify-center">
        <p className="text-sm text-muted-foreground">読み込み中...</p>
      </div>
    );
  }

  if (callLoadError || !callInfo?.call) {
    return (
      <div className="w-full min-h-full flex flex-col items-center justify-center px-6 gap-2">
        <p className="text-sm text-destructive">
          {callLoadError ?? "通話を取得できませんでした"}
        </p>
      </div>
    );
  }

  const call = callInfo.call;
  const participants: CallParticipant[] = callInfo.participants;
  const roomParticipantCount = roomRef.current
    ? roomRef.current.remoteParticipants.size + 1
    : 0;
  const isHost =
    myUserName !== null && call.host !== undefined && myUserName === call.host.name;

  return (
    <div className="w-full min-h-full flex flex-col items-center justify-center px-6 py-8">
      <div className="w-full max-w-md flex flex-col gap-6">
        <div className="text-center">
          <p className="text-xs text-muted-foreground">
            {VISIBILITY_LABELS[call.visibility]}
          </p>
          <p className="text-sm font-mono break-all">{callId}</p>
          <p className="text-[10px] text-muted-foreground truncate">
            host: {call.host?.displayName ?? "unknown"}
            {call.host?.customId ? ` @${call.host.customId}` : ""}
          </p>
        </div>

        {!connected ? (
          call.visibility === CallVisibility.USERS_ONLY && !isAuthenticated ? (
            <div className="flex flex-col gap-3">
              <p className="text-sm text-muted-foreground text-center">
                この通話はログインしているユーザーのみ参加できます
              </p>
              <Link
                to={`/login?returnTo=${encodeURIComponent(`/calls/${callId}`)}`}
                className="w-full"
              >
                <Button className="w-full h-11">ログインして参加</Button>
              </Link>
            </div>
          ) : (
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
              {!isAuthenticated && (
                <Link
                  to={`/login?returnTo=${encodeURIComponent(`/calls/${callId}`)}`}
                  className="text-xs text-muted-foreground text-center underline underline-offset-2 hover:text-foreground"
                >
                  ログインして参加
                </Link>
              )}
            </div>
          )
        ) : (
          <Button
            variant="destructive"
            onClick={handleLeave}
            className="w-full h-11"
          >
            退出
          </Button>
        )}

        {isHost && (
          <div className="rounded-md border p-4 flex flex-col gap-3">
            <div className="flex flex-col gap-2">
              <p className="text-xs text-muted-foreground">
                可視性（ホストのみ変更可能）
              </p>
              <div className="flex flex-col gap-1">
                {UPDATABLE_VISIBILITIES.map((v) => (
                  <label
                    key={v}
                    className="flex items-center gap-2 text-sm cursor-pointer"
                  >
                    <input
                      type="radio"
                      name="call-visibility"
                      checked={call.visibility === v}
                      disabled={updatingVisibility}
                      onChange={() => void handleUpdateVisibility(v)}
                    />
                    {VISIBILITY_LABELS[v]}
                  </label>
                ))}
              </div>
            </div>
            {call.visibility !== CallVisibility.LOCKED && (
              <Button
                variant="outline"
                size="sm"
                onClick={() => void copyInviteLink()}
              >
                {inviteCopied ? "コピーしました" : "招待リンクをコピー"}
              </Button>
            )}
          </div>
        )}

        <div className="rounded-md border p-4">
          <p className="text-xs text-muted-foreground mb-2">
            参加者（{connected ? `接続中: ${roomParticipantCount}` : "未接続"}）
          </p>
          {participants.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              まだ誰も参加していません
            </p>
          ) : (
            <ul className="flex flex-col gap-1 text-sm">
              {participants.map((p) => (
                <li key={p.name} className="flex items-center gap-2">
                  <span
                    className={`size-1.5 rounded-full ${
                      p.isCurrentlyConnected
                        ? "bg-green-500"
                        : "bg-muted-foreground/30"
                    }`}
                  />
                  <span
                    className={
                      p.isCurrentlyConnected ? "" : "text-muted-foreground"
                    }
                  >
                    {p.user?.displayName || p.displayName}
                  </span>
                  {p.user && (
                    <span className="text-xs text-muted-foreground">
                      @{p.user.customId}
                    </span>
                  )}
                </li>
              ))}
            </ul>
          )}
        </div>

        {connected && (
          <div className="rounded-md border p-4 flex flex-col gap-2">
            <p className="text-xs text-muted-foreground">
              チャット（通話中のみ、履歴なし）
            </p>
            <div
              ref={chatScrollRef}
              className="flex flex-col gap-1 max-h-48 overflow-y-auto text-sm"
            >
              {chatMessages.length === 0 ? (
                <p className="text-xs text-muted-foreground">
                  まだメッセージはありません
                </p>
              ) : (
                chatMessages.map((m) => (
                  <div key={m.id} className="leading-snug">
                    <span className="text-xs font-medium">{m.displayName}</span>
                    <span className="ml-2 whitespace-pre-wrap wrap-break-word">
                      {m.text}
                    </span>
                  </div>
                ))
              )}
            </div>
            <div className="flex gap-2">
              <input
                className="flex-1 h-9 px-2 rounded-md border border-input bg-background text-sm"
                placeholder="メッセージを入力"
                value={chatInput}
                onChange={(e) => setChatInput(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter" && !e.shiftKey) {
                    e.preventDefault();
                    void sendChat();
                  }
                }}
              />
              <Button
                size="sm"
                onClick={() => void sendChat()}
                disabled={!chatInput.trim()}
              >
                送信
              </Button>
            </div>
          </div>
        )}

        {joinError && (
          <p className="text-sm text-destructive text-center">{joinError}</p>
        )}
      </div>
      <div ref={audioContainerRef} className="hidden" aria-hidden />
    </div>
  );
}
