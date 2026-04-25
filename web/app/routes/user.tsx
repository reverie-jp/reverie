import { useState, useEffect, useCallback } from "react";
import { Link, useNavigate } from "react-router";
import {
  getUser,
  getMe,
  listUserPosts,
  createPost,
  likePost,
  unlikePost,
  deletePost,
  isLoggedIn,
  createDirectRoom,
  type User as ApiUser,
} from "~/lib/api";
import { useCurrentUser } from "~/lib/use-current-user";
import { apiPostToUiPost } from "~/lib/utils";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
  AvatarGroup,
  AvatarGroupCount,
} from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { FollowButton } from "~/components/follow-button";
import { BottomNav } from "~/components/bottom-nav";
import { ComposeFab } from "~/components/compose-fab";
import { PostCard, type Post } from "~/components/post-card";
import {
  ComposePostDialog,
  type ComposeMode,
  type PostOptions,
} from "~/components/compose-post-dialog";
import { GroupAvatar, type Call } from "~/components/call-list";
import { JoinCallDialog } from "~/components/join-call-dialog";
import {
  ArrowLeft,
  CalendarDays,
  Link2,
  MapPin,
  MessageSquare,
  Phone,
  Settings,
  ShieldBan,
  Video,
  Crown,
  Image as ImageIcon,
} from "lucide-react";
import type { Route } from "./+types/user";

interface UserProfile {
  name: string;
  customId: string;
  avatarUrl?: string;
  bannerUrl?: string;
  bio: string;
  location?: string;
  website?: string;
  joinedAt: Date;
  followingCount: number;
  followerCount: number;
  isFollowing: boolean;
  followsYou?: boolean;
  isMe: boolean;
  blockedByThem?: boolean;
  onlineStatus?: "online" | "idle" | "offline";
}

function apiUserToProfile(u: ApiUser): UserProfile {
  return {
    name: u.displayName ?? "",
    customId: u.customId,
    bio: u.biography ?? "",
    joinedAt: new Date(u.createTime),
    followingCount: u.followingCount ?? 0,
    followerCount: u.followerCount ?? 0,
    isFollowing: u.isFollowing ?? false,
    followsYou: u.isFollowedBy ?? false,
    isMe: u.isMe,
  };
}

interface MutualFollower {
  name: string;
  customId: string;
  avatarUrl?: string;
}

const users: Record<string, UserProfile> = {
  me: {
    name: "自分",
    customId: "me",
    bio: "ソフトウェアエンジニア。TypeScript / React が好きです。",
    location: "東京",
    website: "example.com",
    joinedAt: new Date(2024, 3, 1),
    followingCount: 128,
    followerCount: 256,
    isFollowing: false,
    isMe: true,
    onlineStatus: "online",
  },
  tanaka: {
    name: "田中太郎",
    customId: "tanaka",
    bio: "エンジニア兼ゲーマー。TypeScriptの型パズルが趣味。",
    location: "大阪",
    joinedAt: new Date(2024, 0, 15),
    followingCount: 85,
    followerCount: 342,
    isFollowing: true,
    followsYou: true,
    isMe: false,
    onlineStatus: "online",
  },
  hanako_s: {
    name: "佐藤花子",
    customId: "hanako_s",
    bio: "カフェ巡りと読書が好き。デザイナーやってます。",
    location: "京都",
    website: "hanako.design",
    joinedAt: new Date(2024, 1, 20),
    followingCount: 200,
    followerCount: 510,
    isFollowing: true,
    followsYou: true,
    isMe: false,
    onlineStatus: "idle",
  },
  ichiro_dev: {
    name: "鈴木一郎",
    customId: "ichiro_dev",
    bio: "",
    joinedAt: new Date(2024, 2, 10),
    followingCount: 0,
    followerCount: 0,
    isFollowing: false,
    isMe: false,
    blockedByThem: true,
    onlineStatus: "offline",
  },
};

const defaultUser: UserProfile = {
  name: "ユーザー",
  customId: "unknown",
  bio: "",
  joinedAt: new Date(2024, 6, 1),
  followingCount: 0,
  followerCount: 0,
  isFollowing: false,
  isMe: false,
};

const mutualFollowers: Record<string, MutualFollower[]> = {
  tanaka: [
    { name: "佐藤花子", customId: "hanako_s" },
    { name: "山田美咲", customId: "misaki_y" },
    { name: "渡辺大輔", customId: "daisuke_w" },
    { name: "小林あおい", customId: "aoi_kb" },
    { name: "木村拓也", customId: "takuya_k" },
  ],
  hanako_s: [
    { name: "田中太郎", customId: "tanaka" },
    { name: "山田美咲", customId: "misaki_y" },
  ],
};

function getUserPosts(customId: string): Post[] {
  return [
    {
      id: `${customId}-1`,
      author: {
        name: users[customId]?.name ?? "ユーザー",
        customId,
        avatarUrl: "",
      },
      content: "今日はとてもいい天気ですね。散歩に行ってきました！",
      createdAt: new Date(Date.now() - 3 * 60_000),
      replyCount: 2,
      repostCount: 1,
      likeCount: 5,
    },
    {
      id: `${customId}-2`,
      author: {
        name: users[customId]?.name ?? "ユーザー",
        customId,
        avatarUrl: "",
      },
      content:
        "新しいプロジェクトを始めました。React Router v7を使ってSNSアプリを作っています。",
      createdAt: new Date(Date.now() - 5 * 3_600_000),
      replyCount: 4,
      repostCount: 3,
      likeCount: 15,
    },
    {
      id: `${customId}-3`,
      author: {
        name: users[customId]?.name ?? "ユーザー",
        customId,
        avatarUrl: "",
      },
      content: "週末に読んだ本がとても面白かった。おすすめです！",
      createdAt: new Date(Date.now() - 2 * 86_400_000),
      replyCount: 1,
      repostCount: 0,
      likeCount: 8,
    },
  ];
}

interface ProfileCall extends Call {
  status: "live" | "ended";
  duration?: string;
  hostId: string;
}

function getUserCalls(customId: string): ProfileCall[] {
  const name = users[customId]?.name ?? "ユーザー";
  return [
    {
      id: `${customId}-call-1`,
      name: "雑談部屋",
      type: "audio",
      host: name,
      hostId: customId,
      participants: [
        { name, customId, avatarUrl: "" },
        { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
        { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
        { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
        { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
        { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
        { name: "高橋真一", customId: "shin_t", avatarUrl: "" },
      ],
      status: "live",
    },
    {
      id: `${customId}-call-2`,
      name: "デザインレビュー",
      type: "video",
      host: name,
      hostId: customId,
      participants: [
        { name, customId, avatarUrl: "" },
        { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
        { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
      ],
      status: "ended",
      duration: "45:32",
    },
    {
      id: `${customId}-call-3`,
      name: "チーム定例",
      type: "audio",
      host: "小林あおい",
      hostId: "aoi_kb",
      participants: [
        { name, customId, avatarUrl: "" },
        { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
        { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
        { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
      ],
      status: "ended",
      duration: "1:12:05",
    },
  ];
}

function getUserLikedPosts(): Post[] {
  return [
    {
      id: "liked-1",
      author: { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
      content:
        "React Routerの新しいバージョンを試してみたけど、かなり使いやすくなってる。",
      createdAt: new Date(Date.now() - 1 * 86_400_000),
      replyCount: 8,
      repostCount: 15,
      likeCount: 42,
    },
    {
      id: "liked-2",
      author: { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
      content:
        "Rustでウェブサーバーを書いてみた。所有権の概念、最初は戸惑ったけどコンパイラに怒られながら学ぶのが逆に楽しい。",
      createdAt: new Date(Date.now() - 3 * 86_400_000),
      replyCount: 15,
      repostCount: 20,
      likeCount: 78,
    },
    {
      id: "liked-3",
      author: { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
      content: "今日の夕焼けが本当にきれいだった。写真では伝わらないくらい。",
      createdAt: new Date(Date.now() - 45 * 60_000),
      replyCount: 2,
      repostCount: 5,
      likeCount: 28,
    },
  ];
}

interface MediaItem {
  id: string;
  url: string;
  postId: string;
}

function getUserMedia(customId: string): MediaItem[] {
  return [
    { id: `${customId}-m1`, url: "", postId: `${customId}-1` },
    { id: `${customId}-m2`, url: "", postId: `${customId}-2` },
    { id: `${customId}-m3`, url: "", postId: `${customId}-1` },
    { id: `${customId}-m4`, url: "", postId: `${customId}-3` },
    { id: `${customId}-m5`, url: "", postId: `${customId}-2` },
    { id: `${customId}-m6`, url: "", postId: `${customId}-1` },
  ];
}

function ProfileCallItem({
  call,
  onTap,
}: {
  call: ProfileCall;
  onTap: (call: ProfileCall) => void;
}) {
  const TypeIcon = call.type === "video" ? Video : Phone;
  const participantNames = call.participants.map((p) => p.name).join("、");
  const isLive = call.status === "live";

  return (
    <button
      onClick={() => onTap(call)}
      className="flex items-center gap-3 w-full px-4 py-3 hover:bg-muted/50 transition-colors"
    >
      <GroupAvatar participants={call.participants} className="size-12" />
      <div className="flex-1 text-left min-w-0">
        <div className="flex items-center gap-2">
          <span className="font-medium text-sm">{call.name}</span>
          {isLive && (
            <div className="flex gap-1 items-center">
              <span className="size-1.5 rounded-full bg-green-500 animate-pulse" />
              <span className="text-xs text-green-500 font-medium">通話中</span>
            </div>
          )}
        </div>
        <div className="text-xs text-muted-foreground truncate">
          {isLive ? (
            <>{call.participants.length}人参加中</>
          ) : (
            <>
              {call.duration} · {participantNames}
            </>
          )}
        </div>
      </div>
      <div
        className={`size-8 rounded-full flex items-center justify-center shrink-0 ${
          isLive
            ? call.type === "video"
              ? "bg-blue-500/10 text-blue-500"
              : "bg-green-500/10 text-green-500"
            : call.type === "video"
              ? "bg-blue-500/5 text-blue-500/40"
              : "bg-green-500/5 text-green-500/40"
        }`}
      >
        <TypeIcon className="size-4" />
      </div>
    </button>
  );
}

function FloatingCallBadge({
  call,
  onClick,
}: {
  call: ProfileCall;
  onClick: () => void;
}) {
  const TypeIcon = call.type === "video" ? Video : Phone;
  const maxAvatars = 5;
  const shown = call.participants.slice(0, maxAvatars);
  const overflow = call.participants.length - maxAvatars;
  const hostIndex = call.participants.findIndex(
    (p) => p.customId === call.hostId,
  );

  return (
    <button
      onClick={onClick}
      className="fixed bottom-20 right-4 z-40 flex items-center gap-2.5 rounded-full border bg-background pl-3 pr-4 py-2 shadow-lg hover:bg-muted/50 transition-colors max-w-50"
    >
      <div
        className={`size-8 rounded-full flex items-center justify-center shrink-0 ${call.type === "video" ? "bg-blue-500/10" : "bg-green-500/10"}`}
      >
        <TypeIcon
          className={`size-4 ${call.type === "video" ? "text-blue-500" : "text-green-500"}`}
        />
      </div>
      <div className="flex flex-col gap-1 min-w-0">
        <div className="flex items-center gap-1.5">
          <span className="text-sm font-medium truncate">通話に参加中</span>
          <span className="size-1.5 rounded-full bg-green-500 animate-pulse shrink-0" />
        </div>
        <AvatarGroup className="-space-x-1.5 *:data-[slot=avatar]:ring-background *:data-[slot=avatar]:ring-[1.5px]">
          {shown.map((p, i) => (
            <div key={i} className="relative">
              <Avatar className="size-5!">
                <AvatarImage src={p.avatarUrl} alt={p.name} />
                <AvatarFallback className="text-[8px]">
                  {p.name.slice(0, 1)}
                </AvatarFallback>
              </Avatar>
              {i === hostIndex && (
                <div className="absolute -top-px -left-0.5 size-2 rounded-full bg-yellow-400 flex items-center justify-center z-10">
                  <Crown className="size-1 text-yellow-800" />
                </div>
              )}
            </div>
          ))}
          {overflow > 0 && (
            <AvatarGroupCount className="size-5! text-[8px] ring-[1.5px]! ring-background!">
              +{overflow}
            </AvatarGroupCount>
          )}
        </AvatarGroup>
      </div>
    </button>
  );
}

export default function User({ params }: Route.ComponentProps) {
  const navigate = useNavigate();
  const currentUser = useCurrentUser();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [apiUserId, setApiUserId] = useState<string | null>(null);
  const [loadingProfile, setLoadingProfile] = useState(true);
  const [posts, setPosts] = useState<Post[]>([]);
  const [isBlocked, setIsBlocked] = useState(false);
  const [showBlockedProfile, setShowBlockedProfile] = useState(false);
  const [composeMode, setComposeMode] = useState<ComposeMode | null>(null);
  const [selectedCall, setSelectedCall] = useState<ProfileCall | null>(null);
  const [creatingRoom, setCreatingRoom] = useState(false);

  const calls = getUserCalls(params.id);
  const likedPosts = getUserLikedPosts();
  const media = getUserMedia(params.id);
  const mutuals: MutualFollower[] = [];

  useEffect(() => {
    if (!isLoggedIn()) return;
    setLoadingProfile(true);
    const fetcher = params.id === "me" ? getMe() : getUser(params.id);
    fetcher
      .then(({ user }) => {
        setProfile(apiUserToProfile(user));
        setApiUserId(user.id);
        setLoadingProfile(false);
      })
      .catch(() => {
        setProfile(null);
        setLoadingProfile(false);
      });
  }, [params.id]);

  useEffect(() => {
    if (!isLoggedIn() || !apiUserId) return;
    listUserPosts(apiUserId, { pageSize: 20 })
      .then(({ posts: p }) => setPosts((p ?? []).map(apiPostToUiPost)))
      .catch(() => setPosts([]));
  }, [apiUserId]);

  const isMe = params.id === "me" || (profile?.isMe ?? false);

  // For other users, show their live call as a floating badge
  const otherUserLiveCall = !isMe
    ? calls.find((c) => c.status === "live")
    : null;

  if (loadingProfile) {
    return (
      <div className="w-full min-h-full flex items-center justify-center">
        <p className="text-sm text-muted-foreground">読み込み中...</p>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="w-full min-h-full flex items-center justify-center">
        <p className="text-sm text-muted-foreground">ユーザーが見つかりません</p>
      </div>
    );
  }

  const handlePost = async (content: string, options?: PostOptions) => {
    try {
      await createPost({ text: content, replyToId: options?.replyToId, repostId: options?.repostId });
      const res = await listUserPosts(apiUserId ?? params.id, { pageSize: 20 });
      setPosts((res.posts ?? []).map(apiPostToUiPost));
    } catch (e: any) {
      alert(`投稿に失敗しました: ${e?.message ?? "不明なエラー"}`);
    }
  };

  const handleLike = async (postId: string) => { try { await likePost(postId); } catch {} };
  const handleUnlike = async (postId: string) => { try { await unlikePost(postId); } catch {} };
  const handleDelete = async (postId: string) => {
    try {
      await deletePost(postId);
      setPosts((prev) => prev.filter((p) => p.id !== postId));
    } catch {}
  };

  const handleReply = (post: Post) => {
    setComposeMode({ type: "reply", post });
  };

  const handleRepost = (post: Post) => {
    setComposeMode({ type: "repost", post });
  };

  const handleOpenChat = async () => {
    if (!apiUserId || creatingRoom) return;
    setCreatingRoom(true);
    try {
      const res = await createDirectRoom(apiUserId);
      navigate(`/chat/${res.room.id}`);
    } catch {
      setCreatingRoom(false);
    }
  };

  const formatJoinDate = (date: Date) => {
    return `${date.getFullYear()}年${date.getMonth() + 1}月`;
  };

  return (
    <div className="w-full min-h-full flex flex-col">
      {/* Header */}
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-14">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <div className="min-w-0 flex-1">
            <h1 className="text-base font-bold truncate">{profile.name}</h1>
          </div>
          {isMe && (
            <Link
              to="/settings"
              className="inline-flex items-center justify-center size-9 rounded-md text-muted-foreground hover:bg-accent hover:text-foreground transition-colors"
            >
              <Settings className="size-5" />
            </Link>
          )}
        </div>
      </div>

      {/* Banner */}
      <div className="w-full h-32 bg-muted" />

      {/* Profile info */}
      <div className="px-4 relative">
        {/* Action button (top-right of profile section) */}
        <div className="absolute top-3 right-4 flex items-center gap-2">
          {isMe ? (
            <Button
              variant="outline"
              className="rounded-full h-9 px-4"
              onClick={() => navigate("/settings/profile")}
            >
              プロフィールを編集
            </Button>
          ) : profile.blockedByThem ? null : (
            <>
              <Button
                variant="outline"
                size="icon"
                className="rounded-full size-9"
                onClick={handleOpenChat}
                disabled={creatingRoom}
              >
                <MessageSquare className="size-4" />
              </Button>
              <FollowButton
                userId={apiUserId ?? ""}
                customId={profile.customId}
                initialFollowing={profile.isFollowing}
                followsYou={profile.followsYou}
                size="md"
                onBlockChange={setIsBlocked}
                onFollowChange={(following) =>
                  setPosts((prev) =>
                    prev.map((p) => ({
                      ...p,
                      author: { ...p.author, isFollowing: following },
                    }))
                  )
                }
              />
            </>
          )}
        </div>

        {/* Avatar */}
        <div className="-mt-12 relative w-fit">
          <Avatar className="size-20 ring-4 ring-background">
            {!profile.blockedByThem && (
              <AvatarImage src={profile.avatarUrl} alt={profile.name} />
            )}
            <AvatarFallback className="text-2xl">
              {profile.name.slice(0, 2)}
            </AvatarFallback>
          </Avatar>
          {profile.onlineStatus && profile.onlineStatus !== "offline" && (
            <span
              className={`absolute bottom-0.5 right-0.5 size-4 rounded-full ring-[3px] ring-background ${
                profile.onlineStatus === "online"
                  ? "bg-green-500"
                  : "bg-yellow-500"
              }`}
            />
          )}
        </div>

        {/* Name & ID */}
        <div className="mt-3">
          <div className="flex items-center gap-2">
            <p className="text-lg font-bold">{profile.name}</p>
            {profile.followsYou && (
              <span className="text-[10px] bg-muted text-muted-foreground rounded px-1.5 py-0.5 leading-none">
                フォローされています
              </span>
            )}
          </div>
          <p className="text-sm text-muted-foreground">@{profile.customId}</p>
        </div>

        {!profile.blockedByThem && (
          <>
            {/* Bio */}
            {profile.bio && (
              <p className="mt-2 text-sm whitespace-pre-wrap">{profile.bio}</p>
            )}

            {/* Meta info */}
            <div className="flex flex-wrap gap-x-4 gap-y-1 mt-2 text-sm text-muted-foreground">
              {profile.location && (
                <span className="flex items-center gap-1">
                  <MapPin className="size-3.5" />
                  {profile.location}
                </span>
              )}
              {profile.website && (
                <a
                  href={
                    profile.website.startsWith("http")
                      ? profile.website
                      : `https://${profile.website}`
                  }
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-1 hover:underline"
                  onClick={(e) => e.stopPropagation()}
                >
                  <Link2 className="size-3.5" />
                  <span className="text-primary">{profile.website}</span>
                </a>
              )}
              <span className="flex items-center gap-1">
                <CalendarDays className="size-3.5" />
                {formatJoinDate(profile.joinedAt)}に登録
              </span>
            </div>

            {/* Following / Followers */}
            <div className="flex gap-4 mt-3 text-sm">
              <Link
                to={`/users/${profile.customId}/connections?tab=following`}
                className="hover:underline"
              >
                <span className="font-bold">{profile.followingCount}</span>{" "}
                <span className="text-muted-foreground">フォロー中</span>
              </Link>
              <Link
                to={`/users/${profile.customId}/connections?tab=followers`}
                className="hover:underline"
              >
                <span className="font-bold">{profile.followerCount}</span>{" "}
                <span className="text-muted-foreground">フォロワー</span>
              </Link>
            </div>

            {/* Mutual followers */}
            {!isMe && mutuals.length > 0 && (
              <Link
                to={`/users/${profile.customId}/connections?tab=mutual`}
                className="flex items-center gap-2 mt-3 group"
              >
                <AvatarGroup className="-space-x-2 *:data-[slot=avatar]:ring-background *:data-[slot=avatar]:ring-[1.5px]">
                  {mutuals.slice(0, 3).map((u) => (
                    <Avatar key={u.customId} className="size-5!">
                      <AvatarImage src={u.avatarUrl} alt={u.name} />
                      <AvatarFallback className="text-[8px]">
                        {u.name.slice(0, 1)}
                      </AvatarFallback>
                    </Avatar>
                  ))}
                </AvatarGroup>
                <p className="text-xs text-muted-foreground group-hover:underline">
                  {mutuals.length <= 2
                    ? `${mutuals.map((u) => u.name).join("、")}がフォロー中`
                    : `${mutuals
                        .slice(0, 2)
                        .map((u) => u.name)
                        .join("、")}、他${mutuals.length - 2}人がフォロー中`}
                </p>
              </Link>
            )}
          </>
        )}
      </div>

      {profile.blockedByThem ? (
        <div className="flex-1 flex flex-col items-center justify-center px-6 py-16 text-center">
          <ShieldBan className="size-12 text-muted-foreground/30 mb-4" />
          <p className="text-lg font-bold">あなたをブロックしました</p>
          <p className="text-sm text-muted-foreground mt-2">
            @{profile.customId}{" "}
            さんにブロックされているため、投稿や通話などを閲覧できません。
          </p>
        </div>
      ) : isBlocked && !showBlockedProfile ? (
        <div className="flex-1 flex flex-col items-center justify-center px-6 py-16 text-center">
          <ShieldBan className="size-12 text-muted-foreground/30 mb-4" />
          <p className="text-lg font-bold">
            @{profile.customId} さんをブロックしています
          </p>
          <p className="text-sm text-muted-foreground mt-2">
            投稿を見るだけなら @{profile.customId}{" "}
            さんのブロックを解除しなくても確認することができます。
          </p>
          <Button
            variant="outline"
            size="sm"
            className="mt-4 rounded-full"
            onClick={() => setShowBlockedProfile(true)}
          >
            投稿を表示
          </Button>
        </div>
      ) : (
        <>
          {/* Tabs */}
          <Tabs defaultValue="posts" className="gap-0 flex-1 mt-4">
            <div className="sticky top-14 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
              <TabsList variant="line" className="w-full h-12">
                <TabsTrigger value="posts">投稿</TabsTrigger>
                <TabsTrigger value="calls">通話</TabsTrigger>
                <TabsTrigger value="media">メディア</TabsTrigger>
                <TabsTrigger value="likes">いいね</TabsTrigger>
              </TabsList>
            </div>

            <TabsContent value="posts">
              {posts.length === 0 ? (
                <div className="py-12 text-center text-muted-foreground text-sm">
                  まだ投稿がありません
                </div>
              ) : (
                posts.map((post) => (
                  <PostCard
                    key={post.id}
                    post={post}
                    currentUserId={apiUserId ?? undefined}
                    onReply={handleReply}
                    onRepost={handleRepost}
                    onLike={handleLike}
                    onUnlike={handleUnlike}
                    onDelete={handleDelete}
                  />
                ))
              )}
            </TabsContent>

            <TabsContent value="calls">
              {calls.length === 0 ? (
                <div className="py-12 text-center text-muted-foreground text-sm">
                  通話履歴はありません
                </div>
              ) : (
                <div className="divide-y">
                  {calls.map((call) => (
                    <ProfileCallItem
                      key={call.id}
                      call={call}
                      onTap={setSelectedCall}
                    />
                  ))}
                </div>
              )}
            </TabsContent>

            <TabsContent value="media">
              {media.length === 0 ? (
                <div className="py-12 text-center text-muted-foreground text-sm">
                  メディアはありません
                </div>
              ) : (
                <div className="grid grid-cols-3 gap-0.5">
                  {media.map((item) => (
                    <Link
                      key={item.id}
                      to={`/posts/${item.postId}`}
                      className="aspect-square bg-muted flex items-center justify-center hover:opacity-80 transition-opacity"
                    >
                      {item.url ? (
                        <img
                          src={item.url}
                          alt=""
                          className="w-full h-full object-cover"
                        />
                      ) : (
                        <ImageIcon className="size-8 text-muted-foreground/40" />
                      )}
                    </Link>
                  ))}
                </div>
              )}
            </TabsContent>

            <TabsContent value="likes">
              {likedPosts.length === 0 ? (
                <div className="py-12 text-center text-muted-foreground text-sm">
                  いいねした投稿はありません
                </div>
              ) : (
                likedPosts.map((post) => (
                  <PostCard
                    key={post.id}
                    post={post}
                    onReply={handleReply}
                    onRepost={handleRepost}
                  />
                ))
              )}
            </TabsContent>
          </Tabs>
        </>
      )}

      <ComposePostDialog
        open={composeMode !== null}
        onClose={() => setComposeMode(null)}
        onPost={async (content, options) => { await handlePost(content, options); setComposeMode(null); }}
        mode={composeMode ?? undefined}
      />
      <JoinCallDialog
        call={selectedCall}
        onClose={() => setSelectedCall(null)}
        ended={selectedCall?.status === "ended"}
      />
      {/* Floating call badge for other users */}
      {otherUserLiveCall && (
        <FloatingCallBadge
          call={otherUserLiveCall}
          onClick={() => setSelectedCall(otherUserLiveCall)}
        />
      )}
      {isMe && (
        <ComposeFab
          onPost={handlePost}
          currentUser={currentUser ? { name: currentUser.displayName, avatarUrl: undefined } : undefined}
        />
      )}
      <BottomNav />
    </div>
  );
}
