import { useEffect, useRef, useState } from "react";
import { useNavigate, useSearchParams } from "react-router";
import { accountClient, tokenStore } from "~/lib/api-client";
import { AuthProvider } from "~/lib/gen/account/v1/account_pb";

export default function AuthCallback() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const triggered = useRef(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (triggered.current) return;
    triggered.current = true;

    const code = searchParams.get("code");
    const oauthError = searchParams.get("error");

    if (oauthError) {
      setError(`Google 認証が拒否されました: ${oauthError}`);
      return;
    }
    if (!code) {
      navigate("/login", { replace: true });
      return;
    }

    accountClient
      .socialLogin({ provider: AuthProvider.GOOGLE, code })
      .then((res) => {
        if (!res.tokenPair) {
          setError("トークンの取得に失敗しました");
          return;
        }
        tokenStore.setTokens(
          res.tokenPair.accessToken,
          res.tokenPair.refreshToken,
        );
        navigate("/", { replace: true });
      })
      .catch((err) => {
        console.error("SocialLogin failed:", err);
        setError(err instanceof Error ? err.message : String(err));
      });
  }, [searchParams, navigate]);

  return (
    <div className="w-full min-h-full flex items-center justify-center px-6">
      {error ? (
        <div className="max-w-sm text-center flex flex-col gap-4">
          <p className="text-sm text-destructive">ログインに失敗しました</p>
          <p className="text-xs text-muted-foreground wrap-break-word">
            {error}
          </p>
          <button
            type="button"
            onClick={() => navigate("/login", { replace: true })}
            className="text-sm underline"
          >
            ログイン画面に戻る
          </button>
        </div>
      ) : (
        <p className="text-sm text-muted-foreground">認証中...</p>
      )}
    </div>
  );
}
