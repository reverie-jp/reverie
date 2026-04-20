import { redirect } from "react-router";
import { Button } from "~/components/ui/button";
import { tokenStore } from "~/lib/api-client";
import type { Route } from "./+types/login";

export const clientLoader = (_args: Route.ClientLoaderArgs) => {
  if (tokenStore.getAccessToken()) {
    return redirect("/");
  }
  return null;
};

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

export default function Login() {
  const handleGoogleLogin = () => {
    const clientId = import.meta.env.VITE_GOOGLE_CLIENT_ID;
    const redirectUri = import.meta.env.VITE_GOOGLE_REDIRECT_URI;
    if (!clientId || !redirectUri) {
      console.error(
        "VITE_GOOGLE_CLIENT_ID / VITE_GOOGLE_REDIRECT_URI are not set",
      );
      return;
    }
    const params = new URLSearchParams({
      client_id: clientId,
      redirect_uri: redirectUri,
      response_type: "code",
      scope: "openid email profile",
      prompt: "select_account",
    });
    window.location.href = `https://accounts.google.com/o/oauth2/v2/auth?${params}`;
  };

  return (
    <div className="w-full min-h-full flex flex-col items-center justify-center px-6">
      <div className="w-full max-w-sm flex flex-col gap-10">
        <div className="text-center">
          <h1 className="text-3xl font-bold tracking-tight">Reverie</h1>
          <p className="text-sm text-muted-foreground mt-2">
            ログインまたは新規登録
          </p>
        </div>

        <div className="flex flex-col gap-3">
          <Button
            variant="outline"
            className="w-full h-11 gap-3 text-sm font-medium"
            onClick={handleGoogleLogin}
          >
            <GoogleIcon className="size-5" />
            Googleで続ける
          </Button>
        </div>

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
