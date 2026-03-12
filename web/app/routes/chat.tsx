import { BottomNav } from "~/components/bottom-nav";

export default function Chat() {
  return (
    <div className="w-full h-full">
      <div className="sticky top-0 left-0 w-full h-12 border-b bg-background flex items-center px-4 font-medium">
        チャット
      </div>
      <div className="grid place-items-center h-200">チャットのコンテンツ</div>
      <BottomNav />
    </div>
  );
}
