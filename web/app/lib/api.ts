const API_BASE = "http://localhost:50051";

// Token management
const ACCESS_TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";

export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function setTokens(accessToken: string, refreshToken: string) {
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
}

export function clearTokens() {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
}

export function isLoggedIn(): boolean {
  return !!getAccessToken();
}

// Fetch wrapper
async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getAccessToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ message: res.statusText }));
    throw new ApiError(res.status, err.message ?? res.statusText);
  }

  return res.json();
}

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message);
  }
}

// Types
export interface User {
  id: string;
  custom_id: string;
  display_name: string;
  biography?: string;
  is_private: boolean;
  is_me: boolean;
  create_time: string;
}

export interface Post {
  id: string;
  text: string;
  author: User;
  reply_to_id?: string;
  repost_id?: string;
  reply_count: number;
  repost_count: number;
  favorite_count: number;
  is_favorited: boolean;
  create_time: string;
}

// Account API
export async function socialLogin(
  provider: "GOOGLE" | "LINE",
  code: string
): Promise<{ token_pair: { access_token: string; refresh_token: string }; is_new_account: boolean }> {
  return request("/v1/account:socialLogin", {
    method: "POST",
    body: JSON.stringify({ provider: provider === "GOOGLE" ? 1 : 2, code }),
  });
}

// User API
export async function getUser(userId: string): Promise<{ user: User }> {
  return request(`/v1/users/${userId}`);
}

export interface Account {
  id: string;
  custom_id: string;
  display_name: string;
  create_time: string;
}

export async function getAccount(): Promise<{ account: Account }> {
  return request("/v1/account");
}

export async function getMe(): Promise<{ user: User }> {
  const { account } = await getAccount();
  return getUser(account.id);
}

export async function updateUser(params: {
  display_name?: string;
  biography?: string;
  is_private?: boolean;
}): Promise<{ user: User }> {
  return request("/v1/users/me", {
    method: "PATCH",
    body: JSON.stringify({ user: params }),
  });
}

// Post API
export async function createPost(params: {
  text: string;
  reply_to_id?: string;
  repost_id?: string;
}): Promise<{ post: Post }> {
  return request("/v1/posts", {
    method: "POST",
    body: JSON.stringify(params),
  });
}

export async function deletePost(postId: string): Promise<void> {
  return request(`/v1/posts/${postId}`, { method: "DELETE" });
}

export async function listTimeline(params?: {
  limit?: number;
  cursor?: string;
}): Promise<{ posts: Post[]; next_cursor?: string }> {
  const qs = new URLSearchParams();
  if (params?.limit) qs.set("limit", String(params.limit));
  if (params?.cursor) qs.set("cursor", params.cursor);
  const query = qs.toString() ? `?${qs}` : "";
  return request(`/v1/timeline${query}`);
}

export async function listUserPosts(
  userId: string,
  params?: { limit?: number; cursor?: string }
): Promise<{ posts: Post[]; next_cursor?: string }> {
  const qs = new URLSearchParams();
  if (params?.limit) qs.set("limit", String(params.limit));
  if (params?.cursor) qs.set("cursor", params.cursor);
  const query = qs.toString() ? `?${qs}` : "";
  return request(`/v1/users/${userId}/posts${query}`);
}
