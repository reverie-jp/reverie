import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("search", "routes/search.tsx"),
  route("chat", "routes/chat.tsx"),
  route("notifications", "routes/notifications.tsx"),
  route("users/:id", "routes/user.tsx"),
  route("calls", "routes/calls.tsx"),
] satisfies RouteConfig;
