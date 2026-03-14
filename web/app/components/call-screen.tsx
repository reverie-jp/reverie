import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import {
  Minimize2,
  Mic,
  MicOff,
  Volume2,
  VolumeOff,
  LogOut,
  Phone,
  Video,
  Globe,
  Users,
  UserCheck,
  Lock,
  Gavel,
  Timer,
  Flag,
  SlidersHorizontal,
  Shield,
  Crown,
  User,
  UserPlus,
  UserCheck as UserCheckIcon,
  Settings,
  ScreenShare,
  ScreenShareOff,
  Search,
} from "lucide-react";
import { useState, useRef, useCallback } from "react";
import { useCall, type CallVisibility } from "~/components/call-context";
import { GroupAvatar } from "~/components/call-list";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "~/components/ui/alert-dialog";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "~/components/ui/dialog";
import { Input } from "~/components/ui/input";
import { Slider } from "~/components/ui/slider";
import { Switch } from "~/components/ui/switch";

function DraggableBubble({
  onTap,
  children,
}: {
  onTap: () => void;
  children: React.ReactNode;
}) {
  const ref = useRef<HTMLDivElement>(null);
  const dragState = useRef<{
    startX: number;
    startY: number;
    originX: number;
    originY: number;
    moved: boolean;
  } | null>(null);
  const [pos, setPos] = useState({ x: 16, y: window.innerHeight - 140 });

  const clamp = useCallback((x: number, y: number) => {
    const size = 56;
    return {
      x: Math.max(0, Math.min(window.innerWidth - size, x)),
      y: Math.max(0, Math.min(window.innerHeight - size, y)),
    };
  }, []);

  const handlePointerDown = (e: React.PointerEvent) => {
    dragState.current = {
      startX: e.clientX,
      startY: e.clientY,
      originX: pos.x,
      originY: pos.y,
      moved: false,
    };
    (e.target as HTMLElement).setPointerCapture(e.pointerId);
  };

  const handlePointerMove = (e: React.PointerEvent) => {
    const ds = dragState.current;
    if (!ds) return;
    const dx = e.clientX - ds.startX;
    const dy = e.clientY - ds.startY;
    if (Math.abs(dx) > 3 || Math.abs(dy) > 3) {
      ds.moved = true;
    }
    setPos(clamp(ds.originX + dx, ds.originY + dy));
  };

  const handlePointerUp = () => {
    const ds = dragState.current;
    if (ds && !ds.moved) {
      onTap();
    }
    dragState.current = null;
  };

  return (
    <div
      ref={ref}
      onPointerDown={handlePointerDown}
      onPointerMove={handlePointerMove}
      onPointerUp={handlePointerUp}
      className="fixed z-50 touch-none select-none cursor-grab active:cursor-grabbing"
      style={{ left: pos.x, top: pos.y }}
    >
      {children}
    </div>
  );
}

type OnlineStatus = "online" | "idle" | "offline";

interface FollowingUser {
  name: string;
  customId: string;
  avatarUrl: string;
  status: OnlineStatus;
}

const statusColor: Record<OnlineStatus, string> = {
  online: "bg-green-500",
  idle: "bg-yellow-500",
  offline: "bg-muted-foreground/50",
};

const statusOrder: Record<OnlineStatus, number> = {
  online: 0,
  idle: 1,
  offline: 2,
};

const followingUsers: FollowingUser[] = [
  { name: "田中太郎", customId: "tanaka", avatarUrl: "", status: "online" },
  { name: "佐藤花子", customId: "sato", avatarUrl: "", status: "online" },
  { name: "鈴木一郎", customId: "suzuki", avatarUrl: "", status: "idle" },
  { name: "高橋美咲", customId: "takahashi", avatarUrl: "", status: "online" },
  { name: "伊藤健太", customId: "ito", avatarUrl: "", status: "offline" },
  { name: "渡辺さくら", customId: "watanabe", avatarUrl: "", status: "idle" },
  { name: "山本翔太", customId: "yamamoto", avatarUrl: "", status: "offline" },
  { name: "中村あおい", customId: "nakamura", avatarUrl: "", status: "online" },
  { name: "小林悠", customId: "kobayashi", avatarUrl: "", status: "offline" },
  { name: "加藤莉子", customId: "kato", avatarUrl: "", status: "idle" },
] satisfies FollowingUser[];

followingUsers.sort((a, b) => statusOrder[a.status] - statusOrder[b.status]);

const visibilityOrder: CallVisibility[] = [
  "すべてのユーザー",
  "フォロワーのみ",
  "相互フォローのみ",
  "招待した人のみ",
];

const visibilityIcon: Record<CallVisibility, typeof Globe> = {
  すべてのユーザー: Globe,
  フォロワーのみ: Users,
  相互フォローのみ: UserCheck,
  招待した人のみ: Lock,
};

export function CallScreen() {
  const {
    activeCall,
    isMinimized,
    minimize,
    maximize,
    leaveCall,
    setVisibility,
  } = useCall();
  const [isMuted, setIsMuted] = useState(false);
  const [isSpeaker, setIsSpeaker] = useState(true);
  const [isClosing, setIsClosing] = useState(false);
  const [isScreenSharing, setIsScreenSharing] = useState(false);
  const [screenShareDialog, setScreenShareDialog] = useState<"start" | "stop" | null>(null);
  const [showSettings, setShowSettings] = useState(false);
  const [showInvite, setShowInvite] = useState(false);
  const [invitedUsers, setInvitedUsers] = useState<Set<string>>(new Set());
  const [inviteSearch, setInviteSearch] = useState("");
  const [outputVolume, setOutputVolume] = useState(80);
  const [inputVolume, setInputVolume] = useState(80);
  const [joinSound, setJoinSound] = useState(true);
  const [highQuality, setHighQuality] = useState(false);

  const handleMinimize = () => {
    setIsClosing(true);
    setTimeout(() => {
      minimize();
      setIsClosing(false);
    }, 250);
  };

  if (!activeCall) return null;

  const TypeIcon = activeCall.type === "video" ? Video : Phone;

  if (isMinimized) {
    return (
      <DraggableBubble onTap={maximize}>
        <div className="relative size-14 shadow-lg rounded-full ring-2 ring-primary ring-offset-2 ring-offset-background">
          <GroupAvatar
            participants={activeCall.participants}
            className="size-14"
          />
          <div className="absolute -bottom-0.5 -left-0.5 size-5 rounded-full bg-primary flex items-center justify-center">
            <TypeIcon className="size-3 text-primary-foreground" />
          </div>
        </div>
      </DraggableBubble>
    );
  }

  return (
    <div
      className={`fixed inset-0 z-50 bg-background flex flex-col duration-300 ${isClosing ? "animate-out slide-out-to-bottom fill-mode-forwards" : "animate-in slide-in-from-bottom"}`}
    >
      <div className="flex items-center justify-between px-6 py-8 shrink-0">
        <button
          className="flex items-center gap-3 min-w-0 hover:opacity-70 transition-opacity"
          onClick={() => {
            const idx = visibilityOrder.indexOf(activeCall.visibility);
            setVisibility(visibilityOrder[(idx + 1) % visibilityOrder.length]);
          }}
        >
          {(() => {
            const Icon = visibilityIcon[activeCall.visibility];
            return <Icon className="size-5 text-muted-foreground shrink-0" />;
          })()}
          <div className="min-w-0 text-left">
            <p className="text-sm font-medium truncate">
              {activeCall.visibility}
            </p>
            <p className="text-xs text-muted-foreground">
              {activeCall.participants.length}人が参加中
            </p>
          </div>
        </button>
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={() => setShowInvite(true)}>
            <UserPlus className="size-5" />
          </Button>
          <Button variant="ghost" size="icon" onClick={() => setShowSettings(true)}>
            <Settings className="size-5" />
          </Button>
          <Button variant="ghost" size="icon" onClick={handleMinimize}>
            <Minimize2 className="size-5" />
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-6">
        {isScreenSharing && (
          <div className="mb-6 rounded-lg border bg-muted/30 aspect-video flex flex-col items-center justify-center gap-2">
            <ScreenShare className="size-8 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">画面を共有中</span>
          </div>
        )}
        <div className="grid grid-cols-3 gap-6 justify-items-center">
          {activeCall.participants.map((p, i) => (
            <DropdownMenu key={i}>
              <DropdownMenuTrigger className="flex flex-col items-center gap-2 outline-none hover:opacity-70 transition-opacity">
                <Avatar className="size-16">
                  <AvatarImage src={p.avatarUrl} alt={p.name} />
                  <AvatarFallback>{p.name.slice(0, 2)}</AvatarFallback>
                </Avatar>
                <span className="text-xs truncate max-w-20 text-center">
                  {p.name}
                </span>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="min-w-44">
                <DropdownMenuGroup>
                  <DropdownMenuLabel>{p.name}</DropdownMenuLabel>
                </DropdownMenuGroup>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <User className="size-4" />
                  プロフィール
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Shield className="size-4" />
                  サブホストにする
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Crown className="size-4" />
                  ホスト権限を譲渡
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <MicOff className="size-4" />
                  ミュートする
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <SlidersHorizontal className="size-4" />
                  音量調整
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <Flag className="size-4" />
                  通報する
                </DropdownMenuItem>
                <DropdownMenuItem className="text-destructive">
                  <Timer className="size-4" />
                  一時追放
                </DropdownMenuItem>
                <DropdownMenuItem className="text-destructive">
                  <Gavel className="size-4" />
                  永久追放
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ))}
        </div>
      </div>

      <div className="shrink-0 pb-safe px-6 py-8">
        <div className="flex items-center justify-center gap-6">
          <Button
            variant={isSpeaker ? "outline" : "secondary"}
            size="icon-lg"
            className="rounded-full size-14"
            onClick={() => setIsSpeaker(!isSpeaker)}
          >
            {isSpeaker ? (
              <Volume2 className="size-5" />
            ) : (
              <VolumeOff className="size-5" />
            )}
          </Button>
          <Button
            variant={isMuted ? "destructive" : "outline"}
            size="icon-lg"
            className="rounded-full size-14"
            onClick={() => setIsMuted(!isMuted)}
          >
            {isMuted ? (
              <MicOff className="size-5" />
            ) : (
              <Mic className="size-5" />
            )}
          </Button>
          <Button
            variant={isScreenSharing ? "destructive" : "outline"}
            size="icon-lg"
            className="rounded-full size-14"
            onClick={() => setScreenShareDialog(isScreenSharing ? "stop" : "start")}
          >
            {isScreenSharing ? (
              <ScreenShareOff className="size-5" />
            ) : (
              <ScreenShare className="size-5" />
            )}
          </Button>
          <Button
            variant="destructive"
            size="icon-lg"
            className="rounded-full size-14"
            onClick={leaveCall}
          >
            <LogOut className="size-5" />
          </Button>
        </div>
      </div>

      <AlertDialog
        open={screenShareDialog === "start"}
        onOpenChange={(open) => { if (!open) setScreenShareDialog(null); }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>画面を共有</AlertDialogTitle>
            <AlertDialogDescription>
              通話の参加者全員にあなたの画面が表示されます。共有を開始しますか？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>キャンセル</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                setIsScreenSharing(true);
                setScreenShareDialog(null);
              }}
            >
              共有を開始
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog
        open={screenShareDialog === "stop"}
        onOpenChange={(open) => { if (!open) setScreenShareDialog(null); }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>画面共有を停止</AlertDialogTitle>
            <AlertDialogDescription>
              画面の共有を停止しますか？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>キャンセル</AlertDialogCancel>
            <AlertDialogAction
              variant="destructive"
              onClick={() => {
                setIsScreenSharing(false);
                setScreenShareDialog(null);
              }}
            >
              共有を停止
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <Dialog open={showInvite} onOpenChange={(open) => { setShowInvite(open); if (!open) setInviteSearch(""); }}>
        <DialogContent showCloseButton={false} className="max-w-80 sm:max-w-80">
          <DialogHeader>
            <DialogTitle>ユーザーを招待</DialogTitle>
          </DialogHeader>
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
            <Input
              placeholder="ユーザーを検索"
              value={inviteSearch}
              onChange={(e) => setInviteSearch(e.target.value)}
              className="pl-9 h-9"
            />
          </div>
          <div className="max-h-72 overflow-y-auto -mx-2">
            {followingUsers
              .filter((u) => {
                if (!inviteSearch) return true;
                const q = inviteSearch.toLowerCase();
                return u.name.toLowerCase().includes(q) || u.customId.toLowerCase().includes(q);
              })
              .map((user) => {
                const invited = invitedUsers.has(user.customId);
                return (
                  <button
                    key={user.customId}
                    className="flex items-center gap-3 w-full px-3 py-2 rounded-md hover:bg-muted transition-colors"
                    onClick={() => {
                      setInvitedUsers((prev) => {
                        const next = new Set(prev);
                        if (invited) next.delete(user.customId);
                        else next.add(user.customId);
                        return next;
                      });
                    }}
                  >
                    <div className="relative">
                      <Avatar className="size-10">
                        <AvatarImage src={user.avatarUrl} alt={user.name} />
                        <AvatarFallback>{user.name.slice(0, 2)}</AvatarFallback>
                      </Avatar>
                      <div className={`absolute -bottom-0.5 -right-0.5 size-3 rounded-full border-2 border-background ${statusColor[user.status]}`} />
                    </div>
                    <div className="flex-1 min-w-0 text-left">
                      <p className="text-sm font-medium truncate">{user.name}</p>
                      <p className="text-xs text-muted-foreground">@{user.customId}</p>
                    </div>
                    {invited ? (
                      <UserCheckIcon className="size-4 text-primary shrink-0" />
                    ) : (
                      <UserPlus className="size-4 text-muted-foreground shrink-0" />
                    )}
                  </button>
                );
              })}
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={showSettings} onOpenChange={setShowSettings}>
        <DialogContent showCloseButton={false} className="max-w-80 sm:max-w-80">
          <DialogHeader>
            <DialogTitle>通話設定</DialogTitle>
          </DialogHeader>
          <div className="space-y-6">
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm">出力音量</span>
                <span className="text-xs text-muted-foreground">{outputVolume}%</span>
              </div>
              <Slider
                value={[outputVolume]}
                onValueChange={(v) => setOutputVolume(Array.isArray(v) ? v[0] : v)}
                max={100}
                step={1}
              />
              <p className="text-xs text-muted-foreground">自分に聞こえる音量</p>
            </div>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm">入力音量</span>
                <span className="text-xs text-muted-foreground">{inputVolume}%</span>
              </div>
              <Slider
                value={[inputVolume]}
                onValueChange={(v) => setInputVolume(Array.isArray(v) ? v[0] : v)}
                max={100}
                step={1}
              />
              <p className="text-xs text-muted-foreground">相手に聞こえる音量</p>
            </div>
            <div className="flex items-center justify-between">
              <div>
                <span className="text-sm">入室効果音</span>
                <p className="text-xs text-muted-foreground">誰かが入室した時に効果音を鳴らす</p>
              </div>
              <Switch checked={joinSound} onCheckedChange={setJoinSound} />
            </div>
            <div className="flex items-center justify-between gap-4">
              <div>
                <span className="text-sm">高音質モード</span>
                <p className="text-xs text-muted-foreground">音質が向上しますが、バッテリー消費と通信量が増加します</p>
              </div>
              <Switch checked={highQuality} onCheckedChange={setHighQuality} />
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
