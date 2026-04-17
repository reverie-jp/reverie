import { NavLink } from "react-router";
import { Home, Search, MessageCircle, Bell, CircleUserRound, Settings } from "lucide-react";

export function BottomNav() {
  return (
    <nav className="sticky bottom-0 left-0 w-full h-16 border-t bg-background">
      <div className="flex items-center justify-around h-full">
        <NavLink to="/" className={navLinkClass}>
          <Home className="size-6" />
        </NavLink>
        <NavLink to="/search" className={navLinkClass}>
          <Search className="size-6" />
        </NavLink>
        <NavLink to="/chat" className={navLinkClass}>
          <MessageCircle className="size-6" />
        </NavLink>
        <NavLink to="/notifications" className={navLinkClass}>
          <Bell className="size-6" />
        </NavLink>
        <NavLink to="/users/me" className={navLinkClass}>
          <CircleUserRound className="size-6" />
        </NavLink>
        <NavLink to="/settings" className={navLinkClass}>
          <Settings className="size-6" />
        </NavLink>
      </div>
    </nav>
  );
}

function navLinkClass({ isActive }: { isActive: boolean }) {
  return `flex items-center justify-center p-2 transition-colors ${
    isActive ? "text-foreground" : "text-muted-foreground hover:text-foreground"
  }`;
}
