import { useRef, useState, useEffect } from "react";
import { Link, useSearchParams } from "react-router";
import {
  Search,
  X,
  TrendingUp,
  UserPlus,
  Clock,
  FileText,
  Users,
  Image,
  Phone,
} from "lucide-react";
import { Input } from "~/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { BottomNav } from "~/components/bottom-nav";
import { PostCard, type Post } from "~/components/post-card";
import { type Call } from "~/components/call-list";
import { CallListItem } from "~/routes/calls";
import { JoinCallDialog } from "~/components/join-call-dialog";
import { FollowButton } from "~/components/follow-button";
import {
  ComposePostDialog,
  type ComposeMode,
} from "~/components/compose-post-dialog";
import { searchUsers, type User as ApiUser } from "~/lib/api";

// --- Types ---

type SearchType = "posts" | "users" | "media" | "calls";

interface TrendItem {
  id: string;
  keyword: string;
  postCount: number;
}

interface SuggestedUser {
  id: string;
  name: string;
  customId: string;
  avatarUrl?: string;
  bio: string;
  isFollowing: boolean;
  followsYou: boolean;
}

interface HistoryEntry {
  id: string;
  keyword: string;
  type: SearchType;
}

function apiUserToSuggested(u: ApiUser): SuggestedUser {
  return {
    id: u.id,
    name: u.displayName,
    customId: u.customId,
    bio: u.biography ?? "",
    isFollowing: u.isFollowing ?? false,
    followsYou: u.isFollowedBy ?? false,
  };
}

// --- Sample data ---

const sampleTrends: TrendItem[] = [
  { id: "t1", keyword: "React Router v7", postCount: 1842 },
  { id: "t2", keyword: "桜の開花", postCount: 5621 },
  { id: "t3", keyword: "新生活", postCount: 3204 },
  { id: "t4", keyword: "TypeScript 6.0", postCount: 982 },
  { id: "t5", keyword: "花見スポット", postCount: 2847 },
];

const sampleSuggestedUsers: SuggestedUser[] = [
  {
    id: "user-tanaka",
    name: "田中太郎",
    customId: "tanaka",
    avatarUrl: "",
    bio: "フロントエンドエンジニア。React / TypeScript が好きです。",
    isFollowing: false,
    followsYou: false,
  },
  {
    id: "user-hanako_s",
    name: "佐藤花子",
    customId: "hanako_s",
    avatarUrl: "",
    bio: "デザイナー兼コーヒー愛好家。UI/UXについて発信中",
    isFollowing: false,
    followsYou: true,
  },
  {
    id: "user-ichiro_dev",
    name: "鈴木一郎",
    customId: "ichiro_dev",
    avatarUrl: "",
    bio: "バックエンドエンジニア。Go / Rust で開発しています。",
    isFollowing: false,
    followsYou: false,
  },
  {
    id: "user-misaki_y",
    name: "山田美咲",
    customId: "misaki_y",
    avatarUrl: "",
    bio: "映画とアニメが大好き。週末はいつも映画館にいます。",
    isFollowing: false,
    followsYou: true,
  },
  {
    id: "user-kenta_t",
    name: "高橋健太",
    customId: "kenta_t",
    avatarUrl: "",
    bio: "プログラミング学習中の大学生。毎日コツコツ頑張ってます。",
    isFollowing: false,
    followsYou: false,
  },
  {
    id: "user-yu_nkmr",
    name: "中村悠",
    customId: "yu_nkmr",
    avatarUrl: "",
    bio: "写真が趣味。風景と街スナップを撮っています。",
    isFollowing: false,
    followsYou: false,
  },
  {
    id: "user-aoi_kb",
    name: "小林あおい",
    customId: "aoi_kb",
    avatarUrl: "",
    bio: "イラストレーター。お仕事募集中です。",
    isFollowing: false,
    followsYou: true,
  },
  {
    id: "user-daisuke_w",
    name: "渡辺大輔",
    customId: "daisuke_w",
    avatarUrl: "",
    bio: "ガジェットレビューと自作キーボードにハマってます。",
    isFollowing: false,
    followsYou: false,
  },
  {
    id: "user-sakura_ito",
    name: "伊藤さくら",
    customId: "sakura_ito",
    avatarUrl: "",
    bio: "読書と散歩が日課。月10冊ペースで読んでます。",
    isFollowing: false,
    followsYou: false,
  },
  {
    id: "user-rina_m",
    name: "松本りな",
    customId: "rina_m",
    avatarUrl: "",
    bio: "料理研究家見習い。毎日のレシピを投稿中。",
    isFollowing: false,
    followsYou: true,
  },
];

const sampleHistory: HistoryEntry[] = [
  { id: "h1", keyword: "React hooks", type: "posts" },
  { id: "h2", keyword: "tanaka", type: "users" },
  { id: "h3", keyword: "夕焼け", type: "media" },
  { id: "h4", keyword: "雑談", type: "calls" },
  { id: "h5", keyword: "TypeScript", type: "posts" },
  { id: "h6", keyword: "カフェ", type: "posts" },
];

const samplePostResults: Post[] = [
  {
    id: "s1",
    author: { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
    content:
      "React Routerの新しいバージョンを試してみたけど、かなり使いやすくなってる。特にローダーの仕組みが良い。",
    createdAt: new Date(Date.now() - 1 * 86_400_000),
    replyCount: 8,
    repostCount: 15,
    likeCount: 42,
  },
  {
    id: "s2",
    author: { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
    content:
      "Rustでウェブサーバーを書いてみた。所有権の概念、最初は戸惑ったけどコンパイラに怒られながら学ぶのが逆に楽しい。",
    createdAt: new Date(Date.now() - 3 * 86_400_000),
    replyCount: 15,
    repostCount: 20,
    likeCount: 78,
  },
  {
    id: "s3",
    author: { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
    content:
      "TypeScriptの型パズル、難しいけど楽しい。最近はConditional Typesにハマってます。",
    createdAt: new Date(Date.now() - 2 * 86_400_000),
    replyCount: 6,
    repostCount: 10,
    likeCount: 35,
  },
  {
    id: "s4",
    author: { name: "井上翔", customId: "sho_inoue", avatarUrl: "" },
    content:
      "デザインシステムを一から構築中。コンポーネントの粒度をどこまで細かくするか、チームで議論が白熱してる。",
    createdAt: new Date(Date.now() - 4 * 86_400_000),
    replyCount: 7,
    repostCount: 12,
    likeCount: 44,
  },
];

const sampleUserResults: SuggestedUser[] = [
  {
    id: "user-tanaka",
    name: "田中太郎",
    customId: "tanaka",
    avatarUrl: "",
    bio: "フロントエンドエンジニア。React / TypeScript が好きです。",
    isFollowing: true,
    followsYou: false,
  },
  {
    id: "user-ichiro_dev",
    name: "鈴木一郎",
    customId: "ichiro_dev",
    avatarUrl: "",
    bio: "バックエンドエンジニア。Go / Rust で開発しています。",
    isFollowing: false,
    followsYou: true,
  },
  {
    id: "user-kenta_t",
    name: "高橋健太",
    customId: "kenta_t",
    avatarUrl: "",
    bio: "プログラミング学習中の大学生。毎日コツコツ頑張ってます。",
    isFollowing: false,
    followsYou: false,
  },
];

const sampleMediaPosts: Post[] = [
  {
    id: "m1",
    author: { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
    content:
      "今日の夕焼けが本当にきれいだった。写真では伝わらないくらい。自然の美しさに感動する日々。",
    createdAt: new Date(Date.now() - 45 * 60_000),
    replyCount: 2,
    repostCount: 5,
    likeCount: 28,
  },
  {
    id: "m2",
    author: { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
    content:
      "新しいキーボードを買いました。打鍵感が最高すぎて仕事が捗る。静電容量無接点方式、一度使うと戻れない。",
    createdAt: new Date(Date.now() - 5 * 3_600_000),
    replyCount: 12,
    repostCount: 8,
    likeCount: 56,
  },
];

const sampleCallResults: Call[] = [
  {
    id: "cr1",
    name: "雑談部屋",
    type: "audio",
    host: "田中太郎",
    participants: [
      { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
      { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
    ],
  },
  {
    id: "cr2",
    name: "デザインレビュー",
    type: "video",
    host: "佐藤花子",
    participants: [
      { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
      { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
      { name: "井上翔", customId: "sho_inoue", avatarUrl: "" },
    ],
  },
  {
    id: "cr3",
    name: "開発ミーティング",
    type: "video",
    host: "鈴木一郎",
    participants: [
      { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
      { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
    ],
  },
];

// --- Helpers ---

function formatCount(n: number): string {
  if (n >= 10_000) return `${(n / 10_000).toFixed(1)}万`;
  if (n >= 1_000) return `${(n / 1_000).toFixed(1)}千`;
  return `${n}`;
}

const searchTypeLabels: Record<
  SearchType,
  { label: string; icon: typeof FileText }
> = {
  posts: { label: "投稿", icon: FileText },
  users: { label: "ユーザー", icon: Users },
  media: { label: "メディア", icon: Image },
  calls: { label: "通話", icon: Phone },
};

// --- Components ---

function SearchTypeChips({
  selected,
  onChange,
}: {
  selected: SearchType;
  onChange: (type: SearchType) => void;
}) {
  return (
    <div className="flex gap-2 px-4 py-2.5">
      {(Object.keys(searchTypeLabels) as SearchType[]).map((type) => {
        const { label, icon: Icon } = searchTypeLabels[type];
        const isActive = selected === type;
        return (
          <button
            key={type}
            className={`flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-medium transition-colors ${
              isActive
                ? "bg-foreground text-background"
                : "bg-muted text-muted-foreground hover:bg-muted/80"
            }`}
            onClick={() => onChange(type)}
          >
            <Icon className="size-3.5" />
            {label}
          </button>
        );
      })}
    </div>
  );
}

function SearchHistory({
  history,
  onSelect,
  onRemove,
}: {
  history: HistoryEntry[];
  onSelect: (entry: HistoryEntry) => void;
  onRemove: (id: string) => void;
}) {
  if (history.length === 0) return null;

  return (
    <div>
      <div className="flex items-center justify-between px-4 py-2.5">
        <h3 className="text-sm font-bold">検索履歴</h3>
      </div>
      {history.map((entry) => {
        const { icon: Icon } = searchTypeLabels[entry.type];
        return (
          <button
            key={entry.id}
            className="w-full flex items-center gap-3 px-4 py-2.5 hover:bg-muted/30 transition-colors text-left"
            onClick={() => onSelect(entry)}
          >
            <Clock className="size-4 text-muted-foreground shrink-0" />
            <span className="flex-1 text-sm truncate">{entry.keyword}</span>
            <span className="flex items-center gap-1 text-xs text-muted-foreground">
              <Icon className="size-3" />
              {searchTypeLabels[entry.type].label}
            </span>
            <button
              className="text-muted-foreground hover:text-foreground transition-colors p-0.5"
              onClick={(e) => {
                e.stopPropagation();
                onRemove(entry.id);
              }}
            >
              <X className="size-3.5" />
            </button>
          </button>
        );
      })}
    </div>
  );
}

function TrendList({ trends }: { trends: TrendItem[] }) {
  return (
    <div>
      {trends.map((trend) => (
        <Link
          key={trend.id}
          to={`/trends/${encodeURIComponent(trend.keyword)}`}
          className="block px-4 py-3 border-b hover:bg-muted/30 transition-colors"
        >
          <p className="text-sm font-bold">{trend.keyword}</p>
          <p className="text-xs text-muted-foreground mt-0.5">
            {formatCount(trend.postCount)}件の投稿
          </p>
        </Link>
      ))}
    </div>
  );
}

function UserRow({ user }: { user: SuggestedUser }) {
  return (
    <div className="flex items-center gap-3 px-4 py-3 border-b hover:bg-muted/30 transition-colors">
      <Link to={`/users/${user.customId}`} className="shrink-0">
        <Avatar className="size-11">
          <AvatarImage src={user.avatarUrl} alt={user.name} />
          <AvatarFallback>{user.name.slice(0, 2)}</AvatarFallback>
        </Avatar>
      </Link>
      <div className="flex-1 min-w-0">
        <Link to={`/users/${user.customId}`}>
          <p className="text-sm font-bold truncate">{user.name}</p>
          <p className="text-xs text-muted-foreground truncate">
            @{user.customId}
          </p>
        </Link>
        <p className="text-xs text-muted-foreground mt-0.5 line-clamp-1">
          {user.bio}
        </p>
      </div>
      <FollowButton
        userId={user.id}
        customId={user.customId}
        initialFollowing={user.isFollowing}
        followsYou={user.followsYou}
      />
    </div>
  );
}

function SuggestedUserList({ users }: { users: SuggestedUser[] }) {
  return (
    <div>
      <div className="px-4 py-3 border-b">
        <h2 className="text-base font-bold flex items-center gap-2">
          <UserPlus className="size-4" />
          おすすめのユーザー
        </h2>
      </div>
      {users.map((user) => (
        <UserRow key={user.customId} user={user} />
      ))}
    </div>
  );
}

function EmptyResult({ query }: { query: string }) {
  return (
    <div className="flex flex-col items-center justify-center py-20 px-4 text-center">
      <Search className="size-10 text-muted-foreground/50 mb-3" />
      <p className="text-sm font-medium">「{query}」の検索結果はありません</p>
      <p className="text-xs text-muted-foreground mt-1">
        別のキーワードで検索してみてください
      </p>
    </div>
  );
}

// --- Main ---

export default function SearchPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const initialQuery = searchParams.get("q") ?? "";
  const initialType = (searchParams.get("type") as SearchType) || "posts";

  const [query, setQuery] = useState(initialQuery);
  const [searchType, setSearchType] = useState<SearchType>(initialType);
  const [isFocused, setIsFocused] = useState(false);
  const [committedQuery, setCommittedQuery] = useState(initialQuery);
  const [committedType, setCommittedType] = useState<SearchType>(initialType);
  const [history, setHistory] = useState(sampleHistory);
  const [userResults, setUserResults] = useState<SuggestedUser[]>([]);
  const [userResultsLoading, setUserResultsLoading] = useState(false);
  const [selectedCall, setSelectedCall] = useState<Call | null>(null);
  const [composeMode, setComposeMode] = useState<ComposeMode | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const overlayRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!committedQuery.trim()) {
      setUserResults([]);
      return;
    }
    setUserResultsLoading(true);
    searchUsers(committedQuery, { pageSize: 20 })
      .then((res) => setUserResults((res.users ?? []).map(apiUserToSuggested)))
      .catch(console.error)
      .finally(() => setUserResultsLoading(false));
  }, [committedQuery]);

  const hasResults = committedQuery.trim().length > 0;
  const showOverlay = isFocused && !hasResults;

  const commitSearch = (keyword?: string, type?: SearchType) => {
    const q = keyword ?? query;
    const t = type ?? searchType;
    if (q.trim().length === 0) return;

    setCommittedQuery(q);
    setCommittedType(t);
    setQuery(q);
    setSearchType(t);
    setIsFocused(false);
    inputRef.current?.blur();

    // Persist to URL so back navigation restores state
    setSearchParams({ q, type: t }, { replace: true });

    // Add to history (deduplicate)
    setHistory((prev) => {
      const filtered = prev.filter((h) => !(h.keyword === q && h.type === t));
      return [
        { id: `h-${Date.now()}`, keyword: q, type: t },
        ...filtered,
      ].slice(0, 10);
    });
  };

  const clearSearch = () => {
    setQuery("");
    setCommittedQuery("");
    setSearchParams({}, { replace: true });
    inputRef.current?.focus();
  };

  const handleReply = (post: Post) => {
    setComposeMode({ type: "reply", post });
  };

  const handleRepost = (post: Post) => {
    setComposeMode({ type: "repost", post });
  };

  const handleHistorySelect = (entry: HistoryEntry) => {
    setQuery(entry.keyword);
    setSearchType(entry.type);
    commitSearch(entry.keyword, entry.type);
  };

  const handleHistoryRemove = (id: string) => {
    setHistory((prev) => prev.filter((h) => h.id !== id));
  };



  return (
    <div className="w-full min-h-full flex flex-col">
      {/* Header */}
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="px-4 py-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
            <Input
              ref={inputRef}
              placeholder="キーワードで検索"
              value={query}
              onChange={(e) => {
                setQuery(e.target.value);
                if (committedQuery) {
                  setCommittedQuery("");
                }
              }}
              onFocus={() => setIsFocused(true)}
              onBlur={(e) => {
                // Keep overlay open if clicking inside the overlay area
                if (
                  overlayRef.current &&
                  e.relatedTarget instanceof Node &&
                  overlayRef.current.contains(e.relatedTarget)
                ) {
                  return;
                }
                setIsFocused(false);
              }}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  commitSearch();
                }
              }}
              className="pl-9 pr-9 h-10"
            />
            {(query.length > 0 || hasResults) && (
              <button
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                onClick={clearSearch}
              >
                <X className="size-4" />
              </button>
            )}
          </div>
        </div>

        {/* Search type chips when focused */}
        {showOverlay && (
          <div onMouseDown={(e) => e.preventDefault()}>
            <SearchTypeChips selected={searchType} onChange={setSearchType} />
          </div>
        )}

        {/* Result tabs */}
        {hasResults && (
          <Tabs
            value={committedType}
            onValueChange={(v) => {
              const t = v as SearchType;
              setCommittedType(t);
              setSearchType(t);
            }}
            className="gap-0"
          >
            <TabsList variant="line" className="w-full h-11">
              <TabsTrigger value="posts">投稿</TabsTrigger>
              <TabsTrigger value="users">ユーザー</TabsTrigger>
              <TabsTrigger value="media">メディア</TabsTrigger>
              <TabsTrigger value="calls">通話</TabsTrigger>
            </TabsList>
          </Tabs>
        )}
      </div>

      {/* Content */}
      <div className="flex-1">
        {showOverlay ? (
          /* Focused: search history */
          <div ref={overlayRef} onMouseDown={(e) => e.preventDefault()}>
            <SearchHistory
              history={history}
              onSelect={handleHistorySelect}
              onRemove={handleHistoryRemove}
            />
          </div>
        ) : hasResults ? (
          /* Search results */
          <Tabs value={committedType} className="gap-0">
            <TabsContent value="posts">
              {samplePostResults.length === 0 ? (
                <EmptyResult query={committedQuery} />
              ) : (
                <div>
                  {samplePostResults.map((post) => (
                    <PostCard key={post.id} post={post} onReply={handleReply} onRepost={handleRepost} />
                  ))}
                </div>
              )}
            </TabsContent>
            <TabsContent value="users">
              {userResultsLoading ? (
                <div className="flex justify-center items-center py-12 text-muted-foreground text-sm">読み込み中...</div>
              ) : userResults.length === 0 ? (
                <EmptyResult query={committedQuery} />
              ) : (
                <div>
                  {userResults.map((user) => (
                    <UserRow key={user.customId} user={user} />
                  ))}
                </div>
              )}
            </TabsContent>
            <TabsContent value="media">
              {sampleMediaPosts.length === 0 ? (
                <EmptyResult query={committedQuery} />
              ) : (
                <div>
                  {sampleMediaPosts.map((post) => (
                    <PostCard key={post.id} post={post} onReply={handleReply} onRepost={handleRepost} />
                  ))}
                </div>
              )}
            </TabsContent>
            <TabsContent value="calls">
              {sampleCallResults.length === 0 ? (
                <EmptyResult query={committedQuery} />
              ) : (
                <div className="divide-y">
                  {sampleCallResults.map((call) => (
                    <CallListItem key={call.id} call={call} onTap={setSelectedCall} />
                  ))}
                </div>
              )}
            </TabsContent>
          </Tabs>
        ) : (
          /* Default: trends + suggested users */
          <div>
            <div className="px-4 py-3 border-b">
              <h2 className="text-base font-bold flex items-center gap-2">
                <TrendingUp className="size-4" />
                トレンド
              </h2>
            </div>
            <TrendList trends={sampleTrends} />
            <SuggestedUserList users={sampleSuggestedUsers} />
          </div>
        )}
      </div>

      <ComposePostDialog
        open={composeMode !== null}
        onClose={() => setComposeMode(null)}
        onPost={() => setComposeMode(null)}
        mode={composeMode ?? undefined}
      />
      <JoinCallDialog call={selectedCall} onClose={() => setSelectedCall(null)} />
      <BottomNav />
    </div>
  );
}