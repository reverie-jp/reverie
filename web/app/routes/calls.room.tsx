import { useEffect, useRef, useState } from "react";
import { Link, useNavigate, useParams } from "react-router";
import type {
  CallParticipant,
  GetCallResponse,
} from "~/lib/gen/call/v1/call_pb";
import { CallVisibility } from "~/lib/gen/call/v1/call_pb";
import { callClient } from "~/lib/api-client";
import { formatCall, parseCallParticipant } from "~/lib/resource-name";
import {
  CALL_DEFAULT_VOLUME,
  CALL_UPDATED_EVENT,
  useCall,
} from "~/lib/call-context";
import { Button } from "~/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { ParticipantTile } from "~/components/participant-tile";
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
} from "~/components/ui/drawer";
import { Slider } from "~/components/ui/slider";
import { GoogleLoginButton } from "~/components/google-login-button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "~/components/ui/dialog";
import type { CallBan } from "~/lib/gen/call/v1/call_pb";

const GUEST_DISPLAY_NAME_KEY = "reverie.guest_display_name";

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

function formatDuration(ms: number): string {
  const total = Math.max(0, Math.floor(ms / 1000));
  const h = Math.floor(total / 3600);
  const m = Math.floor((total % 3600) / 60);
  const s = total % 60;
  const mm = String(m).padStart(2, "0");
  const ss = String(s).padStart(2, "0");
  return h > 0 ? `${h}:${mm}:${ss}` : `${m}:${ss}`;
}

function formatClockTime(ms: number): string {
  const d = new Date(ms);
  const h = String(d.getHours()).padStart(2, "0");
  const m = String(d.getMinutes()).padStart(2, "0");
  return `${h}:${m}`;
}

export default function CallRoomRoute() {
  const { callId } = useParams<{ callId: string }>();
  const callName = callId ? formatCall(callId) : "";
  const navigate = useNavigate();
  const call = useCall();

  const [guestDisplayName, setGuestDisplayName] = useState("");
  const [callInfo, setCallInfo] = useState<GetCallResponse | null>(null);
  const [loadingCall, setLoadingCall] = useState(true);
  const [callLoadError, setCallLoadError] = useState<string | null>(null);
  const [updatingVisibility, setUpdatingVisibility] = useState(false);
  const [inviteCopied, setInviteCopied] = useState(false);
  const [moderatingIdentity, setModeratingIdentity] = useState<string | null>(
    null,
  );
  const [chatInput, setChatInput] = useState("");
  const [drawerIdentity, setDrawerIdentity] = useState<string | null>(null);
  const [now, setNow] = useState<number>(() => Date.now());
  const [leaveDialogOpen, setLeaveDialogOpen] = useState(false);
  const [transferMode, setTransferMode] = useState(false);
  const [leavingAction, setLeavingAction] = useState(false);
  const [banListOpen, setBanListOpen] = useState(false);
  const [bans, setBans] = useState<CallBan[]>([]);
  const [bansLoading, setBansLoading] = useState(false);
  const [bansError, setBansError] = useState<string | null>(null);
  const [unbanningName, setUnbanningName] = useState<string | null>(null);
  const [resolvedUnmuteMsgIds, setResolvedUnmuteMsgIds] = useState<Set<string>>(
    new Set(),
  );
  const chatScrollRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const id = window.setInterval(() => setNow(Date.now()), 1000);
    return () => window.clearInterval(id);
  }, []);

  const connected = call.connected && call.callId === callId;
  const connecting = call.connecting && call.callId === callId;
  const joinError = call.callId === callId ? call.joinError : null;

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
        guestIdentity: connected && call.identity ? call.identity : "",
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
    setGuestDisplayName(localStorage.getItem(GUEST_DISPLAY_NAME_KEY) ?? "");
  }, []);

  useEffect(() => {
    void fetchCallInfo();
    const onUpdated = () => void fetchCallInfo({ silent: true });
    window.addEventListener(CALL_UPDATED_EVENT, onUpdated);
    return () => window.removeEventListener(CALL_UPDATED_EVENT, onUpdated);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [callName]);

  // Refetch when participants connect/disconnect in LiveKit.
  useEffect(() => {
    if (!connected) return;
    void fetchCallInfo({ silent: true });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [call.participantTick]);

  useEffect(() => {
    if (chatScrollRef.current) {
      chatScrollRef.current.scrollTop = chatScrollRef.current.scrollHeight;
    }
  }, [call.chatMessages]);

  const handleJoin = async () => {
    if (!callId) return;
    if (!call.isAuthenticated && !guestDisplayName.trim()) return;
    if (call.callId && call.callId !== callId) {
      // User is already in a different call. Leave current first.
      await call.leave();
    }
    const res = await call.join(callId, guestDisplayName);
    if (res.ok) {
      await fetchCallInfo({ silent: true });
    }
  };

  const handleLeave = async () => {
    await call.leave();
    await fetchCallInfo({ silent: true });
  };

  const handleLeaveClick = () => {
    if (!isHost) {
      void handleLeave();
      return;
    }
    const activeOthers = participants.filter((p) => {
      if (!p.isCurrentlyConnected) return false;
      const isSelf =
        (call.myUserName && p.user?.name === call.myUserName) ||
        (!call.isAuthenticated &&
          call.identity ===
            parseCallParticipant(p.name).identity);
      return !isSelf;
    });
    if (activeOthers.length === 0) {
      void handleLeave();
      return;
    }
    setTransferMode(false);
    setLeaveDialogOpen(true);
  };

  const handleEndAndLeave = async () => {
    if (!callName) return;
    setLeavingAction(true);
    try {
      await callClient.endCall({ name: callName });
      // Host's own LiveKit connection receives ROOM_DELETED; the call
      // context's Disconnected handler performs cleanup (no explicit
      // call.leave() needed — that would race with the disconnect event).
      await fetchCallInfo({ silent: true });
      setLeaveDialogOpen(false);
    } catch (err) {
      console.error("EndCall failed:", err);
    } finally {
      setLeavingAction(false);
    }
  };

  const handleTransferAndLeave = async (
    newHostCustomID: string,
    newHostDisplayName: string,
  ) => {
    if (!callName) return;
    setLeavingAction(true);
    try {
      await callClient.transferCallHost({
        name: callName,
        newHost: `users/${newHostCustomID}`,
      });
      // Broadcast before leaving; leave disconnects the LiveKit room so
      // publishing afterwards wouldn't reach anyone.
      await call.publishData({
        type: "host_transferred",
        newHostName: `users/${newHostCustomID}`,
        newHostDisplayName,
      });
      await call.leave();
      await fetchCallInfo({ silent: true });
      setLeaveDialogOpen(false);
    } catch (err) {
      console.error("TransferCallHost failed:", err);
    } finally {
      setLeavingAction(false);
    }
  };

  const fetchBans = async () => {
    if (!callName) return;
    setBansLoading(true);
    setBansError(null);
    try {
      const res = await callClient.listCallBans({
        parent: callName,
        pageSize: 100,
      });
      setBans(res.bans);
    } catch (err) {
      console.error("ListCallBans failed:", err);
      setBansError(err instanceof Error ? err.message : String(err));
    } finally {
      setBansLoading(false);
    }
  };

  const handleOpenBanList = () => {
    setBanListOpen(true);
    void fetchBans();
  };

  const handleTransferHost = async (
    newHostCustomID: string,
    newHostDisplayName: string,
  ) => {
    if (!callName) return;
    try {
      await callClient.transferCallHost({
        name: callName,
        newHost: `users/${newHostCustomID}`,
      });
      await fetchCallInfo({ silent: true });
      call.addSystemMessage(
        `${newHostDisplayName} が新しいホストになりました`,
      );
      await call.publishData({
        type: "host_transferred",
        newHostName: `users/${newHostCustomID}`,
        newHostDisplayName,
      });
    } catch (err) {
      console.error("TransferCallHost failed:", err);
    }
  };

  const handleUnban = async (banName: string) => {
    setUnbanningName(banName);
    try {
      await callClient.unbanCallParticipant({ name: banName });
      setBans((prev) => prev.filter((b) => b.name !== banName));
    } catch (err) {
      console.error("UnbanCallParticipant failed:", err);
    } finally {
      setUnbanningName(null);
    }
  };

  const sendChat = async () => {
    const t = chatInput.trim();
    if (!t) return;
    await call.sendChat(t);
    setChatInput("");
  };

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
      await call.publishData({ type: "call_updated" });
    } catch (err) {
      console.error("UpdateCall failed:", err);
    } finally {
      setUpdatingVisibility(false);
    }
  };

  const handleMute = async (
    identity: string,
    displayName: string,
    muted: boolean,
  ) => {
    if (!callName) return;
    setModeratingIdentity(identity);
    try {
      const name = `${callName}/participants/${identity}`;
      if (muted) {
        await callClient.muteCallParticipant({ name });
      } else {
        await callClient.unmuteCallParticipant({ name });
      }
      const kind = muted ? "host_muted" : "host_unmuted";
      await call.publishData({
        type: kind,
        targetIdentity: identity,
        targetDisplayName: displayName,
      });
      call.addSystemMessage(
        muted
          ? `ホストが ${displayName} をミュートしました`
          : `ホストが ${displayName} のミュートを解除しました`,
      );
      if (!muted) {
        // Mark any existing unmute-request chat entries for this identity as
        // resolved so their "解除" buttons don't reappear on a subsequent
        // re-mute.
        setResolvedUnmuteMsgIds((prev) => {
          const next = new Set(prev);
          for (const m of call.chatMessages) {
            if (m.unmuteRequest?.targetIdentity === identity) {
              next.add(m.id);
            }
          }
          return next;
        });
      }
    } catch (err) {
      console.error("MuteCallParticipant failed:", err);
    } finally {
      setModeratingIdentity(null);
    }
  };

  const handleKick = async (identity: string, displayName: string) => {
    if (!callName) return;
    if (!window.confirm("この参加者を一時的にキックしますか？")) return;
    setModeratingIdentity(identity);
    try {
      await callClient.kickCallParticipant({
        name: `${callName}/participants/${identity}`,
      });
      call.addSystemMessage(`ホストが ${displayName} をキックしました`);
      await call.publishData({
        type: "host_kicked",
        targetDisplayName: displayName,
      });
      await fetchCallInfo({ silent: true });
    } catch (err) {
      console.error("KickCallParticipant failed:", err);
    } finally {
      setModeratingIdentity(null);
    }
  };

  const handleBan = async (identity: string, displayName: string) => {
    if (!callName) return;
    if (!window.confirm("この参加者を永久に追放します。よろしいですか？"))
      return;
    setModeratingIdentity(identity);
    try {
      await callClient.banCallParticipant({
        name: `${callName}/participants/${identity}`,
      });
      call.addSystemMessage(`ホストが ${displayName} を追放しました`);
      await call.publishData({
        type: "host_banned",
        targetDisplayName: displayName,
      });
      await fetchCallInfo({ silent: true });
    } catch (err) {
      console.error("BanCallParticipant failed:", err);
    } finally {
      setModeratingIdentity(null);
    }
  };

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

  const callData = callInfo.call;
  const participants: CallParticipant[] = callInfo.participants;
  const activeParticipantCount = participants.filter(
    (p) => p.isCurrentlyConnected,
  ).length;
  const isHost =
    call.myUserName !== null &&
    callData.host !== undefined &&
    call.myUserName === callData.host.name;
  const isSelfMuted = call.isSelfMuted();
  const callEnded = callData.endTime !== undefined;
  const durationClockMs = callEnded
    ? callData.endTime!.toDate().getTime()
    : now;

  return (
    <div className="w-full min-h-full flex flex-col items-center justify-center px-6 py-8">
      <div className="w-full max-w-md flex flex-col gap-6">
        <div className="text-center">
          <p className="text-xs text-muted-foreground">
            {VISIBILITY_LABELS[callData.visibility]}
          </p>
          <p className="text-sm font-mono break-all">{callId}</p>
          <p className="text-[10px] text-muted-foreground truncate">
            host: {callData.host?.displayName ?? "unknown"}
            {callData.host?.customId ? ` @${callData.host.customId}` : ""}
          </p>
          {callData.createTime && (
            <p className="text-xs text-muted-foreground font-mono mt-1">
              ⏱{" "}
              {formatDuration(
                durationClockMs - callData.createTime.toDate().getTime(),
              )}
              {callEnded && " (終了)"}
            </p>
          )}
        </div>

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
                      checked={callData.visibility === v}
                      disabled={updatingVisibility}
                      onChange={() => void handleUpdateVisibility(v)}
                    />
                    {VISIBILITY_LABELS[v]}
                  </label>
                ))}
              </div>
            </div>
            {callData.visibility !== CallVisibility.LOCKED && (
              <Button
                variant="outline"
                size="sm"
                onClick={() => void copyInviteLink()}
              >
                {inviteCopied ? "コピーしました" : "招待リンクをコピー"}
              </Button>
            )}
            <Button
              variant="outline"
              size="sm"
              onClick={handleOpenBanList}
            >
              追放リストを管理
            </Button>
          </div>
        )}

        <div className="rounded-2xl border border-white/10 bg-gradient-to-b from-white/[0.05] to-white/[0.02] backdrop-blur-xl p-4">
          <div className="mb-4 flex items-baseline justify-between">
            <span className="text-[10px] font-medium uppercase tracking-widest text-muted-foreground/80">
              参加者 {activeParticipantCount}
            </span>
            {participants.length > activeParticipantCount && (
              <span className="text-[10px] text-muted-foreground/60">
                +{participants.length - activeParticipantCount} オフライン
              </span>
            )}
          </div>
          {participants.length === 0 ? (
            <p className="text-sm text-muted-foreground text-center py-4">
              まだ誰も参加していません
            </p>
          ) : (
            <div className="grid grid-cols-4 gap-x-2 gap-y-4">
              {participants.map((p) => {
                const identity = parseCallParticipant(p.name).identity;
                const muted = call.isParticipantMuted(identity);
                const displayName = p.user?.displayName || p.displayName;
                const participantIsHost =
                  !!callData.host &&
                  !!p.user &&
                  p.user.name === callData.host.name;
                return (
                  <ParticipantTile
                    key={p.name}
                    displayName={displayName}
                    avatarUrl={p.user?.avatarUrl}
                    connected={p.isCurrentlyConnected}
                    muted={muted}
                    isHost={participantIsHost}
                    participant={
                      p.isCurrentlyConnected && connected
                        ? call.getLKParticipant(identity)
                        : null
                    }
                    onClick={() => {
                      if (!connected && p.user) {
                        navigate(`/@${p.user.customId}`);
                      } else {
                        setDrawerIdentity(identity);
                      }
                    }}
                  />
                );
              })}
            </div>
          )}
        </div>

        {connected && (
          <div className="rounded-md border p-4 flex flex-col gap-2">
            <p className="text-xs text-muted-foreground">
              チャット（通話中のみ、履歴なし）
            </p>
            <div
              ref={chatScrollRef}
              className="flex flex-col gap-1 h-48 overflow-y-auto text-sm"
            >
              {call.chatMessages.length === 0 ? (
                <p className="text-xs text-muted-foreground">
                  まだメッセージはありません
                </p>
              ) : (
                call.chatMessages.map((m) =>
                  m.system ? (
                    <div
                      key={m.id}
                      className="text-xs text-muted-foreground italic leading-snug flex items-center gap-2 flex-wrap"
                    >
                      <span>
                        <span className="font-mono mr-1">
                          {formatClockTime(m.time)}
                        </span>
                        {m.text}
                      </span>
                      {m.unmuteRequest &&
                        isHost &&
                        call.isParticipantMuted(
                          m.unmuteRequest.targetIdentity,
                        ) &&
                        !resolvedUnmuteMsgIds.has(m.id) && (
                          <Button
                            size="xs"
                            variant="outline"
                            onClick={() => {
                              const req = m.unmuteRequest!;
                              void handleMute(
                                req.targetIdentity,
                                req.targetDisplayName,
                                false,
                              );
                            }}
                          >
                            解除
                          </Button>
                        )}
                    </div>
                  ) : (
                    <div key={m.id} className="leading-snug">
                      <span className="text-[10px] text-muted-foreground font-mono mr-1">
                        {formatClockTime(m.time)}
                      </span>
                      <span className="text-xs font-medium">
                        {m.displayName}
                      </span>
                      <span className="ml-2 whitespace-pre-wrap wrap-break-word">
                        {m.text}
                      </span>
                    </div>
                  ),
                )
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

        {callEnded ? (
          <div className="rounded-md border border-muted-foreground/30 bg-muted/40 px-3 py-3 text-sm text-muted-foreground text-center">
            この通話は終了しました
          </div>
        ) : !connected ? (
          callData.visibility === CallVisibility.USERS_ONLY &&
          !call.isAuthenticated ? (
            <div className="flex flex-col gap-3">
              <p className="text-sm text-muted-foreground text-center">
                この通話はログインしているユーザーのみ参加できます
              </p>
              <GoogleLoginButton returnTo={`/calls/${callId}`} />
            </div>
          ) : (
            <div className="flex flex-col gap-3">
              {!call.isAuthenticated && (
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
                  connecting ||
                  (!call.isAuthenticated && !guestDisplayName.trim())
                }
                className="w-full h-11"
              >
                {connecting ? "参加中..." : "通話に参加"}
              </Button>
              {!call.isAuthenticated && (
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
          <div className="flex flex-col gap-2">
            {isSelfMuted && call.hostMutedMe && (
              <div className="rounded-md border border-destructive/40 bg-destructive/10 px-3 py-2 text-sm text-destructive flex items-center gap-2">
                <span aria-hidden>🔇</span>
                <span>あなたはミュートされています</span>
              </div>
            )}
            <Button
              variant="outline"
              onClick={() => void call.toggleSelfMic()}
              className="w-full h-11"
            >
              {isSelfMuted
                ? call.hostMutedMe
                  ? "解除をリクエスト"
                  : "マイクをオンにする"
                : "マイクをミュート"}
            </Button>
            <Button
              variant="destructive"
              onClick={handleLeaveClick}
              className="w-full h-11"
            >
              退出
            </Button>
          </div>
        )}

        {joinError && (
          <p className="text-sm text-destructive text-center">{joinError}</p>
        )}
      </div>

      <Drawer
        open={drawerIdentity !== null}
        onOpenChange={(open) => {
          if (!open) setDrawerIdentity(null);
        }}
      >
        <DrawerContent>
          {(() => {
            if (!drawerIdentity) return null;
            const target = participants.find(
              (pp) => parseCallParticipant(pp.name).identity === drawerIdentity,
            );
            if (!target) return null;
            const targetIsHost = callData.host?.name === target.user?.name;
            const targetIsSelf =
              (call.myUserName && target.user?.name === call.myUserName) ||
              (!call.isAuthenticated && call.identity === drawerIdentity);
            const canModerateTarget =
              isHost &&
              !targetIsHost &&
              !targetIsSelf &&
              target.isCurrentlyConnected;
            const targetMuted = call.isParticipantMuted(drawerIdentity);
            const targetBusy = moderatingIdentity === drawerIdentity;
            const u = target.user;
            const close = () => setDrawerIdentity(null);
            const joinedAt = u?.createTime?.toDate();
            const displayName = u?.displayName || target.displayName;
            return (
              <>
                <DrawerHeader>
                  <DrawerTitle className="sr-only">
                    {displayName} のプロフィール
                  </DrawerTitle>
                  <DrawerDescription className="sr-only">
                    参加者の詳細とモデレーション操作
                  </DrawerDescription>
                </DrawerHeader>
                <div className="px-4 flex flex-col gap-4">
                  <div className="flex items-center gap-3">
                    <Avatar className="size-14">
                      {u?.avatarUrl && (
                        <AvatarImage src={u.avatarUrl} alt={displayName} />
                      )}
                      <AvatarFallback className="text-lg">
                        {displayName.slice(0, 2)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="min-w-0 flex-1">
                      <p className="font-bold truncate">{displayName}</p>
                      {u ? (
                        <p className="text-sm text-muted-foreground truncate">
                          @{u.customId}
                        </p>
                      ) : (
                        <p className="text-xs text-muted-foreground">ゲスト</p>
                      )}
                    </div>
                    {targetMuted && (
                      <span
                        className="text-xs text-muted-foreground"
                        title="ミュート中"
                      >
                        🔇
                      </span>
                    )}
                  </div>

                  {u?.biography && (
                    <p className="text-sm whitespace-pre-wrap">{u.biography}</p>
                  )}

                  {u && (
                    <div className="flex flex-wrap gap-x-4 gap-y-1 text-sm">
                      <span>
                        <span className="font-bold">{u.followingCount}</span>{" "}
                        <span className="text-muted-foreground">
                          フォロー中
                        </span>
                      </span>
                      <span>
                        <span className="font-bold">{u.followerCount}</span>{" "}
                        <span className="text-muted-foreground">
                          フォロワー
                        </span>
                      </span>
                      {joinedAt && (
                        <span className="text-muted-foreground">
                          {joinedAt.getFullYear()}年{joinedAt.getMonth() + 1}
                          月に登録
                        </span>
                      )}
                    </div>
                  )}

                  {connected &&
                    !targetIsSelf &&
                    target.isCurrentlyConnected && (
                      <div className="flex flex-col gap-2">
                        <div className="flex justify-between text-xs">
                          <span className="text-muted-foreground">
                            音量（あなただけに適用）
                          </span>
                          <span>
                            {Math.round(
                              (call.volumes[drawerIdentity] ??
                                CALL_DEFAULT_VOLUME) * 100,
                            )}
                            %
                          </span>
                        </div>
                        <Slider
                          value={[
                            (call.volumes[drawerIdentity] ??
                              CALL_DEFAULT_VOLUME) * 100,
                          ]}
                          min={0}
                          max={100}
                          step={5}
                          onValueChange={(v) => {
                            const n = Array.isArray(v) ? v[0] : v;
                            call.setVolume(drawerIdentity, n / 100);
                          }}
                        />
                      </div>
                    )}
                </div>
                <DrawerFooter>
                  {u && (
                    <Button
                      variant="outline"
                      onClick={() => {
                        close();
                        navigate(`/@${u.customId}`);
                      }}
                    >
                      プロフィールを見る
                    </Button>
                  )}
                  {canModerateTarget && (
                    <>
                      {!targetMuted && (
                        <Button
                          variant="outline"
                          disabled={targetBusy}
                          onClick={async () => {
                            await handleMute(
                              drawerIdentity,
                              displayName,
                              true,
                            );
                            close();
                          }}
                        >
                          ミュート
                        </Button>
                      )}
                      {u && (
                        <Button
                          variant="outline"
                          disabled={targetBusy}
                          onClick={async () => {
                            if (
                              !window.confirm(
                                `${displayName} にホストを委譲しますか？`,
                              )
                            )
                              return;
                            await handleTransferHost(u.customId, displayName);
                            close();
                          }}
                        >
                          ホストを委譲
                        </Button>
                      )}
                      <Button
                        variant="outline"
                        disabled={targetBusy}
                        onClick={async () => {
                          await handleKick(drawerIdentity, displayName);
                          close();
                        }}
                      >
                        キック
                      </Button>
                      {u && (
                        <Button
                          variant="destructive"
                          disabled={targetBusy}
                          onClick={async () => {
                            await handleBan(drawerIdentity, displayName);
                            close();
                          }}
                        >
                          追放
                        </Button>
                      )}
                    </>
                  )}
                  <Button variant="ghost" onClick={close}>
                    閉じる
                  </Button>
                </DrawerFooter>
              </>
            );
          })()}
        </DrawerContent>
      </Drawer>

      <Dialog
        open={leaveDialogOpen}
        onOpenChange={(open) => {
          if (!open) {
            setLeaveDialogOpen(false);
            setTransferMode(false);
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>ホストの退出</DialogTitle>
            <DialogDescription>
              ホストが退出すると通話は終了します。他のユーザーに委譲することもできます。
            </DialogDescription>
          </DialogHeader>
          {!transferMode ? (
            <div className="flex flex-col gap-2">
              <Button
                variant="outline"
                disabled={leavingAction}
                onClick={() => setTransferMode(true)}
              >
                ホストを委譲して退出
              </Button>
              <Button
                variant="destructive"
                disabled={leavingAction}
                onClick={() => void handleEndAndLeave()}
              >
                通話を終了する
              </Button>
              <Button
                variant="ghost"
                disabled={leavingAction}
                onClick={() => setLeaveDialogOpen(false)}
              >
                キャンセル
              </Button>
            </div>
          ) : (
            <div className="flex flex-col gap-2">
              <p className="text-xs text-muted-foreground">
                新しいホストを選択してください
              </p>
              {(() => {
                const eligible = participants.filter((p) => {
                  if (!p.isCurrentlyConnected) return false;
                  if (!p.user) return false; // guests can't be host
                  if (p.user.name === callData.host?.name) return false; // current host
                  return true;
                });
                if (eligible.length === 0) {
                  return (
                    <p className="text-sm text-muted-foreground">
                      委譲可能な認証済み参加者がいません
                    </p>
                  );
                }
                return eligible.map((p) => (
                  <Button
                    key={p.user!.name}
                    variant="outline"
                    disabled={leavingAction}
                    onClick={() =>
                      void handleTransferAndLeave(
                        p.user!.customId,
                        p.user!.displayName,
                      )
                    }
                    className="justify-between"
                  >
                    <span>{p.user!.displayName}</span>
                    <span className="text-xs text-muted-foreground">
                      @{p.user!.customId}
                    </span>
                  </Button>
                ));
              })()}
              <Button
                variant="ghost"
                disabled={leavingAction}
                onClick={() => setTransferMode(false)}
              >
                戻る
              </Button>
            </div>
          )}
        </DialogContent>
      </Dialog>

      <Dialog
        open={call.becameHost}
        onOpenChange={(open) => {
          if (!open) call.dismissBecameHost();
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>ホストになりました</DialogTitle>
            <DialogDescription>
              あなたがこの通話の新しいホストになりました。
            </DialogDescription>
          </DialogHeader>
          <div className="flex flex-col gap-2">
            <Button onClick={() => call.dismissBecameHost()}>OK</Button>
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={banListOpen} onOpenChange={setBanListOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>追放リスト</DialogTitle>
            <DialogDescription>
              通話から追放されているユーザーです。解除すると再参加できるようになります。
            </DialogDescription>
          </DialogHeader>
          {bansLoading ? (
            <p className="text-sm text-muted-foreground">読み込み中...</p>
          ) : bansError ? (
            <p className="text-sm text-destructive">{bansError}</p>
          ) : bans.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              追放されているユーザーはいません
            </p>
          ) : (
            <ul className="flex flex-col gap-2 max-h-80 overflow-y-auto">
              {bans.map((b) => (
                <li
                  key={b.name}
                  className="flex items-center gap-3 rounded-md border p-2"
                >
                  <Avatar className="size-8">
                    {b.user?.avatarUrl && (
                      <AvatarImage
                        src={b.user.avatarUrl}
                        alt={b.user.displayName}
                      />
                    )}
                    <AvatarFallback className="text-xs">
                      {(b.user?.displayName ?? "??").slice(0, 2)}
                    </AvatarFallback>
                  </Avatar>
                  <div className="min-w-0 flex-1">
                    <p className="text-sm font-medium truncate">
                      {b.user?.displayName ?? "unknown"}
                    </p>
                    {b.user?.customId && (
                      <p className="text-xs text-muted-foreground truncate">
                        @{b.user.customId}
                      </p>
                    )}
                  </div>
                  <Button
                    size="sm"
                    variant="outline"
                    disabled={unbanningName === b.name}
                    onClick={() => void handleUnban(b.name)}
                  >
                    解除
                  </Button>
                </li>
              ))}
            </ul>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
