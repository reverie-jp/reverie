import { useState, useRef, useCallback, useEffect } from "react";
import { Link, useNavigate } from "react-router";
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
  MessageSquarePlus,
} from "lucide-react";
import {
  listRooms,
  markRoomAsRead,
  pinRoom,
  unpinRoom,
  muteRoom,
  unmuteRoom,
  leaveRoom,
  type Room,
} from "~/lib/api";

function formatChatTime(dateStr: string): string {
  const date = new Date(dateStr);
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

function getRoomDisplayName(room: Room): string {
  if (room.roomType === "direct" && room.otherUser) {
    return room.otherUser.displayName;
  }
  return room.name ?? "グループ";
}

function getRoomCustomId(room: Room): string | null {
  if (room.roomType === "direct" && room.otherUser) {
    return room.otherUser.customId;
  }
  return null;
}

// --- Swipeable row ---

const SWIPE_THRESHOLD = 60;
const BUTTON_WIDTH = 64;

function SwipeableChatRow({
  room,
  openSwipeId,
  onSwipeOpen,
  onMarkAsRead,
  onTogglePin,
  onToggleMute,
  onDelete,
}: {
  room: Room;
  openSwipeId: string | null;
  onSwipeOpen: (id: string | null) => void;
  onMarkAsRead: (id: string) => void;
  onTogglePin: (id: string) => void;
  onToggleMute: (id: string) => void;
  onDelete: (id: string) => void;
}) {
  const hasUnread = room.unreadCount > 0;
  const actionWidth = BUTTON_WIDTH * (hasUnread ? 4 : 3);

  const startX = useRef(0);
  const currentX = useRef(0);
  const rowRef = useRef<HTMLDivElement>(null);
  const isDragging = useRef(false);

  const isOpen = openSwipeId === room.id;

  const setTranslate = useCallback((x: number) => {
    if (rowRef.current) {
      rowRef.current.style.transform = `translateX(${x}px)`;
    }
  }, []);

  const handleTouchStart = (e: React.TouchEvent) => {
    startX.current = e.touches[0].clientX;
    isDragging.current = false;
    if (rowRef.current) rowRef.current.style.transition = "none";
  };

  const handleTouchMove = (e: React.TouchEvent) => {
    const diff = e.touches[0].clientX - startX.current;
    const base = isOpen ? -actionWidth : 0;
    const next = Math.min(0, Math.max(-actionWidth, base + diff));
    currentX.current = next;
    setTranslate(next);
    if (Math.abs(diff) > 10) isDragging.current = true;
  };

  const handleTouchEnd = () => {
    if (rowRef.current) rowRef.current.style.transition = "transform 0.2s ease-out";
    if (currentX.current < -SWIPE_THRESHOLD) {
      setTranslate(-actionWidth);
      onSwipeOpen(room.id);
    } else {
      setTranslate(0);
      onSwipeOpen(null);
    }
  };

  const displayX = isOpen ? -actionWidth : 0;
  const displayName = getRoomDisplayName(room);
  const customId = getRoomCustomId(room);

  return (
    <div className="relative overflow-hidden border-b">
      <div
        className="absolute inset-y-0 right-0 flex items-stretch"
        style={{ width: actionWidth }}
      >
        {hasUnread && (
          <button
            className="flex-1 flex flex-col items-center justify-center gap-1 bg-blue-500 text-white text-xs"
            onClick={() => { onMarkAsRead(room.id); onSwipeOpen(null); }}
          >
            <MailCheck className="size-5" />
            既読
          </button>
        )}
        <button
          className="flex-1 flex flex-col items-center justify-center gap-1 bg-violet-500 text-white text-xs"
          onClick={() => { onToggleMute(room.id); onSwipeOpen(null); }}
        >
          <BellOff className="size-5" />
          {room.isMuted ? "通知ON" : "ミュート"}
        </button>
        <button
          className="flex-1 flex flex-col items-center justify-center gap-1 bg-amber-500 text-white text-xs"
          onClick={() => { onTogglePin(room.id); onSwipeOpen(null); }}
        >
          <Pin className="size-5" />
          {room.isPinned ? "解除" : "ピン"}
        </button>
        <button
          className="flex-1 flex flex-col items-center justify-center gap-1 bg-red-500 text-white text-xs"
          onClick={() => { onDelete(room.id); onSwipeOpen(null); }}
        >
          <Trash2 className="size-5" />
          退出
        </button>
      </div>

      <div
        ref={rowRef}
        className="relative bg-background"
        style={{ transform: `translateX(${displayX}px)`, transition: "transform 0.2s ease-out" }}
        onTouchStart={handleTouchStart}
        onTouchMove={handleTouchMove}
        onTouchEnd={handleTouchEnd}
      >
        <div className="flex items-center gap-3 px-4 py-3 hover:bg-muted/30 transition-colors">
          {customId ? (
            <Link
              to={`/users/${customId}`}
              className="relative shrink-0"
              onClick={(e) => { if (isDragging.current) e.preventDefault(); e.stopPropagation(); }}
            >
              <Avatar className="size-12">
                <AvatarFallback>{displayName.slice(0, 2)}</AvatarFallback>
              </Avatar>
            </Link>
          ) : (
            <div className="relative shrink-0">
              <Avatar className="size-12">
                <AvatarFallback>{displayName.slice(0, 2)}</AvatarFallback>
              </Avatar>
            </div>
          )}
          <Link
            to={`/chat/${room.id}`}
            className="flex-1 min-w-0"
            onClick={(e) => { if (isDragging.current) e.preventDefault(); }}
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-1.5 min-w-0">
                <span className="text-sm font-bold truncate">{displayName}</span>
                {room.isPinned && <Pin className="size-3 text-muted-foreground shrink-0 fill-current" />}
                {room.isMuted && <BellOff className="size-3 text-muted-foreground shrink-0" />}
              </div>
              {room.lastMessageAt && (
                <span className="text-xs text-muted-foreground whitespace-nowrap ml-2">
                  {formatChatTime(room.lastMessageAt)}
                </span>
              )}
            </div>
            <div className="flex items-center justify-between mt-0.5">
              <p className="text-sm text-muted-foreground truncate">
                {room.isLastMessageMine && <span className="text-muted-foreground">自分: </span>}
                {room.lastMessageText ?? ""}
              </p>
              {room.unreadCount > 0 && (
                <span className="ml-2 shrink-0 size-5 rounded-full bg-blue-500 text-white text-xs flex items-center justify-center">
                  {room.unreadCount}
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
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState("");
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [openSwipeId, setOpenSwipeId] = useState<string | null>(null);

  useEffect(() => {
    listRooms().then((res) => {
      setRooms(res.rooms ?? []);
    }).finally(() => setLoading(false));
  }, []);

  const filteredRooms = rooms.filter((room) => {
    const name = getRoomDisplayName(room).toLowerCase();
    const q = searchQuery.toLowerCase();
    return name.includes(q) || (room.otherUser?.customId ?? "").toLowerCase().includes(q);
  });

  const sortedRooms = [...filteredRooms].sort((a, b) => {
    if (a.isPinned && !b.isPinned) return -1;
    if (!a.isPinned && b.isPinned) return 1;
    return 0;
  });

  const handleMarkAsRead = async (id: string) => {
    await markRoomAsRead(id).catch(() => {});
    setRooms((prev) => prev.map((r) => r.id === id ? { ...r, unreadCount: 0 } : r));
  };

  const handleTogglePin = async (id: string) => {
    const room = rooms.find((r) => r.id === id);
    if (!room) return;
    try {
      const res = room.isPinned ? await unpinRoom(id) : await pinRoom(id);
      setRooms((prev) => prev.map((r) => r.id === id ? res.room : r));
    } catch {}
  };

  const handleToggleMute = async (id: string) => {
    const room = rooms.find((r) => r.id === id);
    if (!room) return;
    try {
      const res = room.isMuted ? await unmuteRoom(id) : await muteRoom(id);
      setRooms((prev) => prev.map((r) => r.id === id ? res.room : r));
    } catch {}
  };

  const handleDelete = async (id: string) => {
    await leaveRoom(id).catch(() => {});
    setRooms((prev) => prev.filter((r) => r.id !== id));
  };

  const toggleSelect = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handleBulkAction = async (action: "read" | "pin" | "mute" | "delete") => {
    const ids = Array.from(selected);
    if (ids.length === 0) return;
    switch (action) {
      case "read":
        await Promise.all(ids.map((id) => markRoomAsRead(id).catch(() => {})));
        setRooms((prev) => prev.map((r) => ids.includes(r.id) ? { ...r, unreadCount: 0 } : r));
        break;
      case "delete":
        await Promise.all(ids.map((id) => leaveRoom(id).catch(() => {})));
        setRooms((prev) => prev.filter((r) => !ids.includes(r.id)));
        break;
      default:
        break;
    }
    setSelected(new Set());
  };

  const exitEditing = () => { setEditing(false); setSelected(new Set()); };

  return (
    <div className="w-full min-h-full flex flex-col">
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center justify-between px-4 h-14">
          {editing ? (
            <>
              <Button variant="ghost" size="sm" onClick={exitEditing}>キャンセル</Button>
              <h1 className="text-base font-bold">
                {selected.size > 0 ? `${selected.size}件選択中` : "チャットを選択"}
              </h1>
              <Button variant="ghost" size="sm" onClick={() => {
                if (selected.size === sortedRooms.length) setSelected(new Set());
                else setSelected(new Set(sortedRooms.map((r) => r.id)));
              }}>
                {selected.size === sortedRooms.length ? "全解除" : "全選択"}
              </Button>
            </>
          ) : (
            <>
              <h1 className="text-lg font-bold">チャット</h1>
              <div className="flex items-center gap-1">
                <Button variant="ghost" size="icon" onClick={() => navigate("/chat/new")}>
                  <MessageSquarePlus className="size-5" />
                </Button>
                <Button variant="ghost" size="sm" onClick={() => setEditing(true)}>
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

      {!editing && (
        <Link
          to="/chat/requests"
          className="flex items-center justify-between px-4 py-2.5 border-b hover:bg-muted/30 transition-colors"
        >
          <span className="text-sm text-blue-500">メッセージリクエスト</span>
          <ChevronRight className="size-4 text-blue-500" />
        </Link>
      )}

      <div className="flex-1">
        {loading ? (
          <div className="flex items-center justify-center py-12 text-muted-foreground text-sm">読み込み中...</div>
        ) : sortedRooms.length === 0 ? (
          <div className="flex items-center justify-center py-12 text-muted-foreground text-sm">
            チャットがありません
          </div>
        ) : editing ? (
          sortedRooms.map((room) => {
            const displayName = getRoomDisplayName(room);
            return (
              <button
                key={room.id}
                className="flex items-center gap-3 px-4 py-3 border-b w-full text-left hover:bg-muted/30 transition-colors"
                onClick={() => toggleSelect(room.id)}
              >
                <div className={`shrink-0 size-5 rounded-full border-2 flex items-center justify-center transition-colors ${selected.has(room.id) ? "bg-blue-500 border-blue-500" : "border-muted-foreground"}`}>
                  {selected.has(room.id) && (
                    <svg viewBox="0 0 12 12" className="size-3 text-white" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round">
                      <polyline points="2,6 5,9 10,3" />
                    </svg>
                  )}
                </div>
                <div className="relative shrink-0">
                  <Avatar className="size-12">
                    <AvatarFallback>{displayName.slice(0, 2)}</AvatarFallback>
                  </Avatar>
                </div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-1.5 min-w-0">
                      <span className="text-sm font-bold truncate">{displayName}</span>
                      {room.isPinned && <Pin className="size-3 text-muted-foreground shrink-0 fill-current" />}
                      {room.isMuted && <BellOff className="size-3 text-muted-foreground shrink-0" />}
                    </div>
                    {room.lastMessageAt && (
                      <span className="text-xs text-muted-foreground whitespace-nowrap ml-2">
                        {formatChatTime(room.lastMessageAt)}
                      </span>
                    )}
                  </div>
                  <p className="text-sm text-muted-foreground truncate mt-0.5">
                    {room.isLastMessageMine && "自分: "}
                    {room.lastMessageText ?? ""}
                  </p>
                </div>
              </button>
            );
          })
        ) : (
          sortedRooms.map((room) => (
            <SwipeableChatRow
              key={room.id}
              room={room}
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

      {editing ? (
        <div className="sticky bottom-0 left-0 w-full border-t bg-background px-2 py-2">
          <div className="flex items-center justify-around">
            <button className="flex flex-col items-center gap-0.5 text-xs text-muted-foreground disabled:opacity-30" disabled={selected.size === 0} onClick={() => handleBulkAction("read")}>
              <MailCheck className="size-5" />
              既読
            </button>
            <button className="flex flex-col items-center gap-0.5 text-xs text-destructive disabled:opacity-30" disabled={selected.size === 0} onClick={() => handleBulkAction("delete")}>
              <Trash2 className="size-5" />
              退出
            </button>
          </div>
        </div>
      ) : (
        <BottomNav />
      )}
    </div>
  );
}
