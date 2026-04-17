import { Link } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { FollowButton } from "~/components/follow-button";
import { BottomNav } from "~/components/bottom-nav";
import { ArrowLeft, Eye, Clock } from "lucide-react";
import { formatRelativeTime } from "~/components/post-card";

interface FootprintUser {
  id: string;
  name: string;
  customId: string;
  avatarUrl?: string;
  isFollowing: boolean;
  followsYou: boolean;
}

interface FootprintItem {
  id: string;
  user: FootprintUser;
  visitedAt: Date;
}

const receivedFootprints: FootprintItem[] = [
  {
    id: "f1",
    user: {
      id: "user-tanaka",
      name: "田中太郎",
      customId: "tanaka",
      isFollowing: true,
      followsYou: true,
    },
    visitedAt: new Date(Date.now() - 3 * 60_000),
  },
  {
    id: "f2",
    user: {
      id: "user-ichiro_dev",
      name: "鈴木一郎",
      customId: "ichiro_dev",
      isFollowing: false,
      followsYou: true,
    },
    visitedAt: new Date(Date.now() - 25 * 60_000),
  },
  {
    id: "f3",
    user: {
      id: "user-hanako_s",
      name: "佐藤花子",
      customId: "hanako_s",
      isFollowing: true,
      followsYou: true,
    },
    visitedAt: new Date(Date.now() - 2 * 3_600_000),
  },
  {
    id: "f4",
    user: {
      id: "user-kenta_t",
      name: "高橋健太",
      customId: "kenta_t",
      isFollowing: false,
      followsYou: true,
    },
    visitedAt: new Date(Date.now() - 4 * 3_600_000),
  },
  {
    id: "f5",
    user: {
      id: "user-takuya_k",
      name: "木村拓也",
      customId: "takuya_k",
      isFollowing: true,
      followsYou: true,
    },
    visitedAt: new Date(Date.now() - 8 * 3_600_000),
  },
  {
    id: "f6",
    user: {
      id: "user-yosuke_m",
      name: "森田陽介",
      customId: "yosuke_m",
      isFollowing: false,
      followsYou: false,
    },
    visitedAt: new Date(Date.now() - 1 * 86_400_000),
  },
  {
    id: "f7",
    user: {
      id: "user-mari_k",
      name: "川口真理",
      customId: "mari_k",
      isFollowing: false,
      followsYou: false,
    },
    visitedAt: new Date(Date.now() - 1.5 * 86_400_000),
  },
  {
    id: "f8",
    user: {
      id: "user-yu_nkmr",
      name: "中村悠",
      customId: "yu_nkmr",
      isFollowing: false,
      followsYou: true,
    },
    visitedAt: new Date(Date.now() - 2 * 86_400_000),
  },
];

const sentFootprints: FootprintItem[] = [
  {
    id: "s1",
    user: {
      id: "user-hanako_s",
      name: "佐藤花子",
      customId: "hanako_s",
      isFollowing: true,
      followsYou: true,
    },
    visitedAt: new Date(Date.now() - 10 * 60_000),
  },
  {
    id: "s2",
    user: {
      id: "user-misaki_y",
      name: "山田美咲",
      customId: "misaki_y",
      isFollowing: true,
      followsYou: false,
    },
    visitedAt: new Date(Date.now() - 1 * 3_600_000),
  },
  {
    id: "s3",
    user: {
      id: "user-daisuke_w",
      name: "渡辺大輔",
      customId: "daisuke_w",
      isFollowing: true,
      followsYou: true,
    },
    visitedAt: new Date(Date.now() - 5 * 3_600_000),
  },
  {
    id: "s4",
    user: {
      id: "user-rina_m",
      name: "松本りな",
      customId: "rina_m",
      isFollowing: false,
      followsYou: false,
    },
    visitedAt: new Date(Date.now() - 12 * 3_600_000),
  },
  {
    id: "s5",
    user: {
      id: "user-sho_inoue",
      name: "井上翔",
      customId: "sho_inoue",
      isFollowing: false,
      followsYou: false,
    },
    visitedAt: new Date(Date.now() - 2 * 86_400_000),
  },
  {
    id: "s6",
    user: {
      id: "user-tanaka",
      name: "田中太郎",
      customId: "tanaka",
      isFollowing: true,
      followsYou: true,
    },
    visitedAt: new Date(Date.now() - 3 * 86_400_000),
  },
];

function FootprintRow({ item }: { item: FootprintItem }) {
  const { user } = item;

  return (
    <Link
      to={`/users/${user.customId}`}
      className="flex items-center gap-3 px-4 py-3 border-b hover:bg-muted/30 transition-colors"
    >
      <Avatar className="size-11 shrink-0">
        <AvatarImage src={user.avatarUrl} alt={user.name} />
        <AvatarFallback>{user.name.slice(0, 2)}</AvatarFallback>
      </Avatar>

      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-1.5">
          <p className="text-sm font-bold truncate">{user.name}</p>
        </div>
        <p className="text-xs text-muted-foreground">
          訪れた日時: {formatRelativeTime(item.visitedAt)}
        </p>
      </div>

      <FollowButton
        userId={user.id}
        customId={user.customId}
        initialFollowing={user.isFollowing}
        followsYou={user.followsYou}
      />
    </Link>
  );
}

function EmptyState({ type }: { type: "received" | "sent" }) {
  return (
    <div className="flex flex-col items-center justify-center py-20 px-4 text-center">
      <div className="size-16 rounded-full bg-muted flex items-center justify-center mb-4">
        <Eye className="size-8 text-muted-foreground" />
      </div>
      <p className="text-sm font-medium">
        {type === "received"
          ? "まだ足あとはありません"
          : "まだ誰のプロフィールも見ていません"}
      </p>
      <p className="text-xs text-muted-foreground mt-1">
        {type === "received"
          ? "あなたのプロフィールを訪れた人がここに表示されます"
          : "他のユーザーのプロフィールを訪れると記録されます"}
      </p>
    </div>
  );
}

export default function FootprintsPage() {
  return (
    <div className="w-full min-h-full flex flex-col">
      {/* Header */}
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-14">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <h1 className="text-base font-bold">足あと</h1>
        </div>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="received" className="gap-0 flex-1">
        <div className="sticky top-14 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
          <TabsList variant="line" className="w-full h-12">
            <TabsTrigger value="received">
              <Eye className="size-3.5" />
              あなたへの足あと
            </TabsTrigger>
            <TabsTrigger value="sent">
              <Clock className="size-3.5" />
              あなたがつけた足あと
            </TabsTrigger>
          </TabsList>
        </div>

        <TabsContent value="received">
          {receivedFootprints.length === 0 ? (
            <EmptyState type="received" />
          ) : (
            receivedFootprints.map((item) => (
              <FootprintRow key={item.id} item={item} />
            ))
          )}
        </TabsContent>

        <TabsContent value="sent">
          {sentFootprints.length === 0 ? (
            <EmptyState type="sent" />
          ) : (
            sentFootprints.map((item) => (
              <FootprintRow key={item.id} item={item} />
            ))
          )}
        </TabsContent>
      </Tabs>

      <BottomNav />
    </div>
  );
}