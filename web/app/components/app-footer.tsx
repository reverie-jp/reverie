import { useEffect, useState } from "react";
import { Link, useLocation } from "react-router";
import { Home, User as UserIcon } from "lucide-react";
import { tokenStore } from "~/lib/api-client";
import { cn } from "~/lib/utils";
import { NotificationNavItem } from "~/components/notification-nav-item";

function NavItem({
  to,
  icon: Icon,
  label,
  active,
}: {
  to: string;
  icon: typeof Home;
  label: string;
  active: boolean;
}) {
  return (
    <Link
      to={to}
      className={cn(
        "flex flex-col items-center gap-1 px-3 py-1 transition-colors",
        active
          ? "text-foreground"
          : "text-muted-foreground/70 hover:text-foreground",
      )}
    >
      <span
        className={cn(
          "grid place-items-center w-10 h-7 rounded-xl transition-all",
          active && "bg-[var(--reverie-accent)]/15 text-[var(--reverie-accent)]",
        )}
      >
        <Icon className="size-4" />
      </span>
      <span className="text-[10px] leading-none">{label}</span>
    </Link>
  );
}

export function AppFooter() {
  const location = useLocation();
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

  // Hide on pages where a persistent footer gets in the way.
  if (
    location.pathname === "/login" ||
    location.pathname === "/auth/callback" ||
    location.pathname.startsWith("/calls/")
  ) {
    return null;
  }

  const onHome = location.pathname === "/";
  const onMe =
    location.pathname === "/me" || location.pathname.startsWith("/me/");

  return (
    <nav
      aria-label="Global navigation"
      className="shrink-0 relative border-t border-white/10 backdrop-blur-xl bg-black/30"
    >
      <div className="mx-auto flex max-w-xl items-end justify-around px-6 pt-2 pb-3">
        <NavItem to="/" icon={Home} label="ホーム" active={onHome} />
        {isAuthenticated && <NotificationNavItem />}
        <NavItem
          to={isAuthenticated ? "/me" : "/login"}
          icon={UserIcon}
          label={isAuthenticated ? "自分" : "ログイン"}
          active={onMe}
        />
      </div>
    </nav>
  );
}
