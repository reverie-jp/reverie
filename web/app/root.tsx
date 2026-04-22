import {
  isRouteErrorResponse,
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
} from "react-router";

import type { Route } from "./+types/root";
import "./app.css";
import { CallProvider } from "~/lib/call-context";
import { NotificationProvider } from "~/lib/notification-context";
import { usePresenceHeartbeat } from "~/lib/use-presence-heartbeat";
import { CallHeader } from "~/components/call-header";
import { AppFooter } from "~/components/app-footer";
import { FloatingCreateButton } from "~/components/floating-create-button";
import { Toaster } from "~/components/ui/sonner";

export const links: Route.LinksFunction = () => [
  { rel: "preconnect", href: "https://fonts.googleapis.com" },
  {
    rel: "preconnect",
    href: "https://fonts.gstatic.com",
    crossOrigin: "anonymous",
  },
  {
    rel: "stylesheet",
    href: "https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap",
  },
];

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="dark">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body>
        <div className="w-screen h-screen flex flex-col overflow-hidden">
          {children}
        </div>
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

function GlobalEffects() {
  usePresenceHeartbeat();
  return null;
}

export default function App() {
  return (
    <NotificationProvider>
      <CallProvider>
        <GlobalEffects />
        {/* Scroll container holds the mini call bar as a sticky overlay so
            the bar's backdrop-filter blur reacts to the content scrolling
            behind it, instead of sitting on a flat row. */}
        <div className="flex-1 overflow-x-hidden overflow-y-auto">
          <CallHeader />
          <Outlet />
        </div>
        <AppFooter />
        <FloatingCreateButton />
        <Toaster position="top-right" />
      </CallProvider>
    </NotificationProvider>
  );
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
  let message = "Oops!";
  let details = "An unexpected error occurred.";
  let stack: string | undefined;

  if (isRouteErrorResponse(error)) {
    message = error.status === 404 ? "404" : "Error";
    details =
      error.status === 404
        ? "The requested page could not be found."
        : error.statusText || details;
  } else if (import.meta.env.DEV && error && error instanceof Error) {
    details = error.message;
    stack = error.stack;
  }

  return (
    <main className="pt-16 p-4 container mx-auto">
      <h1>{message}</h1>
      <p>{details}</p>
      {stack && (
        <pre className="w-full p-4 overflow-x-auto">
          <code>{stack}</code>
        </pre>
      )}
    </main>
  );
}
