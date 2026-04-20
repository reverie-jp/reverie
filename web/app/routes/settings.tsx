import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { clearTokens } from "~/lib/api";
import { useTheme, THEME_OPTIONS } from "~/lib/theme-context";
import { Input } from "~/components/ui/input";
import { Switch } from "~/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "~/components/ui/alert-dialog";
import {
  ArrowLeft,
  Bell,
  ChevronRight,
  Eye,
  Globe,
  Lock,
  LogOut,
  Palette,
  Shield,
  Smartphone,
  Trash2,
  User,
  Volume2,
  UserX,
} from "lucide-react";

type SettingsSection = {
  label: string;
  items: SettingsItem[];
};

type SettingsItem =
  | { type: "toggle"; icon: React.ElementType; label: string; description?: string; key: string }
  | { type: "select"; icon: React.ElementType; label: string; description?: string; key: string; options: { value: string; label: string }[] }
  | { type: "link"; icon: React.ElementType; label: string; description?: string; to: string }
  | { type: "action"; icon: React.ElementType; label: string; description?: string; variant?: "destructive"; action: string };

const sections: SettingsSection[] = [
  {
    label: "アカウント",
    items: [
      { type: "link", icon: User, label: "プロフィールを編集", to: "/settings/profile" },
      {
        type: "toggle",
        icon: Lock,
        label: "非公開アカウント",
        description: "承認したユーザーだけがあなたの投稿を見ることができます",
        key: "privateAccount",
      },
      { type: "link", icon: UserX, label: "ブロックしたアカウント", to: "/settings/blocked" },
    ],
  },
  {
    label: "通知",
    items: [
      {
        type: "toggle",
        icon: Bell,
        label: "プッシュ通知",
        description: "いいね、返信、フォローなどの通知を受け取ります",
        key: "pushNotifications",
      },
      {
        type: "toggle",
        icon: Volume2,
        label: "通知サウンド",
        key: "notificationSound",
      },
    ],
  },
  {
    label: "プライバシー",
    items: [
      {
        type: "toggle",
        icon: Eye,
        label: "オンラインステータスを表示",
        description: "他のユーザーにオンライン状態を表示します",
        key: "showOnlineStatus",
      },
      {
        type: "select",
        icon: Globe,
        label: "投稿のデフォルト公開範囲",
        key: "defaultVisibility",
        options: [
          { value: "public", label: "全体公開" },
          { value: "followers", label: "フォロワーのみ" },
          { value: "mutual", label: "相互フォローのみ" },
        ],
      },
    ],
  },
  {
    label: "表示",
    items: [
      {
        type: "select",
        icon: Palette,
        label: "テーマ",
        key: "theme",
        options: THEME_OPTIONS.map((o) => ({
          value: o.value,
          label: o.premium ? `${o.label} ★` : o.label,
        })),
      },
    ],
  },
  {
    label: "セキュリティ",
    items: [
      { type: "link", icon: Shield, label: "パスワードを変更", to: "/settings/password" },
      { type: "link", icon: Lock, label: "二段階認証", to: "/settings/2fa" },
      { type: "link", icon: Smartphone, label: "ログイン中のデバイス", to: "/settings/sessions" },
    ],
  },
  {
    label: "その他",
    items: [
      { type: "action", icon: LogOut, label: "ログアウト", action: "logout" },
      {
        type: "action",
        icon: Trash2,
        label: "アカウントを削除",
        description: "この操作は取り消せません",
        variant: "destructive",
        action: "deleteAccount",
      },
    ],
  },
];

const defaultToggles: Record<string, boolean> = {
  privateAccount: false,
  pushNotifications: true,
  notificationSound: true,
  showOnlineStatus: true,
};

const defaultSelects: Record<string, string> = {
  defaultVisibility: "public",
};

export default function Settings() {
  const navigate = useNavigate();
  const { theme, setTheme } = useTheme();
  const [toggles, setToggles] = useState(defaultToggles);
  const [selects, setSelects] = useState(defaultSelects);
  const [confirmDialog, setConfirmDialog] = useState<"logout" | "deleteAccount" | null>(null);
  const [deleteInput, setDeleteInput] = useState("");

  const handleToggle = (key: string) => {
    setToggles((prev) => ({ ...prev, [key]: !prev[key] }));
  };

  const handleSelect = (key: string, value: string | null) => {
    if (!value) return;
    if (key === "theme") {
      setTheme(value as "light" | "dark" | "system");
    } else {
      setSelects((prev) => ({ ...prev, [key]: value }));
    }
  };

  return (
    <div className="w-full min-h-full flex flex-col">
      {/* Header */}
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center gap-3 px-4 h-14">
          <button
            onClick={() => history.back()}
            className="text-muted-foreground hover:text-foreground transition-colors"
          >
            <ArrowLeft className="size-5" />
          </button>
          <h1 className="font-bold text-base">設定</h1>
        </div>
      </div>

      <div className="flex-1 pb-8">
        {sections.map((section) => (
          <div key={section.label}>
            <div className="px-4 pt-6 pb-2">
              <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                {section.label}
              </h2>
            </div>
            <div className="divide-y">
              {section.items.map((item) => {
                if (item.type === "toggle") {
                  const Icon = item.icon;
                  return (
                    <div
                      key={item.key}
                      className="flex items-center gap-3 px-4 py-3"
                    >
                      <Icon className="size-5 text-muted-foreground shrink-0" />
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium">{item.label}</p>
                        {item.description && (
                          <p className="text-xs text-muted-foreground mt-0.5">
                            {item.description}
                          </p>
                        )}
                      </div>
                      <Switch
                        checked={toggles[item.key] ?? false}
                        onCheckedChange={() => handleToggle(item.key)}
                      />
                    </div>
                  );
                }

                if (item.type === "select") {
                  const Icon = item.icon;
                  return (
                    <div
                      key={item.key}
                      className="flex items-center gap-3 px-4 py-3"
                    >
                      <Icon className="size-5 text-muted-foreground shrink-0" />
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium">{item.label}</p>
                        {item.description && (
                          <p className="text-xs text-muted-foreground mt-0.5">
                            {item.description}
                          </p>
                        )}
                      </div>
                      <Select
                        value={item.key === "theme" ? theme : selects[item.key]}
                        onValueChange={(v) => handleSelect(item.key, v)}
                      >
                        <SelectTrigger className="w-auto min-w-28 h-8 text-xs">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          {item.options.map((opt) => (
                            <SelectItem key={opt.value} value={opt.value}>
                              {opt.label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  );
                }

                if (item.type === "link") {
                  const Icon = item.icon;
                  return (
                    <Link
                      key={item.to}
                      to={item.to}
                      className="flex items-center gap-3 px-4 py-3 hover:bg-muted/50 transition-colors"
                    >
                      <Icon className="size-5 text-muted-foreground shrink-0" />
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium">{item.label}</p>
                        {item.description && (
                          <p className="text-xs text-muted-foreground mt-0.5">
                            {item.description}
                          </p>
                        )}
                      </div>
                      <ChevronRight className="size-4 text-muted-foreground shrink-0" />
                    </Link>
                  );
                }

                if (item.type === "action") {
                  const Icon = item.icon;
                  const isDestructive = item.variant === "destructive";
                  return (
                    <button
                      key={item.action}
                      onClick={() => setConfirmDialog(item.action as "logout" | "deleteAccount")}
                      className={`flex items-center gap-3 px-4 py-3 w-full hover:bg-muted/50 transition-colors ${
                        isDestructive ? "text-destructive" : ""
                      }`}
                    >
                      <Icon className="size-5 shrink-0" />
                      <div className="flex-1 text-left min-w-0">
                        <p className="text-sm font-medium">{item.label}</p>
                        {item.description && (
                          <p className="text-xs opacity-70 mt-0.5">
                            {item.description}
                          </p>
                        )}
                      </div>
                    </button>
                  );
                }

                return null;
              })}
            </div>
          </div>
        ))}
      </div>

      {/* Logout confirmation */}
      <AlertDialog
        open={confirmDialog === "logout"}
        onOpenChange={(open) => !open && setConfirmDialog(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>ログアウト</AlertDialogTitle>
            <AlertDialogDescription>
              ログアウトしてもよろしいですか？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel size="lg">キャンセル</AlertDialogCancel>
            <AlertDialogAction
              size="lg"
              onClick={() => {
                clearTokens();
                setConfirmDialog(null);
                navigate("/login");
              }}
            >
              ログアウト
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Delete account confirmation */}
      <AlertDialog
        open={confirmDialog === "deleteAccount"}
        onOpenChange={(open) => {
          if (!open) {
            setConfirmDialog(null);
            setDeleteInput("");
          }
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>アカウントを削除</AlertDialogTitle>
            <AlertDialogDescription>
              アカウントを削除すると、すべてのデータが永久に失われます。この操作は取り消せません。確認のため、あなたのID「<span className="font-mono font-semibold text-foreground">me</span>」を入力してください。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <Input
            placeholder="IDを入力"
            value={deleteInput}
            onChange={(e) => setDeleteInput(e.target.value)}
          />
          <AlertDialogFooter>
            <AlertDialogCancel size="lg" onClick={() => setDeleteInput("")}>
              キャンセル
            </AlertDialogCancel>
            <AlertDialogAction
              size="lg"
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90 disabled:opacity-50 disabled:pointer-events-none"
              disabled={deleteInput !== "me"}
              onClick={() => {
                setConfirmDialog(null);
                setDeleteInput("");
                navigate("/login");
              }}
            >
              削除する
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
