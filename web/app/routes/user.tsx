import { BottomNav } from "~/components/bottom-nav";
import type { Route } from "./+types/user";

export default function User({ params }: Route.ComponentProps) {
  return (
    <div className="w-full h-full">
      <div className="sticky top-0 left-0 w-full h-12 border-b bg-background flex items-center px-4 font-medium">
        プロフィール
      </div>
      <div className="grid place-items-center h-200">
        ユーザー: {params.id}
      </div>
      <BottomNav />
    </div>
  );
}
