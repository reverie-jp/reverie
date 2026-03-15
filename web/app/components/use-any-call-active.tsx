import { useCall } from "~/components/call-context";
import { usePrivateCall } from "~/components/private-call-context";

interface AnyCallState {
  isInCall: boolean;
  currentCallName: string | null;
  endCurrentCall: () => void;
}

export function useAnyCallActive(): AnyCallState {
  const { activeCall: groupCall, leaveCall } = useCall();
  const { activeCall: privateCall, endCall } = usePrivateCall();

  if (groupCall) {
    return {
      isInCall: true,
      currentCallName: groupCall.name,
      endCurrentCall: leaveCall,
    };
  }
  if (privateCall) {
    return {
      isInCall: true,
      currentCallName: `${privateCall.participant.name}との通話`,
      endCurrentCall: endCall,
    };
  }
  return {
    isInCall: false,
    currentCallName: null,
    endCurrentCall: () => {},
  };
}
