import { useState } from "react";
import { LogOut } from "lucide-react";
import { useCall } from "~/components/call-context";
import { formatCallDuration } from "~/components/private-call-context";
import { CallBar } from "./call-bar";

export function GroupCallBar() {
  const { activeCall, isMinimized, maximize, leaveCall } = useCall();
  const [isMuted, setIsMuted] = useState(false);

  if (!activeCall || !isMinimized) return null;

  return (
    <CallBar
      title={activeCall.name}
      subtitle={`${activeCall.participants.length}人参加中 · ${formatCallDuration(activeCall.elapsed)}`}
      micOn={!isMuted}
      onToggleMic={() => setIsMuted(!isMuted)}
      onTapCenter={maximize}
      onEnd={leaveCall}
      EndIcon={LogOut}
    />
  );
}
