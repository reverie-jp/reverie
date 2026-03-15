import { useState, useRef, useCallback } from "react";
import { Link } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Input } from "~/components/ui/input";
import { Button } from "~/components/ui/button";
import { BottomNav } from "~/components/bottom-nav";
import {
  Search,
  SquarePen,
  Pin,
  BellOff,
  Trash2,
  MailCheck,
  ChevronRight,
} from "lucide-react";

type OnlineStatus = "online" | "idle" | "offline";

interface ChatThread {
  id: string;
  participant: {
    name: string;
    customId: string;
    avatarUrl: string;
    status: OnlineStatus;
  };
  lastMessage: string;
  lastMessageAt: Date;
  unreadCount: number;
  isLastMessageMine: boolean;
  pinned?: boolean;
  muted?: boolean;
}

const statusColor: Record<OnlineStatus, string> = {
  online: "bg-green-500",
  idle: "bg-yellow-500",
  offline: "bg-gray-400",
};

const chatThreads: ChatThread[] = [
  {
    id: "chat1",
    participant: {
      name: "佐藤花子",
      customId: "hanako_s",
      avatarUrl: "",
      status: "online",
    },
    lastMessage: "了解です！明日の予定について確認しますね。",
    lastMessageAt: new Date(Date.now() - 5 * 60_000),
    unreadCount: 2,
    isLastMessageMine: false,
  },
  {
    id: "chat2",
    participant: {
      name: "鈴木一郎",
      customId: "ichiro_dev",
      avatarUrl: "",
      status: "online",
    },
    lastMessage: "コードレビューありがとうございました！",
    lastMessageAt: new Date(Date.now() - 30 * 60_000),
    unreadCount: 0,
    isLastMessageMine: true,
  },
  {
    id: "chat3",
    participant: {
      name: "山田美咲",
      customId: "misaki_y",
      avatarUrl: "",
      status: "idle",
    },
    lastMessage: "写真送りました！見てみてください 📷",
    lastMessageAt: new Date(Date.now() - 2 * 3_600_000),
    unreadCount: 1,
    isLastMessageMine: false,
  },
  {
    id: "chat4",
    participant: {
      name: "高橋健太",
      customId: "kenta_t",
      avatarUrl: "",
      status: "offline",
    },
    lastMessage: "来週のミーティング、水曜日に変更できますか？",
    lastMessageAt: new Date(Date.now() - 5 * 3_600_000),
    unreadCount: 0,
    isLastMessageMine: false,
  },
  {
    id: "chat5",
    participant: {
      name: "中村悠",
      customId: "yu_nkmr",
      avatarUrl: "",
      status: "offline",
    },
    lastMessage: "おつかれさまです！",
    lastMessageAt: new Date(Date.now() - 86_400_000),
    unreadCount: 0,
    isLastMessageMine: true,
  },
  {
    id: "chat6",
    participant: {
      name: "田中太郎",
      customId: "tanaka",
      avatarUrl: "",
      status: "idle",
    },
    lastMessage: "新しいプロジェクトの件、相談させてください。",
    lastMessageAt: new Date(Date.now() - 2 * 86_400_000),
    unreadCount: 0,
    isLastMessageMine: false,
  },
];

function formatChatTime(date: Date): string {
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMinutes = Math.floor(diffMs / 60_000);
  const diffHours = Math.floor(diffMs / 3_600_000);
  const diffDays = Math.floor(diffMs / 86_400_000);

  if (diffMinutes < 1) return "たった今";
  if (diffMinutes < 60) return `${diffMinutes}分前`;
  if (diffHours < 24) return `${diffHours}時間前`;
  if (diffDays < 7) return `${diffDays}日前`;
  return `${date.getMonth() + 1}/${date.getDate()}`;
}

// --- Swipeable row ---

const SWIPE_THRESHOLD = 60;
const BUTTON_WIDTH = 64;

function SwipeableChatRow({
  thread,
  openSwipeId,
  onSwipeOpen,
  onMarkAsRead,
  onTogglePin,
  onToggleMute,
  onDelete,
}: {
  thread: ChatThread;
  openSwipeId: string | null;
  onSwipeOpen: (id: string | null) => void;
  onMarkAsRead: (id: string) => void;
  onTogglePin: (id: string) => void;
  onToggleMute: (id: string) => void;
  onDelete: (id: string) => void;
}) {
  const hasUnread = thread.unreadCount > 0;
  const actionWidth = BUTTON_WIDTH * (hasUnread ? 4 : 3);

  const startX = useRef(0);
  const currentX = useRef(0);
  const rowRef = useRef<HTMLDivElement>(null);
  const isDragging = useRef(false);

  const isOpen = openSwipeId === thread.id;

  const setTranslate = useCallback((x: number) => {
    if (rowRef.current) {
      rowRef.current.style.transform = `translateX(${x}px)`;
    }
  }, []);

  const handleTouchStart = (e: React.TouchEvent) => {
    startX.current = e.touches[0].clientX;
    isDragging.current = false;
    if (rowRef.current) {
      rowRef.current.style.transition = "none";
    }
  };

  const handleTouchMove = (e: React.TouchEvent) => {
    const diff = e.touches[0].clientX - startX.current;
    const base = isOpen ? -actionWidth : 0;
    const next = Math.min(0, Math.max(-actionWidth, base + diff));
    currentX.current = next;
    setTranslate(next);
    if (Math.abs(diff) > 10) {
      isDragging.current = true;
    }
  };

  const handleTouchEnd = () => {
    if (rowRef.current) {
      rowRef.current.style.transition = "transform 0.2s ease-out";
    }
    if (currentX.current < -SWIPE_THRESHOLD) {
      setTranslate(-actionWidth);
      onSwipeOpen(thread.id);
    } else {
      setTranslate(0);
      onSwipeOpen(null);
    }
  };

  const displayX = isOpen ? -actionWidth : 0;

  return (
    <div className="relative overflow-hidden border-b">
      {/* Action buttons behind */}
      <div
        className="absolute inset-y-0 right-0 flex items-stretch"
        style={{ width: actionWidth }}
      >
        {hasUnread && (
          <button
            className="flex-1 flex flex-col items-center justify-center gap-1 bg-blue-500 text-white text-xs"
            onClick={() => {
              onMarkAsRead(thread.id);
              onSwipeOpen(null);
            }}
          >
            <MailCheck className="size-5" />
            既読
          </button>
        )}
        <button
          className="flex-1 flex flex-col items-center justify-center gap-1 bg-violet-500 text-white text-xs"
          onClick={() => {
            onToggleMute(thread.id);
            onSwipeOpen(null);
          }}
        >
          <BellOff className="size-5" />
          {thread.muted ? "通知ON" : "ミュート"}
        </button>
        <button
          className="flex-1 flex flex-col items-center justify-center gap-1 bg-amber-500 text-white text-xs"
          onClick={() => {
            onTogglePin(thread.id);
            onSwipeOpen(null);
          }}
        >
          <Pin className="size-5" />
          {thread.pinned ? "解除" : "ピン"}
        </button>
        <button
          className="flex-1 flex flex-col items-center justify-center gap-1 bg-red-500 text-white text-xs"
          onClick={() => {
            onDelete(thread.id);
            onSwipeOpen(null);
          }}
        >
          <Trash2 className="size-5" />
          削除
        </button>
      </div>

      {/* Foreground row */}
      <div
        ref={rowRef}
        className="relative bg-background"
        style={{
          transform: `translateX(${displayX}px)`,
          transition: "transform 0.2s ease-out",
        }}
        onTouchStart={handleTouchStart}
        onTouchMove={handleTouchMove}
        onTouchEnd={handleTouchEnd}
      >
        <div className="flex items-center gap-3 px-4 py-3 hover:bg-muted/30 transition-colors">
          <Link
            to={`/users/${thread.participant.customId}`}
            className="relative shrink-0"
            onClick={(e) => {
              if (isDragging.current) e.preventDefault();
              e.stopPropagation();
            }}
          >
            <Avatar className="size-12">
              <AvatarImage
                src={thread.participant.avatarUrl}
                alt={thread.participant.name}
              />
              <AvatarFallback>
                {thread.participant.name.slice(0, 2)}
              </AvatarFallback>
            </Avatar>
            <span
              className={`absolute bottom-0 right-0 size-3.5 rounded-full border-2 border-background ${statusColor[thread.participant.status]}`}
            />
          </Link>
          <Link
            to={`/chat/${thread.id}`}
            className="flex-1 min-w-0"
            onClick={(e) => {
              if (isDragging.current) e.preventDefault();
            }}
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-1.5 min-w-0">
                <span className="text-sm font-bold truncate">
                  {thread.participant.name}
                </span>
                {thread.pinned && (
                  <Pin className="size-3 text-muted-foreground shrink-0 fill-current" />
                )}
                {thread.muted && (
                  <BellOff className="size-3 text-muted-foreground shrink-0" />
                )}
              </div>
              <span className="text-xs text-muted-foreground whitespace-nowrap ml-2">
                {formatChatTime(thread.lastMessageAt)}
              </span>
            </div>
            <div className="flex items-center justify-between mt-0.5">
              <p className="text-sm text-muted-foreground truncate">
                {thread.isLastMessageMine && (
                  <span className="text-muted-foreground">自分: </span>
                )}
                {thread.lastMessage}
              </p>
              {thread.unreadCount > 0 && (
                <span className="ml-2 shrink-0 size-5 rounded-full bg-blue-500 text-white text-xs flex items-center justify-center">
                  {thread.unreadCount}
                </span>
              )}
            </div>
          </Link>
        </div>
      </div>
    </div>
  );
}

// --- Main ---

export default function Chat() {
  const [searchQuery, setSearchQuery] = useState("");
  const [threads, setThreads] = useState(chatThreads);
  const [editing, setEditing] = useState(false);
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [openSwipeId, setOpenSwipeId] = useState<string | null>(null);

  const filteredThreads = threads.filter(
    (thread) =>
      thread.participant.name
        .toLowerCase()
        .includes(searchQuery.toLowerCase()) ||
      thread.participant.customId
        .toLowerCase()
        .includes(searchQuery.toLowerCase()),
  );

  const sortedThreads = [...filteredThreads].sort((a, b) => {
    if (a.pinned && !b.pinned) return -1;
    if (!a.pinned && b.pinned) return 1;
    return 0;
  });

  const handleMarkAsRead = (id: string) => {
    setThreads((prev) =>
      prev.map((t) => (t.id === id ? { ...t, unreadCount: 0 } : t)),
    );
  };

  const handleTogglePin = (id: string) => {
    setThreads((prev) =>
      prev.map((t) => (t.id === id ? { ...t, pinned: !t.pinned } : t)),
    );
  };

  const handleToggleMute = (id: string) => {
    setThreads((prev) =>
      prev.map((t) => (t.id === id ? { ...t, muted: !t.muted } : t)),
    );
  };

  const handleDelete = (id: string) => {
    setThreads((prev) => prev.filter((t) => t.id !== id));
  };

  const toggleSelect = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handleBulkAction = (action: "read" | "pin" | "mute" | "delete") => {
    const ids = selected;
    if (ids.size === 0) return;
    switch (action) {
      case "read":
        setThreads((prev) =>
          prev.map((t) => (ids.has(t.id) ? { ...t, unreadCount: 0 } : t)),
        );
        break;
      case "pin":
        setThreads((prev) =>
          prev.map((t) => (ids.has(t.id) ? { ...t, pinned: !t.pinned } : t)),
        );
        break;
      case "mute":
        setThreads((prev) =>
          prev.map((t) => (ids.has(t.id) ? { ...t, muted: !t.muted } : t)),
        );
        break;
      case "delete":
        setThreads((prev) => prev.filter((t) => !ids.has(t.id)));
        break;
    }
    setSelected(new Set());
  };

  const exitEditing = () => {
    setEditing(false);
    setSelected(new Set());
  };

  return (
    <div className="w-full min-h-full flex flex-col">
      {/* Header */}
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center justify-between px-4 h-14">
          {editing ? (
            <>
              <Button variant="ghost" size="sm" onClick={exitEditing}>
                キャンセル
              </Button>
              <h1 className="text-base font-bold">
                {selected.size > 0
                  ? `${selected.size}件選択中`
                  : "チャットを選択"}
              </h1>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => {
                  if (selected.size === sortedThreads.length) {
                    setSelected(new Set());
                  } else {
                    setSelected(new Set(sortedThreads.map((t) => t.id)));
                  }
                }}
              >
                {selected.size === sortedThreads.length ? "全解除" : "全選択"}
              </Button>
            </>
          ) : (
            <>
              <h1 className="text-lg font-bold">チャット</h1>
              <div className="flex items-center gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setEditing(true)}
                >
                  <SquarePen className="size-5 mr-1" />
                  編集
                </Button>
              </div>
            </>
          )}
        </div>
        {!editing && (
          <div className="px-4 pb-3">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
              <Input
                placeholder="検索"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-9 h-10"
              />
            </div>
          </div>
        )}
      </div>

      {/* Message requests */}
      {!editing && (
        <Link
          to="/chat/requests"
          className="flex items-center justify-between px-4 py-2.5 border-b hover:bg-muted/30 transition-colors"
        >
          <span className="text-sm text-blue-500">メッセージリクエスト (3)</span>
          <ChevronRight className="size-4 text-blue-500" />
        </Link>
      )}

      {/* Chat list */}
      <div className="flex-1">
        {sortedThreads.length === 0 ? (
          <div className="flex items-center justify-center py-12 text-muted-foreground text-sm">
            チャットが見つかりません
          </div>
        ) : editing ? (
          sortedThreads.map((thread) => (
            <button
              key={thread.id}
              className="flex items-center gap-3 px-4 py-3 border-b w-full text-left hover:bg-muted/30 transition-colors"
              onClick={() => toggleSelect(thread.id)}
            >
              <div
                className={`shrink-0 size-5 rounded-full border-2 flex items-center justify-center transition-colors ${
                  selected.has(thread.id)
                    ? "bg-blue-500 border-blue-500"
                    : "border-muted-foreground"
                }`}
              >
                {selected.has(thread.id) && (
                  <svg
                    viewBox="0 0 12 12"
                    className="size-3 text-white"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth={2}
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <polyline points="2,6 5,9 10,3" />
                  </svg>
                )}
              </div>
              <div className="relative shrink-0">
                <Avatar className="size-12">
                  <AvatarImage
                    src={thread.participant.avatarUrl}
                    alt={thread.participant.name}
                  />
                  <AvatarFallback>
                    {thread.participant.name.slice(0, 2)}
                  </AvatarFallback>
                </Avatar>
                <span
                  className={`absolute bottom-0 right-0 size-3.5 rounded-full border-2 border-background ${statusColor[thread.participant.status]}`}
                />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-1.5 min-w-0">
                    <span className="text-sm font-bold truncate">
                      {thread.participant.name}
                    </span>
                    {thread.pinned && (
                      <Pin className="size-3 text-muted-foreground shrink-0 fill-current" />
                    )}
                    {thread.muted && (
                      <BellOff className="size-3 text-muted-foreground shrink-0" />
                    )}
                  </div>
                  <span className="text-xs text-muted-foreground whitespace-nowrap ml-2">
                    {formatChatTime(thread.lastMessageAt)}
                  </span>
                </div>
                <p className="text-sm text-muted-foreground truncate mt-0.5">
                  {thread.isLastMessageMine && "自分: "}
                  {thread.lastMessage}
                </p>
              </div>
            </button>
          ))
        ) : (
          sortedThreads.map((thread) => (
            <SwipeableChatRow
              key={thread.id}
              thread={thread}
              openSwipeId={openSwipeId}
              onSwipeOpen={setOpenSwipeId}
              onMarkAsRead={handleMarkAsRead}
              onTogglePin={handleTogglePin}
              onToggleMute={handleToggleMute}
              onDelete={handleDelete}
            />
          ))
        )}
      </div>

      {/* Bulk action bar (editing mode) */}
      {editing ? (
        <div className="sticky bottom-0 left-0 w-full border-t bg-background px-2 py-2">
          <div className="flex items-center justify-around">
            <button
              className="flex flex-col items-center gap-0.5 text-xs text-muted-foreground disabled:opacity-30"
              disabled={selected.size === 0}
              onClick={() => handleBulkAction("read")}
            >
              <MailCheck className="size-5" />
              既読
            </button>
            <button
              className="flex flex-col items-center gap-0.5 text-xs text-muted-foreground disabled:opacity-30"
              disabled={selected.size === 0}
              onClick={() => handleBulkAction("pin")}
            >
              <Pin className="size-5" />
              ピン
            </button>
            <button
              className="flex flex-col items-center gap-0.5 text-xs text-muted-foreground disabled:opacity-30"
              disabled={selected.size === 0}
              onClick={() => handleBulkAction("mute")}
            >
              <BellOff className="size-5" />
              ミュート
            </button>
            <button
              className="flex flex-col items-center gap-0.5 text-xs text-destructive disabled:opacity-30"
              disabled={selected.size === 0}
              onClick={() => handleBulkAction("delete")}
            >
              <Trash2 className="size-5" />
              削除
            </button>
          </div>
        </div>
      ) : (
        <BottomNav />
      )}
    </div>
  );
}
