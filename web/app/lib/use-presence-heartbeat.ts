import { useEffect } from "react";
import {
  AUTH_CHANGED_EVENT,
  presenceClient,
  tokenStore,
} from "~/lib/api-client";

// Heartbeat cadence. Server treats users as offline if last_seen_time is
// older than entity.UserPresenceStaleSeconds (60s). 20s interval = 3 beats
// within that window, so two missed beats (transient network / slow RTT)
// won't flip the dot to offline.
const HEARTBEAT_INTERVAL_MS = 20_000;

// usePresenceHeartbeat keeps the caller's last_seen_time fresh while
// authenticated and the tab is visible. Mount once at the app root.
export function usePresenceHeartbeat() {
  useEffect(() => {
    let intervalId: number | null = null;

    const sendHeartbeat = () => {
      void presenceClient.heartbeat({}).catch((err) => {
        console.warn("Heartbeat failed:", err);
      });
    };

    const start = () => {
      if (intervalId !== null) return;
      if (!tokenStore.getAccessToken()) return;
      if (typeof document !== "undefined" && document.hidden) return;
      sendHeartbeat();
      intervalId = window.setInterval(sendHeartbeat, HEARTBEAT_INTERVAL_MS);
    };

    const stop = () => {
      if (intervalId !== null) {
        window.clearInterval(intervalId);
        intervalId = null;
      }
    };

    const onVisibilityChange = () => {
      if (document.hidden) {
        stop();
      } else {
        start();
      }
    };

    const onAuthChange = () => {
      stop();
      start();
    };

    start();
    document.addEventListener("visibilitychange", onVisibilityChange);
    window.addEventListener(AUTH_CHANGED_EVENT, onAuthChange);
    window.addEventListener("storage", onAuthChange);

    return () => {
      stop();
      document.removeEventListener("visibilitychange", onVisibilityChange);
      window.removeEventListener(AUTH_CHANGED_EVENT, onAuthChange);
      window.removeEventListener("storage", onAuthChange);
    };
  }, []);
}
