import { useState } from "react";
import { Link } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { ArrowLeft, Check, X } from "lucide-react";

interface MessageRequest {
  id: string;
  chatId: string;
  sender: {
    name: string;
    customId: string;
    avatarUrl: string;
  };
  lastMessage: string;
  receivedAt: Date;
}

const initialRequests: MessageRequest[] = [
  {
    id: "req1",
    chatId: "chat-req1",
    sender: {
      name: "小林さくら",
      customId: "sakura_k",
      avatarUrl: "",
    },
    lastMessage: "はじめまして！共通の趣味がありそうなのでメッセージしました。",
    receivedAt: new Date(Date.now() - 30 * 60_000),
  },
  {
    id: "req2",
    chatId: "chat-req2",
    sender: {
      name: "渡辺誠",
      customId: "makoto_w",
      avatarUrl: "",
    },
    lastMessage: "先日のイベントでお話した渡辺です。よろしくお願いします！",
    receivedAt: new Date(Date.now() - 3 * 3_600_000),
  },
  {
    id: "req3",
    chatId: "chat-req3",
    sender: {
      name: "松本あい",
      customId: "ai_matsumoto",
      avatarUrl: "",
    },
    lastMessage: "こんにちは、あなたの投稿をよく見ています。",
    receivedAt: new Date(Date.now() - 86_400_000),
  },
];

function formatRequestTime(date: Date): string {
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMinutes = Math.floor(diffMs / 60_000);
  const diffHours = Math.floor(diffMs / 3_600_000);
  const diffDays = Math.floor(diffMs / 86_400_000);

  if (diffMinutes < 1) return "たった今";
  if (diffMinutes < 60) return `${diffMinutes}分前`;
  if (diffHours < 24) return `${diffHours}時間前`;
  if (diffDays < 7) return `${diffDays}日前`;
  return `${date.getMonth() + 1}/${date.getDate()}`;
}

export default function ChatRequests() {
  const [requests, setRequests] = useState(initialRequests);

  const handleAccept = (id: string) => {
    setRequests((prev) => prev.filter((r) => r.id !== id));
  };

  const handleReject = (id: string) => {
    setRequests((prev) => prev.filter((r) => r.id !== id));
  };

  return (
    <div className="w-full min-h-full flex flex-col">
      {/* Header */}
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-2 h-14">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <h1 className="text-base font-bold">メッセージリクエスト</h1>
        </div>
      </div>

      {/* Info */}
      <div className="px-4 py-3 border-b bg-muted/30">
        <p className="text-xs text-muted-foreground">
          フォローしていないユーザーからのメッセージです。承認するとチャットが開始されます。
        </p>
      </div>

      {/* Request list */}
      <div className="flex-1">
        {requests.length === 0 ? (
          <div className="flex items-center justify-center py-16 text-muted-foreground text-sm">
            メッセージリクエストはありません
          </div>
        ) : (
          requests.map((req) => (
            <div
              key={req.id}
              className="flex items-start gap-3 px-4 py-4 border-b"
            >
              <Link
                to={`/users/${req.sender.customId}`}
                className="shrink-0"
              >
                <Avatar className="size-12">
                  <AvatarImage
                    src={req.sender.avatarUrl}
                    alt={req.sender.name}
                  />
                  <AvatarFallback>
                    {req.sender.name.slice(0, 2)}
                  </AvatarFallback>
                </Avatar>
              </Link>
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between">
                  <Link
                    to={`/users/${req.sender.customId}`}
                    className="min-w-0"
                  >
                    <span className="text-sm font-bold truncate">
                      {req.sender.name}
                    </span>
                    <span className="text-xs text-muted-foreground ml-1.5">
                      @{req.sender.customId}
                    </span>
                  </Link>
                  <span className="text-xs text-muted-foreground whitespace-nowrap ml-2">
                    {formatRequestTime(req.receivedAt)}
                  </span>
                </div>
                <Link
                  to={`/chat/${req.chatId}`}
                  className="block mt-1 hover:underline"
                >
                  <p className="text-sm text-muted-foreground line-clamp-2">
                    {req.lastMessage}
                  </p>
                </Link>
                <div className="flex items-center gap-2 mt-3">
                  <Button
                    size="sm"
                    onClick={() => handleAccept(req.id)}
                    className="flex-1"
                  >
                    <Check className="size-4 mr-1" />
                    承認
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleReject(req.id)}
                    className="flex-1"
                  >
                    <X className="size-4 mr-1" />
                    拒否
                  </Button>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
