import { Link } from "react-router";
import { ArrowLeft, Phone, Video } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import type { Call } from "~/components/call-list";

const allCalls: Call[] = [
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
];

export default function Calls() {
  return (
    <div className="w-full min-h-full">
      <div className="sticky top-0 w-full h-12 border-b bg-background z-10 flex items-center px-4 gap-3">
        <Link
          to="/"
          className="text-muted-foreground hover:text-foreground transition-colors"
        >
          <ArrowLeft className="size-5" />
        </Link>
        <span className="font-medium">通話</span>
      </div>
      <div className="divide-y">
        {allCalls.map((call) => {
          const TypeIcon = call.type === "video" ? Video : Phone;
          const participantNames = call.participants
            .map((p) => p.name)
            .join("、");
          return (
            <button
              key={call.id}
              className="flex items-center gap-3 w-full px-4 py-3 hover:bg-muted/50 transition-colors"
            >
              <Avatar className="size-12">
                <AvatarImage
                  src={call.participants[0]?.avatarUrl}
                  alt={call.name}
                />
                <AvatarFallback>
                  {call.name.slice(0, 2)}
                </AvatarFallback>
              </Avatar>
              <div className="flex-1 text-left">
                <div className="font-medium text-sm">{call.name}</div>
                <div className="text-xs text-muted-foreground">
                  {participantNames}
                </div>
              </div>
              <TypeIcon className="size-4 text-muted-foreground shrink-0" />
            </button>
          );
        })}
      </div>
    </div>
  );
}
