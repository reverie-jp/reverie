import { useState, useRef, useEffect } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import {
  Phone,
  PhoneOff,
  Video,
  VideoOff,
  Mic,
  MicOff,
  Volume2,
  VolumeOff,
  Minimize2,
} from "lucide-react";
import { usePrivateCall, formatCallDuration } from "./private-call-context";
import { Button } from "./ui/button";

export function PrivateCallScreen() {
  const {
    activeCall,
    isMinimized,
    endCall,
    minimize,
    setMicOn,
    setSpeakerOn,
    setCameraOn,
  } = usePrivateCall();

  const [swapped, setSwapped] = useState(false);

  // Draggable PiP state
  const pipRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [pipPos, setPipPos] = useState({ x: -1, y: -1 });
  const dragState = useRef({
    dragging: false,
    startX: 0,
    startY: 0,
    startPosX: 0,
    startPosY: 0,
    moved: false,
    usedTouch: false,
  });

  const PIP_W = 112;
  const PIP_H = 160;
  const PIP_MARGIN = 16;

  // Initialize pip position to bottom-right
  useEffect(() => {
    if (!containerRef.current || pipPos.x !== -1) return;
    const el = containerRef.current;
    requestAnimationFrame(() => {
      const rect = el.getBoundingClientRect();
      setPipPos({
        x: Math.max(PIP_MARGIN, rect.width - PIP_W - PIP_MARGIN),
        y: Math.max(PIP_MARGIN, rect.height - PIP_H - PIP_MARGIN),
      });
    });
  }, [activeCall?.cameraOn]);

  // Reset swapped & pip when call changes
  useEffect(() => {
    setSwapped(false);
    setPipPos({ x: -1, y: -1 });
  }, [activeCall?.chatId]);

  const handlePipTouchStart = (e: React.TouchEvent) => {
    const touch = e.touches[0];
    dragState.current = {
      dragging: true,
      startX: touch.clientX,
      startY: touch.clientY,
      startPosX: pipPos.x,
      startPosY: pipPos.y,
      moved: false,
      usedTouch: true,
    };
  };

  const handlePipTouchMove = (e: React.TouchEvent) => {
    if (!dragState.current.dragging || !containerRef.current) return;
    const touch = e.touches[0];
    const dx = touch.clientX - dragState.current.startX;
    const dy = touch.clientY - dragState.current.startY;
    if (Math.abs(dx) > 5 || Math.abs(dy) > 5) {
      dragState.current.moved = true;
    }
    const container = containerRef.current.getBoundingClientRect();
    const newX = Math.max(
      PIP_MARGIN,
      Math.min(
        container.width - PIP_W - PIP_MARGIN,
        dragState.current.startPosX + dx,
      ),
    );
    const newY = Math.max(
      PIP_MARGIN,
      Math.min(
        container.height - PIP_H - PIP_MARGIN,
        dragState.current.startPosY + dy,
      ),
    );
    setPipPos({ x: newX, y: newY });
  };

  const handlePipTouchEnd = () => {
    if (!dragState.current.moved) {
      setSwapped((s) => !s);
    }
    dragState.current.dragging = false;
  };

  if (!activeCall || isMinimized) return null;

  const { participant, micOn, speakerOn, cameraOn, elapsed } = activeCall;

  return (
    <div className="fixed inset-0 z-50 bg-background flex flex-col">
      {/* Header */}
      <div className="flex items-center justify-between px-6 py-8">
        <div className="size-9" />
        <span className="text-sm text-muted-foreground">
          {formatCallDuration(elapsed)}
        </span>
        <Button variant="ghost" size="icon" onClick={minimize}>
          <Minimize2 className="size-5" />
        </Button>
      </div>

      {/* Main area */}
      <div className="flex-1 flex flex-col items-center justify-center gap-4">
        {cameraOn ? (
          <div
            ref={containerRef}
            className="w-full flex-1 relative overflow-hidden"
          >
            {/* Full-screen view */}
            <div className="absolute inset-0 bg-muted/30 flex items-center justify-center">
              <div className="flex flex-col items-center gap-3">
                {swapped ? (
                  <span className="text-lg text-muted-foreground">自分</span>
                ) : (
                  <>
                    <Avatar className="size-24">
                      <AvatarImage
                        src={participant.avatarUrl}
                        alt={participant.name}
                      />
                      <AvatarFallback className="text-2xl">
                        {participant.name.slice(0, 2)}
                      </AvatarFallback>
                    </Avatar>
                    <p className="text-lg font-bold">{participant.name}</p>
                  </>
                )}
              </div>
            </div>
            {/* PiP (draggable, tappable to swap) */}
            <div
              ref={pipRef}
              className="absolute rounded-xl bg-muted border shadow-lg flex items-center justify-center touch-none select-none cursor-grab active:cursor-grabbing"
              style={{
                width: PIP_W,
                height: PIP_H,
                left: pipPos.x,
                top: pipPos.y,
              }}
              onTouchStart={handlePipTouchStart}
              onTouchMove={handlePipTouchMove}
              onTouchEnd={handlePipTouchEnd}
              onMouseDown={(e) => {
                if (dragState.current.usedTouch) return;
                dragState.current = {
                  dragging: true,
                  startX: e.clientX,
                  startY: e.clientY,
                  startPosX: pipPos.x,
                  startPosY: pipPos.y,
                  moved: false,
                  usedTouch: false,
                };
              }}
              onMouseUp={() => {
                if (dragState.current.usedTouch) {
                  dragState.current.usedTouch = false;
                  return;
                }
                if (!dragState.current.moved) {
                  setSwapped((s) => !s);
                }
                dragState.current.dragging = false;
              }}
            >
              {swapped ? (
                <Avatar className="size-12">
                  <AvatarImage
                    src={participant.avatarUrl}
                    alt={participant.name}
                  />
                  <AvatarFallback className="text-sm">
                    {participant.name.slice(0, 2)}
                  </AvatarFallback>
                </Avatar>
              ) : (
                <span className="text-xs text-muted-foreground">自分</span>
              )}
            </div>
          </div>
        ) : (
          <>
            <Avatar className="size-28">
              <AvatarImage src={participant.avatarUrl} alt={participant.name} />
              <AvatarFallback className="text-3xl">
                {participant.name.slice(0, 2)}
              </AvatarFallback>
            </Avatar>
            <p className="text-lg font-bold">{participant.name}</p>
          </>
        )}
      </div>

      {/* Controls */}
      <div className="shrink-0 px-6 py-8">
        <div className="flex items-center justify-center gap-6">
          <Button
            variant="secondary"
            size="icon-lg"
            className={`size-14! rounded-full ${!speakerOn ? "opacity-80 text-muted-foreground" : ""}`}
            onClick={() => setSpeakerOn(!speakerOn)}
          >
            {speakerOn ? (
              <Volume2 className="size-6" />
            ) : (
              <VolumeOff className="size-6" />
            )}
          </Button>
          <Button
            variant="secondary"
            size="icon-lg"
            className={`size-14! rounded-full ${!cameraOn ? "opacity-80 text-muted-foreground" : ""}`}
            onClick={() => setCameraOn(!cameraOn)}
          >
            {cameraOn ? (
              <Video className="size-6" />
            ) : (
              <VideoOff className="size-6" />
            )}
          </Button>
          <Button
            variant="secondary"
            size="icon-lg"
            className={`size-14! rounded-full ${!micOn ? "opacity-80 text-muted-foreground" : ""}`}
            onClick={() => setMicOn(!micOn)}
          >
            {micOn ? <Mic className="size-6" /> : <MicOff className="size-6" />}
          </Button>
          <Button
            size="icon-lg"
            className="size-14! rounded-full bg-red-500 text-white hover:bg-red-600"
            onClick={endCall}
          >
            <PhoneOff className="size-6" />
          </Button>
        </div>
      </div>
    </div>
  );
}
