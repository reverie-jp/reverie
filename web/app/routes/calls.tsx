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
    host: "田中太郎",
    participants: [{ name: "田中太郎", customId: "tanaka", avatarUrl: "" }],
  },
  {
    id: "c2",
    name: "デザインレビュー",
    type: "video",
    host: "佐藤花子",
    participants: [
      { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
      { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
    ],
  },
  {
    id: "c3",
    name: "開発ミーティング",
    type: "video",
    host: "鈴木一郎",
    participants: [
      { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
      { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
      { name: "中村悠", customId: "yu_nkmr", avatarUrl: "" },
      { name: "森田陽介", customId: "yosuke_m", avatarUrl: "" },
      { name: "川口真理", customId: "mari_k", avatarUrl: "" },
    ],
  },
  {
    id: "c4",
    name: "チーム定例",
    type: "audio",
    host: "小林あおい",
    participants: [
      { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
      { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
      { name: "伊藤さくら", customId: "sakura_ito", avatarUrl: "" },
      { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
    ],
  },
  {
    id: "c5",
    name: "スプリントレビュー",
    type: "video",
    host: "鈴木一郎",
    participants: [
      { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
      { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
      { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
      { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
      { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
      { name: "中村悠", customId: "yu_nkmr", avatarUrl: "" },
      { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
      { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
      { name: "森田陽介", customId: "yosuke_m", avatarUrl: "" },
      { name: "川口真理", customId: "mari_k", avatarUrl: "" },
      { name: "藤原誠", customId: "makoto_f", avatarUrl: "" },
    ],
  },
  {
    id: "c6",
    name: "作業通話",
    type: "audio",
    host: "松本りな",
    participants: [
      { name: "松本りな", customId: "rina_m", avatarUrl: "" },
      { name: "井上翔", customId: "sho_inoue", avatarUrl: "" },
    ],
  },
  {
    id: "c13",
    name: "朝会",
    type: "video",
    host: "森田陽介",
    participants: [
      { name: "森田陽介", customId: "yosuke_m", avatarUrl: "" },
      { name: "川口真理", customId: "mari_k", avatarUrl: "" },
      { name: "藤原誠", customId: "makoto_f", avatarUrl: "" },
      { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
      { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
      { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
      { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
      { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
      { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
    ],
  },
];

const publicCalls: Call[] = [
  {
    id: "c7",
    name: "企画ブレスト",
    type: "video",
    host: "田中太郎",
    participants: [
      { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
      { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
      { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
    ],
  },
  {
    id: "c8",
    name: "もくもく会",
    type: "audio",
    host: "高橋健太",
    participants: [
      { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
      { name: "中村悠", customId: "yu_nkmr", avatarUrl: "" },
      { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
      { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
    ],
  },
  {
    id: "c9",
    name: "1on1",
    type: "video",
    host: "伊藤さくら",
    participants: [{ name: "伊藤さくら", customId: "sakura_ito", avatarUrl: "" }],
  },
  {
    id: "c10",
    name: "読書会",
    type: "audio",
    host: "山田美咲",
    participants: [
      { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
      { name: "松本りな", customId: "rina_m", avatarUrl: "" },
    ],
  },
  {
    id: "c11",
    name: "LT大会",
    type: "video",
    host: "井上翔",
    participants: [
      { name: "井上翔", customId: "sho_inoue", avatarUrl: "" },
      { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
      { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
      { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
      { name: "中村悠", customId: "yu_nkmr", avatarUrl: "" },
      { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
      { name: "森田陽介", customId: "yosuke_m", avatarUrl: "" },
      { name: "川口真理", customId: "mari_k", avatarUrl: "" },
    ],
  },
  {
    id: "c12",
    name: "全社キックオフ",
    type: "video",
    host: "渡辺大輔",
    participants: [
      { name: "渡辺大輔", customId: "daisuke_w", avatarUrl: "" },
      { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
      { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
      { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
      { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
      { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
      { name: "中村悠", customId: "yu_nkmr", avatarUrl: "" },
      { name: "小林あおい", customId: "aoi_kb", avatarUrl: "" },
      { name: "伊藤さくら", customId: "sakura_ito", avatarUrl: "" },
      { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
      { name: "松本りな", customId: "rina_m", avatarUrl: "" },
      { name: "井上翔", customId: "sho_inoue", avatarUrl: "" },
      { name: "森田陽介", customId: "yosuke_m", avatarUrl: "" },
      { name: "川口真理", customId: "mari_k", avatarUrl: "" },
      { name: "藤原誠", customId: "makoto_f", avatarUrl: "" },
    ],
  },
  {
    id: "c14",
    name: "ゲーム配信",
    type: "audio",
    host: "藤原誠",
    participants: [
      { name: "藤原誠", customId: "makoto_f", avatarUrl: "" },
      { name: "田中太郎", customId: "tanaka", avatarUrl: "" },
      { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
      { name: "井上翔", customId: "sho_inoue", avatarUrl: "" },
      { name: "森田陽介", customId: "yosuke_m", avatarUrl: "" },
      { name: "木村拓也", customId: "takuya_k", avatarUrl: "" },
      { name: "中村悠", customId: "yu_nkmr", avatarUrl: "" },
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
      <div className={`size-8 rounded-full flex items-center justify-center shrink-0 ${call.type === "video" ? "bg-blue-500/10 text-blue-500" : "bg-green-500/10 text-green-500"}`}>
        <TypeIcon className="size-4" />
      </div>
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
