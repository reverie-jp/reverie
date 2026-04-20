import { useEffect } from "react";
import { useNavigate } from "react-router";
import { userClient } from "~/lib/api-client";

export default function Me() {
  const navigate = useNavigate();

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const res = await userClient.getMyUser({});
        if (cancelled) return;
        if (!res.user) {
          navigate("/login", { replace: true });
          return;
        }
        navigate(`/@${res.user.customId}`, { replace: true });
      } catch {
        if (cancelled) return;
        navigate("/login", { replace: true });
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [navigate]);

  return (
    <div className="w-full min-h-full flex items-center justify-center">
      <p className="text-sm text-muted-foreground">読み込み中...</p>
    </div>
  );
}
