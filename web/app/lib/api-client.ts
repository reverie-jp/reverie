import {
  createClient,
  ConnectError,
  Code,
  type Interceptor,
} from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { AccountService } from "./gen/account/v1/account_connect";
import { CallService } from "./gen/call/v1/call_connect";
import { EventService } from "./gen/event/v1/event_connect";
import { FollowService } from "./gen/follow/v1/follow_connect";
import { NotificationService } from "./gen/notification/v1/notification_connect";
import { PresenceService } from "./gen/presence/v1/presence_connect";
import { UserService } from "./gen/user/v1/user_connect";

const ACCESS_TOKEN_KEY = "reverie.access_token";
const REFRESH_TOKEN_KEY = "reverie.refresh_token";

// AUTH_CHANGED_EVENT lets components outside React's tree (NotificationProvider,
// AppHeader) react to login/logout without threading a full auth context.
export const AUTH_CHANGED_EVENT = "reverie:auth_changed";

function fireAuthChanged() {
  if (typeof window === "undefined") return;
  window.dispatchEvent(new CustomEvent(AUTH_CHANGED_EVENT));
}

export const tokenStore = {
  getAccessToken(): string | null {
    if (typeof window === "undefined") return null;
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  },
  getRefreshToken(): string | null {
    if (typeof window === "undefined") return null;
    return localStorage.getItem(REFRESH_TOKEN_KEY);
  },
  setTokens(accessToken: string, refreshToken: string) {
    const wasAuthed = Boolean(localStorage.getItem(ACCESS_TOKEN_KEY));
    localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
    localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
    // Only fire on login transitions, not silent refresh — otherwise the
    // stream would reconnect on every token rotation.
    if (!wasAuthed) fireAuthChanged();
  },
  clear() {
    const wasAuthed = Boolean(localStorage.getItem(ACCESS_TOKEN_KEY));
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    if (wasAuthed) fireAuthChanged();
  },
};

export const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? "http://localhost:50051";

// Procedures that don't need Authorization and must not trigger the 401-retry
// path (otherwise a failed SocialLogin would try to refresh with no refresh
// token and redirect to /login from the callback page).
const PUBLIC_PROCEDURES = new Set<string>([
  `${AccountService.typeName}/SocialLogin`,
  `${AccountService.typeName}/RefreshToken`,
]);

// Bare transport (no interceptors) used inside the auth interceptor for the
// refresh call, so the refresh itself can never re-enter the interceptor.
const bareTransport = createConnectTransport({ baseUrl: API_BASE_URL });
const bareAccountClient = createClient(AccountService, bareTransport);

// In-flight refresh promise — concurrent 401s share a single refresh call.
let refreshInFlight: Promise<boolean> | null = null;

async function refreshTokens(): Promise<boolean> {
  if (refreshInFlight) return refreshInFlight;
  refreshInFlight = (async () => {
    const refreshToken = tokenStore.getRefreshToken();
    if (!refreshToken) return false;
    try {
      const res = await bareAccountClient.refreshToken({ refreshToken });
      if (!res.tokenPair) return false;
      tokenStore.setTokens(
        res.tokenPair.accessToken,
        res.tokenPair.refreshToken,
      );
      return true;
    } catch {
      return false;
    }
  })().finally(() => {
    refreshInFlight = null;
  });
  return refreshInFlight;
}

function redirectToLogin() {
  tokenStore.clear();
  if (typeof window === "undefined") return;
  const { pathname, search } = window.location;
  if (pathname === "/login" || pathname === "/auth/callback") return;
  const returnTo = `${pathname}${search}`;
  window.location.href = `/login?returnTo=${encodeURIComponent(returnTo)}`;
}

const authInterceptor: Interceptor = (next) => async (req) => {
  const procedure = `${req.service.typeName}/${req.method.name}`;
  if (PUBLIC_PROCEDURES.has(procedure)) {
    return next(req);
  }

  const token = tokenStore.getAccessToken();
  if (token) {
    req.header.set("Authorization", `Bearer ${token}`);
  }

  try {
    return await next(req);
  } catch (err) {
    if (!(err instanceof ConnectError) || err.code !== Code.Unauthenticated) {
      throw err;
    }
    const refreshed = await refreshTokens();
    if (!refreshed) {
      redirectToLogin();
      throw err;
    }
    const newToken = tokenStore.getAccessToken();
    if (newToken) {
      req.header.set("Authorization", `Bearer ${newToken}`);
    }
    return await next(req);
  }
};

const transport = createConnectTransport({
  baseUrl: API_BASE_URL,
  interceptors: [authInterceptor],
});

export const accountClient = createClient(AccountService, transport);
export const userClient = createClient(UserService, transport);
export const callClient = createClient(CallService, transport);
export const followClient = createClient(FollowService, transport);
export const notificationClient = createClient(
  NotificationService,
  transport,
);
export const presenceClient = createClient(PresenceService, transport);
export const eventClient = createClient(EventService, transport);
