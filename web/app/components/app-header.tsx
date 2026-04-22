import { Link } from "react-router";
import { NotificationBell } from "~/components/notification-bell";
import { tokenStore } from "~/lib/api-client";
import { useEffect, useState } from "react";

export function AppHeader() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

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

  return (
    <header className="shrink-0 h-12 px-3 border-b border-foreground/10 flex items-center justify-between">
      <Link to="/" className="text-sm font-semibold tracking-wide">
        reverie
      </Link>
      {isAuthenticated && (
        <div className="flex items-center gap-1">
          <NotificationBell />
        </div>
      )}
    </header>
  );
}
