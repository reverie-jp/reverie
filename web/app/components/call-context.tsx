import { createContext, useContext, useState } from "react";
import type { Call, CallParticipant } from "~/components/call-list";

export type CallVisibility =
  | "すべてのユーザー"
  | "フォロワーのみ"
  | "相互フォローのみ"
  | "招待した人のみ";

interface ActiveCall {
  name: string;
  type: "audio" | "video";
  visibility: CallVisibility;
  participants: CallParticipant[];
}

interface CallContextValue {
  activeCall: ActiveCall | null;
  isMinimized: boolean;
  joinCall: (call: Call) => void;
  createCall: (name: string, type: "audio" | "video", visibility: CallVisibility) => void;
  setVisibility: (visibility: CallVisibility) => void;
  minimize: () => void;
  maximize: () => void;
  leaveCall: () => void;
}

const CallContext = createContext<CallContextValue | null>(null);

export function CallProvider({ children }: { children: React.ReactNode }) {
  const [activeCall, setActiveCall] = useState<ActiveCall | null>(null);
  const [isMinimized, setIsMinimized] = useState(false);

  const joinCall = (call: Call) => {
    setActiveCall({
      name: call.name,
      type: call.type,
      visibility: "すべてのユーザー",
      participants: [
        ...call.participants,
        { name: "自分", customId: "me", avatarUrl: "" },
      ],
    });
    setIsMinimized(false);
  };

  const createCall = (name: string, type: "audio" | "video", visibility: CallVisibility) => {
    setActiveCall({
      name,
      type,
      visibility,
      participants: [{ name: "自分", customId: "me", avatarUrl: "" }],
    });
    setIsMinimized(false);
  };

  const setVisibility = (visibility: CallVisibility) => {
    setActiveCall((prev) => prev ? { ...prev, visibility } : null);
  };

  const minimize = () => setIsMinimized(true);
  const maximize = () => setIsMinimized(false);
  const leaveCall = () => {
    setActiveCall(null);
    setIsMinimized(false);
  };

  return (
    <CallContext value={{
      activeCall,
      isMinimized,
      joinCall,
      createCall,
      setVisibility,
      minimize,
      maximize,
      leaveCall,
    }}>
      {children}
    </CallContext>
  );
}

export function useCall() {
  const ctx = useContext(CallContext);
  if (!ctx) throw new Error("useCall must be used within CallProvider");
  return ctx;
}
