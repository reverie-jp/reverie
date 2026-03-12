import { BottomNav } from "~/components/bottom-nav";
import { ComposeFab } from "~/components/compose-fab";

export default function Notifications() {
  return (
    <div className="w-full h-full">
      <div className="sticky top-0 left-0 w-full h-12 border-b bg-background flex items-center px-4 font-medium">
        通知
      </div>
      <div className="grid place-items-center h-200">通知のコンテンツ</div>
      <ComposeFab />
      <BottomNav />
    </div>
  );
}
