import { useState, useRef, useEffect } from "react";
import { useNavigate } from "react-router";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Textarea } from "~/components/ui/textarea";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { ArrowLeft, Camera, ImagePlus, Trash2 } from "lucide-react";
import { getMe, updateUser, isLoggedIn, type User } from "~/lib/api";

export default function ProfileEdit() {
  const navigate = useNavigate();

  const [user, setUser] = useState<User | null>(null);
  const [name, setName] = useState("");
  const [bio, setBio] = useState("");
  const [isPrivate, setIsPrivate] = useState(false);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
  const [bannerPreview, setBannerPreview] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const avatarInputRef = useRef<HTMLInputElement>(null);
  const bannerInputRef = useRef<HTMLInputElement>(null);

  const maxBioLength = 160;

  useEffect(() => {
    if (!isLoggedIn()) {
      navigate("/login");
      return;
    }
    getMe().then(({ user }) => {
      setUser(user);
      setName(user.display_name);
      setBio(user.biography ?? "");
      setIsPrivate(user.is_private);
    }).catch(() => navigate("/login"));
  }, []);

  const handleImageSelect = (
    e: React.ChangeEvent<HTMLInputElement>,
    setter: (url: string | null) => void,
  ) => {
    const file = e.target.files?.[0];
    if (!file) return;
    const url = URL.createObjectURL(file);
    setter(url);
    e.target.value = "";
  };

  const handleSave = async () => {
    if (!name.trim()) {
      setError("名前を入力してください");
      return;
    }
    setSaving(true);
    setError(null);
    try {
      await updateUser({
        display_name: name,
        biography: bio,
        is_private: isPrivate,
      });
      navigate(-1);
    } catch (e: any) {
      setError(e.message ?? "保存に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="w-full min-h-full flex flex-col">
      {/* Hidden file inputs */}
      <input
        ref={avatarInputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={(e) => handleImageSelect(e, setAvatarPreview)}
      />
      <input
        ref={bannerInputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={(e) => handleImageSelect(e, setBannerPreview)}
      />

      {/* Header */}
      <div className="sticky top-0 left-0 w-full border-b bg-background/60 backdrop-blur-lg z-10">
        <div className="flex items-center justify-between px-4 h-14">
          <div className="flex items-center gap-3">
            <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
              <ArrowLeft className="size-5" />
            </Button>
            <h1 className="text-base font-bold">プロフィールを編集</h1>
          </div>
          <Button
            className="rounded-full h-8 px-4"
            onClick={handleSave}
            disabled={saving}
          >
            {saving ? "保存中..." : "保存"}
          </Button>
        </div>
      </div>

      {error && (
        <div className="mx-4 mt-3 px-3 py-2 bg-destructive/10 text-destructive text-sm rounded-md">
          {error}
        </div>
      )}

      {/* Banner */}
      <DropdownMenu>
        <DropdownMenuTrigger
          render={
            <button className="relative w-full h-32 bg-muted group overflow-hidden" />
          }
        >
          {bannerPreview && (
            <img
              src={bannerPreview}
              alt=""
              className="absolute inset-0 w-full h-full object-cover"
            />
          )}
          <div className="absolute inset-0 flex items-center justify-center bg-black/0 group-hover:bg-black/30 transition-colors">
            <Camera className="size-6 text-white opacity-0 group-hover:opacity-100 transition-opacity" />
          </div>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem onClick={() => bannerInputRef.current?.click()}>
            <ImagePlus className="size-4" />
            画像をアップロード
          </DropdownMenuItem>
          {bannerPreview && (
            <DropdownMenuItem
              className="text-destructive"
              onClick={() => setBannerPreview(null)}
            >
              <Trash2 className="size-4" />
              画像を削除
            </DropdownMenuItem>
          )}
        </DropdownMenuContent>
      </DropdownMenu>

      {/* Avatar */}
      <div className="px-4 -mt-12 relative z-1">
        <DropdownMenu>
          <DropdownMenuTrigger
            render={<button className="relative group w-fit rounded-full" />}
          >
            <Avatar className="size-20 ring-4 ring-background">
              <AvatarImage src={avatarPreview ?? undefined} alt={name} />
              <AvatarFallback className="text-2xl">
                {name.slice(0, 2)}
              </AvatarFallback>
            </Avatar>
            <div className="absolute inset-0 rounded-full flex items-center justify-center bg-black/0 group-hover:bg-black/40 transition-colors">
              <Camera className="size-5 text-white opacity-0 group-hover:opacity-100 transition-opacity" />
            </div>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem onClick={() => avatarInputRef.current?.click()}>
              <ImagePlus className="size-4" />
              画像をアップロード
            </DropdownMenuItem>
            {avatarPreview && (
              <DropdownMenuItem
                className="text-destructive"
                onClick={() => setAvatarPreview(null)}
              >
                <Trash2 className="size-4" />
                画像を削除
              </DropdownMenuItem>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {/* Form */}
      <div className="px-4 mt-6 space-y-5 pb-8">
        <div className="space-y-1.5">
          <span className="text-xs text-muted-foreground">名前</span>
          <Input
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="h-10"
            maxLength={20}
          />
        </div>

        <div className="space-y-1.5">
          <div className="flex items-center justify-between">
            <span className="text-xs text-muted-foreground">ユーザーID</span>
            <span className="text-xs text-muted-foreground">
              @{user?.custom_id ?? ""}
            </span>
          </div>
        </div>

        <div className="space-y-1.5">
          <div className="flex items-center justify-between">
            <span className="text-xs text-muted-foreground">自己紹介</span>
            <span
              className={`text-xs ${bio.length > maxBioLength ? "text-destructive" : "text-muted-foreground"}`}
            >
              {bio.length}/{maxBioLength}
            </span>
          </div>
          <Textarea
            value={bio}
            onChange={(e) => setBio(e.target.value)}
            placeholder="自己紹介を入力"
            className="min-h-20"
            maxLength={maxBioLength}
          />
        </div>
      </div>
    </div>
  );
}
