import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { socialLogin, setTokens, isLoggedIn } from "~/lib/api";

const IS_DEV = import.meta.env.DEV;

const GOOGLE_CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID ?? "";
const REDIRECT_URI = "http://localhost:5173/login";

function buildGoogleOAuthUrl() {
  const params = new URLSearchParams({
    client_id: GOOGLE_CLIENT_ID,
    redirect_uri: REDIRECT_URI,
    response_type: "code",
    scope: "openid email profile",
    access_type: "offline",
  });
  return `https://accounts.google.com/o/oauth2/v2/auth?${params}`;
}

function GoogleIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" className={className}>
      <path
        d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"
        fill="#4285F4"
      />
      <path
        d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
        fill="#34A853"
      />
      <path
        d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
        fill="#FBBC05"
      />
      <path
        d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
        fill="#EA4335"
      />
    </svg>
  );
}

function LineIcon({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 24 24" className={className} fill="#06C755">
      <path d="M24 10.304c0-5.369-5.383-9.738-12-9.738S0 4.935 0 10.304c0 4.814 4.27 8.846 10.035 9.608.391.084.922.258 1.057.592.121.303.079.778.039 1.084l-.171 1.027c-.053.303-.242 1.186 1.039.647 1.281-.54 6.911-4.069 9.428-6.967C23.309 14.253 24 12.38 24 10.304zm-16.66 3.088a.348.348 0 0 1-.348.348H4.51a.348.348 0 0 1-.348-.348V8.585a.348.348 0 0 1 .348-.348h.695a.348.348 0 0 1 .348.348v4.112h1.938a.348.348 0 0 1 .348.348v.347zm2.097 0a.348.348 0 0 1-.348.348h-.694a.348.348 0 0 1-.348-.348V8.585a.348.348 0 0 1 .348-.348h.694a.348.348 0 0 1 .348.348v4.807zm5.628 0a.348.348 0 0 1-.348.348h-.694a.35.35 0 0 1-.273-.132l-1.98-2.678v2.462a.348.348 0 0 1-.348.348h-.694a.348.348 0 0 1-.349-.348V8.585a.348.348 0 0 1 .349-.348h.694a.35.35 0 0 1 .272.131l1.981 2.679V8.585a.348.348 0 0 1 .348-.348h.694a.348.348 0 0 1 .348.348v4.807zm3.839-3.764a.348.348 0 0 1-.348.348h-1.938v.96h1.938a.348.348 0 0 1 .348.348v.695a.348.348 0 0 1-.348.348h-1.938v.96h1.938a.348.348 0 0 1 .348.348v.347a.348.348 0 0 1-.348.348h-2.98a.348.348 0 0 1-.349-.348V8.585a.348.348 0 0 1 .349-.348h2.98a.348.348 0 0 1 .348.348v.695z" />
    </svg>
  );
}

export default function Login() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [devToken, setDevToken] = useState("");

  useEffect(() => {
    if (isLoggedIn()) {
      navigate("/");
      return;
    }

    const code = searchParams.get("code");
    if (code) {
      socialLogin("GOOGLE", code)
        .then(({ token_pair }) => {
          setTokens(token_pair.access_token, token_pair.refresh_token);
          navigate("/");
        })
        .catch((e) => {
          console.error("login failed:", e);
        });
    }
  }, [searchParams]);

  const handleGoogleLogin = () => {
    if (!GOOGLE_CLIENT_ID) {
      alert("Google Client ID が設定されていません。\n.env.development.local に VITE_GOOGLE_CLIENT_ID を設定してください。");
      return;
    }
    window.location.href = buildGoogleOAuthUrl();
  };

  return (
    <div className="w-full min-h-full flex flex-col items-center justify-center px-6">
      <div className="w-full max-w-sm flex flex-col gap-10">
        {/* Logo / App name */}
        <div className="text-center">
          <h1 className="text-3xl font-bold tracking-tight">Reverie</h1>
          <p className="text-sm text-muted-foreground mt-2">
            ログインまたは新規登録
          </p>
        </div>

        {/* Social login buttons */}
        <div className="flex flex-col gap-3">
          <Button
            variant="outline"
            className="w-full h-11 gap-3 text-sm font-medium"
            onClick={handleGoogleLogin}
          >
            <GoogleIcon className="size-5" />
            Googleで続ける
          </Button>

          <Button
            variant="outline"
            className="w-full h-11 gap-3 text-sm font-medium"
            disabled
          >
            <LineIcon className="size-5" />
            LINEで続ける（準備中）
          </Button>
        </div>

        {/* Dev login (development only) */}
        {IS_DEV && (
          <div className="border rounded-lg p-4 space-y-2 bg-muted/30">
            <p className="text-xs font-medium text-muted-foreground">開発用ログイン</p>
            <p className="text-xs text-muted-foreground">
              <code>go run ./cmd/seed</code> で取得したトークンを入力
            </p>
            <Input
              placeholder="ACCESS_TOKEN=..."
              value={devToken}
              onChange={(e) => setDevToken(e.target.value)}
              className="h-8 text-xs font-mono"
            />
            <Button
              size="sm"
              className="w-full h-8 text-xs"
              onClick={() => {
                const token = devToken.replace(/^ACCESS_TOKEN=/, "").trim();
                if (token) {
                  setTokens(token, "");
                  navigate("/");
                }
              }}
            >
              トークンでログイン
            </Button>
          </div>
        )}

        {/* Terms */}
        <p className="text-center text-xs text-muted-foreground leading-relaxed">
          続行することで、
          <button
            type="button"
            className="underline underline-offset-2 hover:text-foreground"
          >
            利用規約
          </button>
          および
          <button
            type="button"
            className="underline underline-offset-2 hover:text-foreground"
          >
            プライバシーポリシー
          </button>
          に同意したことになります。
        </p>
      </div>
    </div>
  );
}
