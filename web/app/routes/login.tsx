import { redirect } from "react-router";
import { tokenStore } from "~/lib/api-client";
import { GoogleLoginButton } from "~/components/google-login-button";
import type { Route } from "./+types/login";

export const clientLoader = ({ request }: Route.ClientLoaderArgs) => {
  const url = new URL(request.url);
  const returnTo = url.searchParams.get("returnTo");
  if (tokenStore.getAccessToken()) {
    return redirect(returnTo && returnTo.startsWith("/") ? returnTo : "/");
  }
  return null;
};

export default function Login() {
  const returnTo =
    typeof window !== "undefined"
      ? new URLSearchParams(window.location.search).get("returnTo") ?? undefined
      : undefined;

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
          <GoogleLoginButton returnTo={returnTo} />
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
