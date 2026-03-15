import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("search", "routes/search.tsx"),
  route("chat", "routes/chat.tsx"),
  route("chat/requests", "routes/chat-requests.tsx"),
  route("chat/:id", "routes/chat-detail.tsx"),
  route("notifications", "routes/notifications.tsx"),
  route("posts/:id", "routes/post.tsx"),
  route("posts/:id/likes", "routes/post-likes.tsx"),
  route("posts/:id/reposts", "routes/post-reposts.tsx"),
  route("users/:id", "routes/user.tsx"),
  route("users/:id/connections", "routes/user-connections.tsx"),
  route("settings/profile", "routes/profile-edit.tsx"),
  route("calls", "routes/calls.tsx"),
] satisfies RouteConfig;
