import { Link } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { BottomNav } from "~/components/bottom-nav";
import { ArrowLeft } from "lucide-react";

const likedUsers = [
  { name: "佐藤花子", customId: "hanako_s", avatarUrl: "" },
  { name: "鈴木一郎", customId: "ichiro_dev", avatarUrl: "" },
  { name: "山田美咲", customId: "misaki_y", avatarUrl: "" },
  { name: "高橋健太", customId: "kenta_t", avatarUrl: "" },
  { name: "中村悠", customId: "yu_nkmr", avatarUrl: "" },
];

export default function PostLikes() {
  return (
    <div className="w-full min-h-full flex flex-col">
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-14">
          <Button variant="ghost" size="icon" onClick={() => history.back()}>
            <ArrowLeft className="size-5" />
          </Button>
          <h1 className="text-base font-bold">いいねしたユーザー</h1>
        </div>
      </div>

      <div className="flex-1">
        {likedUsers.map((user) => (
          <Link
            key={user.customId}
            to={`/users/${user.customId}`}
            className="flex items-center gap-3 px-4 py-3 border-b hover:bg-muted/30 transition-colors"
          >
            <Avatar className="size-10">
              <AvatarImage src={user.avatarUrl} alt={user.name} />
              <AvatarFallback>{user.name.slice(0, 2)}</AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-bold truncate">{user.name}</p>
              <p className="text-sm text-muted-foreground">@{user.customId}</p>
            </div>
          </Link>
        ))}
      </div>

      <BottomNav />
    </div>
  );
}
