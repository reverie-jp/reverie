import {
  createPromiseClient,
  ConnectError,
  Code,
  type Interceptor,
} from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { AccountService } from "./gen/account/v1/account_connect";
import { CallService } from "./gen/call/v1/call_connect";
import { FollowService } from "./gen/follow/v1/follow_connect";
import { UserService } from "./gen/user/v1/user_connect";

const ACCESS_TOKEN_KEY = "reverie.access_token";
const REFRESH_TOKEN_KEY = "reverie.refresh_token";

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
    localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
    localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
  },
  clear() {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
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
const bareAccountClient = createPromiseClient(AccountService, bareTransport);

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

export const accountClient = createPromiseClient(AccountService, transport);
export const userClient = createPromiseClient(UserService, transport);
export const callClient = createPromiseClient(CallService, transport);
export const followClient = createPromiseClient(FollowService, transport);
