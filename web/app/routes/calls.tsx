import { useState } from "react";
import { Link, useSearchParams } from "react-router";
import { ArrowLeft, Phone, Video } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { JoinCallDialog } from "~/components/join-call-dialog";
import { GroupAvatar, type Call } from "~/components/call-list";

const followingCalls: Call[] = [
  {
    id: "c1",
    name: "雑談部屋",
    type: "audio",
    participants: [{ name: "田中太郎", avatarUrl: "" }],
  },
  {
    id: "c2",
    name: "デザインレビュー",
    type: "video",
    participants: [
      { name: "佐藤花子", avatarUrl: "" },
      { name: "山田美咲", avatarUrl: "" },
    ],
  },
  {
    id: "c3",
    name: "開発ミーティング",
    type: "video",
    participants: [
      { name: "鈴木一郎", avatarUrl: "" },
      { name: "高橋健太", avatarUrl: "" },
      { name: "中村悠", avatarUrl: "" },
    ],
  },
  {
    id: "c4",
    name: "チーム定例",
    type: "audio",
    participants: [
      { name: "小林あおい", avatarUrl: "" },
      { name: "渡辺大輔", avatarUrl: "" },
      { name: "伊藤さくら", avatarUrl: "" },
      { name: "木村拓也", avatarUrl: "" },
    ],
  },
  {
    id: "c5",
    name: "作業通話",
    type: "audio",
    participants: [
      { name: "松本りな", avatarUrl: "" },
      { name: "井上翔", avatarUrl: "" },
    ],
  },
];

const publicCalls: Call[] = [
  {
    id: "c6",
    name: "企画ブレスト",
    type: "video",
    participants: [
      { name: "田中太郎", avatarUrl: "" },
      { name: "佐藤花子", avatarUrl: "" },
      { name: "鈴木一郎", avatarUrl: "" },
    ],
  },
  {
    id: "c7",
    name: "もくもく会",
    type: "audio",
    participants: [
      { name: "高橋健太", avatarUrl: "" },
      { name: "中村悠", avatarUrl: "" },
      { name: "小林あおい", avatarUrl: "" },
      { name: "渡辺大輔", avatarUrl: "" },
    ],
  },
  {
    id: "c8",
    name: "1on1",
    type: "video",
    participants: [{ name: "伊藤さくら", avatarUrl: "" }],
  },
  {
    id: "c9",
    name: "読書会",
    type: "audio",
    participants: [
      { name: "山田美咲", avatarUrl: "" },
      { name: "松本りな", avatarUrl: "" },
    ],
  },
  {
    id: "c10",
    name: "LT大会",
    type: "video",
    participants: [
      { name: "井上翔", avatarUrl: "" },
      { name: "田中太郎", avatarUrl: "" },
      { name: "佐藤花子", avatarUrl: "" },
      { name: "高橋健太", avatarUrl: "" },
    ],
  },
];

function CallListItem({
  call,
  onTap,
}: {
  call: Call;
  onTap: (call: Call) => void;
}) {
  const TypeIcon = call.type === "video" ? Video : Phone;
  const participantNames = call.participants.map((p) => p.name).join("、");

  return (
    <button
      onClick={() => onTap(call)}
      className="flex items-center gap-3 w-full px-4 py-3 hover:bg-muted/50 transition-colors"
    >
      <GroupAvatar participants={call.participants} className="size-12" />
      <div className="flex-1 text-left min-w-0">
        <div className="font-medium text-sm">{call.name}</div>
        <div className="text-xs text-muted-foreground truncate">
          {participantNames}
        </div>
      </div>
      <TypeIcon className="size-4 text-muted-foreground shrink-0" />
    </button>
  );
}

export default function Calls() {
  const [searchParams] = useSearchParams();
  const defaultTab = searchParams.get("tab") === "public" ? "public" : "following";
  const [selectedCall, setSelectedCall] = useState<Call | null>(null);

  return (
    <div className="w-full min-h-full">
      <Tabs defaultValue={defaultTab} className="gap-0">
        <div className="sticky top-0 w-full border-b bg-background z-10">
          <div className="flex items-center px-4 h-12 gap-3">
            <Link
              to="/"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              <ArrowLeft className="size-5" />
            </Link>
            <span className="font-medium">通話</span>
          </div>
          <TabsList variant="line" className="w-full h-12">
            <TabsTrigger value="following">フォロー中</TabsTrigger>
            <TabsTrigger value="public">オープン</TabsTrigger>
          </TabsList>
        </div>
        <TabsContent value="following">
          <div className="divide-y">
            {followingCalls.map((call) => (
              <CallListItem key={call.id} call={call} onTap={setSelectedCall} />
            ))}
          </div>
        </TabsContent>
        <TabsContent value="public">
          <div className="divide-y">
            {publicCalls.map((call) => (
              <CallListItem key={call.id} call={call} onTap={setSelectedCall} />
            ))}
          </div>
        </TabsContent>
      </Tabs>

      <JoinCallDialog call={selectedCall} onClose={() => setSelectedCall(null)} />
    </div>
  );
}
