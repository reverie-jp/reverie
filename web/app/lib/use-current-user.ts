import { useState, useEffect } from "react";
import { getMe, isLoggedIn, type User } from "~/lib/api";

export function useCurrentUser(): User | null {
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    if (!isLoggedIn()) return;
    getMe()
      .then(({ user }) => setUser(user))
      .catch(() => {});
  }, []);

  return user;
}
