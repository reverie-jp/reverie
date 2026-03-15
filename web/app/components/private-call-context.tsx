import { createContext, useContext, useState, useEffect, useRef } from "react";

interface PrivateCallParticipant {
  name: string;
  customId: string;
  avatarUrl: string;
}

interface ActivePrivateCall {
  chatId: string;
  participant: PrivateCallParticipant;
  type: "audio" | "video";
  micOn: boolean;
  speakerOn: boolean;
  cameraOn: boolean;
  elapsed: number;
}

interface PrivateCallContextValue {
  activeCall: ActivePrivateCall | null;
  isMinimized: boolean;
  startCall: (
    chatId: string,
    participant: PrivateCallParticipant,
    type: "audio" | "video",
  ) => void;
  endCall: () => void;
  minimize: () => void;
  maximize: () => void;
  setMicOn: (on: boolean) => void;
  setSpeakerOn: (on: boolean) => void;
  setCameraOn: (on: boolean) => void;
}

const PrivateCallContext = createContext<PrivateCallContextValue | null>(null);

export function PrivateCallProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [activeCall, setActiveCall] = useState<ActivePrivateCall | null>(null);
  const [isMinimized, setIsMinimized] = useState(false);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    if (activeCall) {
      timerRef.current = setInterval(() => {
        setActiveCall((prev) =>
          prev ? { ...prev, elapsed: prev.elapsed + 1 } : null,
        );
      }, 1000);
    }
    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
  }, [activeCall?.chatId]);

  const startCall = (
    chatId: string,
    participant: PrivateCallParticipant,
    type: "audio" | "video",
  ) => {
    setActiveCall({
      chatId,
      participant,
      type,
      micOn: true,
      speakerOn: true,
      cameraOn: type === "video",
      elapsed: 0,
    });
    setIsMinimized(false);
  };

  const endCall = () => {
    setActiveCall(null);
    setIsMinimized(false);
    if (timerRef.current) clearInterval(timerRef.current);
  };

  const minimize = () => setIsMinimized(true);
  const maximize = () => setIsMinimized(false);

  const setMicOn = (on: boolean) =>
    setActiveCall((prev) => (prev ? { ...prev, micOn: on } : null));
  const setSpeakerOn = (on: boolean) =>
    setActiveCall((prev) => (prev ? { ...prev, speakerOn: on } : null));
  const setCameraOn = (on: boolean) =>
    setActiveCall((prev) => (prev ? { ...prev, cameraOn: on } : null));

  return (
    <PrivateCallContext value={{
      activeCall,
      isMinimized,
      startCall,
      endCall,
      minimize,
      maximize,
      setMicOn,
      setSpeakerOn,
      setCameraOn,
    }}>
      {children}
    </PrivateCallContext>
  );
}

export function usePrivateCall() {
  const ctx = useContext(PrivateCallContext);
  if (!ctx)
    throw new Error("usePrivateCall must be used within PrivateCallProvider");
  return ctx;
}

export function formatCallDuration(seconds: number): string {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
}
