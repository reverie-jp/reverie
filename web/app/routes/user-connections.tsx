import { useState } from "react";
import { Link, useSearchParams } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { ConfirmActionDialog } from "~/components/confirm-action-dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { BottomNav } from "~/components/bottom-nav";
import { ArrowLeft } from "lucide-react";
import type { Route } from "./+types/user-connections";

interface UserItem {
  name: string;
  customId: string;
  avatarUrl?: string;
  bio?: string;
  isFollowing: boolean;
  followsYou?: boolean;
  isMe: boolean;
}

const followingUsers: UserItem[] = [
  {
    name: "田中太郎",
    customId: "tanaka",
    bio: "エンジニア兼ゲーマー。TypeScriptの型パズルが趣味。",
    isFollowing: true,
    followsYou: true,
    isMe: false,
  },
  {
    name: "佐藤花子",
    customId: "hanako_s",
    bio: "カフェ巡りと読書が好き。デザイナーやってます。",
    isFollowing: true,
    followsYou: true,
    isMe: false,
  },
  {
    name: "山田美咲",
    customId: "misaki_y",
    bio: "フロントエンドエンジニア。猫が好き。",
    isFollowing: true,
    followsYou: false,
    isMe: false,
  },
  {
    name: "渡辺大輔",
    customId: "daisuke_w",
    bio: "バックエンドエンジニア。Go / Rust。",
    isFollowing: true,
    followsYou: true,
    isMe: false,
  },
  {
    name: "小林あおい",
    customId: "aoi_kb",
    bio: "PM やってます。スクラム好き。",
    isFollowing: true,
    followsYou: false,
    isMe: false,
  },
];

const followerUsers: UserItem[] = [
  {
    name: "田中太郎",
    customId: "tanaka",
    bio: "エンジニア兼ゲーマー。TypeScriptの型パズルが趣味。",
    isFollowing: true,
    followsYou: true,
    isMe: false,
  },
  {
    name: "佐藤花子",
    customId: "hanako_s",
    bio: "カフェ巡りと読書が好き。デザイナーやってます。",
    isFollowing: true,
    followsYou: true,
    isMe: false,
  },
  {
    name: "鈴木一郎",
    customId: "ichiro_dev",
    bio: "",
    isFollowing: false,
    followsYou: true,
    isMe: false,
  },
  {
    name: "高橋健太",
    customId: "kenta_t",
    bio: "データサイエンティスト。Python / R。",
    isFollowing: false,
    followsYou: true,
    isMe: false,
  },
  {
    name: "木村拓也",
    customId: "takuya_k",
    bio: "インフラエンジニア。Kubernetes。",
    isFollowing: true,
    followsYou: true,
    isMe: false,
  },
  {
    name: "中村悠",
    customId: "yu_nkmr",
    bio: "モバイルエンジニア。Swift / Kotlin。",
    isFollowing: false,
    followsYou: true,
    isMe: false,
  },
];

const userNames: Record<string, string> = {
  me: "自分",
  tanaka: "田中太郎",
  hanako_s: "佐藤花子",
  ichiro_dev: "鈴木一郎",
};

function UserListItem({ user }: { user: UserItem }) {
  const [following, setFollowing] = useState(user.isFollowing);
  const [showUnfollowConfirm, setShowUnfollowConfirm] = useState(false);

  return (
    <>
      <div className="flex items-start gap-3 px-4 py-3 border-b">
        <Link to={`/users/${user.customId}`} className="shrink-0">
          <Avatar className="size-10">
            <AvatarImage src={user.avatarUrl} alt={user.name} />
            <AvatarFallback>{user.name.slice(0, 2)}</AvatarFallback>
          </Avatar>
        </Link>
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between gap-2">
            <Link to={`/users/${user.customId}`} className="min-w-0">
              <div className="flex items-center gap-1.5">
                <p className="text-sm font-bold truncate">{user.name}</p>
                {user.followsYou && (
                  <span className="shrink-0 text-[10px] bg-muted text-muted-foreground rounded px-1 py-0.5 leading-none">
                    フォローされています
                  </span>
                )}
              </div>
              <p className="text-xs text-muted-foreground">@{user.customId}</p>
            </Link>
            {!user.isMe && (
              <Button
                variant={following ? "outline" : "default"}
                className="rounded-full h-8 px-3 shrink-0 text-xs"
                onClick={() => {
                  if (following) {
                    setShowUnfollowConfirm(true);
                  } else {
                    setFollowing(true);
                  }
                }}
              >
                {following ? "フォロー中" : "フォローする"}
              </Button>
            )}
          </div>
          {user.bio && (
            <p className="mt-1 text-sm text-muted-foreground line-clamp-2">
              {user.bio}
            </p>
          )}
        </div>
      </div>

      <ConfirmActionDialog
        action={showUnfollowConfirm ? "unfollow" : null}
        customId={user.customId}
        onConfirm={() => {
          setFollowing(false);
          setShowUnfollowConfirm(false);
        }}
        onCancel={() => setShowUnfollowConfirm(false)}
      />
    </>
  );
}

export default function UserConnections({ params }: Route.ComponentProps) {
  const [searchParams] = useSearchParams();
  const defaultTab = searchParams.get("tab") === "followers" ? "followers" : "following";
  const displayName = userNames[params.id] ?? params.id;

  return (
    <div className="w-full min-h-full flex flex-col">
      {/* Header */}
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-14">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <div className="min-w-0">
            <h1 className="text-base font-bold truncate">{displayName}</h1>
            <p className="text-xs text-muted-foreground">@{params.id}</p>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <Tabs defaultValue={defaultTab} className="gap-0 flex-1">
        <div className="sticky top-14 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
          <TabsList variant="line" className="w-full h-12">
            <TabsTrigger value="following">フォロー中</TabsTrigger>
            <TabsTrigger value="followers">フォロワー</TabsTrigger>
          </TabsList>
        </div>

        <TabsContent value="following">
          {followingUsers.length === 0 ? (
            <div className="py-12 text-center text-muted-foreground text-sm">
              フォロー中のユーザーはいません
            </div>
          ) : (
            followingUsers.map((user) => (
              <UserListItem key={user.customId} user={user} />
            ))
          )}
        </TabsContent>

        <TabsContent value="followers">
          {followerUsers.length === 0 ? (
            <div className="py-12 text-center text-muted-foreground text-sm">
              フォロワーはいません
            </div>
          ) : (
            followerUsers.map((user) => (
              <UserListItem key={user.customId} user={user} />
            ))
          )}
        </TabsContent>
      </Tabs>

      <BottomNav />
    </div>
  );
}
