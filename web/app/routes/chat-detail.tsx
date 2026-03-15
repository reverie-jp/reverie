import { useState, useRef, useEffect } from "react";
import { useParams, Link } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import {
  ArrowLeft,
  SendHorizontal,
  Phone,
  Video,
  ImagePlus,
} from "lucide-react";

type OnlineStatus = "online" | "idle" | "offline";

interface ChatMessage {
  id: string;
  content: string;
  sentAt: Date;
  isMine: boolean;
}

interface ChatParticipant {
  name: string;
  customId: string;
  avatarUrl: string;
  status: OnlineStatus;
}

const statusColor: Record<OnlineStatus, string> = {
  online: "bg-green-500",
  idle: "bg-yellow-500",
  offline: "bg-gray-400",
};

const statusLabel: Record<OnlineStatus, string> = {
  online: "オンライン",
  idle: "離席中",
  offline: "オフライン",
};

const chatData: Record<
  string,
  { participant: ChatParticipant; messages: ChatMessage[] }
> = {
  chat1: {
    participant: {
      name: "佐藤花子",
      customId: "hanako_s",
      avatarUrl: "",
      status: "online",
    },
    messages: [
      {
        id: "m1",
        content: "こんにちは！元気ですか？",
        sentAt: new Date(Date.now() - 60 * 60_000),
        isMine: false,
      },
      {
        id: "m2",
        content: "元気ですよ！最近どうですか？",
        sentAt: new Date(Date.now() - 55 * 60_000),
        isMine: true,
      },
      {
        id: "m3",
        content: "明日の予定ってもう決まりましたか？",
        sentAt: new Date(Date.now() - 50 * 60_000),
        isMine: false,
      },
      {
        id: "m4",
        content: "まだ決まってないです。何かありますか？",
        sentAt: new Date(Date.now() - 45 * 60_000),
        isMine: true,
      },
      {
        id: "m5",
        content:
          "新しくオープンしたカフェに行きませんか？すごくおしゃれらしいですよ。",
        sentAt: new Date(Date.now() - 40 * 60_000),
        isMine: false,
      },
      {
        id: "m6",
        content: "いいですね！行きましょう！何時がいいですか？",
        sentAt: new Date(Date.now() - 35 * 60_000),
        isMine: true,
      },
      {
        id: "m7",
        content: "14時くらいはどうですか？",
        sentAt: new Date(Date.now() - 10 * 60_000),
        isMine: false,
      },
      {
        id: "m8",
        content: "了解です！明日の予定について確認しますね。",
        sentAt: new Date(Date.now() - 5 * 60_000),
        isMine: false,
      },
    ],
  },
  chat2: {
    participant: {
      name: "鈴木一郎",
      customId: "ichiro_dev",
      avatarUrl: "",
      status: "online",
    },
    messages: [
      {
        id: "m1",
        content: "PRのレビューお願いできますか？",
        sentAt: new Date(Date.now() - 3 * 3_600_000),
        isMine: false,
      },
      {
        id: "m2",
        content: "もちろんです！見てみますね。",
        sentAt: new Date(Date.now() - 2.5 * 3_600_000),
        isMine: true,
      },
      {
        id: "m3",
        content: "いくつかコメント入れました。全体的にはとても良いと思います！",
        sentAt: new Date(Date.now() - 2 * 3_600_000),
        isMine: true,
      },
      {
        id: "m4",
        content: "ありがとうございます！修正しました。",
        sentAt: new Date(Date.now() - 60 * 60_000),
        isMine: false,
      },
      {
        id: "m5",
        content: "コードレビューありがとうございました！",
        sentAt: new Date(Date.now() - 30 * 60_000),
        isMine: true,
      },
    ],
  },
  chat3: {
    participant: {
      name: "山田美咲",
      customId: "misaki_y",
      avatarUrl: "",
      status: "idle",
    },
    messages: [
      {
        id: "m1",
        content: "昨日の写真、整理できました！",
        sentAt: new Date(Date.now() - 3 * 3_600_000),
        isMine: false,
      },
      {
        id: "m2",
        content: "おー！見たい見たい！",
        sentAt: new Date(Date.now() - 2.5 * 3_600_000),
        isMine: true,
      },
      {
        id: "m3",
        content: "写真送りました！見てみてください 📷",
        sentAt: new Date(Date.now() - 2 * 3_600_000),
        isMine: false,
      },
    ],
  },
};

function formatMessageTime(date: Date): string {
  return `${date.getHours()}:${String(date.getMinutes()).padStart(2, "0")}`;
}

function shouldShowTimestamp(
  current: ChatMessage,
  previous?: ChatMessage,
): boolean {
  if (!previous) return true;
  const diff = current.sentAt.getTime() - previous.sentAt.getTime();
  return diff > 10 * 60_000;
}

export default function ChatDetail() {
  const { id } = useParams();
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const scrollRef = useRef<HTMLDivElement>(null);

  const chat = id ? chatData[id] : undefined;

  useEffect(() => {
    if (chat) {
      setMessages(chat.messages);
    }
  }, [id]);

  useEffect(() => {
    scrollRef.current?.scrollTo({
      top: scrollRef.current.scrollHeight,
    });
  }, [messages]);

  if (!chat) {
    return (
      <div className="w-full h-full flex items-center justify-center text-muted-foreground">
        チャットが見つかりません
      </div>
    );
  }

  const { participant } = chat;

  const handleSend = () => {
    const trimmed = input.trim();
    if (!trimmed) return;
    setMessages((prev) => [
      ...prev,
      {
        id: `new-${Date.now()}`,
        content: trimmed,
        sentAt: new Date(),
        isMine: true,
      },
    ]);
    setInput("");
  };

  return (
    <div className="w-full h-dvh flex flex-col">
      {/* Header */}
      <div className="shrink-0 border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-16">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <Link
            to={`/users/${participant.customId}`}
            className="flex items-center gap-3 flex-1 min-w-0"
          >
            <div className="relative shrink-0">
              <Avatar className="size-9">
                <AvatarImage
                  src={participant.avatarUrl}
                  alt={participant.name}
                />
                <AvatarFallback>{participant.name.slice(0, 2)}</AvatarFallback>
              </Avatar>
              <span
                className={`absolute bottom-0 right-0 size-2.5 rounded-full border-2 border-background ${statusColor[participant.status]}`}
              />
            </div>
            <div className="min-w-0">
              <p className="text-sm font-bold truncate">{participant.name}</p>
              <p className="text-xs text-muted-foreground">
                {statusLabel[participant.status]}
              </p>
            </div>
          </Link>
          <div className="flex items-center gap-2">
            <Button variant="ghost" size="icon">
              <Phone className="size-5" />
            </Button>
            <Button variant="ghost" size="icon">
              <Video className="size-5" />
            </Button>
          </div>
        </div>
      </div>

      {/* Messages */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto px-4 py-4">
        <div className="flex flex-col gap-1">
          {messages.map((msg, i) => {
            const showTime = shouldShowTimestamp(msg, messages[i - 1]);
            const isConsecutive =
              i > 0 && messages[i - 1].isMine === msg.isMine && !showTime;

            return (
              <div key={msg.id}>
                {showTime && (
                  <div className="flex justify-center my-3">
                    <span className="text-xs text-muted-foreground bg-muted/50 px-2 py-0.5 rounded-full">
                      {formatMessageTime(msg.sentAt)}
                    </span>
                  </div>
                )}
                <div
                  className={`flex items-start ${msg.isMine ? "justify-end" : "justify-start"} ${isConsecutive ? "mt-0.5" : "mt-2"}`}
                >
                  {!msg.isMine && !isConsecutive && (
                    <Avatar className="size-8 mr-2 shrink-0 self-start">
                      <AvatarImage
                        src={participant.avatarUrl}
                        alt={participant.name}
                      />
                      <AvatarFallback className="text-xs">
                        {participant.name.slice(0, 2)}
                      </AvatarFallback>
                    </Avatar>
                  )}
                  {!msg.isMine && isConsecutive && (
                    <div className="size-8 mr-2 shrink-0" />
                  )}
                  <div
                    className={`max-w-[75%] px-3 py-2 text-sm whitespace-pre-wrap wrap-break-word ${
                      msg.isMine
                        ? "bg-blue-500 text-white rounded-2xl rounded-tr-md"
                        : "bg-muted rounded-2xl rounded-tl-md"
                    }`}
                  >
                    {msg.content}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Input */}
      <div className="shrink-0 border-t bg-background px-4 py-3">
        <div className="flex gap-2">
          <button
            className="shrink-0 p-2 text-muted-foreground hover:text-foreground transition-colors"
            onClick={() => {
              // TODO: open file picker for images/videos
            }}
          >
            <ImagePlus className="size-5" />
          </button>
          <textarea
            placeholder="メッセージを入力..."
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => {
              if (
                e.key === "Enter" &&
                !e.shiftKey &&
                !e.nativeEvent.isComposing
              ) {
                e.preventDefault();
                handleSend();
              }
            }}
            rows={1}
            className="flex-1 resize-none rounded-md border border-input bg-background px-3 py-2 text-sm min-h-10 max-h-32 field-sizing-content placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
          <button
            className="shrink-0 px-1 text-blue-500 disabled:text-muted-foreground transition-colors"
            onClick={handleSend}
            disabled={!input.trim()}
          >
            <SendHorizontal className="size-5" />
          </button>
        </div>
      </div>
    </div>
  );
}
