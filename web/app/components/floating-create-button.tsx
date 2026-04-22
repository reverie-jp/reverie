import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router";
import { Plus } from "lucide-react";
import { tokenStore } from "~/lib/api-client";
import { CreateCallDialog } from "~/components/create-call-dialog";

export function FloatingCreateButton() {
  const location = useLocation();
  const navigate = useNavigate();
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [createOpen, setCreateOpen] = useState(false);

  useEffect(() => {
    const sync = () =>
      setIsAuthenticated(Boolean(tokenStore.getAccessToken()));
    sync();
    window.addEventListener("storage", sync);
    window.addEventListener("reverie:auth_changed", sync);
    return () => {
      window.removeEventListener("storage", sync);
      window.removeEventListener("reverie:auth_changed", sync);
    };
  }, []);

  // Same exclusion rules as AppFooter — the FAB belongs to the same
  // "main browse" surface.
  if (
    location.pathname === "/login" ||
    location.pathname === "/auth/callback" ||
    location.pathname.startsWith("/calls/")
  ) {
    return null;
  }

  const handleClick = () => {
    if (!isAuthenticated) {
      navigate("/login");
      return;
    }
    setCreateOpen(true);
  };

  return (
    <>
      <button
        type="button"
        onClick={handleClick}
        aria-label="新しい通話を作成"
        // Sits above AppFooter (~64px tall) with breathing room.
        className="fixed bottom-24 right-5 z-30 grid place-items-center size-14 rounded-2xl text-[#120a2e] font-bold transition-transform active:scale-95"
        style={{
          background: "linear-gradient(160deg, #c9b5ff, var(--reverie-accent))",
          boxShadow:
            "0 10px 28px rgba(0,0,0,0.4), 0 0 26px var(--reverie-accent-glow)",
        }}
      >
        <Plus className="size-6" strokeWidth={2.5} />
      </button>

      <CreateCallDialog open={createOpen} onOpenChange={setCreateOpen} />
    </>
  );
}
