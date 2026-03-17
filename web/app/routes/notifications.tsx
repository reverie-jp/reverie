import { Link } from "react-router";
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
  AvatarGroup,
} from "~/components/ui/avatar";
import { BottomNav } from "~/components/bottom-nav";
import { ComposeFab } from "~/components/compose-fab";
import {
  Heart,
  MessageCircle,
  Repeat2,
  UserPlus,
  Footprints,
  ChevronRight,
  Phone,
  FileText,
} from "lucide-react";
import { formatRelativeTime } from "~/components/post-card";

type NotificationType =
  | "like"
  | "reply"
  | "repost"
  | "follow"
  | "watch-post"
  | "watch-call";

interface UserRef {
  name: string;
  customId: string;
  avatarUrl?: string;
}

interface NotificationItem {
  id: string;
  type: NotificationType;
  users: UserRef[];
  createdAt: Date;
  postSnippet?: string;
  postId?: string;
  replyContent?: string;
  callName?: string;
}

const notificationStyle: Record<
  NotificationType,
  { icon: React.ElementType; color: string; bg: string; fill?: boolean }
> = {
  like: { icon: Heart, color: "text-pink-500", bg: "bg-pink-500/10", fill: true },
  reply: { icon: MessageCircle, color: "text-blue-500", bg: "bg-blue-500/10" },
  repost: { icon: Repeat2, color: "text-green-500", bg: "bg-green-500/10" },
  follow: { icon: UserPlus, color: "text-purple-500", bg: "bg-purple-500/10" },
  "watch-post": { icon: FileText, color: "text-amber-500", bg: "bg-amber-500/10" },
  "watch-call": { icon: Phone, color: "text-emerald-500", bg: "bg-emerald-500/10" },
};

const sampleNotifications: NotificationItem[] = [
  {
    id: "n1",
    type: "like",
    users: [
      { name: "田中太郎", customId: "tanaka" },
      { name: "渡辺大輔", customId: "daisuke_w" },
      { name: "木村拓也", customId: "takuya_k" },
    ],
    createdAt: new Date(Date.now() - 5 * 60_000),
    postSnippet: "今日はとてもいい天気ですね。散歩に行ってきました！",
    postId: "1",
  },
  {
    id: "n2",
    type: "reply",
    users: [{ name: "佐藤花子", customId: "hanako_s" }],
    createdAt: new Date(Date.now() - 15 * 60_000),
    postSnippet: "今日はとてもいい天気ですね。散歩に行ってきました！",
    postId: "1",
    replyContent: "いいですね！どこに行きましたか？",
  },
  {
    id: "n3",
    type: "follow",
    users: [
      { name: "鈴木一郎", customId: "ichiro_dev" },
      { name: "中村悠", customId: "yu_nkmr" },
    ],
    createdAt: new Date(Date.now() - 30 * 60_000),
  },
  {
    id: "n4",
    type: "repost",
    users: [
      { name: "山田美咲", customId: "misaki_y" },
      { name: "佐藤花子", customId: "hanako_s" },
    ],
    createdAt: new Date(Date.now() - 2 * 3_600_000),
    postSnippet:
      "React Routerの新しいバージョンを試してみたけど、かなり使いやすくなってる。",
    postId: "3",
  },
  {
    id: "n5",
    type: "like",
    users: [{ name: "渡辺大輔", customId: "daisuke_w" }],
    createdAt: new Date(Date.now() - 3 * 3_600_000),
    postSnippet:
      "React Routerの新しいバージョンを試してみたけど、かなり使いやすくなってる。",
    postId: "3",
  },
  {
    id: "n6",
    type: "watch-post",
    users: [{ name: "高橋健太", customId: "kenta_t" }],
    createdAt: new Date(Date.now() - 4 * 3_600_000),
    postSnippet:
      "データ分析の結果、今月のアクティブユーザーが20%増加しました。",
    postId: "p2",
  },
  {
    id: "n7",
    type: "watch-call",
    users: [{ name: "小林あおい", customId: "aoi_kb" }],
    createdAt: new Date(Date.now() - 5 * 3_600_000),
    callName: "デザインレビュー",
  },
  {
    id: "n8",
    type: "like",
    users: [
      { name: "木村拓也", customId: "takuya_k" },
      { name: "松本りな", customId: "rina_m" },
      { name: "井上翔", customId: "sho_inoue" },
      { name: "高橋健太", customId: "kenta_t" },
    ],
    createdAt: new Date(Date.now() - 6 * 3_600_000),
    postSnippet: "週末に映画を観に行きました。ストーリーが素晴らしかった！",
    postId: "4",
  },
  {
    id: "n9",
    type: "follow",
    users: [{ name: "森田陽介", customId: "yosuke_m" }],
    createdAt: new Date(Date.now() - 1 * 86_400_000),
  },
  {
    id: "n10",
    type: "reply",
    users: [{ name: "松本りな", customId: "rina_m" }],
    createdAt: new Date(Date.now() - 1 * 86_400_000),
    postSnippet:
      "プログラミングの勉強を始めて半年。少しずつ書けるようになってきた気がする。",
    postId: "5",
    replyContent: "すごいですね！応援してます！",
  },
  {
    id: "n11",
    type: "watch-post",
    users: [{ name: "田中太郎", customId: "tanaka" }],
    createdAt: new Date(Date.now() - 2 * 86_400_000),
    postSnippet:
      "TypeScriptの型パズル、難しいけど楽しい。最近はConditional Typesにハマってます。",
    postId: "7",
  },
  {
    id: "n12",
    type: "repost",
    users: [{ name: "佐藤花子", customId: "hanako_s" }],
    createdAt: new Date(Date.now() - 3 * 86_400_000),
    postSnippet: "週末に映画を観に行きました。ストーリーが素晴らしかった！",
    postId: "4",
  },
];

function NotificationRow({
  notification,
}: {
  notification: NotificationItem;
}) {
  const style = notificationStyle[notification.type];
  const Icon = style.icon;
  const maxAvatars = 5;

  const linkTo =
    notification.type === "follow"
      ? "/users/me/connections?tab=followers"
      : notification.postId
        ? `/posts/${notification.postId}`
        : "#";

  return (
    <Link
      to={linkTo}
      className="flex gap-3 px-4 py-3 border-b hover:bg-muted/30 transition-colors"
    >
      {/* Left: type icon */}
      <div className="shrink-0 pt-0.5">
        <div
          className={`size-8 rounded-full ${style.bg} flex items-center justify-center`}
        >
          <Icon
            className={`size-4 ${style.color} ${style.fill ? "fill-current" : ""}`}
          />
        </div>
      </div>

      {/* Right: content */}
      <div className="flex-1 min-w-0">
        {/* Avatar group + time */}
        <div className="flex items-start justify-between gap-2 mb-1.5">
        <AvatarGroup className="-space-x-1.5">
          {notification.users.slice(0, maxAvatars).map((user) => (
            <Avatar key={user.customId} className="size-7!">
              <AvatarImage src={user.avatarUrl} alt={user.name} />
              <AvatarFallback className="text-[10px]">
                {user.name.slice(0, 1)}
              </AvatarFallback>
            </Avatar>
          ))}
          {notification.users.length > maxAvatars && (
            <div className="relative flex size-7 shrink-0 items-center justify-center rounded-full bg-muted text-[10px] text-muted-foreground ring-2 ring-background">
              +{notification.users.length - maxAvatars}
            </div>
          )}
        </AvatarGroup>
        <span className="shrink-0 text-xs text-muted-foreground">
          {formatRelativeTime(notification.createdAt)}
        </span>
        </div>

        {/* Text */}
        <p className="text-sm leading-snug">
          <span className="font-bold">{notification.users[0].name}</span>
          {notification.users.length === 2 && (
            <>
              <span className="text-muted-foreground">、</span>
              <span className="font-bold">{notification.users[1].name}</span>
            </>
          )}
          {notification.users.length > 2 && (
            <span className="text-muted-foreground">
              、他{notification.users.length - 1}人
            </span>
          )}
          <span className="text-muted-foreground">
            {(() => {
              switch (notification.type) {
                case "like":
                  return "があなたの投稿にいいねしました";
                case "reply":
                  return "があなたの投稿に返信しました";
                case "repost":
                  return "があなたの投稿を再投稿しました";
                case "follow":
                  return "があなたをフォローしました";
                case "watch-post":
                  return "が新しい投稿をしました";
                case "watch-call":
                  return "が通話を開始しました";
              }
            })()}
          </span>
        </p>

        {/* Reply content */}
        {notification.type === "reply" && notification.replyContent && (
          <p className="mt-1 text-sm line-clamp-2">
            {notification.replyContent}
          </p>
        )}

        {/* Post snippet */}
        {notification.postSnippet && (
          <p className="mt-1 text-xs text-muted-foreground line-clamp-1">
            {notification.postSnippet}
          </p>
        )}

        {/* Call name */}
        {notification.type === "watch-call" && notification.callName && (
          <p className="mt-1 text-xs text-muted-foreground">
            「{notification.callName}」
          </p>
        )}
      </div>

    </Link>
  );
}

export default function Notifications() {
  return (
    <div className="w-full min-h-full flex flex-col">
      {/* Header */}
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center justify-between px-4 h-14">
          <h1 className="text-base font-bold">通知</h1>
          <Link
            to="/footprints"
            className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            <Footprints className="size-4" />
            <span>足あと</span>
            <ChevronRight className="size-4" />
          </Link>
        </div>
      </div>

      {/* Notification list */}
      <div className="flex-1">
        {sampleNotifications.map((notification) => (
          <NotificationRow key={notification.id} notification={notification} />
        ))}
      </div>

      <ComposeFab />
      <BottomNav />
    </div>
  );
}
