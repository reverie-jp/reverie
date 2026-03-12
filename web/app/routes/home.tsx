import { Search } from "lucide-react";
import { Input } from "~/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { BottomNav } from "~/components/bottom-nav";
import type { Route } from "./+types/home";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Welcome to React Router!" },
  ];
}

export default function Home() {
  return (
    <div className="w-full h-full">
      <Tabs defaultValue="following" className="gap-0">
        <div className="sticky top-0 left-0 w-full border-b bg-background">
          <div className="px-4 pt-5 pb-2">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
              <Input placeholder="検索" className="pl-9 h-10" />
            </div>
          </div>
          <TabsList variant="line" className="w-full h-12">
            <TabsTrigger value="following">フォロー中</TabsTrigger>
            <TabsTrigger value="public">オープン</TabsTrigger>
          </TabsList>
        </div>
        <TabsContent value="following">
          <div className="grid place-items-center h-200">
            フォロー中のコンテンツ
          </div>
        </TabsContent>
        <TabsContent value="public">
          <div className="grid place-items-center h-200">
            オープンのコンテンツ
          </div>
        </TabsContent>
      </Tabs>
      <BottomNav />
    </div>
  );
}
