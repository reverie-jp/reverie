import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { Avatar, AvatarFallback } from "~/components/ui/avatar";
import { Input } from "~/components/ui/input";
import { Button } from "~/components/ui/button";
import { ArrowLeft, Search, Loader2, Check, X, Users } from "lucide-react";
import {
  searchUsers,
  listFollowing,
  createDirectRoom,
  createGroupRoom,
  getMe,
  type User as ApiUser,
} from "~/lib/api";

type Mode = "direct" | "group";

export default function ChatNew() {
  const navigate = useNavigate();
  const [mode, setMode] = useState<Mode>("direct");
  const [query, setQuery] = useState("");
  const [followingUsers, setFollowingUsers] = useState<ApiUser[]>([]);
  const [searchResults, setSearchResults] = useState<ApiUser[]>([]);
  const [loadingFollowing, setLoadingFollowing] = useState(true);
  const [searching, setSearching] = useState(false);
  const [creating, setCreating] = useState(false);
  const [selectedUsers, setSelectedUsers] = useState<ApiUser[]>([]);
  const [groupName, setGroupName] = useState("");

  useEffect(() => {
    getMe()
      .then(({ user }) => listFollowing(user.id, { pageSize: 50 }))
      .then((res) => setFollowingUsers(res.users ?? []))
      .catch(() => {})
      .finally(() => setLoadingFollowing(false));
  }, []);

  useEffect(() => {
    const q = query.trim();
    if (!q) {
      setSearchResults([]);
      return;
    }
    setSearching(true);
    const timer = setTimeout(() => {
      searchUsers(q, { pageSize: 20 })
        .then((res) => setSearchResults(res.users ?? []))
        .catch(() => setSearchResults([]))
        .finally(() => setSearching(false));
    }, 300);
    return () => clearTimeout(timer);
  }, [query]);

  const handleSelectDirect = async (user: ApiUser) => {
    if (creating) return;
    setCreating(true);
    try {
      const res = await createDirectRoom(user.id);
      navigate(`/chat/${res.room.id}`);
    } catch {
      setCreating(false);
    }
  };

  const toggleSelectGroup = (user: ApiUser) => {
    setSelectedUsers((prev) =>
      prev.find((u) => u.id === user.id)
        ? prev.filter((u) => u.id !== user.id)
        : [...prev, user]
    );
  };

  const handleCreateGroup = async () => {
    if (creating || selectedUsers.length === 0 || !groupName.trim()) return;
    setCreating(true);
    try {
      const res = await createGroupRoom(
        groupName.trim(),
        selectedUsers.map((u) => u.id)
      );
      navigate(`/chat/${res.room.id}`);
    } catch {
      setCreating(false);
    }
  };

  const isSearching = query.trim().length > 0;
  const displayUsers = isSearching ? searchResults : followingUsers;
  const isLoading = isSearching ? searching : loadingFollowing;

  return (
    <div className="w-full min-h-full flex flex-col">
      {/* Header */}
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-14">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <h1 className="text-base font-bold flex-1">
            {mode === "group" ? "グループを作成" : "新しいメッセージ"}
          </h1>
          {mode === "direct" && (
            <Button variant="ghost" size="sm" onClick={() => setMode("group")} className="text-xs gap-1">
              <Users className="size-4" />
              グループ
            </Button>
          )}
          {mode === "group" && (
            <Button variant="ghost" size="sm" onClick={() => { setMode("direct"); setSelectedUsers([]); setGroupName(""); }} className="text-xs">
              キャンセル
            </Button>
          )}
        </div>

        {mode === "group" && selectedUsers.length > 0 && (
          <div className="flex gap-2 px-4 pb-2 overflow-x-auto">
            {selectedUsers.map((u) => (
              <div key={u.id} className="flex items-center gap-1 bg-primary/10 text-primary rounded-full px-2 py-1 text-xs shrink-0">
                <span>{u.displayName}</span>
                <button onClick={() => toggleSelectGroup(u)}>
                  <X className="size-3" />
                </button>
              </div>
            ))}
          </div>
        )}

        {mode === "group" && (
          <div className="px-4 pb-2">
            <Input
              placeholder="グループ名"
              value={groupName}
              onChange={(e) => setGroupName(e.target.value)}
              className="h-9 text-sm"
            />
          </div>
        )}

        <div className="px-4 pb-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
            {(searching || loadingFollowing) && (
              <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground animate-spin" />
            )}
            <Input
              placeholder="名前や@IDで検索"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              className="pl-9 pr-9 h-10"
              autoFocus={mode === "direct"}
            />
          </div>
        </div>
      </div>

      {!isSearching && (
        <div className="px-4 py-2 text-xs font-semibold text-muted-foreground border-b bg-muted/30">
          フォロー中
        </div>
      )}

      <div className="flex-1">
        {isLoading ? (
          <div className="flex items-center justify-center py-16">
            <Loader2 className="size-5 text-muted-foreground animate-spin" />
          </div>
        ) : displayUsers.length === 0 ? (
          <div className="flex items-center justify-center py-16 text-muted-foreground text-sm">
            {isSearching ? "ユーザーが見つかりません" : "フォロー中のユーザーがいません"}
          </div>
        ) : (
          displayUsers.map((user) => {
            const isSelected = selectedUsers.some((u) => u.id === user.id);
            return (
              <button
                key={user.id}
                className="flex items-center gap-3 px-4 py-3 border-b w-full text-left hover:bg-muted/30 transition-colors disabled:opacity-50"
                onClick={() => mode === "group" ? toggleSelectGroup(user) : handleSelectDirect(user)}
                disabled={creating && mode === "direct"}
              >
                <div className="relative shrink-0">
                  <Avatar className="size-11">
                    <AvatarFallback>{user.displayName.slice(0, 2)}</AvatarFallback>
                  </Avatar>
                  {mode === "group" && isSelected && (
                    <div className="absolute -bottom-0.5 -right-0.5 size-5 bg-primary rounded-full flex items-center justify-center">
                      <Check className="size-3 text-primary-foreground" />
                    </div>
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-bold truncate">{user.displayName}</p>
                  <p className="text-xs text-muted-foreground truncate">@{user.customId}</p>
                </div>
              </button>
            );
          })
        )}
      </div>

      {mode === "group" && selectedUsers.length > 0 && (
        <div className="sticky bottom-0 px-4 py-3 border-t bg-background">
          <Button
            className="w-full"
            onClick={handleCreateGroup}
            disabled={creating || !groupName.trim()}
          >
            {creating ? <Loader2 className="size-4 animate-spin mr-2" /> : null}
            グループを作成（{selectedUsers.length}人）
          </Button>
        </div>
      )}
    </div>
  );
}
