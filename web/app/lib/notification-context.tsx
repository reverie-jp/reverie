import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { ConnectError, Code } from "@connectrpc/connect";
import { toast } from "sonner";
import {
  Notification,
  NotificationType,
} from "./gen/notification/v1/notification_pb";
import {
  eventClient,
  notificationClient,
  tokenStore,
} from "./api-client";

type NotificationContextValue = {
  notifications: Notification[];
  unreadCount: number;
  loading: boolean;
  refresh: () => Promise<void>;
  markAllRead: () => Promise<void>;
  markRead: (name: string) => Promise<void>;
};

const NotificationContext = createContext<NotificationContextValue | null>(
  null,
);

export function useNotifications() {
  const ctx = useContext(NotificationContext);
  if (!ctx)
    throw new Error("useNotifications must be used within NotificationProvider");
  return ctx;
}

const STREAM_RECONNECT_MIN_MS = 1_000;
const STREAM_RECONNECT_MAX_MS = 30_000;
// Server sends a KeepAlive envelope every ~15s. If we go 2x that without any
// frame, the connection is effectively dead — force a reconnect.
const STREAM_STALE_MS = 35_000;

function notificationMessage(n: Notification): string {
  const actor = n.actor?.displayName || n.actor?.customId || "誰か";
  switch (n.type) {
    case NotificationType.USER_FOLLOWED:
      return `${actor} があなたをフォローしました`;
    case NotificationType.FOLLOWING_USER_CALL_STARTED:
      return `${actor} が通話を開始しました`;
    default:
      return "新しい通知があります";
  }
}

export function NotificationProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() =>
    Boolean(
      typeof window !== "undefined" ? tokenStore.getAccessToken() : null,
    ),
  );
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(false);

  const seenEventIdsRef = useRef<Set<string>>(new Set());
  const streamAbortRef = useRef<AbortController | null>(null);
  const reconnectTimerRef = useRef<number | null>(null);
  const reconnectAttemptsRef = useRef(0);

  // Keep a ref to refresh so effect dependencies stay minimal.
  const refreshRef = useRef<() => Promise<void>>(async () => {});

  const refresh = useCallback(async () => {
    if (!tokenStore.getAccessToken()) return;
    setLoading(true);
    try {
      const [listRes, countRes] = await Promise.all([
        notificationClient.listNotifications({ pageSize: 30 }),
        notificationClient.getUnreadNotificationCount({}),
      ]);
      setNotifications(listRes.notifications);
      setUnreadCount(countRes.count);
      seenEventIdsRef.current.clear();
      for (const n of listRes.notifications) {
        // Use the notification's trailing ULID (event_id on the envelope
        // matches the notification's DB id) for dedup.
        const id = n.name.split("/").pop();
        if (id) seenEventIdsRef.current.add(id);
      }
    } catch (err) {
      if (err instanceof ConnectError && err.code === Code.Unauthenticated) {
        return;
      }
      console.warn("ListNotifications failed:", err);
    } finally {
      setLoading(false);
    }
  }, []);
  refreshRef.current = refresh;

  const markAllRead = useCallback(async () => {
    if (!tokenStore.getAccessToken()) return;
    try {
      await notificationClient.markNotificationsRead({ names: [] });
      setUnreadCount(0);
      setNotifications((prev) =>
        prev.map((n) =>
          n.readTime
            ? n
            : new Notification({
                ...n,
                readTime: { seconds: BigInt(Math.floor(Date.now() / 1000)), nanos: 0 },
              }),
        ),
      );
    } catch (err) {
      console.warn("MarkNotificationsRead failed:", err);
    }
  }, []);

  const markRead = useCallback(async (name: string) => {
    if (!tokenStore.getAccessToken()) return;
    try {
      await notificationClient.markNotificationsRead({ names: [name] });
      setNotifications((prev) =>
        prev.map((n) =>
          n.name === name && !n.readTime
            ? new Notification({
                ...n,
                readTime: { seconds: BigInt(Math.floor(Date.now() / 1000)), nanos: 0 },
              })
            : n,
        ),
      );
      setUnreadCount((c) => Math.max(0, c - 1));
    } catch (err) {
      console.warn("MarkNotificationsRead(one) failed:", err);
    }
  }, []);

  const handleIncoming = useCallback((n: Notification) => {
    const id = n.name.split("/").pop();
    if (!id) return;
    if (seenEventIdsRef.current.has(id)) return;
    seenEventIdsRef.current.add(id);

    setNotifications((prev) => [n, ...prev].slice(0, 50));
    setUnreadCount((c) => c + 1);
    toast(notificationMessage(n));
  }, []);

  // Auth change detection. api-client's tokenStore uses localStorage directly,
  // so we poll on mount + listen for storage events (for multi-tab) and a
  // custom event we fire from account flows. Keep this small — nothing else
  // needs a full auth context right now.
  useEffect(() => {
    const sync = () => {
      setIsAuthenticated(Boolean(tokenStore.getAccessToken()));
    };
    window.addEventListener("storage", sync);
    window.addEventListener("reverie:auth_changed", sync);
    return () => {
      window.removeEventListener("storage", sync);
      window.removeEventListener("reverie:auth_changed", sync);
    };
  }, []);

  // Open or close the stream in response to auth state.
  useEffect(() => {
    if (!isAuthenticated) {
      setNotifications([]);
      setUnreadCount(0);
      seenEventIdsRef.current.clear();
      if (streamAbortRef.current) {
        streamAbortRef.current.abort();
        streamAbortRef.current = null;
      }
      if (reconnectTimerRef.current !== null) {
        window.clearTimeout(reconnectTimerRef.current);
        reconnectTimerRef.current = null;
      }
      reconnectAttemptsRef.current = 0;
      return;
    }

    let cancelled = false;

    const openStream = async () => {
      if (cancelled) return;
      // Fetch canonical state before subscribing so any messages lost while
      // disconnected are picked up via the DB. Stream is hint-only.
      await refreshRef.current();
      if (cancelled) return;

      const abort = new AbortController();
      streamAbortRef.current = abort;

      // Watchdog: if no frame (data or keepalive) arrives within STREAM_STALE_MS,
      // the connection is silently dead — force-abort so the catch block below
      // triggers a reconnect.
      let watchdog = window.setTimeout(() => abort.abort(), STREAM_STALE_MS);
      const kickWatchdog = () => {
        window.clearTimeout(watchdog);
        watchdog = window.setTimeout(() => abort.abort(), STREAM_STALE_MS);
      };

      try {
        const stream = eventClient.streamEvents(
          {},
          { signal: abort.signal },
        );
        reconnectAttemptsRef.current = 0;
        for await (const envelope of stream) {
          if (cancelled) break;
          kickWatchdog();
          switch (envelope.payload.case) {
            case "notificationCreated": {
              const n = envelope.payload.value.notification;
              if (n) handleIncoming(n);
              break;
            }
            case "keepAlive":
              // No-op: frame arrival alone already reset the watchdog above.
              break;
          }
        }
      } catch (err) {
        if (cancelled || abort.signal.aborted) {
          // aborted intentionally by us (unmount or watchdog). Fall through to
          // the reconnect path only if still authenticated.
        } else {
          console.warn("StreamEvents disconnected:", err);
        }
      } finally {
        window.clearTimeout(watchdog);
      }

      if (cancelled) return;
      // Exponential backoff capped. The stream-as-hint design tolerates gaps
      // because refreshRef runs on reconnect.
      const attempt = reconnectAttemptsRef.current++;
      const delay = Math.min(
        STREAM_RECONNECT_MAX_MS,
        STREAM_RECONNECT_MIN_MS * Math.pow(2, attempt),
      );
      reconnectTimerRef.current = window.setTimeout(() => {
        reconnectTimerRef.current = null;
        void openStream();
      }, delay);
    };

    void openStream();

    return () => {
      cancelled = true;
      if (streamAbortRef.current) {
        streamAbortRef.current.abort();
        streamAbortRef.current = null;
      }
      if (reconnectTimerRef.current !== null) {
        window.clearTimeout(reconnectTimerRef.current);
        reconnectTimerRef.current = null;
      }
    };
  }, [isAuthenticated, handleIncoming]);

  const value: NotificationContextValue = {
    notifications,
    unreadCount,
    loading,
    refresh,
    markAllRead,
    markRead,
  };

  return (
    <NotificationContext.Provider value={value}>
      {children}
    </NotificationContext.Provider>
  );
}
