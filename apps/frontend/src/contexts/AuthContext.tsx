import { authnService } from "@/api/authn.service";
import { setAccessToken } from "@/api/client";
import { jwtDecode } from "jwt-decode";
import { createContext, useContext, useEffect, useRef, useState, type ReactNode } from "react";

type AuthenticatedUser = {
  authenticated: true;
  userId: string;
  sessionId: string;
  name: string;

  jwt: string;
};

type UIUser =
  | AuthenticatedUser
  | {
      authenticated: false;
    };

type AuthContextValue = {
  user: UIUser | null; // null means we haven't checked yet
  setUserFromAuth: (user: Omit<AuthenticatedUser, "userId"> | { authenticated: false }) => void;
  logoutFromAuth: () => void;
};

export const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UIUser | null>(null);
  const userRef = useRef<UIUser | null>(null);

  const updateUser = (newUser: UIUser | null) => {
    userRef.current = newUser;
    setUser(newUser);
  };

  const setUserFromAuth = (
    authUser: Omit<AuthenticatedUser, "userId"> | { authenticated: false },
  ) => {
    if (authUser.authenticated === false) {
      updateUser({
        authenticated: false,
      });
      setAccessToken(null);
      localStorage.removeItem("user");
      return;
    }

    const decoded = jwtDecode(authUser.jwt);

    const user: UIUser = {
      userId: (decoded as { sub?: string }).sub ?? "",
      ...authUser,
    };

    updateUser(user);
    setAccessToken(user.jwt);

    const { jwt, ...userWithoutToken } = user;
    localStorage.setItem("user", JSON.stringify(userWithoutToken));
  };

  const setUserFromRefresh = (user: Omit<AuthenticatedUser, "jwt">, newJWT: string) => {
    const updatedUser = { ...user, jwt: newJWT };

    updateUser(updatedUser);
    setAccessToken(newJWT);
  };

  useEffect(() => {
    const storedUser = localStorage.getItem("user");

    if (storedUser) {
      const user = JSON.parse(storedUser) as Omit<AuthenticatedUser, "jwt">;

      authnService
        .RefreshJWT({
          userId: user.userId,
          sessionId: user.sessionId,
        })
        .then((res) => {
          if (!res.ok) {
            setUserFromAuth({ authenticated: false });
            return;
          }

          setUserFromRefresh(user, res.jwt);
        });
    } else {
      setUserFromAuth({ authenticated: false });
    }
  }, []);

  useEffect(() => {
    const tokenRefresher = setInterval(
      async () => {
        const currentUser = userRef.current;

        if (!currentUser?.authenticated) return;

        const res = await authnService.RefreshJWT({
          userId: currentUser.userId,
          sessionId: currentUser.sessionId,
        });

        if (!res.ok) {
          setUserFromAuth({ authenticated: false });
          return;
        }

        setUserFromRefresh(currentUser, res.jwt);
      },
      10 * 60 * 1000,
    );

    return () => {
      clearInterval(tokenRefresher);
    };
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        setUserFromAuth,
        logoutFromAuth: () => setUserFromAuth({ authenticated: false }),
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }

  return context;
}
