const API_BASE = "http://localhost:50051";

// Token management
const ACCESS_TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";

const isBrowser = typeof window !== "undefined";

export function getAccessToken(): string | null {
  if (!isBrowser) return null;
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function setTokens(accessToken: string, refreshToken: string) {
  if (!isBrowser) return;
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
}

export function clearTokens() {
  if (!isBrowser) return;
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
  customId: string;
  displayName: string;
  biography?: string;
  isPrivate: boolean;
  isMe: boolean;
  isFollowing: boolean;
  isFollowedBy: boolean;
  followerCount: number;
  followingCount: number;
  createTime: string;
}

export interface Post {
  id: string;
  text: string;
  author: User;
  replyToId?: string;
  repostId?: string;
  replyCount: number;
  repostCount: number;
  likeCount: number;
  isLiked: boolean;
  createTime: string;
  repostOf?: Post;
}

// Account API
export async function socialLogin(
  provider: "GOOGLE" | "LINE",
  code: string
): Promise<{ tokenPair: { accessToken: string; refreshToken: string }; isNewAccount: boolean }> {
  return request("/v1/account:socialLogin", {
    method: "POST",
    body: JSON.stringify({ provider: provider === "GOOGLE" ? 1 : 2, code }),
  });
}

export interface Account {
  id: string;
  customId: string;
  displayName: string;
  createTime: string;
}

export async function getAccount(): Promise<{ account: Account }> {
  return request("/v1/account");
}

// User API
export async function getUser(userId: string): Promise<{ user: User }> {
  return request(`/v1/users/${userId}`);
}

export async function searchUsers(
  query: string,
  params?: { pageSize?: number; pageToken?: string }
): Promise<{ users: User[]; nextPageToken?: string }> {
  const qs = new URLSearchParams({ query });
  if (params?.pageSize) qs.set("page_size", String(params.pageSize));
  if (params?.pageToken) qs.set("page_token", params.pageToken);
  return request(`/v1/users?${qs}`);
}

export async function getMe(): Promise<{ user: User }> {
  const { account } = await getAccount();
  return getUser(account.id);
}

export async function followUser(userId: string): Promise<{ user: User }> {
  return request(`/v1/users/${userId}:follow`, { method: "POST", body: "{}" });
}

export async function unfollowUser(userId: string): Promise<{ user: User }> {
  return request(`/v1/users/${userId}:unfollow`, { method: "POST", body: "{}" });
}

export async function updateUser(params: {
  displayName?: string;
  biography?: string;
  isPrivate?: boolean;
}): Promise<{ user: User }> {
  return request("/v1/users/me", {
    method: "PATCH",
    body: JSON.stringify({ user: params }),
  });
}

// Post API
export async function getPost(postId: string): Promise<{ post: Post }> {
  return request(`/v1/posts/${postId}`);
}

export async function createPost(params: {
  text: string;
  replyToId?: string;
  repostId?: string;
}): Promise<{ post: Post }> {
  return request("/v1/posts", {
    method: "POST",
    body: JSON.stringify(params),
  });
}

export async function deletePost(postId: string): Promise<void> {
  return request(`/v1/posts/${postId}`, { method: "DELETE" });
}

export async function likePost(postId: string): Promise<{ post: Post }> {
  return request(`/v1/posts/${postId}:like`, {
    method: "POST",
    body: "{}",
  });
}

export async function unlikePost(postId: string): Promise<{ post: Post }> {
  return request(`/v1/posts/${postId}:unlike`, {
    method: "POST",
    body: "{}",
  });
}

export async function listPostReposts(
  postId: string,
  params?: { pageSize?: number; pageToken?: string }
): Promise<{ posts: Post[]; nextPageToken?: string }> {
  const qs = new URLSearchParams();
  if (params?.pageSize) qs.set("page_size", String(params.pageSize));
  if (params?.pageToken) qs.set("page_token", params.pageToken);
  const query = qs.toString() ? `?${qs}` : "";
  return request(`/v1/posts/${postId}/reposts${query}`);
}

export async function listPostReplies(
  postId: string,
  params?: { pageSize?: number; pageToken?: string }
): Promise<{ posts: Post[]; nextPageToken?: string }> {
  const qs = new URLSearchParams();
  if (params?.pageSize) qs.set("page_size", String(params.pageSize));
  if (params?.pageToken) qs.set("page_token", params.pageToken);
  const query = qs.toString() ? `?${qs}` : "";
  return request(`/v1/posts/${postId}/replies${query}`);
}

export async function listUserPosts(
  userId: string,
  params?: { pageSize?: number; pageToken?: string }
): Promise<{ posts: Post[]; nextPageToken?: string }> {
  const qs = new URLSearchParams();
  if (params?.pageSize) qs.set("page_size", String(params.pageSize));
  if (params?.pageToken) qs.set("page_token", params.pageToken);
  const query = qs.toString() ? `?${qs}` : "";
  return request(`/v1/users/${userId}/posts${query}`);
}

// Timeline API
export async function listPublicTimeline(params?: {
  pageSize?: number;
  pageToken?: string;
}): Promise<{ posts: Post[]; nextPageToken?: string }> {
  const qs = new URLSearchParams();
  if (params?.pageSize) qs.set("page_size", String(params.pageSize));
  if (params?.pageToken) qs.set("page_token", params.pageToken);
  const query = qs.toString() ? `?${qs}` : "";
  return request(`/v1/timeline/public${query}`);
}

export async function listFollowingTimeline(params?: {
  pageSize?: number;
  pageToken?: string;
}): Promise<{ posts: Post[]; nextPageToken?: string }> {
  const qs = new URLSearchParams();
  if (params?.pageSize) qs.set("page_size", String(params.pageSize));
  if (params?.pageToken) qs.set("page_token", params.pageToken);
  const query = qs.toString() ? `?${qs}` : "";
  return request(`/v1/timeline/following${query}`);
}
