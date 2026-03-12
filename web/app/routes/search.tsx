import { useState } from "react";
import { Search } from "lucide-react";
import { Input } from "~/components/ui/input";
import { BottomNav } from "~/components/bottom-nav";

export default function SearchPage() {
  const [query, setQuery] = useState("");

  return (
    <div className="w-full min-h-full flex flex-col">
      <div className="sticky top-0 left-0 w-full border-b bg-background z-10">
        <div className="flex items-center px-4 h-14 gap-3">
          <span className="font-medium">検索</span>
        </div>
        <div className="px-4 pb-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
            <Input
              placeholder="キーワードで検索"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              className="pl-9 h-10"
            />
          </div>
        </div>
      </div>
      <div className="flex-1">
        {query.trim().length === 0 ? (
          <div className="flex items-center justify-center h-60 text-sm text-muted-foreground">
            キーワードを入力して検索
          </div>
        ) : (
          <div className="flex items-center justify-center h-60 text-sm text-muted-foreground">
            「{query}」の検索結果はありません
          </div>
        )}
      </div>
      <BottomNav />
    </div>
  );
}
