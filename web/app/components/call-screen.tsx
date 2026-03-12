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
} from "lucide-react";
import { useState, useRef, useCallback } from "react";
import { useCall, type CallVisibility } from "~/components/call-context";
import { GroupAvatar } from "~/components/call-list";

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
          className="flex items-center gap-2 min-w-0 hover:opacity-70 transition-opacity"
          onClick={() => {
            const idx = visibilityOrder.indexOf(activeCall.visibility);
            setVisibility(visibilityOrder[(idx + 1) % visibilityOrder.length]);
          }}
        >
          {(() => {
            const Icon = visibilityIcon[activeCall.visibility];
            return <Icon className="size-5 text-muted-foreground shrink-0 mr-1" />;
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
        <Button variant="ghost" size="icon" onClick={handleMinimize}>
          <Minimize2 className="size-5" />
        </Button>
      </div>

      <div className="flex-1 overflow-y-auto px-6">
        <div className="grid grid-cols-3 gap-4 justify-items-center">
          {activeCall.participants.map((p, i) => (
            <div key={i} className="flex flex-col items-center gap-2">
              <Avatar className="size-16">
                <AvatarImage src={p.avatarUrl} alt={p.name} />
                <AvatarFallback>{p.name.slice(0, 2)}</AvatarFallback>
              </Avatar>
              <span className="text-xs truncate max-w-20 text-center">
                {p.name}
              </span>
            </div>
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
            variant="destructive"
            size="icon-lg"
            className="rounded-full size-14"
            onClick={leaveCall}
          >
            <LogOut className="size-5" />
          </Button>
        </div>
      </div>
    </div>
  );
}
