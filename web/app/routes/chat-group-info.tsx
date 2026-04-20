import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router";
import { Avatar, AvatarFallback } from "~/components/ui/avatar";
import { Input } from "~/components/ui/input";
import { Button } from "~/components/ui/button";
import { ArrowLeft, Loader2, Pencil, UserPlus, LogOut, X } from "lucide-react";
import {
  listRoomMembers,
  updateRoom,
  addRoomMember,
  removeRoomMember,
  leaveRoom,
  searchUsers,
  getMe,
  type User as ApiUser,
} from "~/lib/api";

export default function ChatGroupInfo() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [members, setMembers] = useState<ApiUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [myId, setMyId] = useState<string | null>(null);
  const [editingName, setEditingName] = useState(false);
  const [groupName, setGroupName] = useState("");
  const [savingName, setSavingName] = useState(false);
  const [addQuery, setAddQuery] = useState("");
  const [addResults, setAddResults] = useState<ApiUser[]>([]);
  const [searching, setSearching] = useState(false);

  useEffect(() => {
    if (!id) return;
    Promise.all([
      listRoomMembers(id).then((res) => setMembers(res.members ?? [])),
      getMe().then(({ user }) => setMyId(user.id)),
    ]).finally(() => setLoading(false));
  }, [id]);

  useEffect(() => {
    const q = addQuery.trim();
    if (!q) { setAddResults([]); return; }
    setSearching(true);
    const t = setTimeout(() => {
      searchUsers(q, { pageSize: 10 })
        .then((res) => setAddResults(res.users ?? []))
        .catch(() => setAddResults([]))
        .finally(() => setSearching(false));
    }, 300);
    return () => clearTimeout(t);
  }, [addQuery]);

  const handleSaveName = async () => {
    if (!id || !groupName.trim()) return;
    setSavingName(true);
    try {
      await updateRoom(id, groupName.trim());
      setEditingName(false);
    } catch {
    } finally {
      setSavingName(false);
    }
  };

  const handleAddMember = async (user: ApiUser) => {
    if (!id) return;
    try {
      await addRoomMember(id, user.id);
      setMembers((prev) => [...prev, user]);
      setAddQuery("");
      setAddResults([]);
    } catch {}
  };

  const handleRemoveMember = async (userId: string) => {
    if (!id) return;
    try {
      await removeRoomMember(id, userId);
      setMembers((prev) => prev.filter((m) => m.id !== userId));
    } catch {}
  };

  const handleLeave = async () => {
    if (!id) return;
    try {
      await leaveRoom(id);
      navigate("/chat");
    } catch {}
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="size-5 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="w-full min-h-full flex flex-col">
      <div className="sticky top-0 border-b bg-background/60 backdrop-blur-lg z-10 flex items-center gap-3 px-4 h-14">
        <Button variant="ghost" size="icon" onClick={() => history.back()}>
          <ArrowLeft className="size-5" />
        </Button>
        <h1 className="text-base font-bold flex-1">グループ情報</h1>
      </div>

      <div className="flex-1 divide-y">
        {/* Group name */}
        <div className="px-4 py-4">
          <p className="text-xs font-semibold text-muted-foreground mb-2">グループ名</p>
          {editingName ? (
            <div className="flex gap-2">
              <Input
                value={groupName}
                onChange={(e) => setGroupName(e.target.value)}
                autoFocus
                className="flex-1 h-9"
              />
              <Button size="sm" onClick={handleSaveName} disabled={savingName || !groupName.trim()}>
                {savingName ? <Loader2 className="size-4 animate-spin" /> : "保存"}
              </Button>
              <Button size="sm" variant="ghost" onClick={() => setEditingName(false)}>キャンセル</Button>
            </div>
          ) : (
            <button
              className="flex items-center gap-2 text-sm hover:opacity-70"
              onClick={() => { setGroupName(""); setEditingName(true); }}
            >
              <span className="font-medium">名前を変更</span>
              <Pencil className="size-3.5 text-muted-foreground" />
            </button>
          )}
        </div>

        {/* Members */}
        <div>
          <p className="px-4 py-3 text-xs font-semibold text-muted-foreground">
            メンバー（{members.length}人）
          </p>
          {members.map((m) => (
            <div key={m.id} className="flex items-center gap-3 px-4 py-3 hover:bg-muted/30">
              <Avatar className="size-10 shrink-0">
                <AvatarFallback>{m.displayName.slice(0, 2)}</AvatarFallback>
              </Avatar>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium truncate">{m.displayName}</p>
                <p className="text-xs text-muted-foreground truncate">@{m.customId}</p>
              </div>
              {m.id !== myId && (
                <button
                  onClick={() => handleRemoveMember(m.id)}
                  className="shrink-0 p-1.5 text-muted-foreground hover:text-destructive transition-colors"
                >
                  <X className="size-4" />
                </button>
              )}
            </div>
          ))}
        </div>

        {/* Add member */}
        <div className="px-4 py-4">
          <p className="text-xs font-semibold text-muted-foreground mb-2">メンバーを追加</p>
          <div className="relative">
            <Input
              placeholder="ユーザーを検索..."
              value={addQuery}
              onChange={(e) => setAddQuery(e.target.value)}
              className="h-9 pr-8"
            />
            {searching && <Loader2 className="absolute right-2 top-2 size-4 text-muted-foreground animate-spin" />}
          </div>
          {addResults.length > 0 && (
            <div className="mt-1 border rounded-md overflow-hidden">
              {addResults
                .filter((u) => !members.find((m) => m.id === u.id))
                .map((u) => (
                  <button
                    key={u.id}
                    className="flex items-center gap-2 px-3 py-2 w-full hover:bg-muted/50 text-left border-b last:border-0"
                    onClick={() => handleAddMember(u)}
                  >
                    <Avatar className="size-7 shrink-0">
                      <AvatarFallback className="text-xs">{u.displayName.slice(0, 2)}</AvatarFallback>
                    </Avatar>
                    <div className="min-w-0">
                      <p className="text-sm font-medium truncate">{u.displayName}</p>
                      <p className="text-xs text-muted-foreground">@{u.customId}</p>
                    </div>
                    <UserPlus className="size-4 text-primary ml-auto shrink-0" />
                  </button>
                ))}
            </div>
          )}
        </div>

        {/* Leave */}
        <div className="px-4 py-4">
          <button
            onClick={handleLeave}
            className="flex items-center gap-2 text-sm text-destructive hover:opacity-70"
          >
            <LogOut className="size-4" />
            グループを退出
          </button>
        </div>
      </div>
    </div>
  );
}
