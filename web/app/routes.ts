import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("login", "routes/login.tsx"),
  route("auth/callback", "routes/auth-callback.tsx"),
  route("me", "routes/me.tsx"),
  route("calls/:callId", "routes/calls.room.tsx"),
  // React Router v7 can't parse literal-prefix params like `@:customId`
  // (the param must follow `/`), so we catch any top-level handle and
  // validate the `@` prefix inside the component.
  route(":handle", "routes/user.tsx"),
] satisfies RouteConfig;
