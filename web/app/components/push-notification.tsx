import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { useLocation, useNavigate } from "react-router";
import {
  Heart,
  MessageCircle,
  Repeat2,
  UserPlus,
  Phone,
  FileText,
  X,
  Footprints,
  Mail,
  LogIn,
  LogOut,
} from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { useCall } from "~/components/call-context";

type PushNotificationType =
  | "like"
  | "reply"
  | "repost"
  | "follow"
  | "chat"
  | "footprint"
  | "watch-post"
  | "watch-call"
  | "call-join"
  | "call-leave";

export interface PushNotification {
  id: string;
  type: PushNotificationType;
  userName: string;
  userAvatarUrl?: string;
  message: string;
  link?: string;
}

const iconMap: Record<
  PushNotificationType,
  { icon: React.ElementType; color: string; bg: string; fill?: boolean }
> = {
  like: {
    icon: Heart,
    color: "text-pink-500",
    bg: "bg-pink-500/10",
    fill: true,
  },
  reply: { icon: MessageCircle, color: "text-blue-500", bg: "bg-blue-500/10" },
  repost: { icon: Repeat2, color: "text-green-500", bg: "bg-green-500/10" },
  follow: { icon: UserPlus, color: "text-purple-500", bg: "bg-purple-500/10" },
  chat: { icon: Mail, color: "text-sky-500", bg: "bg-sky-500/10" },
  footprint: {
    icon: Footprints,
    color: "text-orange-500",
    bg: "bg-orange-500/10",
  },
  "watch-post": {
    icon: FileText,
    color: "text-amber-500",
    bg: "bg-amber-500/10",
  },
  "watch-call": {
    icon: Phone,
    color: "text-emerald-500",
    bg: "bg-emerald-500/10",
  },
  "call-join": { icon: LogIn, color: "text-teal-500", bg: "bg-teal-500/10" },
  "call-leave": { icon: LogOut, color: "text-red-500", bg: "bg-red-500/10" },
};

const DISPLAY_DURATION = 4000;
const DEMO_INTERVAL = 10_000;

const demoNotifications: Omit<PushNotification, "id">[] = [
  {
    type: "like",
    userName: "佐藤花子",
    message: "あなたの投稿にいいねしました",
    link: "/posts/1",
  },
  {
    type: "reply",
    userName: "田中太郎",
    message: "「React Routerの新しいバージョン...」に返信しました",
    link: "/posts/3",
  },
  {
    type: "repost",
    userName: "高橋健太",
    message: "あなたの投稿を再投稿しました",
    link: "/posts/1",
  },
  {
    type: "follow",
    userName: "山田美咲",
    message: "あなたをフォローしました",
    link: "/users/me/connections?tab=followers",
  },
  {
    type: "chat",
    userName: "鈴木一郎",
    message: "お疲れ様です！明日のミーティングの件ですが...",
    link: "/chat/1",
  },
  {
    type: "chat",
    userName: "中村悠",
    message: "写真送りました！確認してください",
    link: "/chat/2",
  },
  {
    type: "footprint",
    userName: "小林あおい",
    message: "あなたのプロフィールを訪問しました",
    link: "/footprints",
  },
  {
    type: "footprint",
    userName: "木村拓也",
    message: "あなたのプロフィールを訪問しました",
    link: "/footprints",
  },
  {
    type: "watch-post",
    userName: "渡辺大輔",
    message: "新しい投稿を作成しました",
    link: "/posts/p3",
  },
  {
    type: "watch-post",
    userName: "伊藤さくら",
    message: "新しい投稿を作成しました",
    link: "/posts/p4",
  },
  {
    type: "watch-call",
    userName: "松本りな",
    message: "通話「作業通話」を開始しました",
    link: "/calls",
  },
  {
    type: "watch-call",
    userName: "井上翔",
    message: "通話「LT大会」を開始しました",
    link: "/calls",
  },
];

interface PushNotificationContextValue {
  push: (notification: Omit<PushNotification, "id">) => void;
}

const PushNotificationContext = createContext<PushNotificationContextValue>({
  push: () => {},
});

export function usePushNotification() {
  return useContext(PushNotificationContext);
}

interface ActiveNotification extends PushNotification {
  visible: boolean;
  startY: number | null;
  currentY: number | null;
}

export function PushNotificationProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const location = useLocation();
  const isLoggedOut = location.pathname === "/login";
  const [notifications, setNotifications] = useState<ActiveNotification[]>([]);
  const idCounter = useRef(0);
  const demoIndex = useRef(0);

  const dismiss = useCallback((id: string) => {
    setNotifications((prev) =>
      prev.map((n) => (n.id === id ? { ...n, visible: false } : n)),
    );
    setTimeout(() => {
      setNotifications((prev) => prev.filter((n) => n.id !== id));
    }, 300);
  }, []);

  const push = useCallback(
    (notification: Omit<PushNotification, "id">) => {
      const id = `push-${++idCounter.current}`;
      setNotifications((prev) => [
        ...prev,
        { ...notification, id, visible: false, startY: null, currentY: null },
      ]);
      requestAnimationFrame(() => {
        requestAnimationFrame(() => {
          setNotifications((prev) =>
            prev.map((n) => (n.id === id ? { ...n, visible: true } : n)),
          );
        });
      });
      setTimeout(() => dismiss(id), DISPLAY_DURATION);
    },
    [dismiss],
  );

  // Global demo: random notification every 10 seconds (only when logged in)
  useEffect(() => {
    if (isLoggedOut) return;

    const initial = setTimeout(() => {
      const idx = Math.floor(Math.random() * demoNotifications.length);
      push(demoNotifications[idx]);
      demoIndex.current = idx;
    }, 3000);

    const interval = setInterval(() => {
      let idx: number;
      do {
        idx = Math.floor(Math.random() * demoNotifications.length);
      } while (idx === demoIndex.current && demoNotifications.length > 1);
      demoIndex.current = idx;
      push(demoNotifications[idx]);
    }, DEMO_INTERVAL);

    return () => {
      clearTimeout(initial);
      clearInterval(interval);
    };
  }, [push, isLoggedOut]);

  return (
    <PushNotificationContext.Provider value={{ push }}>
      {children}
      <div className="fixed top-0 left-0 right-0 z-100 pointer-events-none flex flex-col items-center gap-2 p-3">
        {notifications.map((n) => (
          <PushNotificationToast
            key={n.id}
            notification={n}
            onDismiss={() => dismiss(n.id)}
            onUpdateTouch={(startY, currentY) => {
              setNotifications((prev) =>
                prev.map((item) =>
                  item.id === n.id ? { ...item, startY, currentY } : item,
                ),
              );
            }}
          />
        ))}
      </div>
    </PushNotificationContext.Provider>
  );
}

function PushNotificationToast({
  notification,
  onDismiss,
  onUpdateTouch,
}: {
  notification: ActiveNotification;
  onDismiss: () => void;
  onUpdateTouch: (startY: number | null, currentY: number | null) => void;
}) {
  const navigate = useNavigate();
  const { icon: Icon, color, bg, fill } = iconMap[notification.type];

  const dragOffset =
    notification.startY !== null && notification.currentY !== null
      ? Math.min(0, notification.currentY - notification.startY)
      : 0;

  const translateY = notification.visible ? dragOffset : -100;

  return (
    <div
      className="pointer-events-auto w-full max-w-sm transition-all duration-300 ease-out"
      style={{
        transform: `translateY(${translateY}px)`,
        opacity: notification.visible ? 1 : 0,
      }}
      onTouchStart={(e) => {
        onUpdateTouch(e.touches[0].clientY, e.touches[0].clientY);
      }}
      onTouchMove={(e) => {
        if (notification.startY !== null) {
          onUpdateTouch(notification.startY, e.touches[0].clientY);
        }
      }}
      onTouchEnd={() => {
        if (dragOffset < -50) {
          onDismiss();
        } else {
          onUpdateTouch(null, null);
        }
      }}
    >
      <div
        className="bg-card border border-border rounded-2xl shadow-lg px-4 py-3 flex items-center gap-3 cursor-pointer active:bg-muted/50 transition-colors"
        onClick={() => {
          if (notification.link) {
            navigate(notification.link);
          }
          onDismiss();
        }}
      >
        <div className="relative shrink-0">
          <Avatar size="default">
            <AvatarImage src={notification.userAvatarUrl} />
            <AvatarFallback>{notification.userName[0]}</AvatarFallback>
          </Avatar>
          <div
            className={`absolute -bottom-1 -right-1 size-5 rounded-full flex items-center justify-center ${bg}`}
          >
            <Icon
              className={`size-3 ${color}`}
              fill={fill ? "currentColor" : "none"}
            />
          </div>
        </div>
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium truncate">
            {notification.userName}
          </p>
          <p className="text-xs text-muted-foreground truncate">
            {notification.message}
          </p>
        </div>
        <button
          onClick={(e) => {
            e.stopPropagation();
            onDismiss();
          }}
          className="shrink-0 text-muted-foreground hover:text-foreground transition-colors p-1"
        >
          <X className="size-4" />
        </button>
      </div>
    </div>
  );
}

const demoCallUsers = [
  "田中太郎",
  "佐藤花子",
  "鈴木一郎",
  "高橋健太",
  "中村悠",
  "小林あおい",
  "渡辺大輔",
  "山田美咲",
];

const CALL_NOTIFICATION_INTERVAL = 8000;

export function CallNotificationBridge() {
  const { activeCall, isMinimized } = useCall();
  const { push } = usePushNotification();
  const indexRef = useRef(0);

  useEffect(() => {
    if (!activeCall || !isMinimized) return;

    const callName = activeCall.name;

    const interval = setInterval(() => {
      const user = demoCallUsers[indexRef.current % demoCallUsers.length];
      const isJoin = Math.random() > 0.3;
      indexRef.current++;

      push({
        type: isJoin ? "call-join" : "call-leave",
        userName: user,
        message: isJoin
          ? `「${callName}」に参加しました`
          : `「${callName}」から退出しました`,
      });
    }, CALL_NOTIFICATION_INTERVAL);

    return () => clearInterval(interval);
  }, [activeCall?.name, isMinimized, push]);

  return null;
}
