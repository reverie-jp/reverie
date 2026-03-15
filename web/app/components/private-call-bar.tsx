import { PhoneOff } from "lucide-react";
import { usePrivateCall, formatCallDuration } from "./private-call-context";
import { CallBar } from "./call-bar";

export function PrivateCallBar() {
  const { activeCall, isMinimized, maximize, endCall, setMicOn } =
    usePrivateCall();

  if (!activeCall || !isMinimized) return null;

  return (
    <CallBar
      title={activeCall.participant.name}
      subtitle={formatCallDuration(activeCall.elapsed)}
      micOn={activeCall.micOn}
      onToggleMic={() => setMicOn(!activeCall.micOn)}
      onTapCenter={maximize}
      onEnd={endCall}
      EndIcon={PhoneOff}
    />
  );
}
