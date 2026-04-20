import { useState, useRef, useEffect, useCallback } from "react";
import { useParams, Link } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogAction,
  AlertDialogCancel,
} from "~/components/ui/alert-dialog";
import {
  ArrowLeft,
  SendHorizontal,
  Phone,
  Video,
  ImagePlus,
  Info,
} from "lucide-react";
import { usePrivateCall } from "~/components/private-call-context";
import { useAnyCallActive } from "~/components/use-any-call-active";
import {
  listMessages,
  sendMessage,
  markRoomAsRead,
  addMessageReaction,
  removeMessageReaction,
  type ChatMessage,
  type Room,
  listRooms,
} from "~/lib/api";

function formatMessageTime(dateStr: string): string {
  const date = new Date(dateStr);
  return `${date.getHours()}:${String(date.getMinutes()).padStart(2, "0")}`;
}

function shouldShowTimestamp(current: ChatMessage, previous?: ChatMessage): boolean {
  if (!previous) return true;
  const diff = new Date(current.createTime).getTime() - new Date(previous.createTime).getTime();
  return diff > 10 * 60_000;
}

type HeartAnim = { id: number; x: number; y: number };

const PAGE_SIZE = 20;

export default function ChatDetail() {
  const { id } = useParams<{ id: string }>();
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [room, setRoom] = useState<Room | null>(null);
  const [loading, setLoading] = useState(true);
  const [nextPageToken, setNextPageToken] = useState<string | undefined>();
  const [loadingMore, setLoadingMore] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);
  const topSentinelRef = useRef<HTMLDivElement>(null);
  const isInitialLoad = useRef(true);
  const { startCall } = usePrivateCall();
  const { isInCall, currentCallName, endCurrentCall } = useAnyCallActive();
  const [showConfirm, setShowConfirm] = useState(false);
  const [pendingCallType, setPendingCallType] = useState<"audio" | "video">("audio");
  const [heartAnims, setHeartAnims] = useState<HeartAnim[]>([]);
  const heartIdRef = useRef(0);
  const tapTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const tapCountRef = useRef<Record<string, number>>({});

  useEffect(() => {
    if (!id) return;
    isInitialLoad.current = true;

    Promise.all([
      listRooms().then((res) => {
        const found = res.rooms?.find((r) => r.id === id);
        if (found) setRoom(found);
      }),
      listMessages(id, { pageSize: PAGE_SIZE }).then((res) => {
        const msgs = [...(res.messages ?? [])].reverse();
        setMessages(msgs);
        setNextPageToken(res.nextPageToken || undefined);
      }),
    ]).finally(() => setLoading(false));

    markRoomAsRead(id).catch(() => {});
  }, [id]);

  // 初回ロード完了時のみ最下部へスクロール
  useEffect(() => {
    if (!loading && isInitialLoad.current) {
      isInitialLoad.current = false;
      scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight });
    }
  }, [loading]);

  // 新規送信メッセージのみ最下部へスクロール（messages末尾追加時）
  const prevLengthRef = useRef(0);
  useEffect(() => {
    const cur = messages.length;
    const prev = prevLengthRef.current;
    prevLengthRef.current = cur;
    // 末尾追加（送信）の場合のみスクロール
    if (cur > prev && prev > 0 && !loadingMore) {
      scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight, behavior: "smooth" });
    }
  }, [messages, loadingMore]);

  // 上スクロール検知で古いメッセージをロード
  const loadOlderMessages = useCallback(async () => {
    if (!id || !nextPageToken || loadingMore) return;
    setLoadingMore(true);
    const scrollEl = scrollRef.current;
    const prevScrollHeight = scrollEl?.scrollHeight ?? 0;
    try {
      const res = await listMessages(id, { pageSize: PAGE_SIZE, pageToken: nextPageToken });
      const older = [...(res.messages ?? [])].reverse();
      setMessages((prev) => [...older, ...prev]);
      setNextPageToken(res.nextPageToken || undefined);
      // スクロール位置を保持（新しく追加された分だけ下にずらす）
      requestAnimationFrame(() => {
        if (scrollEl) {
          scrollEl.scrollTop = scrollEl.scrollHeight - prevScrollHeight;
        }
      });
    } catch {
    } finally {
      setLoadingMore(false);
    }
  }, [id, nextPageToken, loadingMore]);

  // sentinel が見えたら古いメッセージをロード
  useEffect(() => {
    if (!nextPageToken || loadingMore) return;
    const el = topSentinelRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      (entries) => { if (entries[0].isIntersecting) loadOlderMessages(); },
      { root: scrollRef.current, rootMargin: "100px" }
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, [nextPageToken, loadingMore, loadOlderMessages]);

  const isGroup = room?.roomType === "group";
  const participantName = room?.otherUser?.displayName ?? room?.name ?? "チャット";
  const participantCustomId = room?.otherUser?.customId;

  const handleDoubleTap = useCallback((msg: ChatMessage, e: React.MouseEvent | React.TouchEvent) => {
    const msgId = msg.id;
    tapCountRef.current[msgId] = (tapCountRef.current[msgId] ?? 0) + 1;

    if (tapTimerRef.current) clearTimeout(tapTimerRef.current);
    tapTimerRef.current = setTimeout(async () => {
      if (tapCountRef.current[msgId] >= 2) {
        const heartReaction = msg.reactions?.find((r) => r.emoji === "❤️");
        const isMine = heartReaction?.isMine ?? false;

        const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
        const clientX = "touches" in e ? e.changedTouches[0]?.clientX ?? rect.left + rect.width / 2 : (e as React.MouseEvent).clientX;
        const clientY = "touches" in e ? e.changedTouches[0]?.clientY ?? rect.top + rect.height / 2 : (e as React.MouseEvent).clientY;
        const x = clientX - rect.left;
        const y = clientY - rect.top;

        const animId = ++heartIdRef.current;
        setHeartAnims((prev) => [...prev, { id: animId, x, y }]);
        setTimeout(() => setHeartAnims((prev) => prev.filter((a) => a.id !== animId)), 800);

        try {
          const apiCall = isMine ? removeMessageReaction : addMessageReaction;
          const res = await apiCall(msgId, "❤️");
          setMessages((prev) => prev.map((m) => m.id === msgId ? res.message : m));
        } catch {}
      }
      tapCountRef.current[msgId] = 0;
    }, 300);
  }, []);

  const handleSend = async () => {
    const trimmed = input.trim();
    if (!trimmed || !id) return;
    setInput("");
    try {
      const res = await sendMessage(id, trimmed);
      setMessages((prev) => [...prev, res.message]);
    } catch {
      setInput(trimmed);
    }
  };

  const handleStartCall = (type: "audio" | "video") => {
    if (!id || !room) return;
    const participant = {
      name: participantName,
      customId: participantCustomId ?? id,
      avatarUrl: "",
      status: "offline" as const,
    };
    if (isInCall) {
      setPendingCallType(type);
      setShowConfirm(true);
      return;
    }
    startCall(id, participant, type);
  };

  const handleConfirmSwitch = () => {
    if (!id || !room) return;
    const participant = {
      name: participantName,
      customId: participantCustomId ?? id,
      avatarUrl: "",
      status: "offline" as const,
    };
    endCurrentCall();
    startCall(id, participant, pendingCallType);
    setShowConfirm(false);
  };

  if (loading) {
    return (
      <div className="w-full h-full flex items-center justify-center text-muted-foreground">
        読み込み中...
      </div>
    );
  }

  return (
    <>
      <div className="w-full h-full flex flex-col">
        {/* Header */}
        <div className="shrink-0 border-b bg-background/60 backdrop-blur-lg z-10">
          <div className="flex items-center gap-3 px-4 h-16">
            <Button variant="ghost" size="icon" onClick={() => history.back()}>
              <ArrowLeft className="size-5" />
            </Button>
            {isGroup ? (
              <Link to={`/chat/${id}/info`} className="flex items-center gap-3 flex-1 min-w-0">
                <Avatar className="size-9">
                  <AvatarFallback><span className="text-xs">G</span></AvatarFallback>
                </Avatar>
                <div className="min-w-0">
                  <p className="text-sm font-bold truncate">{participantName}</p>
                  <p className="text-xs text-muted-foreground">{room?.members?.length ?? 0}人のメンバー</p>
                </div>
              </Link>
            ) : participantCustomId ? (
              <Link to={`/users/${participantCustomId}`} className="flex items-center gap-3 flex-1 min-w-0">
                <Avatar className="size-9">
                  <AvatarFallback>{participantName.slice(0, 2)}</AvatarFallback>
                </Avatar>
                <div className="min-w-0">
                  <p className="text-sm font-bold truncate">{participantName}</p>
                </div>
              </Link>
            ) : (
              <div className="flex items-center gap-3 flex-1 min-w-0">
                <Avatar className="size-9">
                  <AvatarFallback>{participantName.slice(0, 2)}</AvatarFallback>
                </Avatar>
                <p className="text-sm font-bold truncate">{participantName}</p>
              </div>
            )}
            <div className="flex items-center gap-2">
              {isGroup ? (
                <Link to={`/chat/${id}/info`} className="p-2 inline-flex items-center justify-center text-muted-foreground hover:text-foreground transition-colors">
                  <Info className="size-5" />
                </Link>
              ) : (
                <>
                  <Button variant="ghost" size="icon" onClick={() => handleStartCall("audio")}>
                    <Phone className="size-5" />
                  </Button>
                  <Button variant="ghost" size="icon" onClick={() => handleStartCall("video")}>
                    <Video className="size-5" />
                  </Button>
                </>
              )}
            </div>
          </div>
        </div>

        {/* Messages */}
        <div ref={scrollRef} className="flex-1 overflow-y-auto px-4 py-4">
          {/* 上スクロール sentinel */}
          <div ref={topSentinelRef} className="h-1" />
          {loadingMore && (
            <div className="flex justify-center py-2 text-xs text-muted-foreground">読み込み中...</div>
          )}
          <div className="flex flex-col gap-1">
            {messages.map((msg, i) => {
              const showTime = shouldShowTimestamp(msg, messages[i - 1]);
              const prevMsg = messages[i - 1];
              const sameSender = isGroup
                ? prevMsg?.sender?.id === msg.sender?.id
                : prevMsg?.isMine === msg.isMine;
              const isConsecutive = i > 0 && !showTime && prevMsg?.isMine === msg.isMine && sameSender;
              return (
                <div key={msg.id}>
                  {showTime && (
                    <div className="flex justify-center my-3">
                      <span className="text-xs text-muted-foreground bg-muted/50 px-2 py-0.5 rounded-full">
                        {formatMessageTime(msg.createTime)}
                      </span>
                    </div>
                  )}
                  <div className={`group flex items-center ${msg.isMine ? "justify-end" : "justify-start"} ${isConsecutive ? "mt-0.5" : "mt-2"}`}>
                    {!msg.isMine && !isConsecutive && (
                      <Avatar className="size-8 mr-2 shrink-0 self-start">
                        <AvatarFallback className="text-xs">
                          {(msg.sender?.displayName ?? participantName).slice(0, 2)}
                        </AvatarFallback>
                      </Avatar>
                    )}
                    {!msg.isMine && isConsecutive && <div className="size-8 mr-2 shrink-0" />}

                    {/* 自分のメッセージ: ホバーボタンをバブルの左に表示 */}
                    {msg.isMine && (
                      <div className="hidden group-hover:flex items-center gap-0.5 mr-1">
                        <button
                          onClick={async () => {
                            const isMine = msg.reactions?.find((r) => r.emoji === "❤️")?.isMine ?? false;
                            try {
                              const res = await (isMine ? removeMessageReaction : addMessageReaction)(msg.id, "❤️");
                              setMessages((prev) => prev.map((m) => m.id === msg.id ? res.message : m));
                            } catch {}
                          }}
                          className="p-1.5 rounded-full text-muted-foreground hover:text-foreground hover:bg-muted transition-colors text-sm"
                          title="リアクション"
                        >
                          🤍
                        </button>
                        <button className="p-1.5 rounded-full text-muted-foreground hover:text-foreground hover:bg-muted transition-colors text-xs font-bold" title="その他">
                          •••
                        </button>
                      </div>
                    )}

                    {/* バブル + リアクションバッジ */}
                    <div className="flex flex-col max-w-[75%]" style={{ alignItems: msg.isMine ? "flex-end" : "flex-start" }}>
                      {isGroup && !msg.isMine && !isConsecutive && (
                        <p className="text-xs text-muted-foreground mb-0.5 ml-1">{msg.sender?.displayName}</p>
                      )}
                      <div
                        className={`relative w-fit px-3 py-2 text-sm whitespace-pre-wrap break-words cursor-pointer select-none touch-manipulation ${msg.isMine ? "bg-blue-500 text-white rounded-2xl rounded-tr-md" : "bg-muted rounded-2xl rounded-tl-md"}`}
                        onClick={(e) => handleDoubleTap(msg, e)}
                        onTouchEnd={(e) => handleDoubleTap(msg, e)}
                      >
                        {msg.content}
                        {heartAnims.map((a) => (
                          <span
                            key={a.id}
                            className="pointer-events-none absolute text-2xl animate-heart-float"
                            style={{ left: a.x, top: a.y, transform: "translate(-50%, -50%)" }}
                          >
                            ❤️
                          </span>
                        ))}
                      </div>
                      {(msg.reactions?.length ?? 0) > 0 && (
                        <div className="flex gap-1 mt-1 flex-wrap" style={{ justifyContent: msg.isMine ? "flex-end" : "flex-start" }}>
                          {msg.reactions!.map((r) => (
                            <button
                              key={r.emoji}
                              onClick={async () => {
                                try {
                                  const res = await (r.isMine ? removeMessageReaction : addMessageReaction)(msg.id, r.emoji);
                                  setMessages((prev) => prev.map((m) => m.id === msg.id ? res.message : m));
                                } catch {}
                              }}
                              className={`flex items-center gap-0.5 text-xs px-1.5 py-0.5 rounded-full border transition-colors ${r.isMine ? "bg-blue-100 border-blue-300 dark:bg-blue-900/30 dark:border-blue-700" : "bg-muted border-transparent"}`}
                            >
                              <span>{r.emoji}</span>
                              <span className="text-muted-foreground">{r.count}</span>
                            </button>
                          ))}
                        </div>
                      )}
                    </div>

                    {/* 相手のメッセージ: ホバーボタンをバブルの右に表示 */}
                    {!msg.isMine && (
                      <div className="hidden group-hover:flex items-center gap-0.5 ml-1">
                        <button
                          onClick={async () => {
                            const isMine = msg.reactions?.find((r) => r.emoji === "❤️")?.isMine ?? false;
                            try {
                              const res = await (isMine ? removeMessageReaction : addMessageReaction)(msg.id, "❤️");
                              setMessages((prev) => prev.map((m) => m.id === msg.id ? res.message : m));
                            } catch {}
                          }}
                          className="p-1.5 rounded-full text-muted-foreground hover:text-foreground hover:bg-muted transition-colors text-sm"
                          title="リアクション"
                        >
                          🤍
                        </button>
                        <button className="p-1.5 rounded-full text-muted-foreground hover:text-foreground hover:bg-muted transition-colors text-xs font-bold" title="その他">
                          •••
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Input */}
        <div className="shrink-0 border-t bg-background px-4 py-3">
          <div className="flex gap-2">
            <button
              className="shrink-0 p-2 text-muted-foreground hover:text-foreground transition-colors"
              onClick={() => {}}
            >
              <ImagePlus className="size-5" />
            </button>
            <textarea
              placeholder="メッセージを入力..."
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter" && !e.shiftKey && !e.nativeEvent.isComposing) {
                  e.preventDefault();
                  handleSend();
                }
              }}
              rows={1}
              className="flex-1 resize-none rounded-md border border-input bg-background px-3 py-2 text-sm min-h-10 max-h-32 field-sizing-content placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            />
            <button
              className="shrink-0 px-1 text-blue-500 disabled:text-muted-foreground transition-colors"
              onClick={handleSend}
              disabled={!input.trim()}
            >
              <SendHorizontal className="size-5" />
            </button>
          </div>
        </div>
      </div>

      <AlertDialog open={showConfirm} onOpenChange={(open) => { if (!open) setShowConfirm(false); }}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>通話を切り替えますか？</AlertDialogTitle>
            <AlertDialogDescription>
              「{currentCallName}」を終了して、{participantName}との通話を開始しますか？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>キャンセル</AlertDialogCancel>
            <AlertDialogAction onClick={handleConfirmSwitch}>切り替える</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
